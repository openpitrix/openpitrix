// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

// +build k8s

package test

import (
	"io/ioutil"
	"log"
	"testing"
	"time"

	"k8s.io/client-go/util/homedir"

	"openpitrix.io/openpitrix/test/client/app_manager"
	"openpitrix.io/openpitrix/test/client/cluster_manager"
	"openpitrix.io/openpitrix/test/client/repo_indexer"
	"openpitrix.io/openpitrix/test/client/repo_manager"
	"openpitrix.io/openpitrix/test/client/runtime_manager"
	"openpitrix.io/openpitrix/test/models"
)

var (
	clientConfig = &ClientConfig{
		Host:  "192.168.0.3:9100",
		Debug: true,
	}

	RepoNameForTest    = "incubator"
	RuntimeNameForTest = "k8s runtime"
	RepoUrlForTest     = "http://192.168.0.11:8879/"

	KubeConfig string
)

func TestK8s(t *testing.T) {
	log.SetPrefix("[ K8S TEST ] ")

	b, err := ioutil.ReadFile(homedir.HomeDir() + "/.kube/config")
	if err != nil {
		t.Fatal(err)
	}

	KubeConfig = string(b)

	client := GetClient(clientConfig)

	// create repo
	var repo *models.OpenpitrixRepo
	{
		describeParams := repo_manager.NewDescribeReposParams()
		describeParams.SetName([]string{RepoNameForTest})
		describeResp, err := client.RepoManager.DescribeRepos(describeParams)
		if err != nil {
			t.Fatal(err)
		}
		repos := describeResp.Payload.RepoSet

		if len(repos) != 0 {
			repo = repos[0]
		} else {
			createParams := repo_manager.NewCreateRepoParams()
			createParams.SetBody(
				&models.OpenpitrixCreateRepoRequest{
					Name:        RepoNameForTest,
					Description: "incubator charts",
					Type:        "http",
					URL:         RepoUrlForTest,
					Credential:  `{}`,
					Visibility:  "public",
				})
			createResp, err := client.RepoManager.CreateRepo(createParams)
			if err != nil {
				t.Fatal(err)
			}
			repo = createResp.Payload.Repo
		}
	}
	log.Printf("Got repo [%s]\n", repo.Name)

	// index repo
	{
		indexRepoParams := repo_indexer.NewIndexRepoParams()
		indexRepoParams.SetBody(
			&models.OpenpitrixIndexRepoRequest{
				RepoID: repo.RepoID,
			})
		_, err := client.RepoIndexer.IndexRepo(indexRepoParams)
		if err != nil {
			t.Fatal(err)
		}
	}

	time.Sleep(2 * time.Second)

	// waiting for apps indexed by repo indexer
	var app *models.OpenpitrixApp
	{
		for {
			describeParams := app_manager.NewDescribeAppsParams()
			describeParams.WithRepoID([]string{repo.RepoID})
			describeResp, err := client.AppManager.DescribeApps(describeParams)
			if err != nil {
				t.Fatal(err)
			}
			apps := describeResp.Payload.AppSet
			if len(apps) != 0 {
				app = apps[0]
				break
			}
			log.Printf("Waiting for app ...")
			time.Sleep(5 * time.Second)
		}
	}
	log.Printf("Got app [%s]\n", app.Name)

	var appVersion1 *models.OpenpitrixAppVersion
	var appVersion2 *models.OpenpitrixAppVersion
	{
		describeParams := app_manager.NewDescribeAppVersionsParams()
		describeParams.SetAppID([]string{app.AppID})
		describeResp, err := client.AppManager.DescribeAppVersions(describeParams)
		if err != nil {
			t.Fatal(err)
		}
		appVersions := describeResp.Payload.AppVersionSet
		if len(appVersions) == 0 {
			t.Fatal("App has no version released")
		}

		if len(appVersions) != 2 {
			t.Fatal("We need two version to test upgrade")
		}

		appVersion1 = appVersions[0]
		appVersion2 = appVersions[1]
	}
	log.Printf("Got app version [%s] [%s]\n", appVersion1.Name, appVersion2.Name)

	// create runtime
	var runtime *models.OpenpitrixRuntime
	{
		describeParams := runtime_manager.NewDescribeRuntimesParams()
		describeParams.SetSearchWord(&RuntimeNameForTest)
		describeResp, err := client.RuntimeManager.DescribeRuntimes(describeParams)
		if err != nil {
			t.Fatal(err)
		}
		runtimes := describeResp.Payload.RuntimeSet
		if len(runtimes) != 0 {
			runtime = runtimes[0]
		} else {
			createParams := runtime_manager.NewCreateRuntimeParams()
			createParams.SetBody(
				&models.OpenpitrixCreateRuntimeRequest{
					Name:              RuntimeNameForTest,
					Description:       "k8s runtime",
					Provider:          "kubernetes",
					RuntimeURL:        "https://k8s.io",
					RuntimeCredential: KubeConfig,
					Zone:              "default",
				})
			createResp, err := client.RuntimeManager.CreateRuntime(createParams)
			if err != nil {
				t.Fatal(err)
			}
			runtime = createResp.Payload.Runtime
		}
	}
	log.Printf("Got runtime [%s]\n", runtime.Name)

	var clusterId string
	log.Printf("Creating cluster...\n")
	{
		conf := `
Description: "test cluster"
Name: "cluster test"
`

		createParams := cluster_manager.NewCreateClusterParams()
		createParams.SetBody(&models.OpenpitrixCreateClusterRequest{
			AdvancedParam: []string{},
			AppID:         app.AppID,
			Conf:          conf,
			RuntimeID:     runtime.RuntimeID,
			VersionID:     appVersion1.VersionID,
		})

		createResp, err := client.ClusterManager.CreateCluster(createParams)
		if err != nil {
			t.Fatal(err)
		}

		clusterId = createResp.Payload.ClusterID
	}
	log.Printf("Cluster [%s] created \n", clusterId)

	time.Sleep(4 * time.Minute)

	log.Printf("Upgrading cluster [%s]...\n", clusterId)
	{
		upgradeParams := cluster_manager.NewUpgradeClusterParams()
		upgradeParams.SetBody(&models.OpenpitrixUpgradeClusterRequest{
			AdvancedParam: []string{},
			ClusterID:     clusterId,
			VersionID:     appVersion2.VersionID,
		})

		_, err := client.ClusterManager.UpgradeCluster(upgradeParams)
		if err != nil {
			t.Fatal(err)
		}
	}
	log.Printf("Cluster [%s] upgraded \n", clusterId)

	time.Sleep(4 * time.Minute)

	log.Printf("Rolling back cluster [%s]...\n", clusterId)
	{
		rollbackParams := cluster_manager.NewRollbackClusterParams()
		rollbackParams.SetBody(&models.OpenpitrixRollbackClusterRequest{
			AdvancedParam: []string{},
			ClusterID:     clusterId,
		})

		_, err := client.ClusterManager.RollbackCluster(rollbackParams)
		if err != nil {
			t.Fatal(err)
		}
	}
	log.Printf("Cluster [%s] roll back \n", clusterId)

	time.Sleep(4 * time.Minute)

	log.Printf("Deleting cluster [%s]...\n", clusterId)
	{
		deleteParams := cluster_manager.NewDeleteClustersParams()
		deleteParams.SetBody(&models.OpenpitrixDeleteClustersRequest{
			AdvancedParam: []string{},
			ClusterID:     []string{clusterId},
		})

		_, err := client.ClusterManager.DeleteClusters(deleteParams)
		if err != nil {
			t.Fatal(err)
		}
	}
	log.Printf("Cluster [%s] deleted \n", clusterId)

	time.Sleep(4 * time.Minute)

	log.Printf("Purging cluster [%s]...\n", clusterId)
	{
		ceaseParams := cluster_manager.NewCeaseClustersParams()
		ceaseParams.SetBody(&models.OpenpitrixCeaseClustersRequest{
			AdvancedParam: []string{},
			ClusterID:     []string{clusterId},
		})

		_, err := client.ClusterManager.CeaseClusters(ceaseParams)
		if err != nil {
			t.Fatal(err)
		}
	}
	log.Printf("Cluster [%s] purged \n", clusterId)
}
