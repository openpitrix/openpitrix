// Copyright 2017 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

//openpitrix mbing manager
package main

import (
	"openpitrix.io/openpitrix/pkg/config"
	"openpitrix.io/openpitrix/pkg/service/mbing"
)

func main() {
	cfg := config.LoadConf()
	mbing.Serve(cfg)
}
