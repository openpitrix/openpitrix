#!/bin/bash

FILE_NAME="$1"

CONTAINER_NAME=metad

if [ "$FILE_NAME" == "frontgate.conf" ]
then
  # single node metad
  if /usr/bin/docker ps -a -f name=${CONTAINER_NAME} | grep ${CONTAINER_NAME}; then /usr/bin/docker start -a ${CONTAINER_NAME}; else /usr/bin/docker run --name ${CONTAINER_NAME} --network host -v /opt/openpitrix/:/opt/openpitrix/ --privileged openpitrix/metad metad --backend etcdv3 --nodes http://127.0.0.1:2379 --log_level debug --listen :80 -listen_manage 127.0.0.1:9611 --xff; fi
fi
