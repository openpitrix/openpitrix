// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

// +build integration

package test

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"openpitrix.io/openpitrix/pkg/constants"
	"openpitrix.io/openpitrix/pkg/utils/iptool"
)

func TestHelm(t *testing.T) {
	t.Logf("start create repo at [%s]", testRepoDir)

	d := NewDocker(t, "test-helm", "lachlanevenson/k8s-helm:v2.8.2")
	d.Port = testExportPort
	d.WorkDir = testDockerPath
	d.Volume[testRepoDir] = testDockerPath

	t.Log(d.Setup())

	t.Log(d.Exec("helm init --client-only"))

	t.Log(d.Exec("helm create nginx"))

	t.Log(d.Exec("ls nginx"))

	// TODO: write file content to testRepoDir, so that we can test create cluster

	t.Log(d.Exec("helm package nginx"))

	t.Log(d.Exec("helm repo index ./"))

	t.Log(d.Exec("cat index.yaml"))

	ip := strings.TrimSpace(d.Exec("hostname -i"))
	localIp := iptool.GetLocalIP()
	t.Log(d.ExecD(fmt.Sprintf("helm serve --address %s:%d --url http://%s:%d/", ip, testExportPort, localIp, testExportPort)))

	t.Run("create repo", func(t *testing.T) {
		time.Sleep(5 * time.Second)
		testCreateRepo(t, "test-helm-repo-name", constants.ProviderKubernetes, fmt.Sprintf("http://%s:8879/", iptool.GetLocalIP()))
	})

	t.Run("create cluster", func(t *testing.T) {
		t.Log("TODO")
	})

	// cleanup
	t.Log(d.Exec("find . -mindepth 1 -delete"))
	t.Log(d.Teardown())
}
