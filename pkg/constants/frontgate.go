// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package constants

// can use this id for internal test
const FrontgateVersionId = "appv-ABCDEFGHIJKLMNOPQRST"

const FrontgateDefaultConf = `
{
  "app_id": "app-ABCDEFGHIJKLMNOPQRST",
  "version_id": "appv-ABCDEFGHIJKLMNOPQRST",
  "name": "frontgate",
  "description": "OpenPitrix built-in frontgate service",
  "subnet": "",
  "nodes": [{
     "container": {
        "type": "docker",
        "image": "metadata"
     },
     "count": 3,
     "cpu": 1,
     "memory": 1024,
     "volume": {
         "size": 10,
         "mount_point": "/data",
         "filesystem": "ext4"
     }
  }]
}
`
