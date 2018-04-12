package test

import (
	"net/url"
	"testing"

	"openpitrix.io/openpitrix/pkg/constants"
	"openpitrix.io/openpitrix/pkg/utils"
	"openpitrix.io/openpitrix/test/client/runtime_manager"
	"openpitrix.io/openpitrix/test/models"
)

func TestRuntime(t *testing.T) {
	client := GetClient(clientConfig)

	testRuntimeName := "e2e-test-runtime"
	describeParams := runtime_manager.NewDescribeRuntimesParams()
	describeParams.SetSearchWord(&testRuntimeName)
	describeParams.SetStatus([]string{constants.StatusActive})
	describeResp, err := client.RuntimeManager.DescribeRuntimes(describeParams)
	if err != nil {
		t.Fatal(err)
	}
	runtimes := describeResp.Payload.RuntimeSet
	for _, runtime := range runtimes {
		deleteParams := runtime_manager.NewDeleteRuntimeParams()
		deleteParams.SetBody(
			&models.OpenpitrixDeleteRuntimeRequest{
				RuntimeID: runtime.RuntimeID,
			})
		_, err := client.RuntimeManager.DeleteRuntime(deleteParams)
		if err != nil {
			t.Fatal(err)
		}
	}
	// create runtime
	createParams := runtime_manager.NewCreateRuntimeParams()
	createParams.SetBody(
		&models.OpenpitrixCreateRuntimeRequest{
			Name:              testRuntimeName,
			Description:       "description",
			Provider:          "qingcloud",
			RuntimeURL:        "https://github.com/",
			RuntimeCredential: `{}`,
			Zone:              "test",
		})
	createResp, err := client.RuntimeManager.CreateRuntime(createParams)
	if err != nil {
		t.Logf("Create runtime will fail without credential")
		return
	}
	runtimeId := createResp.Payload.Runtime.RuntimeID
	// modify runtime
	modifyParams := runtime_manager.NewModifyRuntimeParams()
	modifyParams.SetBody(
		&models.OpenpitrixModifyRuntimeRequest{
			RuntimeID:   runtimeId,
			Description: "cc",
		})
	modifyResp, err := client.RuntimeManager.ModifyRuntime(modifyParams)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(modifyResp)
	// describe runtime
	describeParams.WithRuntimeID([]string{runtimeId})
	describeResp, err = client.RuntimeManager.DescribeRuntimes(describeParams)
	if err != nil {
		t.Fatal(err)
	}
	runtimes = describeResp.Payload.RuntimeSet
	if len(runtimes) != 1 {
		t.Fatalf("failed to describe runtimes with params [%+v]", describeParams)
	}
	if runtimes[0].Name != testRuntimeName || runtimes[0].Description != "cc" {
		t.Fatalf("failed to modify runtime [%+v]", runtimes[0])
	}
	// delete runtime
	deleteParams := runtime_manager.NewDeleteRuntimeParams()
	deleteParams.WithBody(&models.OpenpitrixDeleteRuntimeRequest{
		RuntimeID: runtimeId,
	})
	deleteResp, err := client.RuntimeManager.DeleteRuntime(deleteParams)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(deleteResp)
	// describe deleted runtime
	describeParams.WithRuntimeID([]string{runtimeId})
	describeParams.WithStatus([]string{constants.StatusDeleted})
	describeParams.WithSearchWord(nil)
	describeResp, err = client.RuntimeManager.DescribeRuntimes(describeParams)
	if err != nil {
		t.Fatal(err)
	}
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

func generateRuntimeLabels() string {
	v := url.Values{}
	v.Add("key1", utils.GetUuid(""))
	v.Add("key2", utils.GetUuid(""))
	v.Add("key3", utils.GetUuid(""))
	v.Add("key4", utils.GetUuid(""))
	v.Add("key5", utils.GetUuid(""))
	return v.Encode()
}

func TestRuntimeLabel(t *testing.T) {
	client := GetClient(clientConfig)
	// Create a test runtime that can attach label on it
	testRuntimeName := "e2e-test-runtime"
	labels := generateRuntimeLabels()
	createParams := runtime_manager.NewCreateRuntimeParams()
	createParams.SetBody(
		&models.OpenpitrixCreateRuntimeRequest{
			Name:              testRuntimeName,
			Description:       "description",
			Provider:          "qingcloud",
			RuntimeURL:        "https://github.com/",
			RuntimeCredential: `{}`,
			Zone:              "test",
			Labels:            labels,
		})
	createResp, err := client.RuntimeManager.CreateRuntime(createParams)
	if err != nil {
		t.Logf("Create runtime will fail without credential")
		return
	}
	runtimeId := createResp.Payload.Runtime.RuntimeID

	describeParams := runtime_manager.NewDescribeRuntimesParams()
	describeParams.Label = &labels
	describeParams.Status = []string{constants.StatusActive}
	describeResp, err := client.RuntimeManager.DescribeRuntimes(describeParams)
	if err != nil {
		t.Fatal(err)
	}
	if len(describeResp.Payload.RuntimeSet) != 1 {
		t.Fatalf("describe runtime with filter failed")
	}
	if describeResp.Payload.RuntimeSet[0].RuntimeID != runtimeId {
		t.Fatalf("describe runtime with filter failed")
	}

	// delete runtime
	deleteParams := runtime_manager.NewDeleteRuntimeParams()
	deleteParams.WithBody(&models.OpenpitrixDeleteRuntimeRequest{
		RuntimeID: runtimeId,
	})
	deleteResp, err := client.RuntimeManager.DeleteRuntime(deleteParams)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(deleteResp)

	t.Log("test runtime label finish, all test is ok")
}
