// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

// +build k8s

package k8s

import (
	"fmt"
	"io/ioutil"
	"testing"

	"openpitrix.io/openpitrix/pkg/constants"
	"openpitrix.io/openpitrix/pkg/util/funcutil"
	apiclient "openpitrix.io/openpitrix/test/client"
	"openpitrix.io/openpitrix/test/client/app_manager"
	"openpitrix.io/openpitrix/test/client/cluster_manager"
	"openpitrix.io/openpitrix/test/client/job_manager"
	"openpitrix.io/openpitrix/test/client/runtime_manager"
	"openpitrix.io/openpitrix/test/models"
	"openpitrix.io/openpitrix/test/testutil"
)

var (
	RuntimeNameForTest = "minikube"

	clientConfig = testutil.GetClientConfig()

	package1, _ = ioutil.ReadFile("./test_data/mysql-1.1.0.tgz")
	package2, _ = ioutil.ReadFile("./test_data/mysql-1.3.3.tgz")

	versionType = "helm"
)

func TestK8S(t *testing.T) {
	client := testutil.GetClient(clientConfig)

	var appId string
	var version1 string
	var version2 string
	// create app and version
	{
		createParams := app_manager.NewCreateAppParams()
		createParams.WithBody(&models.OpenpitrixCreateAppRequest{
			VersionPackage: package1,
			Name:           "mysql",
			VersionType:    versionType,
		})
		createResp, err := client.AppManager.CreateApp(createParams, nil)
		if err != nil {
			t.Fatal(err)
		}
		appId = createResp.Payload.AppID
		version1 = createResp.Payload.VersionID

		createVersionParams := app_manager.NewCreateAppVersionParams()
		createVersionParams.WithBody(&models.OpenpitrixCreateAppVersionRequest{
			AppID:   appId,
			Package: package2,
			Type:    versionType,
		})
		createVersionResp, err := client.AppManager.CreateAppVersion(createVersionParams, nil)
		if err != nil {
			t.Fatal(err)
		}
		version2 = createVersionResp.Payload.VersionID
	}

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
			if err != nil {
				fmt.Print(testutil.ExecCmd(t, "docker-compose logs hyperpitrix"))
				t.Fatal(err)
			}
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
			if err != nil {
				fmt.Print(testutil.ExecCmd(t, "docker-compose logs hyperpitrix"))
				t.Fatal(err)
			}
			runtimeId = createResp.Payload.RuntimeID
		}
	}
	fmt.Printf("Got runtime id [%s]", runtimeId)

	var clusterId string
	{
		fmt.Printf("Creating cluster...")

		conf := `Name: test`

		createParams := cluster_manager.NewCreateClusterParams()
		createParams.SetBody(&models.OpenpitrixCreateClusterRequest{
			AdvancedParam: []string{},
			AppID:         appId,
			Conf:          conf,
			RuntimeID:     runtimeId,
			VersionID:     version1,
		})

		createResp, err := client.ClusterManager.CreateCluster(createParams, nil)
		if err != nil {
			fmt.Print(testutil.ExecCmd(t, "docker-compose logs openpitrix-rp-kubernetes"))
			fmt.Print(testutil.ExecCmd(t, "docker-compose logs hyperpitrix"))
			t.Fatal(err)
		}

		clusterId = createResp.Payload.ClusterID
		jobId := createResp.Payload.JobID

		err = waitJobFinish(t, client, jobId)
		if err != nil {
			fmt.Print(testutil.ExecCmd(t, "docker-compose logs openpitrix-rp-kubernetes"))
			fmt.Print(testutil.ExecCmd(t, "docker-compose logs hyperpitrix"))
			t.Fatal(err)
		}

		fmt.Printf("Cluster [%s] created", clusterId)
	}

	{
		fmt.Printf("Upgrading cluster [%s] ...", clusterId)

		upgradeParams := cluster_manager.NewUpgradeClusterParams()
		upgradeParams.SetBody(&models.OpenpitrixUpgradeClusterRequest{
			AdvancedParam: []string{},
			ClusterID:     clusterId,
			VersionID:     version2,
		})

		upgradeResp, err := client.ClusterManager.UpgradeCluster(upgradeParams, nil)
		if err != nil {
			fmt.Print(testutil.ExecCmd(t, "docker-compose logs openpitrix-rp-kubernetes"))
			fmt.Print(testutil.ExecCmd(t, "docker-compose logs hyperpitrix"))
			t.Fatal(err)
		}

		jobId := upgradeResp.Payload.JobID

		err = waitJobFinish(t, client, jobId)
		if err != nil {
			fmt.Print(testutil.ExecCmd(t, "docker-compose logs openpitrix-rp-kubernetes"))
			fmt.Print(testutil.ExecCmd(t, "docker-compose logs hyperpitrix"))
			t.Fatal(err)
		}

		fmt.Printf("Cluster [%s] upgraded", clusterId)
	}

	{
		fmt.Printf("Rolling back cluster [%s] ...", clusterId)

		rollbackParams := cluster_manager.NewRollbackClusterParams()
		rollbackParams.SetBody(&models.OpenpitrixRollbackClusterRequest{
			AdvancedParam: []string{},
			ClusterID:     clusterId,
		})

		rollbackResp, err := client.ClusterManager.RollbackCluster(rollbackParams, nil)
		if err != nil {
			fmt.Print(testutil.ExecCmd(t, "docker-compose logs openpitrix-rp-kubernetes"))
			fmt.Print(testutil.ExecCmd(t, "docker-compose logs hyperpitrix"))
			t.Fatal(err)
		}

		jobId := rollbackResp.Payload.JobID

		err = waitJobFinish(t, client, jobId)
		if err != nil {
			fmt.Print(testutil.ExecCmd(t, "docker-compose logs openpitrix-rp-kubernetes"))
			fmt.Print(testutil.ExecCmd(t, "docker-compose logs hyperpitrix"))
			t.Fatal(err)
		}

		fmt.Printf("Cluster [%s] roll back", clusterId)
	}

	{
		fmt.Printf("Updating cluster [%s] env ...", clusterId)

		env := `Name: test
Description: test`

		updateEnvParams := cluster_manager.NewUpdateClusterEnvParams()
		updateEnvParams.SetBody(&models.OpenpitrixUpdateClusterEnvRequest{
			ClusterID: clusterId,
			Env:       env,
		})

		updateEnvResp, err := client.ClusterManager.UpdateClusterEnv(updateEnvParams, nil)
		if err != nil {
			fmt.Print(testutil.ExecCmd(t, "docker-compose logs openpitrix-rp-kubernetes"))
			fmt.Print(testutil.ExecCmd(t, "docker-compose logs hyperpitrix"))
			t.Fatal(err)
		}

		jobId := updateEnvResp.Payload.JobID

		err = waitJobFinish(t, client, jobId)
		if err != nil {
			fmt.Print(testutil.ExecCmd(t, "docker-compose logs openpitrix-rp-kubernetes"))
			fmt.Print(testutil.ExecCmd(t, "docker-compose logs hyperpitrix"))
			t.Fatal(err)
		}

		fmt.Printf("Cluster [%s] env updated", clusterId)
	}

	{
		fmt.Printf("Deleting cluster [%s]...", clusterId)

		deleteParams := cluster_manager.NewDeleteClustersParams()
		deleteParams.SetBody(&models.OpenpitrixDeleteClustersRequest{
			AdvancedParam: []string{},
			ClusterID:     []string{clusterId},
		})

		deleteResp, err := client.ClusterManager.DeleteClusters(deleteParams, nil)
		if err != nil {
			fmt.Print(testutil.ExecCmd(t, "docker-compose logs openpitrix-rp-kubernetes"))
			fmt.Print(testutil.ExecCmd(t, "docker-compose logs hyperpitrix"))
			t.Fatal(err)
		}

		jobId := deleteResp.Payload.JobID[0]

		err = waitJobFinish(t, client, jobId)
		if err != nil {
			fmt.Print(testutil.ExecCmd(t, "docker-compose logs openpitrix-rp-kubernetes"))
			fmt.Print(testutil.ExecCmd(t, "docker-compose logs hyperpitrix"))
			t.Fatal(err)
		}

		fmt.Printf("Cluster [%s] deleted", clusterId)
	}

	{
		fmt.Printf("Purging cluster [%s]...", clusterId)

		ceaseParams := cluster_manager.NewCeaseClustersParams()
		ceaseParams.SetBody(&models.OpenpitrixCeaseClustersRequest{
			AdvancedParam: []string{},
			ClusterID:     []string{clusterId},
		})

		ceaseResp, err := client.ClusterManager.CeaseClusters(ceaseParams, nil)
		if err != nil {
			fmt.Print(testutil.ExecCmd(t, "docker-compose logs openpitrix-rp-kubernetes"))
			fmt.Print(testutil.ExecCmd(t, "docker-compose logs hyperpitrix"))
			t.Fatal(err)
		}

		jobId := ceaseResp.Payload.JobID[0]

		err = waitJobFinish(t, client, jobId)
		if err != nil {
			fmt.Print(testutil.ExecCmd(t, "docker-compose logs openpitrix-rp-kubernetes"))
			fmt.Print(testutil.ExecCmd(t, "docker-compose logs hyperpitrix"))
			t.Fatal(err)
		}

		fmt.Printf("Cluster [%s] purged", clusterId)
	}
}

func waitJobFinish(_ *testing.T, client *apiclient.Openpitrix, jobId string) error {
	fmt.Printf("Waiting job [%s]", jobId)

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
		fmt.Printf("Unknown status [%s] for job [%s]", j.Status, jobId)
		return false, nil

	}, constants.WaitTaskTimeout, constants.WaitTaskInterval)
}

func getRuntimeCredential(t *testing.T) string {
	return testutil.ExecCmd(t, "kubectl config view --flatten")
}
