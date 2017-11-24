// Copyright 2017 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

// openpitrix cluster server
package main

import (
	"openpitrix.io/openpitrix/pkg/cmd/cluster"
	"openpitrix.io/openpitrix/pkg/config"
)

func main() {
	if config.IsUserConfigExists() {
		cluster.Main(config.MustLoadUserConfig())
	} else {
		cluster.Main(config.Default())
	}
}
