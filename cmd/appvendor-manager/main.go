// Copyright 2017 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

//openpitrix vendor manager
package main

import (
	"openpitrix.io/openpitrix/pkg/config"
	appvendor "openpitrix.io/openpitrix/pkg/service/appvendor"
)

func main() {
	cfg := config.LoadConf()
	appvendor.Serve(cfg)
}
