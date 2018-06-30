// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

// +build integration

package test

import (
	"fmt"
	"path"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"openpitrix.io/openpitrix/pkg/constants"
	"openpitrix.io/openpitrix/pkg/util/idutil"
	"openpitrix.io/openpitrix/pkg/util/iputil"
	"openpitrix.io/openpitrix/test/client/app_manager"
	"openpitrix.io/openpitrix/test/client/repo_indexer"
	"openpitrix.io/openpitrix/test/client/repo_manager"
	"openpitrix.io/openpitrix/test/models"
)

const testExportPort = 8879
const testDockerPath = "/tmp/openpitrix-test"

var testRepoDir = path.Join(testDockerPath, idutil.GetUuid(""))

func TestDevkit(t *testing.T) {
	t.Logf("start create repo at [%s]", testRepoDir)

	d := NewDocker(t, "test-op", "openpitrix")
	d.Port = testExportPort
	d.WorkDir = testDockerPath
	d.Volume[testRepoDir] = testDockerPath

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

	t.Run("create repo", func(t *testing.T) {
		time.Sleep(5 * time.Second)
		testCreateRepo(t, "test-devkit-repo-name", constants.ProviderQingCloud, fmt.Sprintf("http://%s:8879/", iputil.GetLocalIP()))
	})

	t.Run("create cluster", func(t *testing.T) {
		t.Log("TODO")
	})

	// cleanup
	t.Log(d.Exec("find . -mindepth 1 -delete"))
	t.Log(d.Teardown())
}

func waitRepoEventSuccess(t *testing.T, repoEventId string) {
	client := GetClient(clientConfig)

	for {
		describeEventParams := repo_indexer.NewDescribeRepoEventsParams()
		describeEventParams.RepoEventID = []string{repoEventId}
		describeEventResp, err := client.RepoIndexer.DescribeRepoEvents(describeEventParams)
		require.NoError(t, err)
		require.Equal(t, int64(1), describeEventResp.Payload.TotalCount, "count should be 1")
		require.Equal(t, repoEventId, describeEventResp.Payload.RepoEventSet[0].RepoEventID, "error repo event id")

		status := describeEventResp.Payload.RepoEventSet[0].Status
		require.NotEqual(t, constants.StatusFailed, status, "status should not be failed")

		switch status {
		case constants.StatusSuccessful:
			t.Logf("event [%s] finish with successful", repoEventId)
			return
		case constants.StatusPending, constants.StatusWorking:
			time.Sleep(5 * time.Second)
		}
		continue
	}
}

func testCreateRepo(t *testing.T, name, provider, url string) {
	client := GetClient(clientConfig)

	deleteRepo(t, client, name)

	createParams := repo_manager.NewCreateRepoParams()
	createParams.SetBody(
		&models.OpenpitrixCreateRepoRequest{
			Name:        name,
			Description: "description",
			Type:        "http",
			URL:         url,
			Providers:   []string{provider},
			Credential:  `{}`,
			Visibility:  "public",
		})
	createResp, err := client.RepoManager.CreateRepo(createParams)
	require.NoError(t, err)

	repoId := createResp.Payload.RepoID

	describeParams := repo_indexer.NewDescribeRepoEventsParams()
	describeParams.SetRepoID([]string{repoId})
	describeParams.SetStatus([]string{constants.StatusPending, constants.StatusWorking})
	describeResp, err := client.RepoIndexer.DescribeRepoEvents(describeParams)
	require.NoError(t, err)

	if len(describeResp.Payload.RepoEventSet) < 1 {
		t.Fatal("repo indexer need to working when repo created")
	}

	repoEventId := describeResp.Payload.RepoEventSet[0].RepoEventID
	waitRepoEventSuccess(t, repoEventId)

	describeAppParams := app_manager.NewDescribeAppsParams()
	describeAppParams.WithRepoID([]string{repoId})

	describeAppResp, err := client.AppManager.DescribeApps(describeAppParams)
	require.NoError(t, err)

	t.Logf("success got [%d] apps", describeAppResp.Payload.TotalCount)
	require.Equal(t, int64(1), describeAppResp.Payload.TotalCount, "auto create app more than 1")
	app := describeAppResp.Payload.AppSet[0]

	require.Equal(t, "nginx", app.Name, "app_name not equal nginx")

	t.Logf("got app [%+v]", app)

	deleteRepoParams := repo_manager.NewDeleteReposParams()
	deleteRepoParams.WithBody(&models.OpenpitrixDeleteReposRequest{
		RepoID: []string{repoId},
	})

	_, err = client.RepoManager.DeleteRepos(deleteRepoParams)
	require.NoError(t, err)
}
