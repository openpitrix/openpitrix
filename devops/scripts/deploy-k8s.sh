#!/bin/bash

# Back to the root of the project
cd $(dirname $0)
cd ../..

kubectl apply -f ./devops/kubernetes/openpitrix.yaml
