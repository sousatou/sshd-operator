#! /bin/bash
operator-sdk build sshd-operator
docker tag sshd-operator quay.io/sousatou/sshd-operator
docker push quay.io/sousatou/sshd-operator
