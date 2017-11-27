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
	"fmt"
	"os"

	"openpitrix.io/openpitrix/pkg/cmd/api"
	"openpitrix.io/openpitrix/pkg/config"
	"openpitrix.io/openpitrix/pkg/version"
)

func main() {
	if len(os.Args) == 2 {
		if os.Args[1] == "-v" {
			fmt.Printf("openpitrix-api %s\n", version.ShortVersion)
			os.Exit(0)
		}
		if os.Args[1] == "-version" {
			fmt.Printf("openpitrix-api %s, build date %s\n", version.GitSha1Version, version.BuildDate)
			os.Exit(0)
		}
	}

	if config.IsUserConfigExists() {
		api.Main(config.MustLoadUserConfig())
	} else {
		api.Main(config.Default())
	}
}
