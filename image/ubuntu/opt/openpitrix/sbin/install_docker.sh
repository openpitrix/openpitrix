#!/bin/bash -e

rm -rf /opt/*.deb
wget -P /opt https://openpitrix.pek3a.qingstor.com/image/docker-ce_18.06.0_ce_3-0_ubuntu_amd64.deb
wget -P /opt https://openpitrix.pek3a.qingstor.com/image/libltdl7_2.4.6-0.1_amd64.deb

sudo dpkg -i /opt/libltdl7_2.4.6-0.1_amd64.deb
sudo dpkg -i /opt/docker-ce_18.06.0_ce_3-0_ubuntu_amd64.deb