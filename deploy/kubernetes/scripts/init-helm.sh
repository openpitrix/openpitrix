#!/bin/bash

TILLER_DEPLOY=`kubectl get deploy -n kube-system | grep 'tiller-deploy'`
[ -n "$TILLER_DEPLOY" ] && echo "Helm is already initialized." && exit 1


curl https://raw.githubusercontent.com/helm/helm/master/scripts/get | DESIRED_VERSION=v2.12.0 bash

helm init --tiller-image=gcr.io/kubernetes-helm/tiller:v2.12.0 --upgrade

kubectl create serviceaccount --namespace kube-system tiller
kubectl create clusterrolebinding tiller-cluster-rule --clusterrole=cluster-admin --serviceaccount=kube-system:tiller
kubectl patch deploy --namespace kube-system tiller-deploy -p '{"spec":{"template":{"spec":{"serviceAccount":"tiller"}}}}'

until sudo helm list -a; do sleep 1; done
