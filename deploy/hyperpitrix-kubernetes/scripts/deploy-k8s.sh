#!/bin/bash

echo "Deploying k8s resource..."

[ -z `which kubectl` ] && echo "Deployed failed: You need to install 'kubectl' first." && exit 1

# Back to the root of the project
cd $(dirname $0)
cd ../..

DEFAULT_NAMESPACE=openpitrix-system

NAMESPACE=${DEFAULT_NAMESPACE}
VERSION=""
DBCTRL=0
BASE=0
STORAGE=0
ALL=0
JOB_REPLICA=1
TASK_REPLICA=1
OPENPITRIX_LOG_LEVEL="info"
OPENPITRIX_MYSQL_HOST="openpitrix-db"
OPENPITRIX_ETCD_ENDPOINTS="openpitrix-etcd:2379"
DB_SERVICE="openpitrix-db.${NAMESPACE}.svc"
ETCD_SERVICE="openpitrix-etcd.${NAMESPACE}.svc"
MINIO_SERVICE="minio.kubesphere-system.svc"
DB_LOG_MODE_ENABLE="true"
GRPC_SHOW_ERROR_CAUSE="true"
CPU_REQUESTS=100
MEMORY_REQUESTS=100
CPU_LIMITS=500
MEMORY_LIMITS=500
OPENPITRIX_ATTACHMENT_ENDPOINT='http://minio.openpitrix-system.svc:9000'
OPENPITRIX_ATTACHMENT_BUCKET_NAME='openpitrix-attachment'

REQUESTS=""
LIMITS=""

PROVIDER_PLUGINS=""
REPO_DIR=""
BUSYBOX=busybox:1.28.4
RELEASE_APP_IMAGE="openpitrix/release-app:latest"

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
  echo "        -p PROVIDER     : the runtime provider plugin. such as: qingcloud,aws. such as: all"
  echo "        -e REPO_DIR     : the other repo dir"
  echo "        -b              : base model will be applied."
  echo "        -c              : dbctrl will be applied."
  echo "        -d              : set openpitrix log level to debug."
  echo "        -s              : storage will be applied."
  echo "        -a              : all of base/dbctrl/storage will be applied."
  exit -1
}


while getopts :n:v:r:l:j:t:p:e:m:hbcsa option
do
  case "${option}"
  in
  n) NAMESPACE=${OPTARG};;
  v) VERSION=${OPTARG};;
  r) REQUESTS=${OPTARG};;
  l) LIMITS=${OPTARG};;
  j) JOB_REPLICA=${OPTARG};;
  t) TASK_REPLICA=${OPTARG};;
  p) PROVIDER_PLUGINS=${OPTARG};;
  e) REPO_DIR=${OPTARG};;
  c) DBCTRL=1;;
  d) OPENPITRIX_LOG_LEVEL="debug";;
  b) BASE=1;;
  s) STORAGE=1;;
  a) ALL=1;;
  h) usage ;;
  *) usage ;;
  esac
done

if [ "${DBCTRL}" == "0" ] && \
   [ "${BASE}" == "0" ] && \
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
    replace ./hyperpitrix-kubernetes/openpitrix/${file} | kubectl delete -f - --ignore-not-found=true
    replace ./hyperpitrix-kubernetes/openpitrix/${file} | kubectl apply -f -
  else
    replace ./hyperpitrix-kubernetes/openpitrix/${file} | kubectl apply -f -
    if [ $? -ne 0 ]; then
      replace ./hyperpitrix-kubernetes/openpitrix/${file} | kubectl delete -f - --ignore-not-found=true
      replace ./hyperpitrix-kubernetes/openpitrix/${file} | kubectl apply -f -
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

## export image versions
VERSION_IMAGES=`./version.sh openpitrix-${VERSION}`
if [ $? == 0 ]; then
  export ${VERSION_IMAGES}
else
  # echo error message
  echo ${VERSION_IMAGES}
  exit 1
fi

if [ "x${VERSION}" == "xlatest" ];then
  IMAGE_PULL_POLICY=Always
else
  IMAGE_PULL_POLICY=IfNotPresent
fi

