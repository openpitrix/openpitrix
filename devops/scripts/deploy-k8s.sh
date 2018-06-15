#!/bin/bash

DEFAULT_NAMESPACE=openpitrix-system
DEFAULT_VERSION=latest
DEFAULT_METADATA=off

NAMESPACE=${DEFAULT_NAMESPACE}
VERSION=${DEFAULT_VERSION}

IMAGE="openpitrix/openpitrix:$VERSION"
METADATA_IMAGE="openpitrix/openpitrix:metadata"
FLYWAY_IMAGE="openpitrix/openpitrix:flyway"

if [ "VERSION" != "latest" ];then
  METADATA_IMAGE="openpitrix/openpitrix:metadata-$VERSION"
  FLYWAY_IMAGE="openpitrix/openpitrix:flyway-$VERSION"
fi

while getopts n:v:m:h option
do
  case "${option}"
  in
  n) NAMESPACE=${OPTARG};;
  v) VERSION=${OPTARG};;
  m) METADATA=${OPTARG};;
  h) echo "usage ./deploy-k8s.sh -n namespace -v version -m on" && exit 1 ;;
  *) echo "usage ./deploy-k8s.sh -n namespace -v version -m on" && exit 1 ;;
  esac
done

# Back to the root of the project
cd $(dirname $0)
cd ../..

kubectl create namespace ${NAMESPACE}
kubectl create secret generic mysql-pass --from-file=./devops/kubernetes/password.txt -n ${NAMESPACE}

for FILE in `ls ./devops/kubernetes/db/`;do
  sed -e "s!\${NAMESPACE}!${NAMESPACE}!g" -e "s!\${IMAGE}!${IMAGE}!g" -e "s!\${METADATA_IMAGE}!${METADATA_IMAGE}!g" -e "s!\${FLYWAY_IMAGE}!${FLYWAY_IMAGE}!g" ./devops/kubernetes/db/${FILE} | kubectl apply -f -
done

for FILE in `ls ./devops/kubernetes/etcd/`;do
  sed -e "s!\${NAMESPACE}!${NAMESPACE}!g" -e "s!\${IMAGE}!${IMAGE}!g" -e "s!\${METADATA_IMAGE}!${METADATA_IMAGE}!g" -e "s!\${FLYWAY_IMAGE}!${FLYWAY_IMAGE}!g" ./devops/kubernetes/etcd/${FILE} | kubectl apply -f -
done

for FILE in `ls ./devops/kubernetes/openpitrix/ | grep "^openpitrix-"`;do
  sed -e "s!\${NAMESPACE}!${NAMESPACE}!g" -e "s!\${IMAGE}!${IMAGE}!g" -e "s!\${METADATA_IMAGE}!${METADATA_IMAGE}!g" -e "s!\${FLYWAY_IMAGE}!${FLYWAY_IMAGE}!g" ./devops/kubernetes/openpitrix/${FILE} | kubectl apply -f -
done

if [ "$METADATA" == "on" ];then
  for FILE in `ls ./devops/kubernetes/openpitrix/metadata/`;do
    sed -e "s!\${NAMESPACE}!${NAMESPACE}!g" -e "s!\${IMAGE}!${IMAGE}!g" -e "s!\${METADATA_IMAGE}!${METADATA_IMAGE}!g" -e "s!\${FLYWAY_IMAGE}!${FLYWAY_IMAGE}!g" ./devops/kubernetes/openpitrix/metadata/${FILE} | kubectl apply -f -
  done
fi