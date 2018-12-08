#!/bin/bash

FILE_NAME="$1"

CONTAINER_NAME=etcd

if [ "$FILE_NAME" == "frontgate.conf" ]
then
  # single node etcd
  if /usr/bin/docker ps -a -f name=${CONTAINER_NAME} | grep ${CONTAINER_NAME}; then /usr/bin/docker start -a ${CONTAINER_NAME}; else /usr/bin/docker run -v /opt/openpitrix/etcd/:/opt/openpitrix/etcd/ --name ${CONTAINER_NAME} --network host --privileged appcelerator/etcd --data-dir=/opt/openpitrix/etcd/ --name node1; fi
fi
