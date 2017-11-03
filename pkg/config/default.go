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

[ApiService]
Host = "127.0.0.1"
Port = 8080

[AppService]
Host = "127.0.0.1"
Port = 8081

[AppRuntimeService]
Host = "127.0.0.1"
Port = 8082

[ClusterService]
Host = "127.0.0.1"
Port = 8083

[RepoService]
Host = "127.0.0.1"
Port = 8084

[Database]
Type         = "mysql"
Host         = "127.0.0.1"
Port         = 3306
Encoding     = "utf8"
Engine       = "InnoDB"
DbName       = "openpitrix"
RootPassword = "password"

`
