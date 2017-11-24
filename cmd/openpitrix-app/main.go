// Copyright 2017 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

// openpitrix app server
package main

import (
	"openpitrix.io/openpitrix/pkg/cmd/app"
	"openpitrix.io/openpitrix/pkg/config"
)

func main() {
	if config.IsUserConfigExists() {
		app.Main(config.MustLoadUserConfig())
	} else {
		app.Main(config.Default())
	}
}
