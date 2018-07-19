#!/bin/bash -e

echo y | sudo apt update
echo y | sudo apt install \
    apt-transport-https \
    ca-certificates \
    curl \
    software-properties-common
curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo apt-key add -
echo y | sudo apt-key fingerprint 0EBFCD88
echo y | sudo add-apt-repository \
   "deb [arch=amd64] https://download.docker.com/linux/ubuntu \
   $(lsb_release -cs) \
   stable"
echo y | sudo apt update
echo y | sudo apt install docker-ce
