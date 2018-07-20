#!/bin/bash

FILE_NAME="$1"

if [ "$FILE_NAME" == "frontgate.conf" ]
then
  # single node etcd
  if /usr/bin/docker ps -a | grep etcd; then /usr/bin/docker start -a registry; else /usr/bin/docker run -p 5000:5000 --restart=always --name registry -v /opt/openpitrix/conf/config.yml:/etc/docker/registry/config.yml -v /data:/var/lib/registry registry:2; fi
fi
