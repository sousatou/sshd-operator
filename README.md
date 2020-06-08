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

## How it works
1. Users deploy their sshdservice_cr.yaml.  
They can specify the username to login to the ssh server(omittable).
See: "pkg/apis/sshdoperator/v1alpha1/sshdservice_types.go" for its definition.
Or, "deploy/crds/sshd-operator.sousatou.com_v1alpha1_sshdservice_cr.yaml" for an example.

2. The operator notice the CR deployed, then creates following Pod and Service resources.
The Pod run "fedora" image from dockerhub(with "sleep infinity" command), which is running as root user, and has privilege.
The Service exposes NodePort to connect Pod's 22 port. The operator stores the NodePort number to CR's status.
See: "pkg/controller/sshdservice/sshdservice_controller.go"

3. After the Pod created, the oeprator copies some scripts to the Pod, which are installing and maintaining sshd service.
Then the operator execute those scripts within the Pod.
See: "build/bin/pod_init"

4. Within the Pod, the installer script installs openshift-server and other needed packages, then create user, run sshd service.
The password for the user generated in the operator, and stored in the CR's status.
See: "build/bin/action/sshd_install"

5. After the sshd service running, the CR's status:STAGE changed to "RUNNING", the user can connect to its sshd server through NodePort.
