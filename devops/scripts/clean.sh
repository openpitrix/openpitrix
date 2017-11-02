#!/bin/bash

echo "clean..."

docker rmi openpitrix/openpitrix-dev:latest
docker rmi openpitrix

echo "cleaned successfully"