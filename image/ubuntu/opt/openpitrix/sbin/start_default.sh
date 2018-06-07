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
  START_CMD="frontgate -config=/opt/openpitrix/conf/frontgate.conf serve"
elif [ "$FILE_NAME" == "drone.conf" ]
then
  START_CMD="drone -config=/opt/openpitrix/conf/drone.conf serve"
fi

echo $FILE_CONF > /opt/openpitrix/conf/$FILE_NAME
if ! service docker status|grep running ; then service docker start; fi
if /usr/bin/docker ps -a | grep default; then /usr/bin/docker start -a default; else /usr/bin/docker run $mount -v /opt/openpitrix/:/opt/openpitrix/ -v $CMD_INFO_TOML:$CMD_INFO_TOML -v $CMD_INFO_TMPL:$CMD_INFO_TMPL -v /var/run/docker.sock:/var/run/docker.sock --name default --network host --pid host --privileged $IMAGE $START_CMD; fi
