// Copyright 2017 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

// openpitrix cluster server
package main

import (
	"fmt"
	"os"

	"openpitrix.io/openpitrix/pkg/cmd/cluster"
	"openpitrix.io/openpitrix/pkg/config"
	"openpitrix.io/openpitrix/pkg/version"
)

func main() {
	if len(os.Args) == 2 {
		if os.Args[1] == "-v" {
			fmt.Printf("openpitrix-cluster %s\n", version.ShortVersion)
			os.Exit(0)
		}
		if os.Args[1] == "-version" {
			fmt.Printf("openpitrix-cluster %s, build date %s\n", version.GitSha1Version, version.BuildDate)
			os.Exit(0)
		}
	}

	if config.RunInDocker() {
		config.UseDockerLinkedEnvironmentVariables()
	}

	if config.IsUserConfigExists() {
		cluster.Main(config.MustLoadUserConfig())
	} else {
		cluster.Main(config.Default())
	}
}
