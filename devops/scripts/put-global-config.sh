#!/usr/bin/env bash

POD=`kubectl get pods -l tier=etcd --namespace=openpitrix-system --no-headers -o custom-columns=:metadata.name`


# Back to the root of the project
cd $(dirname $0)
cd ../..

test -s config/global_config.yaml || { echo "[config/global_config.yaml] not exist"; exit 1; }

cat config/global_config.yaml | kubectl run --restart=Never --quiet --rm -i test --image=openpitrix/openpitrix:latest -- opctl validate_global_config

if [[ $? != 0 ]]; then exit 1; fi

cat config/global_config.yaml | kubectl exec --namespace=openpitrix-system -i $POD etcdctl put openpitrix/global_config
