#!/bin/bash

echo "clean..."

docker image rm rayzhou/openpitrix-dev:latest
docker image rm openpitrix

echo "cleaned successfully"