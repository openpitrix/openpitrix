// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

// +build integration

package test

import (
	"log"
	"os"
	"sort"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"openpitrix.io/openpitrix/pkg/constants"
	"openpitrix.io/openpitrix/test/client/app_manager"
	"openpitrix.io/openpitrix/test/models"
)

var clientConfig = &ClientConfig{}

func init() {
	clientConfig = GetClientConfig()
	log.Printf("Got Client Config: %+v", clientConfig)
}

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}

func getSortedString(s []string) string {
	sortedCategoryIds := sort.StringSlice(s)
	sortedCategoryIds.Sort()
	return strings.Join(sortedCategoryIds, ",")
}

func TestApp(t *testing.T) {
	client := GetClient(clientConfig)

	// delete old app
	testAppName := "e2e_test_app"
	testRepoId := "e2e_test_repo"
	testRepoId2 := "e2e_test_repo2"
	describeParams := app_manager.NewDescribeAppsParams()
	describeParams.SetName([]string{testAppName})
	describeParams.SetStatus([]string{constants.StatusActive})
	describeResp, err := client.AppManager.DescribeApps(describeParams)
	if err != nil {
		t.Fatal(err)
	}
	apps := describeResp.Payload.AppSet
	for _, app := range apps {
		deleteParams := app_manager.NewDeleteAppsParams()
		deleteParams.SetBody(
			&models.OpenpitrixDeleteAppsRequest{
				AppID: []string{app.AppID},
			})
		_, err := client.AppManager.DeleteApps(deleteParams)
		if err != nil {
			t.Fatal(err)
		}
	}
	// create app
	createParams := app_manager.NewCreateAppParams()
	createParams.SetBody(
		&models.OpenpitrixCreateAppRequest{
			Name:       testAppName,
			RepoID:     testRepoId,
			CategoryID: "xx,yy,zz",
		})
	createResp, err := client.AppManager.CreateApp(createParams)
	if err != nil {
		t.Fatal(err)
	}
	appId := createResp.Payload.AppID
	// modify app
	modifyParams := app_manager.NewModifyAppParams()
	modifyParams.SetBody(
		&models.OpenpitrixModifyAppRequest{
			AppID:      appId,
			RepoID:     testRepoId2,
			CategoryID: "aa,bb,cc,xx",
		})
	modifyResp, err := client.AppManager.ModifyApp(modifyParams)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(modifyResp)
	// describe app
	describeParams.WithAppID([]string{appId})
	describeResp, err = client.AppManager.DescribeApps(describeParams)
	if err != nil {
		t.Fatal(err)
	}
	apps = describeResp.Payload.AppSet
	if len(apps) != 1 {
		t.Fatalf("failed to describe apps with params [%+v]", describeParams)
	}
	app := apps[0]
	if app.RepoID != testRepoId2 {
		t.Fatalf("failed to modify app, app [%+v] repo is not [%s]", apps[0], testRepoId2)
	}
	var enabledCategoryIds []string
	var disabledCategoryIds []string
	for _, a := range app.CategorySet {
		if a.Status == constants.StatusEnabled {
			enabledCategoryIds = append(enabledCategoryIds, a.CategoryID)
		}
		if a.Status == constants.StatusDisabled {
			disabledCategoryIds = append(disabledCategoryIds, a.CategoryID)
		}
	}

	require.Equal(t, getSortedString(enabledCategoryIds), "aa,bb,cc,xx")
	require.Equal(t, getSortedString(disabledCategoryIds), "yy,zz")
	// delete app
	deleteParams := app_manager.NewDeleteAppsParams()
	deleteParams.WithBody(&models.OpenpitrixDeleteAppsRequest{
		AppID: []string{appId},
	})
	deleteResp, err := client.AppManager.DeleteApps(deleteParams)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(deleteResp)
	// describe deleted app
	describeParams.WithAppID([]string{appId})
	describeParams.WithStatus([]string{constants.StatusDeleted})
	describeParams.WithName(nil)
	describeResp, err = client.AppManager.DescribeApps(describeParams)
	if err != nil {
		t.Fatal(err)
	}
	apps = describeResp.Payload.AppSet
	if len(apps) != 1 {
		t.Fatalf("failed to describe apps with params [%+v]", describeParams)
	}
	app = apps[0]
	if app.AppID != appId {
		t.Fatalf("failed to describe app")
	}
	if app.Status != constants.StatusDeleted {
		t.Fatalf("failed to delete app, got app status [%s]", app.Status)
	}

	getStatisticsResp, err := client.AppManager.GetAppStatistics(nil)
	require.NoError(t, err)
	require.NotEmpty(t, getStatisticsResp.Payload.LastTwoWeekCreated)
	require.NotEmpty(t, getStatisticsResp.Payload.TopTenRepos)
	require.NotEmpty(t, getStatisticsResp.Payload.AppCount)
	require.NotEmpty(t, getStatisticsResp.Payload.RepoCount)

	t.Log("test app finish, all test is ok")
}
