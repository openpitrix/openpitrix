#!/usr/bin/env bash

if [ ! -f "${HOME}/.openpitrix/config.json" ];then
  mkdir -p ${HOME}/.openpitrix/
  cp config.json ${HOME}/.openpitrix/
fi
