// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

// +build integration

package devkit

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

const testExportPort = 9191

var testTmpDir = testutil.GetTmpDir()
var clientConfig = testutil.GetClientConfig()

func TestDevkit(t *testing.T) {
	t.Logf("start create repo at [%s]", testTmpDir)

	d := testutil.NewDocker(t, "test-op", "openpitrix")
	d.Port = testExportPort
	d.WorkDir = testutil.TmpPath
	d.Volume[testTmpDir] = testutil.TmpPath

	t.Log(d.Setup())

	t.Log(d.Exec("op create nginx"))

	t.Log(d.Exec("ls nginx"))

	// TODO: write file content to testRepoDir, so that we can test create cluster

	t.Log(d.Exec("op package nginx"))

	t.Log(d.Exec("op index ./"))

	t.Log(d.Exec("cat index.yaml"))

	ip := strings.TrimSpace(d.Exec("hostname -i"))
	localIp := iputil.GetLocalIP()
	t.Log(d.ExecD(fmt.Sprintf("op serve --address %s:%d --url http://%s:%d/", ip, testExportPort, localIp, testExportPort)))

	client := testutil.GetClient(clientConfig)

	t.Run("create repo", func(t *testing.T) {
		time.Sleep(5 * time.Second)
		repocommon.CreateRepo(t, client, "test-devkit-repo-name", constants.ProviderQingCloud, fmt.Sprintf("http://%s:9191/", iputil.GetLocalIP()))
	})

	t.Run("create cluster", func(t *testing.T) {
		t.Log("TODO")
	})

	// cleanup
	t.Log(d.Exec("find . -mindepth 1 -delete"))
	t.Log(d.Teardown())
}
