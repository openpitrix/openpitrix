// Copyright 2017 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package helm

import (
	"flag"
	"log"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	config "openpitrix.io/openpitrix/pkg/config/unittest"
)

var (
	tShowEnvFlag = flag.Bool("show-env-flag", false, "show env flags")

	tConfig *config.Config
)

func TestMain(m *testing.M) {
	flag.Parse()

	if *tShowEnvFlag {
		config.PrintEnvs()
		os.Exit(0)
	}

	if conf, err := config.LoadConfig(); err == nil {
		tConfig = conf
	} else {
		log.Fatal(err)
	}

	os.Exit(m.Run())
}

func TestHelmRuntime(t *testing.T) {
	if !tConfig.Unittest.K8s.Enabled {
		t.Skip()
	}

	runtime := HelmRuntime{}

	appConf := "~/.helm/cache/archive/zookeeper-0.4.2.tgz"

	_, err := os.Stat(strings.Replace(appConf, "~/", os.Getenv("HOME")+"/", 1))
	if err != nil {
		t.Skipf("Helm runtime test skipped because no [%s], err: %v", appConf, err)
	}

	values := "{servers: 1}"

	clusterId, err := runtime.CreateCluster(appConf, true, values)
	assert.Empty(t, err)
	err = runtime.DeleteClusters(clusterId, true)
	assert.Empty(t, err)
}
