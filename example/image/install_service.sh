#!/bin/bash

chmod +x opt/openpitrix/sbin/*
opt/openpitrix/sbin/install_docker.sh
mv opt/openpitrix/conf/default.service /lib/systemd/system/default.service
mv opt/openpitrix/conf/etcd.service /lib/systemd/system/etcd.service
systemctl enable default.service
systemctl enable etcd.service
mkdir -p /opt/openpitrix/conf
mkdir -p /opt/openpitrix/log
mkdir -p /opt/openpitrix/sbin
mkdir -p /opt/openpitrix/etcd
mkdir -p /etc/confd/conf.d; mv etc/confd/conf.d/cmd.info.toml /etc/confd/conf.d/
mkdir -p /etc/confd/templates; mv etc/confd/templates/cmd.info.tmpl /etc/confd/templates/
docker pull appcelerator/etcd
mv opt/openpitrix/sbin/exec.sh opt/openpitrix/sbin/start_default.sh opt/openpitrix/sbin/start_etcd.sh /opt/openpitrix/sbin/

