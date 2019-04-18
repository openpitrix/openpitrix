#!/usr/bin/env bash

if [ ! -f "~/.openpitrix/config.json" ];then
  mkdir -p ~/.openpitrix/
  cp config.json ~/.openpitrix/
fi