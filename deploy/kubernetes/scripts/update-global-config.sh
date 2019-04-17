#!/bin/bash

set +e

# Back to the root of the project
cd $(dirname $0)
cd ../..

DEFAULT_NAMESPACE=openpitrix-system
NAMESPACE=${DEFAULT_NAMESPACE}
DEFAULT_VERSION=latest
VERSION=${DEFAULT_VERSION}
CONFIG_FILE="./config/global_config.yaml"

usage() {
  echo "Usage:"
  echo "  update-global-config.sh [-n NAMESPACE]"
  echo "Description:"
  echo "    -v VERSION : the version of openpitrix, default: latest."
  echo "    -n NAMESPACE : the namespace of kubernetes, default: openpitrix-system."
  exit -1
}

while getopts v:n:h option
do
  case "${option}"
  in
  v) VERSION=${OPTARG};;
  n) NAMESPACE=${OPTARG};;
  h) usage ;;
  *) usage ;;
  esac
done

if [ ! -f ${CONFIG_FILE} ]; then
  echo "The ${CONFIG_FILE} not exist!"
  exit 1
fi

CONFIG_MAP_FILE=$(sh ./kubernetes/scripts/generate-config-map.sh)

sed -e "s/\${VERSION}/${VERSION}/g" -e "s/\${NAMESPACE}/${NAMESPACE}/g" ${CONFIG_MAP_FILE} | kubectl apply -f -

echo "Global config updated to config_map, please wait for a while to check it."
