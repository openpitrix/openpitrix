#!/bin/bash


TILLER_DEPLOY=`kubectl get deploy -n kube-system | grep 'tiller-deploy'`

[ -n "$TILLER_DEPLOY" ] && echo "Helm is already initialized." && exit 1

helm init --tiller-image=gcr.io/kubernetes-helm/tiller:v2.9.1 --upgrade

kubectl create serviceaccount --namespace kube-system tiller
kubectl create clusterrolebinding tiller-cluster-rule --clusterrole=cluster-admin --serviceaccount=kube-system:tiller
kubectl patch deploy --namespace kube-system tiller-deploy -p '{"spec":{"template":{"spec":{"serviceAccount":"tiller"}}}}'
