// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

// +build k8s

package test

import (
	"fmt"
	"io/ioutil"
	"log"
	"testing"
	"time"

	"k8s.io/client-go/util/homedir"

	"openpitrix.io/openpitrix/test/client/app_manager"
	"openpitrix.io/openpitrix/test/client/repo_indexer"
	"openpitrix.io/openpitrix/test/client/repo_manager"
	//"openpitrix.io/openpitrix/test/client/job_manager"
	"openpitrix.io/openpitrix/test/client/cluster_manager"
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

	KubeConfig string
)

func TestK8s(t *testing.T) {
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
					Type:        "https",
					URL:         "https://kubernetes-charts-incubator.storage.googleapis.com/",
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
			time.Sleep(10 * time.Second)
		}
	}
	log.Printf("Got app [%s]\n", app.Name)

	var appVersion *models.OpenpitrixAppVersion
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

		appVersion = appVersions[0]
	}
	log.Printf("Got app version [%s]\n", appVersion.Name)

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

	// create cluster
	var jobId string
	var clusterId string
	{
		conf := fmt.Sprintf(`ClusterId: "cl-duak9an3f7cjfru"
Description: "test cluster"
Name: "cluster test"
RuntimeId: "%s"
`, runtime.RuntimeID)

		createParams := cluster_manager.NewCreateClusterParams()
		createParams.SetBody(&models.OpenpitrixCreateClusterRequest{
			AdvancedParam: []string{},
			AppID:         app.AppID,
			Conf:          conf,
			RuntimeID:     runtime.RuntimeID,
			VersionID:     appVersion.VersionID,
		})

		createResp, err := client.ClusterManager.CreateCluster(createParams)
		if err != nil {
			t.Fatal(err)
		}

		jobId = createResp.Payload.JobID
		clusterId = createResp.Payload.ClusterID
	}
	log.Printf("Got job id [%s]\n", jobId)
	log.Printf("Got cluster id [%s]\n", clusterId)

	//// check job finish
	//{
	//	describeParams := job_manager.NewDescribeJobsParams()
	//	describeParams.SetJobID([]string{jobId})
	//	job_manager.
	//}

	// check cluster status
	{
		describeParams := cluster_manager.NewDescribeClustersParams()
		describeParams.SetClusterID([]string{clusterId})
		describeResp, err := client.ClusterManager.DescribeClusters(describeParams)
		if err != nil {
			t.Fatal(err)
		}
		clusters := describeResp.Payload.ClusterSet
		if len(clusters) == 0 {
			t.Fatalf("Cluster [%s] not here", clusterId)
		}

		cluster := clusters[0]
		if cluster.Status != "active" {
			t.Fatal("Cluster status must be active when job finished")
		}

	}
}
