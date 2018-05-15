#!/bin/bash

chmod +x ./*
./install_docker.sh
mv ../conf/openpitrix.service /lib/systemd/system/openpitrix.service
mv ../conf/etcd.service /lib/systemd/system/etcd.service
systemctl enable openpitrix.service
systemctl enable etcd.service
mkdir -p /opt/openpitrix/conf
mkdir -p /opt/openpitrix/log
mkdir -p /opt/openpitrix/sbin
mkdir -p /opt/openpitrix/etcd
docker pull openpitrix/openpitrix:metadata
docker pull appcelerator/etcd
