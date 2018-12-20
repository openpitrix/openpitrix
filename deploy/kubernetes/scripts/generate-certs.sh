#!/bin/bash

set -e

[ -z `which make` ] && echo "You need to install 'make' first." && exit 1

DEFAULT_NAMESPACE=openpitrix-system
NAMESPACE=${DEFAULT_NAMESPACE}

usage() {
  echo "Usage:"
  echo "  generate-certs.sh [-n NAMESPACE] [-c HOST]"
  echo "Description:"
  echo "    -n NAMESPACE : the namespace of kubernetes."
  echo "    -c HOST      : the hostname used in certificate."
  exit -1
}

while getopts n:c: option
do
  case "${option}"
  in
    n) NAMESPACE=${OPTARG};;
    c) HOST=${OPTARG};;
  esac
done

cd $(dirname $0)
cd ../..
cd ./kubernetes/tls-config
make ingress.crt-${HOST}

SECRETS=("openpitrix-ca.crt" "pilot-server.crt" "pilot-server.key" "pilot-client.crt" "pilot-client.key")
for i in ${SECRETS[@]}; do
  echo ${i}
  kubectl create secret generic ${i} --from-file=./${i} -n ${NAMESPACE} || true
done

kubectl create secret tls ingress-tls --key ingress.key --cert ingress.crt -n ${NAMESPACE} || true
