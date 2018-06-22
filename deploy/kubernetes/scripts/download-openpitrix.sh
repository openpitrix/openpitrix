#!/bin/bash
#
# This file will be fetched as: curl -L https://git.io/GetOpenPitrix | sh -
#
# The script fetches the latest OpenPitrix release and untars it.
# It's latest stable releases but lets
# users do curl -L https://git.io/GetOpenPitrix | OPENPITRIX_VERSION=0.1.0 sh -
# for instance to change the version fetched.

if [ "x${OPENPITRIX_VERSION}" = "x" ] ; then
  OPENPITRIX_VERSION=$(curl -L -s https://api.github.com/repos/openpitrix/openpitrix/releases/latest | grep tag_name | sed "s/ *\"tag_name\": *\"\(.*\)\",*/\1/")
fi

NAME="openpitrix-${OPENPITRIX_VERSION}-kubernetes"
URL="https://github.com/openpitrix/openpitrix/releases/download/${OPENPITRIX_VERSION}/openpitrix-${OPENPITRIX_VERSION}-kubernetes.tar.gz"
echo "Downloading $NAME from $URL ..."
curl -L "$URL" | tar xz
