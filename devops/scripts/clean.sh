#!/bin/bash

echo "clean..."

docker image rm openpitrix/openpitrix:dev
docker image rm openpitrix

echo "cleaned successfully"