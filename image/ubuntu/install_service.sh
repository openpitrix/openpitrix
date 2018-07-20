#!/bin/bash -e

FILE_NAME="$1"

chmod +x opt/openpitrix/sbin/*

for i in $(seq 1 100)
do
  opt/openpitrix/sbin/install_docker.sh && break || sleep 3
done

mkdir -p /opt/openpitrix/conf
mkdir -p /opt/openpitrix/log
mkdir -p /opt/openpitrix/sbin
mkdir -p /etc/confd/conf.d; cp etc/confd/conf.d/cmd.info.toml /etc/confd/conf.d/
mkdir -p /etc/confd/templates; cp etc/confd/templates/cmd.info.tmpl /etc/confd/templates/
cp opt/openpitrix/sbin/exec.sh opt/openpitrix/sbin/start_default.sh opt/openpitrix/sbin/update_fstab.sh /opt/openpitrix/sbin/
cp opt/openpitrix/conf/default.service /lib/systemd/system/default.service

if [ "$FILE_NAME" == "frontgate.conf" ]
then
  mkdir -p /opt/openpitrix/etcd
  cp opt/openpitrix/sbin/start_etcd.sh /opt/openpitrix/sbin/
  cp opt/openpitrix/conf/etcd.service /lib/systemd/system/etcd.service
  systemctl enable etcd.service && systemctl start etcd.service

  cp opt/openpitrix/sbin/start_registry.sh /opt/openpitrix/sbin/
  cp opt/openpitrix/conf/config.yml /opt/openpitrix/conf/config.yml
  cp opt/openpitrix/conf/registry.service /lib/systemd/system/registry.service
  systemctl enable registry.service && systemctl start registry.service
fi

systemctl enable default.service && systemctl start default.service
