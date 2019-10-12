#!/bin/bash

set -e

mkdir -p /root/.openpitrix
cd /root/.openpitrix
rm -f config.json
touch config.json
echo -e "{" >> ./config.json
echo -e "\t\"client_id\":\"x\"," >> ./config.json
echo -e "\t\"client_secret\":\"y\"," >> ./config.json
echo -e "\t\"endpoint_url\":\"http://localhost:9100\"" >> ./config.json
echo -e "}" >> ./config.json

tars=$(ls -l /data/helm-pkg |grep ".tar" |awk '{print $NF}')
for pkg in $tars
do
    version_id=$(opctl create_app -f pkg |jq '.version_id')
    opctl submit_app_version --version_id ${version_id}
    opctl admin_pass_app_version --version_id ${version_id}
    opctl release_app_version --version_id ${version_id}
done
