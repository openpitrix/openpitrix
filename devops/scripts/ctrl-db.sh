#!/bin/bash
git clone https://github.com/openpitrix/openpitrix.git /opt/openpitrix
# Back to the root of the project
cd $(dirname $0)
cd ../..

kubectl apply -f ./devops/kubernetes/ctrl

