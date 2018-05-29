#!/bin/bash

echo "clean all stuff"

echo "start to clean k8s resource"
# Back to the root of the project
cd $(dirname $0)
cd ../..

kubectl delete namespace openpitrix-system


echo "start to clean docker resource"
docker rmi openpitrix/openpitrix-dev:latest
docker rmi openpitrix

echo "cleaned successfully"
