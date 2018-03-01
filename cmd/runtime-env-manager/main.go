// Copyright 2017 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

// openpitrix runtime-env manager
package main

import (
	"openpitrix.io/openpitrix/pkg/config"
	"openpitrix.io/openpitrix/pkg/manager/runtime_env"
)

func main() {
	cfg := config.LoadConf()
	runtime_env.Serve(cfg)
}
