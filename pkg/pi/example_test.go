// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package pi_test

import (
	"fmt"

	"openpitrix.io/openpitrix/pkg/config"
	"openpitrix.io/openpitrix/pkg/logger"
	"openpitrix.io/openpitrix/pkg/pi"
)

func ExampleNewPi() {
	// TODO: Automatic startup dependent services.
	cfg := config.GetConf()
	logger.SetLevelByString("debug")
	cfg.Mysql.Host = "localhost"
	cfg.Etcd.Endpoints = "localhost:2379"
	p := pi.NewPi(cfg)
	fmt.Println(p.GlobalConfig())
}
