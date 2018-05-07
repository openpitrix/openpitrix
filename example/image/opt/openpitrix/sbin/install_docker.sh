#!/bin/bash

echo y | sudo apt-get remove docker docker-engine docker.io
echo y | sudo apt-get update
echo y | sudo apt-get install \
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
echo y | sudo apt-get update
echo y | sudo apt-get install docker-ce
