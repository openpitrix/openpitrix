// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

// +build integration

package helm

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"openpitrix.io/openpitrix/pkg/constants"
	"openpitrix.io/openpitrix/pkg/util/iputil"
	"openpitrix.io/openpitrix/test/repocommon"
	"openpitrix.io/openpitrix/test/testutil"
)

const testExportPort = 8879

var testTmpDir = testutil.GetTmpDir()
var clientConfig = testutil.GetClientConfig()

func TestHelm(t *testing.T) {
	t.Logf("start create repo at [%s]", testTmpDir)

	d := testutil.NewDocker(t, "test-helm", "lachlanevenson/k8s-helm:v2.11.0")
	d.Port = testExportPort
	d.WorkDir = testutil.TmpPath
	d.Volume[testTmpDir] = testutil.TmpPath

	t.Log(d.Setup())

	t.Log(d.Exec("helm init --client-only"))

	t.Log(d.Exec("helm create nginx"))

	t.Log(d.Exec("ls nginx"))

	// TODO: write file content to testRepoDir, so that we can test create cluster

	t.Log(d.Exec("helm package nginx"))

	t.Log(d.Exec("helm repo index ./"))

	t.Log(d.Exec("cat index.yaml"))

	ip := strings.TrimSpace(d.Exec("hostname -i"))
	localIp := iputil.GetLocalIP()
	t.Log(d.ExecD(fmt.Sprintf("helm serve --address %s:%d --url http://%s:%d/", ip, testExportPort, localIp, testExportPort)))

	client := testutil.GetClient(clientConfig)

	t.Run("create repo", func(t *testing.T) {
		time.Sleep(5 * time.Second)
		repocommon.CreateRepo(t, client, "test-helm-repo-name", constants.ProviderKubernetes, fmt.Sprintf("http://%s:8879/", iputil.GetLocalIP()))
	})

	t.Run("create cluster", func(t *testing.T) {
		t.Log("TODO")
	})

	// cleanup
	t.Log(d.Exec("find . -mindepth 1 -delete"))
	t.Log(d.Teardown())
}
