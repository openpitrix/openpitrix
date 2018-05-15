#!/bin/bash

FILE_NAME="$1"

if [ "$FILE_NAME" == "frontgate.conf" ]
then
  # single node etcd
  if docker ps -a | grep etcd; then docker start -a etcd; else docker run -v /opt/openpitrix/etcd/:/opt/openpitrix/etcd/ --name etcd --network host --privileged appcelerator/etcd --data-dir=/opt/openpitrix/etcd/ --name node1; fi
fi
