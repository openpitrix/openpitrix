#!/bin/bash

DEFAULT_NAMESPACE=openpitrix-system
NAMESPACE=${DEFAULT_NAMESPACE}

usage() {
  echo "Usage:"
  echo "  put-global-config.sh [-n NAMESPACE]"
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

POD=`kubectl get pods -l tier=etcd --namespace=${NAMESPACE} --no-headers -o custom-columns=:metadata.name`

echo "Putting global config..."
# Back to the root of the project
cd $(dirname $0)
cd ../..

test -s config/global_config.yaml || { echo "[config/global_config.yaml] not exist"; exit 1; }

cat config/global_config.yaml | kubectl run --restart=Never --quiet --rm -i test --image=openpitrix/openpitrix:latest -- opctl validate_global_config

if [[ $? != 0 ]]; then exit 1; fi

cat config/global_config.yaml | kubectl exec --namespace=${NAMESPACE} -i ${POD} -- /bin/sh -c "export ETCDCTL_API=3 && etcdctl put openpitrix/global_config"
echo "Put successfully"
