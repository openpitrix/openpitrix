#!/bin/bash

DEFAULT_NAMESPACE=openpitrix-system
NAMESPACE=${DEFAULT_NAMESPACE}
PILOT_IP=
PILOT_PORT=

usage() {
  echo "Usage:"
  echo "  put-global-config.sh [-n NAMESPACE] COMMAND"
  echo "Description:"
  echo "    -n NAMESPACE : the namespace of kubernetes."
  echo "    -i PILOT_IP  : the ip of pilot service."
  echo "    -p PILOT_PORT: the port of pilot service."
  exit -1
}

while getopts n:i:p:h option
do
  case "${option}"
  in
  n) NAMESPACE=${OPTARG};;
  i) PILOT_IP=${OPTARG};;
  p) PILOT_PORT=${OPTARG};;
  h) usage ;;
  *) usage ;;
  esac
done

POD=`kubectl get pods -l tier=etcd --namespace=${NAMESPACE} --no-headers -o custom-columns=:metadata.name`

echo "Putting global config..."
# Back to the root of the project
cd $(dirname $0)
cd ../..

if [ ! -f "config/global_config.yaml" ];then
  cp config/global_config.init.yaml config/global_config.yaml
fi

if [ ! -z "${PILOT_IP}" ]; then
  OLD_PILOT_IP=`cat config/global_config.yaml | grep -e "^  ip:" | awk -F : '{print $2}'`
  sed -i -e "s/${OLD_PILOT_IP}/ ${PILOT_IP}/g" config/global_config.yaml
fi
if [ ! -z "${PILOT_PORT}" ]; then
  OLD_PILOT_PORT=`cat config/global_config.yaml | grep -e "^  port:" | awk -F : '{print $2}'`
  sed -i -e "s/${OLD_PILOT_PORT}/ ${PILOT_PORT}/g" config/global_config.yaml
fi

cat config/global_config.yaml | kubectl run --namespace=${NAMESPACE} --restart=Never --quiet --rm -i test --image=openpitrix/openpitrix:latest -- opctl validate_global_config

if [[ $? != 0 ]]; then exit 1; fi

cat config/global_config.yaml | kubectl exec --namespace=${NAMESPACE} -i ${POD} -- /bin/sh -c "export ETCDCTL_API=3 && etcdctl put openpitrix/global_config"
echo "Put successfully"
