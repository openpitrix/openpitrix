#!/bin/bash

set -e

[ -z `which make` ] && echo "You need to install 'make' first." && exit 1

DEFAULT_NAMESPACE=openpitrix-system
NAMESPACE=${DEFAULT_NAMESPACE}

usage() {
  echo "Usage:"
  echo "  generate-certs.sh [-n NAMESPACE]"
  echo "Description:"
  echo "    -n NAMESPACE: the namespace of kubernetes."
  exit -1
}

while getopts n: option
do
  case "${option}"
  in
    n) NAMESPACE=${OPTARG};;
  esac
done

cd $(dirname $0)
cd ../..
cd ./kubernetes/tls-config
make

SECRETS=("openpitrix-ca.crt" "pilot-server.crt" "pilot-server.key" "pilot-client.crt" "pilot-client.key")
for i in ${SECRETS[@]}; do
  echo ${i}
  kubectl create secret generic ${i} --from-file=./${i} -n ${NAMESPACE} || true
done
