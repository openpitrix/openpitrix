#!/bin/bash

MOUNT_POINT="$1"
FILE_NAME="$2"
IMAGE="$3"

CMD_INFO_TOML="/etc/confd/conf.d/cmd.info.toml"
CMD_INFO_TMPL="/etc/confd/templates/cmd.info.tmpl"
CMD_EXEC="/opt/openpitrix/sbin/exec.sh"

mount=""
if [ "$MOUNT_POINT" != "#" ]
then
  mount="-v $MOUNT_POINT:$MOUNT_POINT"
fi

if ! service docker status|grep running ; then service docker start; fi
if docker ps -a | grep default; then docker start -a default; else docker kill $(docker ps -q); docker run -i $mount -v /opt/openpitrix/log/:/opt/openpitrix/log/ -v $CMD_INFO_TOML:$CMD_INFO_TOML -v $CMD_INFO_TMPL:$CMD_INFO_TMPL -v $CMD_EXEC:$CMD_EXEC --name default --network host --privileged $IMAGE; fi
