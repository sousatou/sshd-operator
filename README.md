# sshd-operator
To run the operator in (normal)Kubernetes clusters:

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
  IP address: node machine's address(any machine in the cluster)
  Port: Nodeport number from status(above)
  User: "user1" (which is specified in the CR)
  Password: Password from status(above, auto-generated)
  ex) # ssh -p 31039 user1@192.168.0.100
  
---

To run the operator in OpenShift clusters:

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
---

When "Stage" field become "RUNNING", then you can connect to the pod using ssh command.
  IP address: node machine's address(any machine in the cluster)
  Port: Nodeport number from status(above)
  User: "user1" (which is specified in the CR)
  Password: Password from status(above, auto-generated)
  ex) # ssh -p 31039 user1@192.168.0.100
  
