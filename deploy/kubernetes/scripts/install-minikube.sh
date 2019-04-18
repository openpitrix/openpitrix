#!/bin/bash

curl -Lo kubectl https://storage.googleapis.com/kubernetes-release/release/v1.10.0/bin/linux/amd64/kubectl \
  && chmod +x kubectl && sudo mv kubectl /usr/local/bin/

curl -Lo minikube https://storage.googleapis.com/minikube/releases/v0.25.2/minikube-linux-amd64 \
  && chmod +x minikube && sudo mv minikube /usr/local/bin/

sudo minikube start --vm-driver=none --kubernetes-version=v1.10.0 --extra-config=apiserver.Authorization.Mode=RBAC

sudo minikube update-context

sudo kubectl create clusterrolebinding add-on-cluster-admin --clusterrole=cluster-admin --serviceaccount=kube-system:default

JSONPATH='{range .items[*]}{@.metadata.name}:{range @.status.conditions[*]}{@.type}={@.status};{end}{end}'
until sudo kubectl get nodes -o jsonpath="$JSONPATH" 2>&1 | grep -q "Ready=True";
do sleep 1; done
