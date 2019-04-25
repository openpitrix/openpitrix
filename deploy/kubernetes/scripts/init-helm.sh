#!/bin/bash

TILLER_DEPLOY=`sudo kubectl get deploy -n kube-system | grep 'tiller-deploy'`
[ -n "$TILLER_DEPLOY" ] && echo "Helm is already initialized." && exit 1

export DESIRED_VERSION=v2.12.0

curl https://raw.githubusercontent.com/helm/helm/${DESIRED_VERSION}/scripts/get | bash

sudo helm init --tiller-image=gcr.io/kubernetes-helm/tiller:${DESIRED_VERSION}

sudo kubectl create serviceaccount --namespace kube-system tiller
sudo kubectl create clusterrolebinding tiller-cluster-rule --clusterrole=cluster-admin --serviceaccount=kube-system:tiller
sudo kubectl patch deploy --namespace kube-system tiller-deploy -p '{"spec":{"template":{"spec":{"serviceAccount":"tiller"}}}}'

until sudo helm list -a; do sleep 1; done
