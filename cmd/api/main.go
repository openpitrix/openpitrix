// Copyright 2017 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

// install httpie:
// - macOS: brew install httpie
// - Ubuntu: apt-get install httpie
// - CentOS: yum install httpie

// http get :8080/v1/apps

// curl http://localhost:8080/v1/apps

// openpitrix server
package main

import (
	"flag"
	"io/ioutil"
	"os"
	"path"
	"strings"

	"openpitrix.io/openpitrix/pkg/config"
	"openpitrix.io/openpitrix/pkg/handlers/apps"
)

var flagConfigFile = flag.String("config-file", "~/.openpitrix/config.yaml", "set config file path")

func main() {
	flag.Parse()

	cfgpath := *flagConfigFile
	if strings.HasPrefix(cfgpath, "~/") || strings.HasPrefix(cfgpath, `~\`) {
		cfgpath = config.GetHomePath() + cfgpath[1:]
	}
	if _, err := os.Stat(cfgpath); err != nil {
		if *flagConfigFile == config.DefaultConfigPath {
			os.MkdirAll(path.Dir(cfgpath), 0755)
			ioutil.WriteFile(cfgpath, []byte(config.DefaultConfigContent), 0644)
		}
	}

	cfg := config.MustLoad(cfgpath)
	apps.ListenAndServeAppsServer(cfg)
}
