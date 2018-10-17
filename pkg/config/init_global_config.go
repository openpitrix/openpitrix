// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package config

const InitialGlobalConfig = `app:
  default_draft_status: true
repo:
  # cron usage: https://godoc.org/github.com/robfig/cron#hdr-Usage
  #
  #   "@every 1h30m" means Every hour thirty
  #   "@hourly" means Every hour
  #   "0 30 * * * *" means Every hour on the half hour
  #
  #	  Field name   | Mandatory? | Allowed values  | Allowed special characters
  #	  ----------   | ---------- | --------------  | --------------------------
  #	  Seconds      | Yes        | 0-59            | * / , -
  #	  Minutes      | Yes        | 0-59            | * / , -
  #	  Hours        | Yes        | 0-23            | * / , -
  #	  Day of month | Yes        | 1-31            | * / , - ?
  #	  Month        | Yes        | 1-12 or JAN-DEC | * / , -
  #	  Day of week  | Yes        | 0-6 or SUN-SAT  | * / , - ?
  #
  cron: "0 30 4 * * *"
  max_repo_events: 20
cluster:
  frontgate_conf: '{"app_id":"app-ABCDEFGHIJKLMNOPQRST","version_id":"appv-ABCDEFGHIJKLMNOPQRST","name":"frontgate","description":"OpenPitrixbuilt-infrontgateservice","subnet":"","nodes":[{"container":{"type":"docker","image":"openpitrix/openpitrix:metadata"},"count":1,"cpu":1,"memory":1024,"volume":{"size":10,"mount_point":"/data","filesystem":"ext4"}}]}'
  frontgate_auto_delete: true
  frontgate_auto_update: false
job:
  max_working_jobs: 20
task:
  max_working_tasks: 20
pilot:
  ip: 127.0.0.1
  port: 9114
runtime:
  qingcloud_provider:
    api_server: api.qingcloud.com
    zone: .*
    image_id: xenial4x64a
    image_url: https://openpitrix.pek3a.qingstor.com/image/ubuntu.tar.gz
  aws_provider:
    api_server: .*.amazonaws.com
    zone: .*
    image_name: amzn2-ami-hvm-2.0.20180622.1-x86_64-gp2
    image_url: https://openpitrix.pek3a.qingstor.com/image/amazon-linux.tar.gz
`
