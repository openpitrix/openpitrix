// Copyright 2017 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

//openpitrix isv manager
package main

import (
	"openpitrix.io/openpitrix/pkg/config"
	isv "openpitrix.io/openpitrix/pkg/service/isv"
)

func main() {
	cfg := config.GetConf()
	isv.Serve(cfg)
}
