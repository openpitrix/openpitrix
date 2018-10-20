#!/bin/sh


COMMAND="/usr/local/bin/drone"
ARGV="-config=/opt/openpitrix/conf/drone.conf"


if [ -f "/opt/openpitrix/conf/pilot-version" ]; then
	PILOT_VERSION=$(cat /opt/openpitrix/conf/pilot-version)

	if [ -f "/opt/openpitrix/bin/${PILOT_VERSION}/drone" ]; then
		COMMAND="/opt/openpitrix/bin/${PILOT_VERSION}/drone"
	fi
fi

exec ${COMMAND} ${ARGV} serve
