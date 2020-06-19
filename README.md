# sshd-operator
## To run the operator in (normal)Kubernetes cluster:

```
# kubectl create -f deploy/crds/sshd-operator.sousatou.com_sshdservices_crd.yaml
# kubectl create -f deploy/
```

Then, add the CR:

```
# kubectl create -f deploy/crds/sshd-operator.sousatou.com_v1alpha1_sshdservice_cr.yaml
```

Check status of the CR:

```
# kubectl describe sshdservice example-sshdservice
...
Spec:
  Username:  user1
Status:
  Nodeport:  31039
  Password:  fgDsc3WD
  Stage:     RUNNING
```

When "Stage" field become "RUNNING", you can connect to the pod by using ssh command.  
- IP address: node machine's IP address(any worker node in the cluster)  
- Port: Nodeport number from status(auto-generated, appeared in Status)  
- Username: "user1" (specified in the CR, appeared in Spec)  
- Password: Password from status(auto-generated, appeared in Status)  
- ex) # ssh -p 31039 user1@192.168.0.100  
  
---

## To run the operator in OpenShift cluster:

```
# oc create -f deploy/crds/sshd-operator.sousatou.com_sshdservices_crd.yaml
# oc create -f deploy/service_account.yaml
# oc adm policy add-cluster-role-to-user cluster-admin -z sshd-operator
# oc create -f deploy/
```

Then, add the CR:

```
# oc create -f deploy/crds/sshd-operator.sousatou.com_v1alpha1_sshdservice_cr.yaml
```

Check status of the CR:

```
# oc describe sshdservice example-sshdservice
...
Spec:
  Username:  user1
Status:
  Nodeport:  31039
  Password:  fgDsc3WD
  Stage:     RUNNING
```

When "Stage" field become "RUNNING", you can connect to the pod by using ssh command.  
- IP address: node machine's IP address(any worker node in the cluster)  
- Port: Nodeport number from status(auto-generated, appeared in Status)  
- Username: "user1" (specified in the CR, appeared in Spec)  
- Password: Password from status(auto-generated, appeared in Status)  
- ex) # ssh -p 31039 user1@192.168.0.100  

---

## How it works
### Custom Resource (SshdService)
In CR, you can specify the username to login to the ssh server(omittable, "user1" is default).  
See: "pkg/apis/sshdoperator/v1alpha1/sshdservice_types.go" for its definition.  
Or, "deploy/crds/sshd-operator.sousatou.com_v1alpha1_sshdservice_cr.yaml" for its example.  
  
### Custom Controller
When the operator notices the CR deployed, it creates following Pod and Service resources.  
The Pod run "fedora" image from dockerhub(with "sleep infinity" command).  
The Pod needs to run as root user, and have privilege, because it installs and run opensshd-server.  
The Service exposes NodePort to connect Pod's 22 port. The operator stores the NodePort number to CR's status.  
The oeprator ganerates password to access sshd, and stores to the CR's status.  
See: "pkg/controller/sshdservice/sshdservice_controller.go".  
Sshd Pod and Service manifests are also within the Go file.  
  
### A script running in the operator
After the Pod created, the operator run "pod_init" script periodically, so it can setup sshd service and check its update.  
"pod_init" script copies some script files to the Sshd Pod.  
Then it executes "sshd_action" script within the Pod.  
See: "build/bin/pod_init"  
  
### Scripts running in the sshd Pod
*/action/sshd_action*  
The script which is called from the operator periodically.  
It run other sub-script according to the Pod status(STAGE).  
See: "build/bin/action/sshd_action"  
  
*/action/sshd_install*  
Installs openssh-server and other packages needed in the Pod.  
Create a user, and start sshd service.  
After the sshd service start, the CR's status(STAGE) will change to "RUNNING".  
See: "build/bin/action/sshd_install"  
  
*/action/sshd_update*  
The script checks if the latest update of openssh-server package exists, and update the service.  
See: "build/bin/action/sshd_update"  
