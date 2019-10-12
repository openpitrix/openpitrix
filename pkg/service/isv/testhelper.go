// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package isv

import (
	"flag"

	"openpitrix.io/openpitrix/pkg/config"
	"openpitrix.io/openpitrix/pkg/pi"
)

var (
	tTestingEnvEnabled = flag.Bool("testing-env-enabled", false, "enable testing env")
	//tTestingEnvEnabled = flag.Bool("testing-env-enabled", true, "enable testing env")
)

func InitGlobelSetting() {
	cfg := config.GetConf()
	pi.SetGlobal(cfg)
}
