#!/bin/bash

set -e

[ -z `which make` ] && echo "You need to install 'make' first." && exit 1

DEFAULT_NAMESPACE=openpitrix-system
NAMESPACE=${DEFAULT_NAMESPACE}

usage() {
  echo "Usage:"
  echo "  generate-certs.sh [-n NAMESPACE] [-o HOST] -t [TYPE] "
  echo "Description:"
  echo "    -n NAMESPACE : the namespace of kubernetes."
  echo "    -o HOST      : the hostname used in certificate."
  echo "    -t TYPE      : the type(all/metadata/ingress) of resource that need to generate cert."
  exit -1
}

while getopts n:o:t: option
do
  case "${option}"
  in
    n) NAMESPACE=${OPTARG};;
    o) HOST=${OPTARG};;
    t) TYPE=${OPTARG};;
  esac
done

cd $(dirname $0)
cd ../..
cd ./kubernetes/tls-config

if [ "x${TYPE}" == "xall" ] || [ "x${TYPE}" == "xmetadata" ] || [ "x${TYPE}" == "x" ]; then
  make pilot-server.crt
  make pilot-client.crt
	SECRETS=("openpitrix-ca.crt" "pilot-server.crt" "pilot-server.key" "pilot-client.crt" "pilot-client.key")
	for i in ${SECRETS[@]}; do
	  echo ${i}
	  kubectl create secret generic ${i} --from-file=./${i} -n ${NAMESPACE} || true
	done
fi

if [ "x${TYPE}" == "xall" ] || [ "x${TYPE}" == "xingress" ] || [ "x${TYPE}" == "x" ]; then
	make ingress.crt-${HOST}
	kubectl create secret tls ingress-tls --key ingress.key --cert ingress.crt -n ${NAMESPACE} || true
fi