#!/usr/bin/env bash

##############################################################################
########## The image versions of OpenPitrix and related services #############
##############################################################################

### OpenPitrix version latest
openpitrix_latest="
	IMAGE=openpitrix/openpitrix:latest
	DASHBOARD_IMAGE=openpitrix/dashboard:latest
	FLYWAY_IMAGE=openpitrix/openpitrix:flyway
	IM_IMAGE=kubespheredev/im:latest
	IM_FLYWAY_IMAGE=kubespheredev/im:flyway
	AM_IMAGE=openpitrix/iam:latest
	AM_FLYWAY_IMAGE=openpitrix/iam:flyway
	NOTIFICATION_IMAGE=openpitrix/notification:latest
	NOTIFICATION_FLYWAY_IMAGE=openpitrix/notification:flyway
	RP_QINGCLOUD_IMAGE=openpitrix/runtime-provider-qingcloud:latest
	RP_AWS_IMAGE=openpitrix/runtime-provider-aws:latest
	RP_ALIYUN_IMAGE=openpitrix/runtime-provider-aliyun:latest
	RP_K8S_IMAGE=openpitrix/runtime-provider-kubernetes:latest
"

### OpenPitrix version v0.4.0
openpitrix_v0_4_0="
	IMAGE=openpitrix/openpitrix:v0.4.0
	DASHBOARD_IMAGE=openpitrix/dashboard:v0.4.1
	FLYWAY_IMAGE=openpitrix/openpitrix:flyway-v0.4.0
	IM_IMAGE=kubespheredev/im:v0.1.0
	IM_FLYWAY_IMAGE=kubespheredev/im:flyway-v0.1.0
	AM_IMAGE=openpitrix/iam:v0.1.0
	AM_FLYWAY_IMAGE=openpitrix/iam:flyway-v0.1.0
	NOTIFICATION_IMAGE=openpitrix/notification:v0.1.0
	NOTIFICATION_FLYWAY_IMAGE=openpitrix/notification:flyway-v0.1.0
	RP_QINGCLOUD_IMAGE=openpitrix/runtime-provider-qingcloud:v0.1.0
	RP_AWS_IMAGE=openpitrix/runtime-provider-aws:v0.1.0
	RP_ALIYUN_IMAGE=openpitrix/runtime-provider-aliyun:v0.1.0
	RP_K8S_IMAGE=openpitrix/runtime-provider-kubernetes:v0.1.0
"

## Usage:
## sh version.sh [openpitrix_${version}]

VERSION=$1
DEFAULT_VERSION="openpitrix_latest"
if [ "x${VERSION}" == "x" ]; then
	VERSION=${DEFAULT_VERSION}
fi

OP_VERSION=${VERSION//[.-]/_}
VERSIONS=`eval echo '$'"${OP_VERSION}"`
# check if the given version exist
if [ "x${VERSIONS}" == "x" ]; then
	echo "The version ${VERSION} of openpitrix not exist!"
	exit 1
fi

for V in ${VERSIONS} ; do
	echo ${V}
done
