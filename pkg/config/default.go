// Copyright 2017 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package config

// DefaultConfigFile is the default config file.
const DefaultConfigPath = "~/.openpitrix/config.yaml"

// DefaultConfigContent is the default config file content.
const DefaultConfigContent = `
# OpenPitrix configuration
# https://openpitrix.io/

db_type: 'mysql'
db_host: 'root:password@tcp(127.0.0.1:3306)/openpitrix'
db_encoding: 'utf8'
db_engine: 'InnoDB'

host: '127.0.0.1'
port: 8080
protocol: 'http'
uri: '/'

# Valid log levels are "debug", "info", "warn", "error", and "fatal".
log_level: 'warn'
`
