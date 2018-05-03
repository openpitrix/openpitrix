#!/bin/bash

MOUNT_POINT="$1"
FILE_NAME="$2"
FILE_CONF="$3"
IMAGE="$4"

mount=""
if [ "$MOUNT_POINT" != "#" ]
then
  mount="-v $MOUNT_POINT:$MOUNT_POINT"
fi

echo $FILE_CONF > /opt/openpitrix/conf/$FILE_NAME
if ! service docker status|grep running ; then service docker start; fi
if docker ps -a | grep default; then docker start -a default; else docker kill $(docker ps -q); docker run $mount -v /opt/openpitrix/conf/:/opt/openpitrix/conf/ --name default --network host --privileged $IMAGE ; fi
