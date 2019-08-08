// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

// +build k8s

package k8s

import (
	"fmt"
	"testing"
	"time"

	"openpitrix.io/openpitrix/pkg/constants"
	log "openpitrix.io/openpitrix/pkg/logger"
	"openpitrix.io/openpitrix/pkg/util/funcutil"
	apiclient "openpitrix.io/openpitrix/test/client"
	"openpitrix.io/openpitrix/test/client/app_manager"
	"openpitrix.io/openpitrix/test/client/cluster_manager"
	"openpitrix.io/openpitrix/test/client/job_manager"
	"openpitrix.io/openpitrix/test/client/repo_manager"
	"openpitrix.io/openpitrix/test/client/runtime_manager"
	"openpitrix.io/openpitrix/test/models"
	"openpitrix.io/openpitrix/test/testutil"
)

var (
	RepoNameForTest    = "google"
	AppNameForTest     = "cerebro"
	RuntimeNameForTest = "minikube"

	clientConfig = testutil.GetClientConfig()

	Service = []string{"hyperpitrix", "openpitrix-rp-kubernetes", "openpitrix-rp-qingcloud"}
)

func TestK8S(t *testing.T) {
	client := testutil.GetClient(clientConfig)

	// create repo
	var repoId string
	{
		describeParams := repo_manager.NewDescribeReposParams()
		describeParams.SetName([]string{RepoNameForTest})
		describeResp, err := client.RepoManager.DescribeRepos(describeParams, nil)
		if err != nil {
			t.Fatal(err)
		}
		repos := describeResp.Payload.RepoSet

		if len(repos) != 0 {
			repoId = repos[0].RepoID
		} else {
			createParams := repo_manager.NewCreateRepoParams()
			createParams.SetBody(
				&models.OpenpitrixCreateRepoRequest{
					Name:        RepoNameForTest,
					Description: "test repo",
					Type:        "https",
					URL:         "https://mirror.azure.cn/kubernetes/charts/",
					Credential:  `{}`,
					Visibility:  "public",
					Providers:   []string{constants.ProviderKubernetes},
				})
			createResp, err := client.RepoManager.CreateRepo(createParams, nil)
			if err != nil {
				t.Fatal(err)
			}
			repoId = createResp.Payload.RepoID
		}
	}
	log.Info(nil, "Got repo id [%s]", repoId)

	// waiting for apps indexed by repo indexer
	var app *models.OpenpitrixApp
	{
		for {
			describeParams := app_manager.NewDescribeAppsParams()
			describeParams.WithRepoID([]string{repoId})
			describeParams.SetSearchWord(&AppNameForTest)
			describeResp, err := client.AppManager.DescribeApps(describeParams, nil)
			if err != nil {
				t.Fatal(err)
			}
			apps := describeResp.Payload.AppSet
			if len(apps) != 0 {
				app = apps[0]
				break
			}
			log.Info(nil, "Waiting for app [%s]...", AppNameForTest)
			time.Sleep(5 * time.Second)
		}
	}
	log.Info(nil, "Got app name [%s] latest version [%s]", app.Name, app.LatestAppVersion.Name)

	var appVersion1 *models.OpenpitrixAppVersion
	var appVersion2 *models.OpenpitrixAppVersion
	{
		for {
			describeParams := app_manager.NewDescribeAppVersionsParams()
			describeParams.SetAppID([]string{app.AppID})
			describeParams.WithPackageName([]string{"cerebro-0.3.0.tgz", "cerebro-0.3.1.tgz"})
			describeResp, err := client.AppManager.DescribeAppVersions(describeParams, nil)
			if err != nil {
				t.Fatal(err)
			}
			appVersions := describeResp.Payload.AppVersionSet

			if len(appVersions) == 2 {
				appVersion1 = appVersions[0]
				appVersion2 = appVersions[1]
				break
			}

			log.Info(nil, "Waiting for app version ...")
			time.Sleep(5 * time.Second)
		}
	}
	log.Info(nil, "Got app version 1, %s", appVersion1.Name)
	log.Info(nil, "Got app version 2, %s", appVersion2.Name)

	// create runtime
	var runtimeId string
	{
		describeParams := runtime_manager.NewDescribeRuntimesParams()
		describeParams.SetSearchWord(&RuntimeNameForTest)
		describeResp, err := client.RuntimeManager.DescribeRuntimes(describeParams, nil)
		if err != nil {
			t.Fatal(err)
		}
		runtimes := describeResp.Payload.RuntimeSet
		if len(runtimes) != 0 {
			runtimeId = runtimes[0].RuntimeID
		} else {
			KubeConfig := getRuntimeCredential(t)

			createCredentialParams := runtime_manager.NewCreateRuntimeCredentialParams()
			createCredentialParams.SetBody(
				&models.OpenpitrixCreateRuntimeCredentialRequest{
					Name:                     RuntimeNameForTest,
					Description:              "minikube",
					Provider:                 constants.ProviderKubernetes,
					RuntimeCredentialContent: KubeConfig,
				})
			createCredentialResp, err := client.RuntimeManager.CreateRuntimeCredential(createCredentialParams, nil)
			testutil.NoError(t, err, Service)
			runtimeCredentialId := createCredentialResp.Payload.RuntimeCredentialID

			createParams := runtime_manager.NewCreateRuntimeParams()
			createParams.SetBody(
				&models.OpenpitrixCreateRuntimeRequest{
					Name:                RuntimeNameForTest,
					Description:         "minikube",
					Provider:            constants.ProviderKubernetes,
					RuntimeCredentialID: runtimeCredentialId,
					Zone:                "test",
				})
			createResp, err := client.RuntimeManager.CreateRuntime(createParams, nil)
			testutil.NoError(t, err, Service)
			runtimeId = createResp.Payload.RuntimeID
		}
	}
	log.Info(nil, "Got runtime id [%s]", runtimeId)

	var clusterId string
	{
		log.Info(nil, "Creating cluster...")

		conf := `Name: test`

		createParams := cluster_manager.NewCreateClusterParams()
		createParams.SetBody(&models.OpenpitrixCreateClusterRequest{
			AdvancedParam: []string{},
			AppID:         app.AppID,
			Conf:          conf,
			RuntimeID:     runtimeId,
			VersionID:     appVersion1.VersionID,
		})

		createResp, err := client.ClusterManager.CreateCluster(createParams, nil)
		testutil.NoError(t, err, Service)

		clusterId = createResp.Payload.ClusterID
		jobId := createResp.Payload.JobID

		err = waitJobFinish(t, client, jobId)
		testutil.NoError(t, err, Service)

		log.Info(nil, "Cluster [%s] created", clusterId)
	}

	{
		log.Info(nil, "Upgrading cluster [%s] ...", clusterId)

		upgradeParams := cluster_manager.NewUpgradeClusterParams()
		upgradeParams.SetBody(&models.OpenpitrixUpgradeClusterRequest{
			AdvancedParam: []string{},
			ClusterID:     clusterId,
			VersionID:     appVersion2.VersionID,
		})

		upgradeResp, err := client.ClusterManager.UpgradeCluster(upgradeParams, nil)
		testutil.NoError(t, err, Service)

		jobId := upgradeResp.Payload.JobID

		err = waitJobFinish(t, client, jobId)
		testutil.NoError(t, err, Service)

		log.Info(nil, "Cluster [%s] upgraded", clusterId)
	}

	{
		log.Info(nil, "Rolling back cluster [%s] ...", clusterId)

		rollbackParams := cluster_manager.NewRollbackClusterParams()
		rollbackParams.SetBody(&models.OpenpitrixRollbackClusterRequest{
			AdvancedParam: []string{},
			ClusterID:     clusterId,
		})

		rollbackResp, err := client.ClusterManager.RollbackCluster(rollbackParams, nil)
		testutil.NoError(t, err, Service)

		jobId := rollbackResp.Payload.JobID

		err = waitJobFinish(t, client, jobId)
		testutil.NoError(t, err, Service)

		log.Info(nil, "Cluster [%s] roll back", clusterId)
	}

	{
		log.Info(nil, "Updating cluster [%s] env ...", clusterId)

		env := `Name: test
Description: test`

		updateEnvParams := cluster_manager.NewUpdateClusterEnvParams()
		updateEnvParams.SetBody(&models.OpenpitrixUpdateClusterEnvRequest{
			ClusterID: clusterId,
			Env:       env,
		})

		updateEnvResp, err := client.ClusterManager.UpdateClusterEnv(updateEnvParams, nil)
		testutil.NoError(t, err, Service)

		jobId := updateEnvResp.Payload.JobID

		err = waitJobFinish(t, client, jobId)
		testutil.NoError(t, err, Service)

		log.Info(nil, "Cluster [%s] env updated", clusterId)
	}

	{
		log.Info(nil, "Deleting cluster [%s]...", clusterId)

		deleteParams := cluster_manager.NewDeleteClustersParams()
		deleteParams.SetBody(&models.OpenpitrixDeleteClustersRequest{
			AdvancedParam: []string{},
			ClusterID:     []string{clusterId},
		})

		deleteResp, err := client.ClusterManager.DeleteClusters(deleteParams, nil)
		testutil.NoError(t, err, Service)

		jobId := deleteResp.Payload.JobID[0]

		err = waitJobFinish(t, client, jobId)
		testutil.NoError(t, err, Service)

		log.Info(nil, "Cluster [%s] deleted", clusterId)
	}

	{
		log.Info(nil, "Purging cluster [%s]...", clusterId)

		ceaseParams := cluster_manager.NewCeaseClustersParams()
		ceaseParams.SetBody(&models.OpenpitrixCeaseClustersRequest{
			AdvancedParam: []string{},
			ClusterID:     []string{clusterId},
		})

		ceaseResp, err := client.ClusterManager.CeaseClusters(ceaseParams, nil)
		testutil.NoError(t, err, Service)

		jobId := ceaseResp.Payload.JobID[0]

		err = waitJobFinish(t, client, jobId)
		testutil.NoError(t, err, Service)

		testutil.NoError(t, err, Service)

		log.Info(nil, "Cluster [%s] purged", clusterId)
	}
}

func waitJobFinish(_ *testing.T, client *apiclient.Openpitrix, jobId string) error {
	log.Info(nil, "Waiting job [%s]", jobId)

	describeParams := job_manager.NewDescribeJobsParams()
	describeParams.WithJobID([]string{jobId})

	return funcutil.WaitForSpecificOrError(func() (bool, error) {
		describeResp, err := client.JobManager.DescribeJobs(describeParams, nil)
		if err != nil {
			//network or api error, not considered job fail.
			return false, nil
		}

		if len(describeResp.Payload.JobSet) == 0 {
			return false, fmt.Errorf("can not find job [%s]", jobId)
		}
		j := describeResp.Payload.JobSet[0]
		if j.Status == "working" || j.Status == "pending" {
			return false, nil
		}
		if j.Status == "successful" {
			return true, nil
		}
		if j.Status == "failed" {
			return false, fmt.Errorf("job [%s] failed", jobId)
		}
		log.Info(nil, "Unknown status [%s] for job [%s]", j.Status, jobId)
		return false, nil

	}, constants.WaitTaskTimeout, constants.WaitTaskInterval)
}

func getRuntimeCredential(t *testing.T) string {
	return testutil.ExecCmd(t, "kubectl config view --flatten")
}
