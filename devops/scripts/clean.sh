#!/bin/bash

echo "clean all stuff"

echo "start to clean k8s resource"
# Back to the root of the project
cd $(dirname $0)
cd ../..

DEFAULT_NAMESPACE=openpitrix-system
NAMESPACE=${DEFAULT_NAMESPACE}

while getopts n:h option
do
 case "${option}"
 in
 n) NAMESPACE=${OPTARG};;
 h) echo "usage ./clean.sh -n namespace" && exit 1 ;;
 *) echo "usage ./clean.sh -n namespace" && exit 1 ;;
 esac
done

kubectl delete namespace ${NAMESPACE}


echo "start to clean docker resource"
docker rmi openpitrix/openpitrix-dev:latest
docker rmi openpitrix

echo "cleaned successfully"
