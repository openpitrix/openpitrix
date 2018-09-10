// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

// +build integration

package test

import (
	"io/ioutil"
	"log"
	"os"
	"path"
	"sort"
	"strings"
	"testing"

	"github.com/go-openapi/strfmt"
	"github.com/stretchr/testify/require"

	"openpitrix.io/openpitrix/pkg/constants"
	"openpitrix.io/openpitrix/pkg/devkit"
	"openpitrix.io/openpitrix/pkg/devkit/opapp"
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

func preparePackage(t *testing.T, v string) strfmt.Base64 {
	var testAppName = "test-app"

	cfile := &opapp.Metadata{
		Name:        testAppName,
		Description: "An OpenPitrix app",
		Version:     v,
		AppVersion:  "1.0",
		ApiVersion:  devkit.ApiVersionV1,
	}

	os.MkdirAll(testTmpDir, 0755)
	_, err := devkit.Create(cfile, testTmpDir)

	require.NoError(t, err)

	ch, err := devkit.LoadDir(path.Join(testTmpDir, testAppName))

	require.NoError(t, err)

	name, err := devkit.Save(ch, testTmpDir)

	require.NoError(t, err)

	t.Logf("save [%s] success", name)

	content, err := ioutil.ReadFile(name)

	require.NoError(t, err)

	require.NoError(t, os.RemoveAll(testTmpDir))

	return strfmt.Base64(content)
}

func testVersionPackage(t *testing.T, appId string) {
	client := GetClient(clientConfig)

	modifyAppParams := app_manager.NewModifyAppParams()
	modifyAppParams.SetBody(
		&models.OpenpitrixModifyAppRequest{
			AppID:  appId,
			RepoID: "repo-vmbased",
			Status: constants.StatusDraft,
		})
	_, err := client.AppManager.ModifyApp(modifyAppParams)

	require.NoError(t, err)

	createAppVersionParams := app_manager.NewCreateAppVersionParams()
	createAppVersionParams.SetBody(
		&models.OpenpitrixCreateAppVersionRequest{
			AppID:   appId,
			Status:  constants.StatusDraft,
			Package: preparePackage(t, "0.0.1"),
		})
	createAppVersionResp, err := client.AppManager.CreateAppVersion(createAppVersionParams)

	require.NoError(t, err)

	versionId1 := createAppVersionResp.Payload.VersionID

	modifyAppVersionParams := app_manager.NewModifyAppVersionParams()
	modifyAppVersionParams.SetBody(
		&models.OpenpitrixModifyAppVersionRequest{
			VersionID: versionId1,
			Package:   preparePackage(t, "0.0.2"),
		})
	_, err = client.AppManager.ModifyAppVersion(modifyAppVersionParams)

	require.NoError(t, err)

	modifyAppVersionParams = app_manager.NewModifyAppVersionParams()
	modifyAppVersionParams.SetBody(
		&models.OpenpitrixModifyAppVersionRequest{
			VersionID: versionId1,
			Package:   preparePackage(t, "0.0.3"),
		})
	_, err = client.AppManager.ModifyAppVersion(modifyAppVersionParams)

	require.NoError(t, err)

	createAppVersionParams = app_manager.NewCreateAppVersionParams()
	createAppVersionParams.SetBody(
		&models.OpenpitrixCreateAppVersionRequest{
			AppID:   appId,
			Status:  constants.StatusDraft,
			Package: preparePackage(t, "0.1.0"),
		})
	createAppVersionResp, err = client.AppManager.CreateAppVersion(createAppVersionParams)

	require.NoError(t, err)

	versionId2 := createAppVersionResp.Payload.VersionID

	modifyAppVersionParams = app_manager.NewModifyAppVersionParams()
	modifyAppVersionParams.SetBody(
		&models.OpenpitrixModifyAppVersionRequest{
			VersionID: versionId2,
			Package:   preparePackage(t, "0.0.3"),
		})
	_, err = client.AppManager.ModifyAppVersion(modifyAppVersionParams)

	require.Error(t, err)

	deleteAppVersionParams := app_manager.NewDeleteAppVersionParams()
	deleteAppVersionParams.SetBody(
		&models.OpenpitrixDeleteAppVersionRequest{
			VersionID: versionId2,
		})
	_, err = client.AppManager.DeleteAppVersion(deleteAppVersionParams)

	require.NoError(t, err)

	deleteAppVersionParams = app_manager.NewDeleteAppVersionParams()
	deleteAppVersionParams.SetBody(
		&models.OpenpitrixDeleteAppVersionRequest{
			VersionID: versionId1,
		})
	_, err = client.AppManager.DeleteAppVersion(deleteAppVersionParams)

	require.NoError(t, err)
}

func testVersionLifeCycle(t *testing.T, appId string) {
	client := GetClient(clientConfig)

	modifyAppParams := app_manager.NewModifyAppParams()
	modifyAppParams.SetBody(
		&models.OpenpitrixModifyAppRequest{
			AppID:  appId,
			Status: constants.StatusDraft,
		})
	_, err := client.AppManager.ModifyApp(modifyAppParams)

	require.NoError(t, err)

	createAppVersionParams := app_manager.NewCreateAppVersionParams()
	createAppVersionParams.SetBody(
		&models.OpenpitrixCreateAppVersionRequest{
			Name:   "test_version",
			AppID:  appId,
			Status: constants.StatusDraft,
		})
	createAppVersionResp, err := client.AppManager.CreateAppVersion(createAppVersionParams)

	require.NoError(t, err)

	versionId := createAppVersionResp.Payload.VersionID

	modifyAppVersionParams := app_manager.NewModifyAppVersionParams()
	modifyAppVersionParams.SetBody(
		&models.OpenpitrixModifyAppVersionRequest{
			VersionID: versionId,
			Name:      "test_version2",
		})
	_, err = client.AppManager.ModifyAppVersion(modifyAppVersionParams)

	require.NoError(t, err)

	submitAppVersionParams := app_manager.NewSubmitAppVersionParams()
	submitAppVersionParams.SetBody(
		&models.OpenpitrixSubmitAppVersionRequest{
			VersionID: versionId,
		})
	_, err = client.AppManager.SubmitAppVersion(submitAppVersionParams)

	require.NoError(t, err)

	rejectAppVersionParams := app_manager.NewRejectAppVersionParams()
	rejectAppVersionParams.SetBody(
		&models.OpenpitrixRejectAppVersionRequest{
			VersionID: versionId,
		})
	_, err = client.AppManager.RejectAppVersion(rejectAppVersionParams)

	require.NoError(t, err)

	_, err = client.AppManager.SubmitAppVersion(submitAppVersionParams)

	require.NoError(t, err)

	passAppVersionParams := app_manager.NewPassAppVersionParams()
	passAppVersionParams.SetBody(
		&models.OpenpitrixPassAppVersionRequest{
			VersionID: versionId,
		})
	_, err = client.AppManager.PassAppVersion(passAppVersionParams)

	require.NoError(t, err)

	releaseAppVersionParams := app_manager.NewReleaseAppVersionParams()
	releaseAppVersionParams.SetBody(
		&models.OpenpitrixReleaseAppVersionRequest{
			VersionID: versionId,
		})
	_, err = client.AppManager.ReleaseAppVersion(releaseAppVersionParams)

	require.NoError(t, err)

	suspendAppVersionParams := app_manager.NewSuspendAppVersionParams()
	suspendAppVersionParams.SetBody(
		&models.OpenpitrixSuspendAppVersionRequest{
			VersionID: versionId,
		})
	_, err = client.AppManager.SuspendAppVersion(suspendAppVersionParams)

	require.NoError(t, err)

	deleteAppVersionParams := app_manager.NewDeleteAppVersionsParams()
	deleteAppVersionParams.SetBody(
		&models.OpenpitrixDeleteAppVersionsRequest{
			VersionID: []string{versionId},
		})
	_, err = client.AppManager.DeleteAppVersions(deleteAppVersionParams)

	require.NoError(t, err)
}

func TestApp(t *testing.T) {
	client := GetClient(clientConfig)

	// delete old app
	testAppName := "e2e_test_app"
	testRepoId := "e2e_test_repo"
	testRepoId2 := "e2e_test_repo2"
	describeParams := app_manager.NewDescribeAppsParams()
	describeParams.SetName([]string{testAppName})
	describeParams.SetStatus([]string{constants.StatusDraft, constants.StatusActive})
	describeResp, err := client.AppManager.DescribeApps(describeParams)

	require.NoError(t, err)

	apps := describeResp.Payload.AppSet
	for _, app := range apps {
		deleteParams := app_manager.NewDeleteAppsParams()
		deleteParams.SetBody(
			&models.OpenpitrixDeleteAppsRequest{
				AppID: []string{app.AppID},
			})
		_, err := client.AppManager.DeleteApps(deleteParams)

		require.NoError(t, err)
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

	require.NoError(t, err)

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

	require.NoError(t, err)

	t.Log(modifyResp)
	// describe app
	describeParams.WithAppID([]string{appId})
	describeResp, err = client.AppManager.DescribeApps(describeParams)

	require.NoError(t, err)

	apps = describeResp.Payload.AppSet

	require.Equal(t, 1, len(apps))

	app := apps[0]

	require.Equal(t, testRepoId2, app.RepoID)

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

	getStatisticsResp, err := client.AppManager.GetAppStatistics(nil)
	require.NoError(t, err)
	require.NotEmpty(t, getStatisticsResp.Payload.LastTwoWeekCreated)
	require.NotEmpty(t, getStatisticsResp.Payload.TopTenRepos)
	require.NotEmpty(t, getStatisticsResp.Payload.AppCount)
	require.NotEmpty(t, getStatisticsResp.Payload.RepoCount)

	testVersionLifeCycle(t, appId)

	testVersionPackage(t, appId)

	// delete app
	deleteParams := app_manager.NewDeleteAppsParams()
	deleteParams.WithBody(&models.OpenpitrixDeleteAppsRequest{
		AppID: []string{appId},
	})
	deleteResp, err := client.AppManager.DeleteApps(deleteParams)

	require.NoError(t, err)

	t.Log(deleteResp)
	// describe deleted app
	describeParams.WithAppID([]string{appId})
	describeParams.WithStatus([]string{constants.StatusDeleted})
	describeParams.WithName(nil)
	describeResp, err = client.AppManager.DescribeApps(describeParams)

	require.NoError(t, err)

	apps = describeResp.Payload.AppSet

	require.Equal(t, 1, len(apps))

	app = apps[0]

	require.Equal(t, appId, app.AppID)

	require.Equal(t, constants.StatusDeleted, app.Status)

	t.Log("test app finish, all test is ok")
}
