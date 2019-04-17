#!/bin/bash

#generate global_config_map yaml file from global_config.yaml and global_config_map.tmpl

set +e

# Back to the root of the project
cd $(dirname $0)
cd ../..

CONFIG_FILE="./config/global_config.yaml"
CONFIG_MAP_TMPL="./kubernetes/etcd/global_config_map.tmpl"
CONFIG_MAP="./kubernetes/etcd/global_config_map.yaml"


usage() {
  echo "Usage:"
  echo "  format-global-config.sh [-v VERSION] [-n NAMESPACE] [-c GLOBAL_CONFIG_MAP]"
  echo "Description:"
  echo "    -c GLOBAL_CONFIG_MAP : the path of global_config_map yaml file."
  exit -1
}

while getopts c:h option
do
  case "${option}"
  in
  c) CONFIG_MAP=${OPTARG};;
  h) usage ;;
  *) usage ;;
  esac
done

if [ ! -f ${CONFIG_FILE} ]; then
  cp ./config/global_config.init.yaml ${CONFIG_FILE}
fi

sed -e "/\${GLOBAL_CONFIG}/d" ${CONFIG_MAP_TMPL} > ${CONFIG_MAP}

prefix="    "
IFS=''
while read line
do
	echo "${prefix}${line}" >> ${CONFIG_MAP}
done < ${CONFIG_FILE}

#return GLOBAL_CONFIG_MAP file
echo ${CONFIG_MAP}
