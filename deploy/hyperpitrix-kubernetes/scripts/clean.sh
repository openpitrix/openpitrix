#!/bin/bash

echo "Cleaning k8s resource..."
# Back to the root of the project
cd $(dirname $0)
cd ../..

DEFAULT_NAMESPACE=openpitrix-system
NAMESPACE=${DEFAULT_NAMESPACE}

usage() {
  echo "Usage:"
  echo "  clean.sh [-n NAMESPACE]"
  echo "Description:"
  echo "    -n NAMESPACE: the namespace of kubernetes."
  exit -1
}

while getopts n:h option
do
 case "${option}"
 in
 n) NAMESPACE=${OPTARG};;
 h) usage ;;
 *) usage ;;
 esac
done

kubectl delete namespace ${NAMESPACE}


echo "Cleaning docker resource..."
docker rmi openpitrix/openpitrix-dev:latest
docker rmi openpitrix/openpitrix-dev:metadata
docker rmi openpitrix/openpitrix-dev:flyway
docker rmi openpitrix

echo "Cleaned successfully"
