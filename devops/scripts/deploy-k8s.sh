#!/bin/bash

# Back to the root of the project
cd $(dirname $0)
cd ../..

kubectl create secret generic mysql-pass --from-file=./devops/kubernetes/password.txt -n default
kubectl apply -f ./devops/kubernetes/openpitrix.yaml
