// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

// +build integration

package test

import (
	"fmt"
	"net"
	"os/exec"
	"path"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"openpitrix.io/openpitrix/pkg/constants"
	"openpitrix.io/openpitrix/pkg/utils/idtool"
	"openpitrix.io/openpitrix/test/client/app_manager"
	"openpitrix.io/openpitrix/test/client/repo_indexer"
	"openpitrix.io/openpitrix/test/client/repo_manager"
	"openpitrix.io/openpitrix/test/models"
)

var testRepoDir = path.Join("/tmp/openpitrix-test", idtool.GetUuid(""))

func execOnTestRepo(t *testing.T, cmd string) string {
	fullCmd := fmt.Sprintf("docker run --rm -i --name='test-op' -p 8879:8879 -v %s:/tmp/openpitrix-test -w /tmp/openpitrix-test openpitrix %s", testRepoDir, cmd)
	t.Logf("run command [%s]", fullCmd)
	c := exec.Command("/bin/sh", "-c", fullCmd)
	output, err := c.CombinedOutput()
	assert.NoError(t, err)
	return string(output)
}

func execOnTestRepoD(t *testing.T, cmd string) string {
	fullCmd := fmt.Sprintf("docker run --rm -i -d --name='test-op' -p 8879:8879 -v %s:/tmp/openpitrix-test -w /tmp/openpitrix-test openpitrix %s", testRepoDir, cmd)
	t.Logf("run command [%s]", fullCmd)
	c := exec.Command("/bin/sh", "-c", fullCmd)
	output, err := c.CombinedOutput()
	assert.NoError(t, err)
	return string(output)
}

func cleanupTestRepoDocker(t *testing.T) string {
	output, err := exec.Command("/bin/sh", "-c", "docker rm -f test-op").CombinedOutput()
	assert.NoError(t, err)
	return string(output)
}

func TestDevkit(t *testing.T) {
	t.Logf("start create repo at [%s]", testRepoDir)

	t.Log(execOnTestRepo(t, "op create nginx"))

	t.Log(execOnTestRepo(t, "ls nginx"))

	t.Log(execOnTestRepo(t, "op package nginx"))

	t.Log(execOnTestRepo(t, "op index ./"))

	t.Log(execOnTestRepo(t, "cat index.yaml"))

	ip := strings.TrimSpace(execOnTestRepo(t, "hostname -i"))
	localIp := GetLocalIP()
	t.Log(execOnTestRepoD(t, fmt.Sprintf("op serve --address %s:8879 --url http://%s:8879/", ip, localIp)))

	t.Run("create repo", func(t *testing.T) {
		time.Sleep(5 * time.Second)
		testCreateRepo(t)
	})

	// cleanup
	t.Log(cleanupTestRepoDocker(t))
	t.Log(execOnTestRepo(t, "find . -mindepth 1 -delete"))
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

// GetLocalIP returns the non loopback local IP of the host
func GetLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}
	for _, address := range addrs {
		// check the address type and if it is not a loopback the display it
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return ""
}

func testCreateRepo(t *testing.T) {
	client := GetClient(clientConfig)
	testRepoName := "test-devkit-repo-name"
	createParams := repo_manager.NewCreateRepoParams()
	createParams.SetBody(
		&models.OpenpitrixCreateRepoRequest{
			Name:        testRepoName,
			Description: "description",
			Type:        "http",
			URL:         fmt.Sprintf("http://%s:8879/", GetLocalIP()),
			Providers:   []string{"qingcloud"},
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
