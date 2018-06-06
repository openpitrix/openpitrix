#!/bin/bash

FILE_NAME="$1"

chmod +x opt/openpitrix/sbin/*
for i in 0 1 2 3 4 5 6 7 8 9 10 11 12 13 14 15 16 17 18 19;
do
  opt/openpitrix/sbin/install_docker.sh && break || sleep 3
done

echo '{
  "registry-mirrors": ["https://registry.docker-cn.com"]
}' > /etc/docker/daemon.json

mkdir -p /opt/openpitrix/conf
mkdir -p /opt/openpitrix/log
mkdir -p /opt/openpitrix/sbin
mkdir -p /etc/confd/conf.d; mv etc/confd/conf.d/cmd.info.toml /etc/confd/conf.d/
mkdir -p /etc/confd/templates; mv etc/confd/templates/cmd.info.tmpl /etc/confd/templates/
mv opt/openpitrix/sbin/exec.sh opt/openpitrix/sbin/start_default.sh opt/openpitrix/sbin/update_fstab.sh /opt/openpitrix/sbin/
mv opt/openpitrix/conf/default.service /lib/systemd/system/default.service

if [ "$FILE_NAME" == "frontgate.conf" ]
then
  mkdir -p /opt/openpitrix/etcd
  mv opt/openpitrix/sbin/start_etcd.sh /opt/openpitrix/sbin/

  mv opt/openpitrix/conf/etcd.service /lib/systemd/system/etcd.service
  systemctl enable etcd.service && systemctl start etcd.service
fi

systemctl enable default.service && systemctl start default.service
