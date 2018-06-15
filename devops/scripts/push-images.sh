#!/usr/bin/env bash

echo "pushing images..."
docker push openpitrix/openpitrix-dev:latest
docker push openpitrix/openpitrix-dev:metadata
docker push openpitrix/openpitrix-dev:flyway
