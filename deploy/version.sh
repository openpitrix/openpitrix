#!/bin/bash

##############################################################################
########## The image versions of OpenPitrix and related services #############
##############################################################################

### OpenPitrix version latest
openpitrix_latest="
  VERSION=latest
  DASHBOARD_VERSION=latest
  IM_VERSION=latest
  AM_VERSION=latest
  NOTIFICATION_VERSION=latest
  RP_QINGCLOUD_VERSION=latest
  RP_AWS_VERSION=latest
  RP_ALIYUN_VERSION=latest
  RP_K8S_VERSION=latest
  WATCHER_VERSION=latest
"

### OpenPitrix version v0.4.0
openpitrix_v0_4_0="
  VERSION=v0.4.0
  DASHBOARD_VERSION=v0.4.1
  IM_VERSION=v0.1.0
  AM_VERSION=v0.1.0
  NOTIFICATION_VERSION=v0.1.0
  RP_QINGCLOUD_VERSION=v0.1.0
  RP_AWS_VERSION=v0.1.0
  RP_ALIYUN_VERSION=v0.1.0
  RP_K8S_VERSION=v0.1.0
  WATCHER_VERSION=v0.1.0
"

### OpenPitrix version v0.4.1
openpitrix_v0_4_1="
  VERSION=v0.4.1
  DASHBOARD_VERSION=v0.4.2
  IM_VERSION=v0.1.0
  AM_VERSION=v0.1.0
  NOTIFICATION_VERSION=v0.1.0
  RP_QINGCLOUD_VERSION=v0.1.0
  RP_AWS_VERSION=v0.1.0
  RP_ALIYUN_VERSION=v0.1.0
  RP_K8S_VERSION=v0.1.0
  WATCHER_VERSION=v0.1.0
"

## Usage:
## sh version.sh [openpitrix_${version}]

VERSION=$1
DEFAULT_VERSION="openpitrix_latest"
if [ "x${VERSION}" == "x" ]; then
  VERSION=${DEFAULT_VERSION}
fi

#get openpitrix version, eg: openpitrix_v0_4_0
OP_VERSION=${VERSION//[.-]/_}
VERSIONS=`eval echo '$'"${OP_VERSION}"`
# check if the given version exist
if [ "x${VERSIONS}" == "x" ]; then
  echo "The version ${VERSION} of openpitrix not exist!"
  exit 1
fi

#echo versions
for V in ${VERSIONS} ; do
  export ${V}
  echo ${V}
done


#openpitrix latest images
openpitrix_latest_images="
  IMAGE=openpitrix/openpitrix:latest
  FLYWAY_IMAGE=openpitrix/openpitrix:flyway
  DASHBOARD_IMAGE=openpitrix/dashboard:latest
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
  WATCHER_IMAGE=openpitrix/watcher:latest
"

#openpitrix images
openpitrix_images="
  IMAGE=openpitrix/openpitrix:${VERSION}
  FLYWAY_IMAGE=openpitrix/openpitrix:flyway-${VERSION}
  DASHBOARD_IMAGE=openpitrix/dashboard:${DASHBOARD_VERSION}
  IM_IMAGE=kubespheredev/im:${IM_VERSION}
  IM_FLYWAY_IMAGE=kubespheredev/im:flyway-${IM_VERSION}
  AM_IMAGE=openpitrix/iam:${AM_VERSION}
  AM_FLYWAY_IMAGE=openpitrix/iam:flyway-${AM_VERSION}
  NOTIFICATION_IMAGE=openpitrix/notification:${NOTIFICATION_VERSION}
  NOTIFICATION_FLYWAY_IMAGE=openpitrix/notification:flyway-${NOTIFICATION_VERSION}
  RP_QINGCLOUD_IMAGE=openpitrix/runtime-provider-qingcloud:${RP_QINGCLOUD_VERSION}
  RP_AWS_IMAGE=openpitrix/runtime-provider-aws:${RP_AWS_VERSION}
  RP_ALIYUN_IMAGE=openpitrix/runtime-provider-aliyun:${RP_ALIYUN_VERSION}
  RP_K8S_IMAGE=openpitrix/runtime-provider-kubernetes:${RP_K8S_VERSION}
  WATCHER_IMAGE=openpitrix/watcher:${WATCHER_VERSION}
"

IMAGES=${openpitrix_images}
if [ "x${OP_VERSION}" == "x${DEFAULT_VERSION}" ]; then
  IMAGES=${openpitrix_latest_images}
fi

for I in ${IMAGES} ; do
  echo ${I}
done
