// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package repocommon

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"openpitrix.io/openpitrix/pkg/constants"
	"openpitrix.io/openpitrix/test/client"
	"openpitrix.io/openpitrix/test/client/app_manager"
	"openpitrix.io/openpitrix/test/client/repo_indexer"
	"openpitrix.io/openpitrix/test/client/repo_manager"
	"openpitrix.io/openpitrix/test/models"
	"openpitrix.io/openpitrix/test/testutil"
)

func DeleteRepo(t *testing.T, c *client.Openpitrix, testRepoName string) {
	describeParams := repo_manager.NewDescribeReposParams()
	describeParams.SetName([]string{testRepoName})
	describeParams.SetStatus([]string{constants.StatusActive})
	describeResp, err := c.RepoManager.DescribeRepos(describeParams, nil)
	if err != nil {
		t.Fatal(err)
	}
	repos := describeResp.Payload.RepoSet
	for _, repo := range repos {
		deleteParams := repo_manager.NewDeleteReposParams()
		deleteParams.SetBody(
			&models.OpenpitrixDeleteReposRequest{
				RepoID: []string{repo.RepoID},
			})
		_, err := c.RepoManager.DeleteRepos(deleteParams, nil)
		if err != nil {
			t.Fatal(err)
		}
	}
}

func waitRepoEventSuccess(t *testing.T, c *client.Openpitrix, repoEventId string) {
	for {
		describeEventParams := repo_indexer.NewDescribeRepoEventsParams()
		describeEventParams.RepoEventID = []string{repoEventId}
		describeEventResp, err := c.RepoIndexer.DescribeRepoEvents(describeEventParams, nil)
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

func CreateRepo(t *testing.T, c *client.Openpitrix, name, provider, url string) {
	DeleteRepo(t, c, name)

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
	createResp, err := c.RepoManager.CreateRepo(createParams, nil)
	if err != nil {
		fmt.Print(testutil.ExecCmd(t, "docker-compose logs openpitrix-repo-manager"))
		fmt.Print(testutil.ExecCmd(t, "docker-compose logs openpitrix-rp-kubernetes"))
		fmt.Print(testutil.ExecCmd(t, "docker-compose logs openpitrix-rp-qingcloud"))
	}
	require.NoError(t, err)

	repoId := createResp.Payload.RepoID

	describeParams := repo_indexer.NewDescribeRepoEventsParams()
	describeParams.SetRepoID([]string{repoId})
	describeParams.SetStatus([]string{constants.StatusPending, constants.StatusWorking})
	describeResp, err := c.RepoIndexer.DescribeRepoEvents(describeParams, nil)
	require.NoError(t, err)

	if len(describeResp.Payload.RepoEventSet) < 1 {
		t.Fatal("repo indexer need to working when repo created")
	}

	repoEventId := describeResp.Payload.RepoEventSet[0].RepoEventID
	waitRepoEventSuccess(t, c, repoEventId)

	describeAppParams := app_manager.NewDescribeAppsParams()
	describeAppParams.WithRepoID([]string{repoId})

	describeAppResp, err := c.AppManager.DescribeApps(describeAppParams, nil)
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

	_, err = c.RepoManager.DeleteRepos(deleteRepoParams, nil)
	require.NoError(t, err)
}
