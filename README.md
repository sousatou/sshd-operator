# sshd-operator
To run the operator in (normal)Kubernetes:

```
# kubectl create -f deploy/crds/sshd-operator.sousatou.com_sshdservices_crd.yaml
# kubectl create -f deploy/
```

Then, add the CR:

```
# kubectl create -f deploy/crds/sshd-operator.sousatou.com_v1alpha1_sshdservice_cr.yaml
```

---


To run the operator in OpenShift:


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

