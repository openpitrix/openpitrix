#!/bin/bash

# Back to the root of the project
cd $(dirname $0)
cd ../../..

echo "Building images..."
docker build -t openpitrix -t openpitrix/openpitrix-dev:latest .
docker build -t openpitrix/openpitrix-dev:metadata -f ./Dockerfile.metadata .
cd ./pkg/db/ && docker build -t openpitrix/openpitrix-dev:flyway -f ./Dockerfile .
echo "Built successfully"
