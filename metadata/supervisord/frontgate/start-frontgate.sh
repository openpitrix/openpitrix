#!/bin/sh


COMMAND="/usr/local/bin/frontgate"
ARGV="-config=/opt/openpitrix/conf/frontgate.conf --openpitrix-ca-crt-file=/opt/openpitrix/conf/openpitrix-ca.crt --pilot-client-crt-file=/opt/openpitrix/conf/pilot-client.crt --pilot-client-key-file=/opt/openpitrix/conf/pilot-client.key"


if [ -f "/opt/openpitrix/conf/pilot-version" ]; then
	PILOT_VERSION=$(cat /opt/openpitrix/conf/pilot-version)

	if [ -f "/opt/openpitrix/bin/${PILOT_VERSION}/frontgate" ]; then
		COMMAND="/opt/openpitrix/bin/${PILOT_VERSION}/frontgate"
	fi
fi

exec ${COMMAND} ${ARGV} serve
