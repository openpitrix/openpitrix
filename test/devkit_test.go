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

	"github.com/stretchr/testify/assert"

	"openpitrix.io/openpitrix/pkg/constants"
	"openpitrix.io/openpitrix/pkg/utils/idtool"
	"openpitrix.io/openpitrix/pkg/utils/iptool"
	"openpitrix.io/openpitrix/test/client/app_manager"
	"openpitrix.io/openpitrix/test/client/repo_indexer"
	"openpitrix.io/openpitrix/test/client/repo_manager"
	"openpitrix.io/openpitrix/test/models"
)

var testRepoDir = path.Join("/tmp/openpitrix-test", idtool.GetUuid(""))

func TestDevkit(t *testing.T) {
	t.Logf("start create repo at [%s]", testRepoDir)

	d := NewDocker(t, "test-op", "openpitrix")
	d.Port = 8879
	d.WorkDir = "/tmp/openpitrix-test"
	d.Volume[testRepoDir] = "/tmp/openpitrix-test"

	t.Log(d.Setup())

	t.Log(d.Exec("op create nginx"))

	t.Log(d.Exec("ls nginx"))

	// TODO: write file content to testRepoDir, so that we can test create cluster

	t.Log(d.Exec("op package nginx"))

	t.Log(d.Exec("op index ./"))

	t.Log(d.Exec("cat index.yaml"))

	ip := strings.TrimSpace(d.Exec("hostname -i"))
	localIp := iptool.GetLocalIP()
	t.Log(d.ExecD(fmt.Sprintf("op serve --address %s:8879 --url http://%s:8879/", ip, localIp)))

	t.Run("create repo", func(t *testing.T) {
		time.Sleep(5 * time.Second)
		testCreateRepo(t, "test-devkit-repo-name", constants.ProviderQingCloud, fmt.Sprintf("http://%s:8879/", iptool.GetLocalIP()))
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
		assert.NoError(t, err)
		assert.Equal(t, int64(1), describeEventResp.Payload.TotalCount, "count should be 1")
		assert.Equal(t, repoEventId, describeEventResp.Payload.RepoEventSet[0].RepoEventID, "error repo event id")

		status := describeEventResp.Payload.RepoEventSet[0].Status
		assert.NotEqual(t, constants.StatusFailed, status, "status should not be failed")

		switch status {
		case constants.StatusSuccessful:
			return
		case constants.StatusPending, constants.StatusWorking:
			time.Sleep(5 * time.Second)
		}
		continue
	}
}

func testCreateRepo(t *testing.T, name, provider, url string) {
	client := GetClient(clientConfig)
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
	assert.NoError(t, err)

	repoId := createResp.Payload.Repo.RepoID

	indexParams := repo_indexer.NewIndexRepoParams()
	indexParams.SetBody(
		&models.OpenpitrixIndexRepoRequest{
			RepoID: repoId,
		})
	indexResp, err := client.RepoIndexer.IndexRepo(indexParams)
	assert.NoError(t, err)
	repoEventId := indexResp.Payload.RepoEvent.RepoEventID
	waitRepoEventSuccess(t, repoEventId)

	describeAppParams := app_manager.NewDescribeAppsParams()
	describeAppParams.WithRepoID([]string{repoId})

	describeAppResp, err := client.AppManager.DescribeApps(describeAppParams)
	assert.NoError(t, err)

	t.Logf("success got [%d] apps", describeAppResp.Payload.TotalCount)
	assert.Equal(t, int64(1), describeAppResp.Payload.TotalCount, "auto create app more than 1")
	app := describeAppResp.Payload.AppSet[0]

	assert.Equal(t, "nginx", app.Name, "app_name not equal nginx")

	t.Logf("got app [%+v]", app)
}
