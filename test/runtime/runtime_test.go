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
	"openpitrix.io/openpitrix/test/testutil"
)

var clientConfig = testutil.GetClientConfig()

func getRuntimeCredential(t *testing.T) string {
	return testutil.ExecCmd(t, "kubectl config view --flatten")
}

func TestRuntime(t *testing.T) {
	credential := getRuntimeCredential(t)

	client := testutil.GetClient(clientConfig)

	// clean runtime
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

	// clean runtime credential
	describeCredentialParams := runtime_manager.NewDescribeRuntimeCredentialsParams()
	describeCredentialParams.SetSearchWord(&testRuntimeName)
	describeCredentialParams.SetStatus([]string{constants.StatusActive})
	describeCredentialResp, err := client.RuntimeManager.DescribeRuntimeCredentials(describeCredentialParams, nil)
	require.NoError(t, err)
	runtimeCredentials := describeCredentialResp.Payload.RuntimeCredentialSet
	for _, runtimeCredential := range runtimeCredentials {
		deleteCredentialParams := runtime_manager.NewDeleteRuntimeCredentialsParams()
		deleteCredentialParams.SetBody(
			&models.OpenpitrixDeleteRuntimeCredentialsRequest{
				RuntimeCredentialID: []string{runtimeCredential.RuntimeCredentialID},
			})
		_, err := client.RuntimeManager.DeleteRuntimeCredentials(deleteCredentialParams, nil)
		require.NoError(t, err)
	}

	// create runtime credential
	createCredentialParams := runtime_manager.NewCreateRuntimeCredentialParams()
	createCredentialParams.SetBody(
		&models.OpenpitrixCreateRuntimeCredentialRequest{
			Name:                     testRuntimeName,
			Description:              "description",
			Provider:                 constants.ProviderKubernetes,
			RuntimeCredentialContent: credential,
		})
	createCredentialResp, err := client.RuntimeManager.CreateRuntimeCredential(createCredentialParams, nil)
	require.NoError(t, err)
	runtimeCredentialId := createCredentialResp.Payload.RuntimeCredentialID

	// create runtime
	createParams := runtime_manager.NewCreateRuntimeParams()
	createParams.SetBody(
		&models.OpenpitrixCreateRuntimeRequest{
			Name:                testRuntimeName,
			Description:         "description",
			Provider:            constants.ProviderKubernetes,
			RuntimeCredentialID: runtimeCredentialId,
			Zone:                idutil.GetUuid36("r-"),
		})
	createResp, err := client.RuntimeManager.CreateRuntime(createParams, nil)
	require.NoError(t, err)
	runtimeId := createResp.Payload.RuntimeID

	// modify runtime credential
	modifyCredentialParams := runtime_manager.NewModifyRuntimeCredentialParams()
	modifyCredentialParams.SetBody(
		&models.OpenpitrixModifyRuntimeCredentialRequest{
			RuntimeCredentialID: runtimeCredentialId,
			Description:         "cc",
		})
	_, err = client.RuntimeManager.ModifyRuntimeCredential(modifyCredentialParams, nil)
	require.NoError(t, err)

	// modify runtime
	modifyParams := runtime_manager.NewModifyRuntimeParams()
	modifyParams.SetBody(
		&models.OpenpitrixModifyRuntimeRequest{
			RuntimeID:   runtimeId,
			Description: "cc",
		})
	_, err = client.RuntimeManager.ModifyRuntime(modifyParams, nil)
	require.NoError(t, err)

	// describe runtime credential
	describeCredentialParams.WithRuntimeCredentialID([]string{runtimeCredentialId})
	describeCredentialResp, err = client.RuntimeManager.DescribeRuntimeCredentials(describeCredentialParams, nil)
	require.NoError(t, err)
	runtimeCredentials = describeCredentialResp.Payload.RuntimeCredentialSet
	if len(runtimeCredentials) != 1 {
		t.Fatalf("failed to describe runtime credentialss with params [%+v]", describeCredentialParams)
	}
	if runtimeCredentials[0].Name != testRuntimeName || runtimeCredentials[0].Description != "cc" {
		t.Fatalf("failed to modify runtime credential [%+v]", runtimeCredentials[0])
	}

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

	// delete runtime credential
	deleteCredentialParams := runtime_manager.NewDeleteRuntimeCredentialsParams()
	deleteCredentialParams.WithBody(&models.OpenpitrixDeleteRuntimeCredentialsRequest{
		RuntimeCredentialID: []string{runtimeCredentialId},
	})
	_, err = client.RuntimeManager.DeleteRuntimeCredentials(deleteCredentialParams, nil)
	require.NoError(t, err)

	// describe deleted runtime
	describeParams.WithRuntimeID([]string{runtimeId})
	describeParams.WithStatus([]string{constants.StatusDeleted})
	describeResp, err = client.RuntimeManager.DescribeRuntimes(describeParams, nil)
	require.NoError(t, err)
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

	// describe deleted runtime credential
	describeCredentialParams.WithRuntimeCredentialID([]string{runtimeCredentialId})
	describeCredentialParams.WithStatus([]string{constants.StatusDeleted})
	describeCredentialResp, err = client.RuntimeManager.DescribeRuntimeCredentials(describeCredentialParams, nil)
	require.NoError(t, err)
	runtimeCredentials = describeCredentialResp.Payload.RuntimeCredentialSet
	if len(runtimeCredentials) != 1 {
		t.Fatalf("failed to describe runtime credentials with params [%+v]", describeCredentialParams)
	}
	runtimeCredential := runtimeCredentials[0]
	if runtimeCredential.RuntimeCredentialID != runtimeCredentialId {
		t.Fatalf("failed to describe runtime credential")
	}
	if runtimeCredential.Status != constants.StatusDeleted {
		t.Fatalf("failed to delete runtime credential, got runtime credential status [%s]", runtimeCredential.Status)
	}

	t.Log("test runtime finish, all test is ok")
}
