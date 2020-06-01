package sshdservice

import (
	"context"
	"reflect"

	sshdoperatorv1alpha1 "github.com/sousatou/sshd-operator/pkg/apis/sshdoperator/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
        "time"
        "os/exec"
	"math/rand"
        "strings"
)

var log = logf.Log.WithName("controller_sshdservice")

/**
* USER ACTION REQUIRED: This is a scaffold file intended for the user to modify with their own Controller
* business logic.  Delete these comments after modifying this file.*
 */

// Add creates a new SshdService Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileSshdService{client: mgr.GetClient(), scheme: mgr.GetScheme()}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("sshdservice-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to primary resource SshdService
	err = c.Watch(&source.Kind{Type: &sshdoperatorv1alpha1.SshdService{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	// TODO(user): Modify this to be the types you create that are owned by the primary resource
	// Watch for changes to secondary resource Pods and requeue the owner SshdService
	err = c.Watch(&source.Kind{Type: &corev1.Pod{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &sshdoperatorv1alpha1.SshdService{},
	})
	if err != nil {
		return err
	}

	err = c.Watch(&source.Kind{Type: &corev1.Service{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &sshdoperatorv1alpha1.SshdService{},
	})
	if err != nil {
		return err
	}

	return nil
}

// blank assignment to verify that ReconcileSshdService implements reconcile.Reconciler
var _ reconcile.Reconciler = &ReconcileSshdService{}

// ReconcileSshdService reconciles a SshdService object
type ReconcileSshdService struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client client.Client
	scheme *runtime.Scheme
}

// Reconcile reads that state of the cluster for a SshdService object and makes changes based on the state read
// and what is in the SshdService.Spec
// TODO(user): Modify this Reconcile function to implement your Controller logic.  This example creates
// a Pod as an example
// Note:
// The Controller will requeue the Request to be processed again if the returned error is non-nil or
// Result.Requeue is true, otherwise upon completion it will remove the work from the queue.
func (r *ReconcileSshdService) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Reconciling SshdService")

	// Fetch the SshdService instance
	instance := &sshdoperatorv1alpha1.SshdService{}
	err := r.client.Get(context.TODO(), request.NamespacedName, instance)
	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return reconcile.Result{}, err
	}

	// Define a new Pod object
	pod := newPodForCR(instance)

	// Set SshdService instance as the owner and controller
	if err := controllerutil.SetControllerReference(instance, pod, r.scheme); err != nil {
		return reconcile.Result{}, err
	}

	// Check if this Pod already exists
	found := &corev1.Pod{}
	err = r.client.Get(context.TODO(), types.NamespacedName{Name: pod.Name, Namespace: pod.Namespace}, found)
	if err != nil && errors.IsNotFound(err) {
		r.updateStage(instance, "INITIALIZING")

		reqLogger.Info("Creating a new Pod", "Pod.Namespace", pod.Namespace, "Pod.Name", pod.Name)
		err = r.client.Create(context.TODO(), pod)
		if err != nil {
			return reconcile.Result{}, err
		}

		reqLogger.Info("Creating a new Service", "Pod.Namespace", pod.Namespace, "Pod.Name", pod.Name)
		err = r.client.Create(context.TODO(), pod)
		if err != nil {
			return reconcile.Result{}, err
		}

                return reconcile.Result{RequeueAfter: time.Second*5}, nil
	} else if err != nil {
		return reconcile.Result{}, err
	}

  // For service
  svc := newServiceForCR(instance)
  if err := controllerutil.SetControllerReference(instance, svc, r.scheme); err != nil {
    return reconcile.Result{}, err
  }
  found_svc := &corev1.Service{}
  err = r.client.Get(context.TODO(), types.NamespacedName{Name: svc.Name, Namespace: svc.Namespace}, found_svc)
  if err != nil && errors.IsNotFound(err) {
    reqLogger.Info("Creating a new Service", "Svc.Namespace", svc.Namespace, "Svc.Name", svc.Name)
    err = r.client.Create(context.TODO(), svc)
    if err != nil {
      return reconcile.Result{}, err
    }
    return reconcile.Result{RequeueAfter: time.Second*5}, nil
  } else if err != nil {
    return reconcile.Result{}, err
  }
  // In next iteration
  nodeport := found_svc.Spec.Ports[0].NodePort
  r.updateNodePort(instance, int(nodeport))
  
  // Do action periodically
  reqLogger.Info("Executing Command")
  username := instance.Spec.UserName
  if username == "" {
    username = "user1"
  }
  password := instance.Status.Password
  if password == "" {
    password = generatePassword()
    r.updatePassword(instance, password)
  }
  out, err := exec.Command("pod_init", pod.Name, username, password).Output()
  r.updateStage(instance, strings.TrimRight(string(out), "\n"))

  return reconcile.Result{RequeueAfter: time.Second*5}, nil
}

// newPodForCR returns a busybox pod with the same name/namespace as the cr
func newPodForCR(cr *sshdoperatorv1alpha1.SshdService) *corev1.Pod {
	labels := map[string]string{
		"app": "sshd",
		"cr": cr.Name,
	}
	return &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cr.Name,
			Namespace: cr.Namespace,
			Labels:    labels,
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name:    "main",
					Image:   "fedora",
					Command: []string{"sleep", "infinity"},
				},
			},
		},
	}
}

func newServiceForCR(cr *sshdoperatorv1alpha1.SshdService) *corev1.Service {
  labels := map[string]string{
    "app": "sshd",
    "cr": cr.Name,
  }
  return &corev1.Service{
    ObjectMeta: metav1.ObjectMeta{
      Name:      cr.Name,
      Namespace: cr.Namespace,
      Labels:    labels,
    },
    Spec: corev1.ServiceSpec{
      Type: "NodePort",
      Selector: labels,
      Ports: []corev1.ServicePort{ { Port: 22, }, },
    },
  }
}

func (r *ReconcileSshdService) updateStage(cr *sshdoperatorv1alpha1.SshdService, stage string) {
  status := sshdoperatorv1alpha1.SshdServiceStatus{
    Stage: stage,
    Password: cr.Status.Password,
    NodePort: cr.Status.NodePort,
  }
  if !reflect.DeepEqual(cr.Status, status) {
    cr.Status = status
    r.client.Status().Update(context.TODO(), cr)
  }
}

func (r *ReconcileSshdService) updateNodePort(cr *sshdoperatorv1alpha1.SshdService, nodeport int) {
  status := sshdoperatorv1alpha1.SshdServiceStatus{
    Stage: cr.Status.Stage,
    Password: cr.Status.Password,
    NodePort: nodeport,
  }
  if !reflect.DeepEqual(cr.Status, status) {
    cr.Status = status
    r.client.Status().Update(context.TODO(), cr)
  }
}

func (r *ReconcileSshdService) updatePassword(cr *sshdoperatorv1alpha1.SshdService, password string) {
  status := sshdoperatorv1alpha1.SshdServiceStatus{
    Stage: cr.Status.Stage,
    Password: password,
    NodePort: cr.Status.NodePort,
  }
  if !reflect.DeepEqual(cr.Status, status) {
    cr.Status = status
    r.client.Status().Update(context.TODO(), cr)
  }
}

func generatePassword() string {
  c := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYG1234567890"
  l := len(c)
  p := ""
  for i	:= 0; i < 8; i++ {
    n := rand.Intn(l)
    p += c[n:n+1]
  }
  return p
}

