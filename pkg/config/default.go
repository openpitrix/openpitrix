// Copyright 2017 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package config

// DefaultConfigFile is the default config file.
const DefaultConfigPath = "~/.openpitrix/config.toml"

// DefaultConfigContent is the default config file content.
const DefaultConfigContent = `
# OpenPitrix configuration
# https://openpitrix.io/

[Glog]
LogToStderr       = false
AlsoLogTostderr   = false
StderrThreshold   = "ERROR" # INFO, WARNING, ERROR, FATAL
LogDir            = ""

LogBacktraceAt    = ""
V                 = 0
VModule           = ""

CopyStandardLogTo = "INFO"

[DB]
Type         = "mysql"
Host         = "openpitrix-db"
Port         = 3306
Encoding     = "utf8"
Engine       = "InnoDB"
DbName       = "openpitrix"
RootPassword = "password"

[Api]
Host = "openpitrix-api"
Port = 9100

[App]
Host = "openpitrix-app"
Port = 9101

[Runtime]
Host = "openpitrix-runtime"
Port = 9102

[Cluster]
Host = "openpitrix-cluster"
Port = 9103

[Repo]
Host = "openpitrix-repo"
Port = 9104


`
