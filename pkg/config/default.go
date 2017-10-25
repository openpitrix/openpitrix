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

database: 'mysql://root:password@127.0.0.1:3306/openpitrix'

host: '127.0.0.1'
port: 443
protocol: 'https'
uri: '/'

# Valid log levels are "debug", "info", "warn", "error", and "fatal".
log_level: 'warn'
`
