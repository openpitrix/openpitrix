// Copyright 2017 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package helm

import (
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHelmRuntime(t *testing.T) {
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