replace() {
  sed -e "s!\${NAMESPACE}!${NAMESPACE}!g" \
	  -e "s!\${VERSION}!${VERSION}!g" \
	  -e "s!\${IMAGE}!${IMAGE}!g" \
	  -e "s!\${FLYWAY_IMAGE}!${FLYWAY_IMAGE}!g" \
	  -e "s!\${RP_K8S_VERSION}!${RP_K8S_VERSION}!g" \
	  -e "s!\${RP_K8S_IMAGE}!${RP_K8S_IMAGE}!g" \
	  -e "s!\${WATCHER_VERSION}!${WATCHER_VERSION}!g" \
	  -e "s!\${WATCHER_IMAGE}!${WATCHER_IMAGE}!g" \
	  -e "s!\${CPU_REQUESTS}!${CPU_REQUESTS}!g" \
	  -e "s!\${MEMORY_REQUESTS}!${MEMORY_REQUESTS}!g" \
	  -e "s!\${CPU_LIMITS}!${CPU_LIMITS}!g" \
	  -e "s!\${MEMORY_LIMITS}!${MEMORY_LIMITS}!g" \
	  -e "s!\${JOB_REPLICA}!${JOB_REPLICA}!g" \
	  -e "s!\${TASK_REPLICA}!${TASK_REPLICA}!g" \
	  -e "s!\${IMAGE_PULL_POLICY}!${IMAGE_PULL_POLICY}!g" \
	  -e "s!\${OPENPITRIX_LOG_LEVEL}!${OPENPITRIX_LOG_LEVEL}!g" \
	  -e "s!\${DB_LOG_MODE_ENABLE}!${DB_LOG_MODE_ENABLE}!g" \
	  -e "s!\${GRPC_SHOW_ERROR_CAUSE}!${GRPC_SHOW_ERROR_CAUSE}!g" \
	  -e "s!\${REPO_DIR}!${REPO_DIR}!g" \
	  -e "s!\${OPENPITRIX_MYSQL_HOST}!${OPENPITRIX_MYSQL_HOST}!g" \
	  -e "s!\${OPENPITRIX_ETCD_ENDPOINTS}!${OPENPITRIX_ETCD_ENDPOINTS}!g" \
	  -e "s!\${DB_SERVICE}!${DB_SERVICE}!g" \
	  -e "s!\${ETCD_SERVICE}!${ETCD_SERVICE}!g" \
	  -e "s!\${MINIO_SERVICE}!${MINIO_SERVICE}!g" \
	  -e "s!\${OPENPITRIX_ATTACHMENT_ENDPOINT}!${OPENPITRIX_ATTACHMENT_ENDPOINT}!g" \
	  -e "s!\${OPENPITRIX_ATTACHMENT_BUCKET_NAME}!${OPENPITRIX_ATTACHMENT_BUCKET_NAME}!g" \
	  -e "s!\${RELEASE_APP_IMAGE}!${RELEASE_APP_IMAGE}!g" \
	  -e "s!\${BUSYBOX}!${BUSYBOX}!g" \
	  $1
}

DELETE_DEPLOYMENT=(
"openpitrix-api-gateway-deployment"
"openpitrix-app-manager-deployment"
"openpitrix-category-manager-deployment"
"openpitrix-cluster-manager-deployment"
"openpitrix-dashboard-deployment"
"openpitrix-iam-service-deployment"
"openpitrix-job-manager-deployment"
"openpitrix-pilot-deployment"
"openpitrix-repo-indexer-deployment"
"openpitrix-repo-manager-deployment"
"openpitrix-runtime-manager-deployment"
"openpitrix-task-manager-deployment"
)

DELETE_SERVICE=(
"openpitrix-api-gateway"
"openpitrix-app-manager"
"openpitrix-category-manager"
"openpitrix-cluster-manager"
"openpitrix-dashboard"
"openpitrix-iam-service"
"openpitrix-job-manager"
"openpitrix-pilot-service"
"openpitrix-repo-indexer"
"openpitrix-repo-manager"
"openpitrix-runtime-manager"
"openpitrix-task-manager"
)


for item in ${DELETE_DEPLOYMENT[@]}
do
    kubectl delete deployment ${item} -n ${NAMESPACE}
    if [ $? -ne 0 ]; then
      echo "Delete deployment ${item} failed."
    fi
done

for item in ${DELETE_SERVICE[@]}
do
    kubectl delete svc ${item} -n ${NAMESPACE}
    if [ $? -ne 0 ]; then
      echo "Delete svc ${item} failed."
    fi
done

[ -z `which make` ] && echo "Deployed failed: You need to install 'make' first." && exit 1


kubectl get ns ${NAMESPACE}
if [ $? != 0 ];then
  kubectl create namespace ${NAMESPACE}
fi

if [ "${STORAGE}" == "1" ] || [ "${ALL}" == "1" ];then
  kubectl create secret generic mysql-pass --from-file=./hyperpitrix-kubernetes/password.txt -n ${NAMESPACE}
  for FILE in `ls ./hyperpitrix-kubernetes/db/`;do
    replace ./hyperpitrix-kubernetes/db/${FILE} | kubectl apply -f -
  done
  ./hyperpitrix-kubernetes/scripts/generate-config-map.sh
  for FILE in `ls ./hyperpitrix-kubernetes/etcd/`;do
    if [ "x${FILE##*.}" == "xyaml" ]; then
      replace ./hyperpitrix-kubernetes/etcd/${FILE} | kubectl apply -f -
    fi
  done
  for FILE in `ls ./hyperpitrix-kubernetes/minio/`;do
    replace ./hyperpitrix-kubernetes/minio/${FILE} | kubectl apply -f -
  done
fi
if [ "${BASE}" == "1" ] || [ "${ALL}" == "1" ];then
  for FILE in `ls ./hyperpitrix-kubernetes/openpitrix/`; do
    apply_yaml ${VERSION} ${FILE}
  done
fi
if [ "${PROVIDER_PLUGINS}" != "" ] || [ "${ALL}" == "1" ];then
  if [ "${PROVIDER_PLUGINS}" == "" ] || [ "${PROVIDER_PLUGINS}" == "all" ];then
    for FILE in `ls ./hyperpitrix-kubernetes/openpitrix/plugin/`;do
      apply_yaml ${VERSION} plugin/${FILE}
    done
  else
    plugin=`echo ${PROVIDER_PLUGINS} | awk -F ',' '{ for(i=1;i<=NF;i++) {print $i}}'`
    for item in `echo ${plugin}`;do
      for FILE in `ls ./hyperpitrix-kubernetes/openpitrix/plugin/ | grep "\-${item}.yaml"`;do
        echo $FILE
        apply_yaml ${VERSION} plugin/${FILE}
      done
    done
  fi
fi

echo "Deployed successfully"
