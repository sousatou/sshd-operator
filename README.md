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
Status:
  Nodeport:  31039
  Password:  fgDsc3WD
  Stage:     RUNNING
```

When "Stage" field become "RUNNING", then you can connect to the pod using ssh command.  
- IP address: node machine's address(any machine in the cluster)  
- Port: Nodeport number from status(above)  
- User: "user1" (which is specified in the CR)  
- Password: Password from status(above, auto-generated)  
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

When "Stage" field become "RUNNING", then you can connect to the pod by using ssh command.  
- IP address: node machine's address(any worker node in the cluster)  
- Port: Nodeport number from status(auto-generated, appeared in Status)  
- Username: "user1" (specified in the CR, appeared in Spec)  
- Password: Password from status(auto-generated, appeared in Status)  
- ex) # ssh -p 31039 user1@192.168.0.100  

---

## How it works
1. *Custom Resource*. Users deploy their Custom Resources, or *CR*.  
In the CR, you can specify the username to login to the ssh server(omittable, "user1" is default).  
See: "pkg/apis/sshdoperator/v1alpha1/sshdservice_types.go" for its definition.  
Or, "deploy/crds/sshd-operator.sousatou.com_v1alpha1_sshdservice_cr.yaml" for its example.  
  
2. *Custom Controller/Operator*. When the operator notices the CR deployed, it creates following Pod and Service resources.  
The Pod run "fedora" image from dockerhub(with "sleep infinity" command).  
The Pod needs to run as root user, and have privilege, because it installs and run opensshd-server.  
The Service exposes NodePort to connect Pod's 22 port. The operator stores the NodePort number to CR's status.  
The oeprator ganerates password to access sshd, and stores to the CR's status.  
See: "pkg/controller/sshdservice/sshdservice_controller.go"  
  
3. After the Pod created, the oeprator copies some scripts to the Pod.  
Then the operator executes "sshd_action" script within the Pod.  
The operator run "pod_init" script periodically, so it can check the service status and openssh-server's latest update.  
See: "build/bin/pod_init", "build/bin/action/sshd_action"  
  
4. "sshd_action" script would run its sub-scripts.  
One of them is "sshd_install", which installs openssh-server and other needed packages in the Pod.  
Then create a user, and start sshd service.  
After the sshd service started, the CR's status:STAGE changed to "RUNNING".  
See: "build/bin/action/sshd_install"  
  
5. Another sub-script is "sshd_update".  
The script checks whether the latest update of openssh-server package exists, and update the service.  
See: "build/bin/action/sshd_update"  
