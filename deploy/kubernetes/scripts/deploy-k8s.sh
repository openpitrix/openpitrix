#!/bin/bash

echo "Deploying k8s resource..."

[ -z `which kubectl` ] && echo "Deployed failed: You need to install 'kubectl' first." && exit 1

# Back to the root of the project
cd $(dirname $0)
cd ../..

DEFAULT_NAMESPACE=openpitrix-system

NAMESPACE=${DEFAULT_NAMESPACE}
VERSION=""
DASHBOARD_VERSION=""
METADATA=0
DBCTRL=0
BASE=0
DASHBOARD=0
INGRESS=0
STORAGE=0
ALL=0
JOB_REPLICA=1
TASK_REPLICA=1
API_NODEPORT=""
PILOT_NODEPORT=""
DASHBOARD_NODEPORT=""

# use nodePort for api/pilot/dashboard service
# $cat ${NODEPORT_FILE}
# API_NODEPORT=31009
# PILOT_NODEPORT=31019
# DASHBOARD_NODEPORT=31029
NODEPORT_FILE="./config/node_port"
if [ -f ${NODEPORT_FILE} ];then
  API_NODEPORT=`cat ${NODEPORT_FILE} | grep API_NODEPORT | awk -F '=' '{print "nodePort: "$2}'`
  PILOT_NODEPORT=`cat ${NODEPORT_FILE} | grep PILOT_NODEPORT | awk -F '=' '{print "nodePort: "$2}'`
  DASHBOARD_NODEPORT=`cat ${NODEPORT_FILE} | grep DASHBOARD_NODEPORT | awk -F '=' '{print "nodePort: "$2}'`
  WEBSOCKET_PORT=`cat ${NODEPORT_FILE} | grep WEBSOCKET_NODEPORT | awk -F '=' '{print $2}'`
fi

if [[ ! -n "${WEBSOCKET_PORT}" ]]; then
  WEBSOCKET_PORT="30300"
fi
WEBSOCKET_NODEPORT="nodePort: "+${WEBSOCKET_PORT}

CPU_REQUESTS=100
MEMORY_REQUESTS=100
CPU_LIMITS=500
MEMORY_LIMITS=500

REQUESTS=""
LIMITS=""

PROVIDER_PLUGINS=""

usage() {
  echo "Usage:"
  echo "  deploy-k8s.sh [-n NAMESPACE] [-v VERSION] COMMAND"
  echo "Description:"
  echo "        -n NAMESPACE    : the namespace of kubernetes."
  echo "        -v VERSION      : the version to be deployed."
  echo "        -r REQUESTS     : the requests of container resources. such as: cpu=100,memory=200, default is: cpu=100,memory=100"
  echo "        -l LIMITS       : the limits of container resources. such as: cpu=100,memory=200, default is: cpu=500,memory=500"
  echo "        -j JOB REPLICA  : the job replica number."
  echo "        -t TASK REPLICA : the task replica number."
  echo "        -c HOST         : the hostname used in ingress."
  echo "        -p PROVIDER     : the runtime provider plugin. such as: qingcloud,aws. such as: all"
  echo "        -b              : base model will be applied."
  echo "        -m              : metadata will be applied."
  echo "        -d              : dbctrl will be applied."
  echo "        -u              : ui/dashboard will be applied."
  echo "        -i              : ingress will be applied."
  echo "        -s              : storage will be applied."
  echo "        -a              : all of base/metadata/dbctrl/dashboard/storage/ingress will be applied."
  exit -1
}


while getopts :c:n:v:r:l:j:t:p:hbdmsuia option
do
  case "${option}"
  in
  c) HOST=${OPTARG};;
  n) NAMESPACE=${OPTARG};;
  v) VERSION=${OPTARG};;
  r) REQUESTS=${OPTARG};;
  l) LIMITS=${OPTARG};;
  j) JOB_REPLICA=${OPTARG};;
  t) TASK_REPLICA=${OPTARG};;
  p) PROVIDER_PLUGINS=${OPTARG};;
  d) DBCTRL=1;;
  m) METADATA=1;;
  b) BASE=1;;
  u) DASHBOARD=1;;
  i) INGRESS=1;;
  s) STORAGE=1;;
  a) ALL=1;;
  h) usage ;;
  *) usage ;;
  esac
