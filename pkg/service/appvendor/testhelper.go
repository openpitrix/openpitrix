package appvendor

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
	cfg := config.LoadConf()
	pi.SetGlobal(cfg)
}
