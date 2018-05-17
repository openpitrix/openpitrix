#!/bin/bash

FILE_NAME="$1"
FILE_CONF="$2"

echo $FILE_CONF > /opt/openpitrix/conf/$FILE_NAME

if [ "$FILE_NAME" == "frontgate.conf" ]
then
  frontgate -config=/opt/openpitrix/conf/frontgate.conf serve
elif [ "$FILE_NAME" == "drone.conf" ]
then
  drone -config=/opt/openpitrix/conf/drone.conf serve
fi
