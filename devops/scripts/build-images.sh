#!/usr/bin/env bash

# Back to the root of the project
cd $(dirname $0)
cd ../..

echo "Building images..."
docker build -t openpitrix -t openpitrix/openpitrix-dev:latest .
