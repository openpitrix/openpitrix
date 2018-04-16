#!/bin/bash

echo "clean all stuff"

echo "start to clean k8s resource"
# Back to the root of the project
cd $(dirname $0)
cd ../..

kubectl delete -f ./devops/kubernetes/openpitrix
kubectl delete -f ./devops/kubernetes/etcd
kubectl delete -f ./devops/kubernetes/ctrl
kubectl delete  -f ./devops/kubernetes/db



echo "start to clean docker resource"
docker rmi openpitrix/openpitrix-dev:latest
docker rmi openpitrix

echo "cleaned successfully"
