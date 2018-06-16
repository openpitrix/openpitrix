#!/bin/bash

DEFAULT_NAMESPACE=openpitrix-system
DEFAULT_MODEL=default

NAMESPACE=${DEFAULT_NAMESPACE}
MODEL=${DEFAULT_MODEL}

while getopts n:v:m:h option
do
  case "${option}"
  in
  n) NAMESPACE=${OPTARG};;
  v) VERSION=${OPTARG};;
  m) MODEL=${OPTARG};;
  h) echo "usage ./deploy-k8s.sh -n namespace -v version -m default/all/dbctrl/metadata" && exit 1 ;;
  *) echo "usage ./deploy-k8s.sh -n namespace -v version -m default/all/dbctrl/metadata" && exit 1 ;;
  esac
done

if [ "${VERSION}" == "" ];then
  IMAGE="openpitrix/openpitrix-dev:latest"
  METADATA_IMAGE="openpitrix/openpitrix-dev:metadata"
  FLYWAY_IMAGE="openpitrix/openpitrix-dev:flyway"
elif [ "$VERSION" == "latest" ];then
  IMAGE="openpitrix/openpitrix:latest"
  METADATA_IMAGE="openpitrix/openpitrix:metadata"
  FLYWAY_IMAGE="openpitrix/openpitrix:flyway"
else
  IMAGE="openpitrix/openpitrix:$VERSION"
  METADATA_IMAGE="openpitrix/openpitrix:metadata-$VERSION"
  FLYWAY_IMAGE="openpitrix/openpitrix:flyway-$VERSION"
fi

# Back to the root of the project
cd $(dirname $0)
cd ../..

kubectl create namespace ${NAMESPACE}
kubectl create secret generic mysql-pass --from-file=./devops/kubernetes/password.txt -n ${NAMESPACE}


if [ "${MODEL}" == "dbctrl" ];then
  for FILE in `ls ./devops/kubernetes/ctrl`;do
    sed -e "s!\${NAMESPACE}!${NAMESPACE}!g" -e "s!\${IMAGE}!${IMAGE}!g" -e "s!\${METADATA_IMAGE}!${METADATA_IMAGE}!g" -e "s!\${FLYWAY_IMAGE}!${FLYWAY_IMAGE}!g" ./devops/kubernetes/ctrl/${FILE} | kubectl apply -f -
  done
  exit 0
elif [ "${MODEL}" == "metadata" ];then
  for FILE in `ls ./devops/kubernetes/openpitrix/metadata/`;do
    sed -e "s!\${NAMESPACE}!${NAMESPACE}!g" -e "s!\${IMAGE}!${IMAGE}!g" -e "s!\${METADATA_IMAGE}!${METADATA_IMAGE}!g" -e "s!\${FLYWAY_IMAGE}!${FLYWAY_IMAGE}!g" ./devops/kubernetes/openpitrix/metadata/${FILE} | kubectl apply -f -
  done
  exit 0
elif [ "${MODEL}" == "default" ];then
  for FILE in `ls ./devops/kubernetes/db/`;do
    sed -e "s!\${NAMESPACE}!${NAMESPACE}!g" -e "s!\${IMAGE}!${IMAGE}!g" -e "s!\${METADATA_IMAGE}!${METADATA_IMAGE}!g" -e "s!\${FLYWAY_IMAGE}!${FLYWAY_IMAGE}!g" ./devops/kubernetes/db/${FILE} | kubectl apply -f -
  done

  for FILE in `ls ./devops/kubernetes/etcd/`;do
    sed -e "s!\${NAMESPACE}!${NAMESPACE}!g" -e "s!\${IMAGE}!${IMAGE}!g" -e "s!\${METADATA_IMAGE}!${METADATA_IMAGE}!g" -e "s!\${FLYWAY_IMAGE}!${FLYWAY_IMAGE}!g" ./devops/kubernetes/etcd/${FILE} | kubectl apply -f -
  done

  for FILE in `ls ./devops/kubernetes/openpitrix/ | grep "^openpitrix-"`;do
    sed -e "s!\${NAMESPACE}!${NAMESPACE}!g" -e "s!\${IMAGE}!${IMAGE}!g" -e "s!\${METADATA_IMAGE}!${METADATA_IMAGE}!g" -e "s!\${FLYWAY_IMAGE}!${FLYWAY_IMAGE}!g" ./devops/kubernetes/openpitrix/${FILE} | kubectl apply -f -
  done
elif [ "${MODEL}" == "all" ];then
  for FILE in `ls ./devops/kubernetes/ctrl`;do
    sed -e "s!\${NAMESPACE}!${NAMESPACE}!g" -e "s!\${IMAGE}!${IMAGE}!g" -e "s!\${METADATA_IMAGE}!${METADATA_IMAGE}!g" -e "s!\${FLYWAY_IMAGE}!${FLYWAY_IMAGE}!g" ./devops/kubernetes/ctrl/${FILE} | kubectl apply -f -
  done

  for FILE in `ls ./devops/kubernetes/db/`;do
    sed -e "s!\${NAMESPACE}!${NAMESPACE}!g" -e "s!\${IMAGE}!${IMAGE}!g" -e "s!\${METADATA_IMAGE}!${METADATA_IMAGE}!g" -e "s!\${FLYWAY_IMAGE}!${FLYWAY_IMAGE}!g" ./devops/kubernetes/db/${FILE} | kubectl apply -f -
  done

  for FILE in `ls ./devops/kubernetes/etcd/`;do
    sed -e "s!\${NAMESPACE}!${NAMESPACE}!g" -e "s!\${IMAGE}!${IMAGE}!g" -e "s!\${METADATA_IMAGE}!${METADATA_IMAGE}!g" -e "s!\${FLYWAY_IMAGE}!${FLYWAY_IMAGE}!g" ./devops/kubernetes/etcd/${FILE} | kubectl apply -f -
  done

  for FILE in `ls ./devops/kubernetes/openpitrix/ | grep "^openpitrix-"`;do
    sed -e "s!\${NAMESPACE}!${NAMESPACE}!g" -e "s!\${IMAGE}!${IMAGE}!g" -e "s!\${METADATA_IMAGE}!${METADATA_IMAGE}!g" -e "s!\${FLYWAY_IMAGE}!${FLYWAY_IMAGE}!g" ./devops/kubernetes/openpitrix/${FILE} | kubectl apply -f -
  done

  for FILE in `ls ./devops/kubernetes/openpitrix/metadata/`;do
    sed -e "s!\${NAMESPACE}!${NAMESPACE}!g" -e "s!\${IMAGE}!${IMAGE}!g" -e "s!\${METADATA_IMAGE}!${METADATA_IMAGE}!g" -e "s!\${FLYWAY_IMAGE}!${FLYWAY_IMAGE}!g" ./devops/kubernetes/openpitrix/metadata/${FILE} | kubectl apply -f -
  done
else
  echo "usage ./deploy-k8s.sh -n namespace -v version -m default/all/dbctrl/metadata" && exit 1
fi

