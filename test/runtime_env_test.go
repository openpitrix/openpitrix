package test

import (
	"reflect"
	"testing"

	"openpitrix.io/openpitrix/pkg/constants"
	"openpitrix.io/openpitrix/test/client/runtime_env_manager"
	"openpitrix.io/openpitrix/test/models"
)

func TestRuntimeEnvCredential(t *testing.T) {
	client := GetClient(clientConfig)
	testRuntimeEnvCredentialName := "e2e_test_runtime_env_credential"
	testRuntimeEnvCredentialName2 := "e2e_test_runtime_env_credential"
	testRuntimeEnvCredentialDescription := "e2e_test_runtime_env_credential_description"
	testRuntimeEnvCredentialContent := map[string]string{
		"access_key": "ak",
		"secret_key": "sk",
	}
	testRuntimeEnvCredentialContent2 := map[string]string{
		"auth": "value",
	}

	describeParams := runtime_env_manager.NewDescribeRuntimeEnvCredentialsParams()
	describeParams.SetSearchWord(&testRuntimeEnvCredentialName)
	describeParams.SetStatus([]string{constants.StatusActive})
	describeResp, err := client.RuntimeEnvManager.DescribeRuntimeEnvCredentials(describeParams)
	if err != nil {
		t.Fatal(err)
	}
	runtimeEnvCredentials := describeResp.Payload.RuntimeEnvCredentialSet

	for _, runtimeEnvCredential := range runtimeEnvCredentials {
		deleteParams := runtime_env_manager.NewDeleteRuntimeEnvCredentialParams()
		deleteParams.SetBody(&models.OpenpitrixDeleteRuntimeEnvCredentialRequset{
			RuntimeEnvCredentialID: runtimeEnvCredential.RuntimeEnvCredentialID,
		})

		_, err := client.RuntimeEnvManager.DeleteRuntimeEnvCredential(deleteParams)
		if err != nil {
			t.Fatal(err)
		}
	}

	createParams := runtime_env_manager.NewCreateRuntimeEnvCredentialParams()
	createParams.SetBody(&models.OpenpitrixCreateRuntimeEnvCredentialRequset{
		Name:        testRuntimeEnvCredentialName,
		Description: testRuntimeEnvCredentialDescription,
		Content:     testRuntimeEnvCredentialContent,
	},
	)
	createResp, err := client.RuntimeEnvManager.CreateRuntimeEnvCredential(createParams)
	if err != nil {
		t.Fatal(err)
	}

	describeParams = runtime_env_manager.NewDescribeRuntimeEnvCredentialsParams()
	describeParams.WithRuntimeEnvCredentialID([]string{createResp.Payload.RuntimeEnvCredential.RuntimeEnvCredentialID})
	describeResp, err = client.RuntimeEnvManager.DescribeRuntimeEnvCredentials(describeParams)
	if err != nil {
		t.Fatal(err)
	}

	runtimeEnvCredentials = describeResp.Payload.RuntimeEnvCredentialSet
	if len(runtimeEnvCredentials) != 1 {
		t.Fatalf("failed to describe runtime_env_credential with params [%+v]", describeParams)
	}
	if runtimeEnvCredentials[0].Name != testRuntimeEnvCredentialName ||
		runtimeEnvCredentials[0].Description != testRuntimeEnvCredentialDescription ||
		!reflect.DeepEqual(runtimeEnvCredentials[0].Content, testRuntimeEnvCredentialContent) {
		t.Fatalf("failed to create runtime_env_credential [%+v]", runtimeEnvCredentials[0])
	}

	modifyParams := runtime_env_manager.NewModifyRuntimeEnvCredentialParams()
	modifyParams.SetBody(&models.OpenpitrixModifyRuntimeEnvCredentialRequest{
		RuntimeEnvCredentialID: createResp.Payload.RuntimeEnvCredential.RuntimeEnvCredentialID,
		Name:    testRuntimeEnvCredentialName2,
		Content: testRuntimeEnvCredentialContent2,
	},
	)
	modifyResp, err := client.RuntimeEnvManager.ModifyRuntimeEnvCredential(modifyParams)
	if err != nil {
		t.Fatal(err)
	}

	describeParams = runtime_env_manager.NewDescribeRuntimeEnvCredentialsParams()
	describeParams.WithRuntimeEnvCredentialID([]string{modifyResp.Payload.RuntimeEnvCredential.RuntimeEnvCredentialID})
	describeResp, err = client.RuntimeEnvManager.DescribeRuntimeEnvCredentials(describeParams)
	if err != nil {
		t.Fatal(err)
	}

	runtimeEnvCredentials = describeResp.Payload.RuntimeEnvCredentialSet
	if len(runtimeEnvCredentials) != 1 {
		t.Fatalf("failed to describe runtime_env_credential with params [%+v]", describeParams)
	}
	if runtimeEnvCredentials[0].Name != testRuntimeEnvCredentialName2 ||
		runtimeEnvCredentials[0].Description != testRuntimeEnvCredentialDescription ||
		!reflect.DeepEqual(runtimeEnvCredentials[0].Content, testRuntimeEnvCredentialContent2) {
		t.Fatalf("failed to modify runtime_env_credential [%+v]", runtimeEnvCredentials[0])
	}

	deleteParams := runtime_env_manager.NewDeleteRuntimeEnvCredentialParams()
	deleteParams.SetBody(&models.OpenpitrixDeleteRuntimeEnvCredentialRequset{
		RuntimeEnvCredentialID: runtimeEnvCredentials[0].RuntimeEnvCredentialID,
	})

	deleteResp, err := client.RuntimeEnvManager.DeleteRuntimeEnvCredential(deleteParams)
	if err != nil {
		t.Fatal(err)
	}

	describeParams = runtime_env_manager.NewDescribeRuntimeEnvCredentialsParams()
	describeParams.WithRuntimeEnvCredentialID([]string{deleteResp.Payload.RuntimeEnvCredential.RuntimeEnvCredentialID})
	describeResp, err = client.RuntimeEnvManager.DescribeRuntimeEnvCredentials(describeParams)
	if err != nil {
		t.Fatal(err)
	}

	runtimeEnvCredentials = describeResp.Payload.RuntimeEnvCredentialSet
	if len(runtimeEnvCredentials) != 1 ||
		runtimeEnvCredentials[0].RuntimeEnvCredentialID != deleteResp.Payload.RuntimeEnvCredential.RuntimeEnvCredentialID {
		t.Fatalf("failed to describe runtime_env_credential with params [%+v]", describeParams)
	}
	if runtimeEnvCredentials[0].Status != constants.StatusDeleted {
		t.Fatalf("failed to delete runtime_env_credential, got runtime_env_credential status [%s]", runtimeEnvCredentials[0].Status)
	}
	t.Log("test runtime_env_credential finish, all test is ok")
}
