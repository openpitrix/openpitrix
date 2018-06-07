#!/bin/bash

FileSystem=$1
Device=$2
MountPoint=$3
MountOptions=$4

uuid=`blkid $Device | awk -F '"' '{print $2}'`
UUID="UUID=$uuid $MountPoint $FileSystem $MountOptions 0 2"
sed -i "s/^UUID=$uuid .*//g" /etc/fstab && echo "$UUID" >> /etc/fstab