done

if [ "${METADATA}" == "0" ] && \
   [ "${DBCTRL}" == "0" ] && \
   [ "${BASE}" == "0" ] && \
   [ "${DASHBOARD}" == "0" ] && \
   [ "${INGRESS}" == "0" ] && \
   [ "${STORAGE}" == "0" ] && \
   [ "${PROVIDER_PLUGINS}" == "" ] && \
   [ "${ALL}" == "0" ]
then
  usage
fi

resource=""
get_resource() {
  key=${2}
  resource=""
  split=`echo ${1} | awk -F ',' '{ for(i=1;i<=NF;i++) {print $i}}'`
  for item in `echo ${split}`;do
    value=`echo ${item} | awk -F '=' '{if ($1=="'${key}'") print $2}'`
    if [ "${value}" != "" ]
    then
      resource=${value}
      break
    fi
  done
}

apply_yaml() {
  version=${1}
  file=${2}
  if [ "${version}" == "latest" ];then
    replace ./kubernetes/openpitrix/${file} | kubectl delete -f - --ignore-not-found=true
    replace ./kubernetes/openpitrix/${file} | kubectl apply -f -
  else
    replace ./kubernetes/openpitrix/${file} | kubectl apply -f -
    if [ $? -ne 0 ]; then
      replace ./kubernetes/openpitrix/${file} | kubectl delete -f - --ignore-not-found=true
      replace ./kubernetes/openpitrix/${file} | kubectl apply -f -
    fi
  fi
}

get_resource ${REQUESTS} "cpu"
if [ "${resource}" != "" ]
then
  CPU_REQUESTS=${resource}
fi

get_resource ${REQUESTS} "memory"
if [ "${resource}" != "" ]
then
  MEMORY_REQUESTS=${resource}
fi

get_resource ${LIMITS} "cpu"
if [ "${resource}" != "" ]
then
  CPU_LIMITS=${resource}
fi

get_resource ${LIMITS} "memory"
if [ "${resource}" != "" ]
then
  MEMORY_LIMITS=${resource}
fi

echo "Deploy resources: CPU_REQUESTS(${CPU_REQUESTS}), MEMORY_REQUESTS(${MEMORY_REQUESTS}), CPU_LIMITS(${CPU_LIMITS}), MEMORY_LIMITS(${MEMORY_LIMITS})"

