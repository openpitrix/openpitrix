#!/bin/bash

MINIKUBE_VERSION=v1.6.2
KUBE_VERSION=v1.15.7

curl -Lo kubectl https://storage.googleapis.com/kubernetes-release/release/${KUBE_VERSION}/bin/linux/amd64/kubectl \
  && chmod +x kubectl && sudo mv kubectl /usr/local/bin/

curl -Lo minikube https://storage.googleapis.com/minikube/releases/${MINIKUBE_VERSION}/minikube-linux-amd64 \
  && chmod +x minikube && sudo mv minikube /usr/local/bin/

sudo minikube start --vm-driver=none --kubernetes-version=${KUBE_VERSION} --extra-config=apiserver.v=10 --extra-config=kubelet.max-pods=100

sudo minikube update-context

sudo chown -R ${USER}:${USER} ${HOME}/.kube
sudo chown -R ${USER}:${USER} ${HOME}/.minikube

kubectl create clusterrolebinding add-on-cluster-admin --clusterrole=cluster-admin --serviceaccount=kube-system:default

JSONPATH='{range .items[*]}{@.metadata.name}:{range @.status.conditions[*]}{@.type}={@.status};{end}{end}'
until kubectl get nodes -o jsonpath="$JSONPATH" 2>&1 | grep -q "Ready=True";
do sleep 1; done
