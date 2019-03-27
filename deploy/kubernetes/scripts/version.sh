#!/usr/bin/env bash

##############################################################################
########## The image versions of openpitrix and related services #############
##############################################################################

### Openpitrix version latest
V_LATEST_OPENPITRIX="IMAGE=openpitrix/openpitrix:latest
	METADATA_IMAGE=openpitrix/openpitrix:metadata
	DASHBOARD_IMAGE=openpitrix/dashboard:latest
	FLYWAY_IMAGE=openpitrix/openpitrix:flyway
	IM_IMAGE=kubespheredev/im:latest
	IM_FLYWAY_IMAGE=kubespheredev/im:flyway
	AM_IMAGE=openpitrix/iam:latest
	AM_FLYWAY_IMAGE=openpitrix/iam:flyway
	NOTIFICATION_IMAGE=openpitrix/notification:latest
	NOTIFICATION_FLYWAY_IMAGE=openpitrix/notification:flyway"


### Openpitrix version 0.4.0
V_0_4_0_OPENPITRIX="IMAGE=openpitrix/openpitrix:0.4.0
 	METADATA_IMAGE=openpitrix/openpitrix:metadata-0.4.0
	FLYWAY_IMAGE=openpitrix/openpitrix:flyway-0.4.0
	IM_IMAGE=kubespheredev/im:0.4.0
	IM_FLYWAY_IMAGE=kubespheredev/im:flyway-0.4.0
	AM_IMAGE=openpitrix/iam:0.4.0
	AM_FLYWAY_IMAGE=openpitrix/iam:flyway-0.4.0
	NOTIFICATION_IMAGE=openpitrix/notification:0.4.0
	NOTIFICATION_FLYWAY_IMAGE=openpitrix/notification:flyway-0.4.0"


## Usage:
## source version.sh
## export_image_version [OPENPITRIX_VERSION]

export_image_version()
{
	VERSION=$1
	OP_VERSION="V_LATEST_OPENPITRIX"
	if [[ ${VERSION} != "" ]]; then
		OP_VERSION="V_${VERSION//./_}_OPENPITRIX"
	fi

	VERSIONS=`eval echo '$'"${OP_VERSION}"`
	# check if the given version exist
	if [[ ${VERSIONS} == "" ]]; then
		echo "The version ${VERSION} of openpitrix not exist!"
		exit 1
	fi

	for V in ${VERSIONS} ; do
		echo $V
		export $V
	done

	# check dashboard_image version
	if [[ "x${VERSION}" != "xlatest" ]];then
		DASHBOARD_VERSION=${VERSION}
	    RELEASES=`curl -L -s https://api.github.com/repos/openpitrix/dashboard/releases`
	    echo ${RELEASES} | grep tag_name | sed "s/ *\"tag_name\": *\"\(.*\)\",*/\1/" | grep ${VERSION}
	    if [[ $? != 0 ]];then
			MAJOR_VERSION=`echo ${VERSION} | awk -F '.' '{print $1}'`
			for version_item in `echo ${RELEASES} | grep tag_name | sed "s/ *\"tag_name\": *\"\(.*\)\",*/\1/"`;do
			  echo version_item | grep ${MAJOR_VERSION}
			  if [ $? == 0 ];then
				DASHBOARD_VERSION=${version_item}
				break
			  fi
			done
			if [ "${DASHBOARD_VERSION}" == "" ];then
			  DASHBOARD_VERSION="latest"
			fi
	    fi
		export DASHBOARD_IMAGE="openpitrix/dashboard:${DASHBOARD_VERSION}"
	fi

}