if [ "${VERSION}" == "" ];then
  VERSION=$(curl -L -s https://api.github.com/repos/openpitrix/openpitrix/releases/latest | grep tag_name | sed "s/ *\"tag_name\": *\"\(.*\)\",*/\1/")
fi

IM_IMAGE="kubespheredev/im:latest"
IM_FLYWAY_IMAGE="kubespheredev/im:flyway"
AM_IMAGE="openpitrix/iam:latest"
AM_FLYWAY_IMAGE="openpitrix/iam:flyway"
NOTIFICATION_IMAGE="openpitrix/notification:latest"
NOTIFICATION_FLYWAY_IMAGE="openpitrix/notification:flyway"

if [ "${VERSION}" == "dev" ];then
  IMAGE="openpitrix/openpitrix-dev:latest"
  METADATA_IMAGE="openpitrix/openpitrix-dev:metadata"
  FLYWAY_IMAGE="openpitrix/openpitrix-dev:flyway"
  DASHBOARD_IMAGE="openpitrix/dashboard:latest"
  IMAGE_PULL_POLICY="Always"
elif [ "${VERSION}" == "latest" ];then
  IMAGE="openpitrix/openpitrix:latest"
  METADATA_IMAGE="openpitrix/openpitrix:metadata"
  FLYWAY_IMAGE="openpitrix/openpitrix:flyway"
  DASHBOARD_IMAGE="openpitrix/dashboard:latest"
  IMAGE_PULL_POLICY="Always"
else
  IMAGE="openpitrix/openpitrix:${VERSION}"
  METADATA_IMAGE="openpitrix/openpitrix:metadata-${VERSION}"
  FLYWAY_IMAGE="openpitrix/openpitrix:flyway-${VERSION}"
  curl -L -s https://api.github.com/repos/openpitrix/dashboard/releases | grep tag_name | sed "s/ *\"tag_name\": *\"\(.*\)\",*/\1/" | grep ${VERSION}
  if [ $? == 0 ];then
    DASHBOARD_IMAGE="openpitrix/dashboard:${VERSION}"
  else
  	MAJOR_VERSION=`echo ${VERSION} | awk -F '.' '{print $1}'`
    for version_item in `curl -L -s https://api.github.com/repos/openpitrix/dashboard/releases | grep tag_name | sed "s/ *\"tag_name\": *\"\(.*\)\",*/\1/"`;do
      echo version_item | grep ${MAJOR_VERSION}
      if [ $? == 0 ];then
        DASHBOARD_VERSION=${version_item}
        break
      fi
    done
    if [ "${DASHBOARD_VERSION}" == "" ];then
      DASHBOARD_VERSION="latest"
    fi
    DASHBOARD_IMAGE="openpitrix/dashboard:${DASHBOARD_VERSION}"
  fi
  IMAGE_PULL_POLICY="IfNotPresent"
fi

replace() {
  sed -e "s!\${NAMESPACE}!${NAMESPACE}!g" \
      -e "s!\${IM_IMAGE}!${IM_IMAGE}!g" \
      -e "s!\${AM_IMAGE}!${AM_IMAGE}!g" \
      -e "s!\${IM_FLYWAY_IMAGE}!${IM_FLYWAY_IMAGE}!g" \
      -e "s!\${IM_FLYWAY_IMAGE}!${AM_FLYWAY_IMAGE}!g" \
      -e "s!\${NOTIFICATION_IMAGE}!${NOTIFICATION_IMAGE}!g" \
      -e "s!\${NOTIFICATION_FLYWAY_IMAGE}!${NOTIFICATION_FLYWAY_IMAGE}!g" \
	  -e "s!\${IMAGE}!${IMAGE}!g" \
	  -e "s!\${DASHBOARD_IMAGE}!${DASHBOARD_IMAGE}!g" \
	  -e "s!\${METADATA_IMAGE}!${METADATA_IMAGE}!g" \
	  -e "s!\${FLYWAY_IMAGE}!${FLYWAY_IMAGE}!g" \
	  -e "s!\${CPU_REQUESTS}!${CPU_REQUESTS}!g" \
	  -e "s!\${MEMORY_REQUESTS}!${MEMORY_REQUESTS}!g" \
	  -e "s!\${CPU_LIMITS}!${CPU_LIMITS}!g" \
	  -e "s!\${MEMORY_LIMITS}!${MEMORY_LIMITS}!g" \
	  -e "s!\${JOB_REPLICA}!${JOB_REPLICA}!g" \
	  -e "s!\${TASK_REPLICA}!${TASK_REPLICA}!g" \
	  -e "s!\${VERSION}!${VERSION}!g" \
	  -e "s!\${IMAGE_PULL_POLICY}!${IMAGE_PULL_POLICY}!g" \
	  -e "s!\${API_NODEPORT}!${API_NODEPORT}!g" \
	  -e "s!\${PILOT_NODEPORT}!${PILOT_NODEPORT}!g" \
	  -e "s!\${DASHBOARD_NODEPORT}!${DASHBOARD_NODEPORT}!g" \
	  -e "s!\${WEBSOCKET_PORT}!${WEBSOCKET_PORT}!g" \
	  -e "s!\${WEBSOCKET_NODEPORT}!${WEBSOCKET_NODEPORT}!g" \
	  -e "s!\${HOST}!${HOST}!g" \
	  $1
}

[ -z `which make` ] && echo "Deployed failed: You need to install 'make' first." && exit 1

[ "${HOST}" == "" ] && HOST=demo.openpitrix.io

kubectl get ns ${NAMESPACE}
if [ $? != 0 ];then
  kubectl create namespace ${NAMESPACE}
fi

if [ "${STORAGE}" == "1" ] || [ "${ALL}" == "1" ];then
  kubectl create secret generic mysql-pass --from-file=./kubernetes/password.txt -n ${NAMESPACE}
  for FILE in `ls ./kubernetes/db/ | grep -v "\-job.yaml$"`;do
    replace ./kubernetes/db/${FILE} | kubectl apply -f -
  done

  for FILE in `ls ./kubernetes/etcd/`;do
    replace ./kubernetes/etcd/${FILE} | kubectl apply -f -
  done

  for FILE in `ls ./kubernetes/minio/`;do
    replace ./kubernetes/minio/${FILE} | kubectl apply -f -
  done
fi
if [ "${DBCTRL}" == "1" ] || [ "${ALL}" == "1" ];then
  for FILE in `ls ./kubernetes/ctrl`;do
    replace ./kubernetes/ctrl/${FILE} | kubectl delete -f - --ignore-not-found=true
    replace ./kubernetes/ctrl/${FILE} | kubectl apply -f -
  done

  for FILE in `ls ./kubernetes/db/ | grep "\-job.yaml$"`;do
    replace ./kubernetes/db/${FILE} | kubectl delete -f - --ignore-not-found=true
    replace ./kubernetes/db/${FILE} | kubectl apply -f -
  done
fi
if [ "${BASE}" == "1" ] || [ "${ALL}" == "1" ];then
  cd ./kubernetes/iam-config
  make
  cd ../..

  kubectl create secret generic iam-secret-key --from-file=./kubernetes/iam-config/secret-key.txt -n ${NAMESPACE}
  kubectl create secret generic iam-client -n ${NAMESPACE} \
  	--from-file=./kubernetes/iam-config/client-id.txt \
  	--from-file=./kubernetes/iam-config/client-secret.txt
  kubectl create secret generic iam-account -n ${NAMESPACE} --from-file=./kubernetes/iam-config/admin-password.txt

  # market service temporary not needed
  for FILE in `ls ./kubernetes/openpitrix/ | grep "^openpitrix-" | grep -v "^openpitrix-market"`; do
    apply_yaml ${VERSION} ${FILE}
  done
fi
if [ "${METADATA}" == "1" ] || [ "${ALL}" == "1" ];then
  ./kubernetes/scripts/generate-certs.sh -n ${NAMESPACE} -c ${HOST}
  if [ $? -ne 0 ]; then
	echo "Deploy failed."
	exit 1
  fi

  for FILE in `ls ./kubernetes/openpitrix/metadata/`;do
    apply_yaml ${VERSION} metadata/${FILE}
  done
fi
if [ "${DASHBOARD}" == "1" ] || [ "${ALL}" == "1" ];then
  for FILE in `ls ./kubernetes/openpitrix/dashboard/`;do
    apply_yaml ${VERSION} dashboard/${FILE}
  done
fi
if [ "${PROVIDER_PLUGINS}" != "" ] || [ "${ALL}" == "1" ];then
  if [ "${PROVIDER_PLUGINS}" == "" ] || [ "${PROVIDER_PLUGINS}" == "all" ];then
    for FILE in `ls ./kubernetes/openpitrix/plugin/`;do
      apply_yaml ${VERSION} plugin/${FILE}
    done
  else
    plugin=`echo ${PROVIDER_PLUGINS} | awk -F ',' '{ for(i=1;i<=NF;i++) {print $i}}'`
    for item in `echo ${plugin}`;do
      for FILE in `ls ./kubernetes/openpitrix/plugin/ | grep "\-${item}.yaml"`;do
        echo $FILE
    	apply_yaml ${VERSION} plugin/${FILE}
  	  done
    done
  fi
fi

if [ "${INGRESS}" == "1" ] || [ "${ALL}" == "1" ];then
  kubectl get ns ingress-nginx
  if [ $? != 0 ];then
    kubectl create namespace ingress-nginx
  fi

  ./kubernetes/scripts/generate-certs.sh -n ${NAMESPACE} -c ${HOST}
  for FILE in `ls ./kubernetes/openpitrix/ingress/`;do
    apply_yaml ${VERSION} ingress/${FILE}
  done
fi

echo "Deployed successfully"
