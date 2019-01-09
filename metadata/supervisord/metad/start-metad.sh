#!/bin/sh


COMMAND="/usr/local/bin/metad"
ARGV="--backend etcdv3 --nodes http://127.0.0.1:2379 --log_level debug --listen :80 -listen_manage 127.0.0.1:9611 --xff"


if [ -f "/opt/openpitrix/conf/pilot-version" ]; then
	PILOT_VERSION=$(cat /opt/openpitrix/conf/pilot-version)

	if [ -f "/opt/openpitrix/bin/${PILOT_VERSION}/metad" ]; then
		COMMAND="/opt/openpitrix/bin/${PILOT_VERSION}/metad"
	fi
fi

exec ${COMMAND} ${ARGV}
