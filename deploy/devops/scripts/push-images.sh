#!/bin/bash

echo "Pushing images..."
docker push openpitrix/openpitrix-dev:latest
docker push openpitrix/openpitrix-dev:metadata
docker push openpitrix/openpitrix-dev:flyway
echo "Pushed successfully"
