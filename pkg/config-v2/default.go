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

Host = "127.0.0.1"
Port = 8080

# Valid log levels are "debug", "info", "warn", "error", and "fatal".
LogLevel = "warn"

[Database]
Type     = "mysql"
Host     = "root:password@tcp(127.0.0.1:3306)"
Encoding = "utf8"
Engine   = "InnoDB"
DbName   = "openpitrix"
`
