// Copyright 2017 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

// start mysql in docker for test

package apps

import (
	"flag"
	"os"
	"testing"

	"openpitrix.io/openpitrix/pkg/config"
)

var (
	tConfigFile = flag.String("config-file", "~/.openpitrix/config.yaml", "set config file path")
	tConfig     *config.Config
)

func TestMain(m *testing.M) {
	flag.Parse()

	tConfig, _ = config.Load(*tConfigFile)

	rv := m.Run()
	os.Exit(rv)
}

func TestAppDatabase(t *testing.T) {
	// TODO
}
