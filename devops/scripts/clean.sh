#!/bin/bash

echo "clean..."

docker image rm openpitrix/openpitrix-dev:latest
docker image rm openpitrix

echo "cleaned successfully"