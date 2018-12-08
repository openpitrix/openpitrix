#!/bin/bash

FILE_NAME="$1"

CONTAINER_NAME=registry

if [ "$FILE_NAME" == "frontgate.conf" ]
then
  # single node registry
  if /usr/bin/docker ps -a -f name=${CONTAINER_NAME} | grep ${CONTAINER_NAME}; then /usr/bin/docker start -a ${CONTAINER_NAME}; else /usr/bin/docker run -p 5000:5000 --restart=always --name ${CONTAINER_NAME} -v /opt/openpitrix/conf/config.yml:/etc/docker/registry/config.yml -v /data:/var/lib/registry registry:2; fi
fi
