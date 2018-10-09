// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

// +build k8s

package runtime

import (
	"testing"

	"github.com/stretchr/testify/require"

	"openpitrix.io/openpitrix/pkg/constants"
	"openpitrix.io/openpitrix/pkg/util/idutil"
	"openpitrix.io/openpitrix/test/client/runtime_manager"
	"openpitrix.io/openpitrix/test/models"
	"openpitrix.io/openpitrix/test/repocommon"
	"openpitrix.io/openpitrix/test/testutil"
)

var clientConfig = testutil.GetClientConfig()

func getRuntimeCredential(t *testing.T) string {
	return testutil.ExecCmd(t, "kubectl config view --flatten")
}

func TestRuntime(t *testing.T) {
	credential := getRuntimeCredential(t)

	client := testutil.GetClient(clientConfig)

	testRuntimeName := "e2e-test-runtime"
	describeParams := runtime_manager.NewDescribeRuntimesParams()
	describeParams.SetSearchWord(&testRuntimeName)
	describeParams.SetStatus([]string{constants.StatusActive})
	describeResp, err := client.RuntimeManager.DescribeRuntimes(describeParams, nil)
	require.NoError(t, err)
	runtimes := describeResp.Payload.RuntimeSet
	for _, runtime := range runtimes {
		deleteParams := runtime_manager.NewDeleteRuntimesParams()
		deleteParams.SetBody(
			&models.OpenpitrixDeleteRuntimesRequest{
				RuntimeID: []string{runtime.RuntimeID},
			})
		_, err := client.RuntimeManager.DeleteRuntimes(deleteParams, nil)
		require.NoError(t, err)
	}
	// create runtime
	createParams := runtime_manager.NewCreateRuntimeParams()
	createParams.SetBody(
		&models.OpenpitrixCreateRuntimeRequest{
			Name:              testRuntimeName,
			Description:       "description",
			Provider:          constants.ProviderKubernetes,
			RuntimeURL:        "",
			RuntimeCredential: credential,
			Zone:              idutil.GetUuid36("r-"),
		})
	createResp, err := client.RuntimeManager.CreateRuntime(createParams, nil)
	require.NoError(t, err)
	runtimeId := createResp.Payload.RuntimeID
	// modify runtime
	modifyParams := runtime_manager.NewModifyRuntimeParams()
	modifyParams.SetBody(
		&models.OpenpitrixModifyRuntimeRequest{
			RuntimeID:   runtimeId,
			Description: "cc",
		})
	modifyResp, err := client.RuntimeManager.ModifyRuntime(modifyParams, nil)
	require.NoError(t, err)
	t.Log(modifyResp)
	// describe runtime
	describeParams.WithRuntimeID([]string{runtimeId})
	describeResp, err = client.RuntimeManager.DescribeRuntimes(describeParams, nil)
	require.NoError(t, err)
	runtimes = describeResp.Payload.RuntimeSet
	if len(runtimes) != 1 {
		t.Fatalf("failed to describe runtimes with params [%+v]", describeParams)
	}
	if runtimes[0].Name != testRuntimeName || runtimes[0].Description != "cc" {
		t.Fatalf("failed to modify runtime [%+v]", runtimes[0])
	}
	// delete runtime
	deleteParams := runtime_manager.NewDeleteRuntimesParams()
	deleteParams.WithBody(&models.OpenpitrixDeleteRuntimesRequest{
		RuntimeID: []string{runtimeId},
	})
	deleteResp, err := client.RuntimeManager.DeleteRuntimes(deleteParams, nil)
	require.NoError(t, err)
	t.Log(deleteResp)
	// describe deleted runtime
	describeParams.WithRuntimeID([]string{runtimeId})
	describeParams.WithStatus([]string{constants.StatusDeleted})
	describeParams.WithSearchWord(nil)
	describeResp, err = client.RuntimeManager.DescribeRuntimes(describeParams, nil)
	require.NoError(t, err)
	runtimes = describeResp.Payload.RuntimeSet
	runtimes = describeResp.Payload.RuntimeSet
	if len(runtimes) != 1 {
		t.Fatalf("failed to describe runtimes with params [%+v]", describeParams)
	}
	runtime := runtimes[0]
	if runtime.RuntimeID != runtimeId {
		t.Fatalf("failed to describe runtime")
	}
	if runtime.Status != constants.StatusDeleted {
		t.Fatalf("failed to delete runtime, got runtime status [%s]", runtime.Status)
	}

	t.Log("test runtime finish, all test is ok")
}

func TestRuntimeLabel(t *testing.T) {
	credential := getRuntimeCredential(t)
	client := testutil.GetClient(clientConfig)
	// Create a test runtime that can attach label on it
	testRuntimeName := "e2e-test-runtime"
	labels := repocommon.GenerateLabels()
	createParams := runtime_manager.NewCreateRuntimeParams()
	createParams.SetBody(
		&models.OpenpitrixCreateRuntimeRequest{
			Name:              testRuntimeName,
			Description:       "description",
			Provider:          constants.ProviderKubernetes,
			RuntimeURL:        "",
			RuntimeCredential: credential,
			Zone:              idutil.GetUuid36("r-"),
			Labels:            labels,
		})
	createResp, err := client.RuntimeManager.CreateRuntime(createParams, nil)
	require.NoError(t, err)
	runtimeId := createResp.Payload.RuntimeID

	describeParams := runtime_manager.NewDescribeRuntimesParams()
	describeParams.Label = &labels
	describeParams.Status = []string{constants.StatusActive}
	describeResp, err := client.RuntimeManager.DescribeRuntimes(describeParams, nil)
	require.NoError(t, err)
	if len(describeResp.Payload.RuntimeSet) != 1 {
		t.Fatalf("describe runtime with filter failed")
	}
	if describeResp.Payload.RuntimeSet[0].RuntimeID != runtimeId {
		t.Fatalf("describe runtime with filter failed")
	}

	// delete runtime
	deleteParams := runtime_manager.NewDeleteRuntimesParams()
	deleteParams.WithBody(&models.OpenpitrixDeleteRuntimesRequest{
		RuntimeID: []string{runtimeId},
	})
	deleteResp, err := client.RuntimeManager.DeleteRuntimes(deleteParams, nil)
	require.NoError(t, err)
	t.Log(deleteResp)

	t.Log("test runtime label finish, all test is ok")
}
