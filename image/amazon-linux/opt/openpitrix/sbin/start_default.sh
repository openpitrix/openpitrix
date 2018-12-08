#!/bin/bash

MOUNT_POINT="$1"
FILE_NAME="$2"
FILE_CONF="$3"
IMAGE="$4"

CMD_INFO_TOML="/etc/confd/conf.d/cmd.info.toml"
CMD_INFO_TMPL="/etc/confd/templates/cmd.info.tmpl"

mount=""
if [ "$MOUNT_POINT" != "#" ]
then
  mount="-v $MOUNT_POINT:$MOUNT_POINT"
fi

START_CMD=""
if [ "$FILE_NAME" == "frontgate.conf" ]
then
  START_CMD="frontgate"
elif [ "$FILE_NAME" == "drone.conf" ]
then
  START_CMD="drone"
fi

CONTAINER_NAME=default

if ! service docker status|grep running ; then service docker start; fi
if /usr/bin/docker ps -a -f name=${CONTAINER_NAME} | grep ${CONTAINER_NAME}; then /usr/bin/docker start -a ${CONTAINER_NAME}; else test -s /opt/openpitrix/conf/$FILE_NAME || echo $FILE_CONF > /opt/openpitrix/conf/$FILE_NAME; /usr/bin/docker run $mount -v /opt/openpitrix/:/opt/openpitrix/ -v $CMD_INFO_TOML:$CMD_INFO_TOML -v $CMD_INFO_TMPL:$CMD_INFO_TMPL -v /var/run/docker.sock:/var/run/docker.sock --name ${CONTAINER_NAME} --network host --pid host --privileged $IMAGE $START_CMD; fi