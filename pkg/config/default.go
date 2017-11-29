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
Host         = "127.0.0.1"
Port         = 3306
Encoding     = "utf8"
Engine       = "InnoDB"
DbName       = "openpitrix"
RootPassword = "password"

[Api]
Host = "127.0.0.1"
Port = 9100

[App]
Host = "127.0.0.1"
Port = 9101

[Runtime]
Host = "127.0.0.1"
Port = 9102

[Cluster]
Host = "127.0.0.1"
Port = 9103

[Repo]
Host = "127.0.0.1"
Port = 9104


`
