// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

// +build aws

// This file maybe discard

package aws

import (
	"log"
	"testing"
	"time"

	"openpitrix.io/openpitrix/test/client/app_manager"
	"openpitrix.io/openpitrix/test/client/cluster_manager"
	"openpitrix.io/openpitrix/test/client/repo_manager"
	"openpitrix.io/openpitrix/test/client/runtime_manager"
	"openpitrix.io/openpitrix/test/models"
	"openpitrix.io/openpitrix/test/testutil"
)

var (
	clientConfig = &testutil.ClientConfig{
		Host:  "192.168.0.6:9100",
		Debug: true,
	}

	RepoNameForTest    = "incubator"
	RuntimeNameForTest = "qingcloud runtime"
	RepoUrlForTest     = "http://139.198.121.182:8879/"

	ClusterConf = `{
    "cluster": {
        "name": "ELK",
        "description": "The description of the ELK service",
        "subnet": "subnet-03f37679",
        "es_node": {
            "cpu": 1,
            "memory": 1024,
            "count": 3,
            "instance_class": 1,
            "volume_size": 10
        },
        "kbn_node": {
            "cpu": 1,
            "memory": 1024,
            "count": 1,
            "instance_class": 1,
            "volume_size": 10
        },
        "lst_node": {
            "cpu": 1,
            "memory": 1024,
            "count": 1,
            "instance_class": 1,
            "volume_size": 10
        }
    },
    "env": {
        "es_node": {
            "action_destructive_requires_name": "true",
            "indices_fielddata_cache_size": "90%",
            "logstash_node_ip": "",
            "discovery_zen_no_master_block": "write",
            "gateway_recover_after_time": "5m",
            "http_cors_enabled": "true",
            "http_cors_allow_origin": "*",
            "indices_queries_cache_size": "10%",
            "indices_memory_index_buffer_size": "10%",
            "indices_requests_cache_size": "2%",
            "script_inline": "true",
            "script_stored": "true",
            "script_file": "false",
            "script_aggs": "true",
            "script_search": "true",
            "script_update": "true",
            "remote_ext_dict": "",
            "remote_ext_stopwords": ""
        },
        "lst_node": {
            "input_conf_content": "http { port => 9700 }",
            "filter_conf_content": "",
            "output_conf_content": "",
            "output_es_content": "",
            "gemfile_append_content": ""
        }
    }
}`
)

func TestAWS(t *testing.T) {
	log.SetPrefix("[ === AWS TEST === ] ")

	client := testutil.GetClient(clientConfig)

	// create repo
	var repoID string
	{
		describeParams := repo_manager.NewDescribeReposParams()
		describeParams.SetName([]string{RepoNameForTest})
		describeResp, err := client.RepoManager.DescribeRepos(describeParams)
		if err != nil {
			t.Fatal(err)
		}
		repos := describeResp.Payload.RepoSet

		if len(repos) != 0 {
			repoID = repos[0].RepoID
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
					Providers:   []string{"qingcloud"},
				})
			createResp, err := client.RepoManager.CreateRepo(createParams)
			if err != nil {
				t.Fatal(err)
			}
			repoID = createResp.Payload.RepoID
		}
	}
	log.Printf("Got repo [%s]\n", repoID)

	// waiting for apps indexed by repo indexer
	var app *models.OpenpitrixApp
	{
		for {
			describeParams := app_manager.NewDescribeAppsParams()
			describeParams.WithRepoID([]string{repoID})
			describeParams.WithName([]string{"elk"})
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

		if len(appVersions) != 1 {
			t.Fatal("We need only one version to test")
		}

		appVersion = appVersions[0]
	}
	log.Printf("Got app version [%s]\n", appVersion.Name)

	// create runtime
	var runtimeID string
	{
		describeParams := runtime_manager.NewDescribeRuntimesParams()
		describeParams.SetSearchWord(&RuntimeNameForTest)
		describeResp, err := client.RuntimeManager.DescribeRuntimes(describeParams)
		if err != nil {
			t.Fatal(err)
		}
		runtimes := describeResp.Payload.RuntimeSet
		if len(runtimes) != 0 {
			runtimeID = runtimes[0].RuntimeID
		} else {
			createParams := runtime_manager.NewCreateRuntimeParams()
			createParams.SetBody(
				&models.OpenpitrixCreateRuntimeRequest{
					Name:              RuntimeNameForTest,
					Description:       "aws runtime",
					Provider:          "aws",
					RuntimeURL:        "https://ec2.us-east-2.amazonaws.com",
					RuntimeCredential: `{"access_key_id": "xxxxxxxxxxxxxxx", "secret_access_key": "xxxxxxxxxxxxxxxx"}`,
					Zone:              "us-east-2",
				})
			createResp, err := client.RuntimeManager.CreateRuntime(createParams)
			if err != nil {
				t.Fatal(err)
			}
			runtimeID = createResp.Payload.RuntimeID
		}
	}
	log.Printf("Got runtime [%s]\n", runtimeID)

	var clusterId string
	log.Printf("Creating cluster...\n")
	{
		createParams := cluster_manager.NewCreateClusterParams()
		createParams.SetBody(&models.OpenpitrixCreateClusterRequest{
			AdvancedParam: []string{},
			AppID:         app.AppID,
			Zone:          "us-east-2b",
			Conf:          ClusterConf,
			RuntimeID:     runtimeID,
			VersionID:     appVersion.VersionID,
		})

		createResp, err := client.ClusterManager.CreateCluster(createParams)
		if err != nil {
			t.Fatal(err)
		}

		clusterId = createResp.Payload.ClusterID
	}
	log.Printf("Cluster [%s] created \n", clusterId)
}
