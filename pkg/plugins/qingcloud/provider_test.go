// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package qingcloud

import (
	"testing"

	runtimeclient "openpitrix.io/openpitrix/pkg/client/runtime"
	"openpitrix.io/openpitrix/pkg/constants"
	"openpitrix.io/openpitrix/pkg/models"
	"openpitrix.io/openpitrix/pkg/plugins/vmbased"
)

type ActionNum struct {
	Action string
	Num    int
}

func testCreateCluster(t *testing.T, frame *vmbased.Frame) {
	rootTaskLayer := frame.CreateClusterLayer()

	expectResult := []ActionNum{
		{vmbased.ActionCreateVolumes, 5},
		{vmbased.ActionRunInstances, 5},
		{vmbased.ActionWaitFrontgateAvailable, 1},
		{vmbased.ActionRegisterMetadata, 1},
		{vmbased.ActionStartConfd, 5},
		{vmbased.ActionRegisterCmd, 1}, // hbase-hdfs-master init
		{vmbased.ActionRegisterCmd, 1}, // hbase-hdfs-master start
		{vmbased.ActionRegisterCmd, 1}, // hbase-master start
		{vmbased.ActionRegisterCmd, 3}, // hbase-slave start
		{vmbased.ActionDeregesterCmd, 5},
	}

	var result []ActionNum
	for rootTaskLayer != nil {
		result = append(result, ActionNum{rootTaskLayer.Tasks[0].TaskAction, len(rootTaskLayer.Tasks)})
		rootTaskLayer = rootTaskLayer.Child
	}

	if len(result) != len(expectResult) {
		t.Errorf("Expect [%d] task layer, while get [%s] task layer", len(expectResult), len(result))
	}

	for index := range result {
		if result[index] != expectResult[index] {
			t.Errorf("Index [%d] expect [%s], while get [%s]", index, expectResult[index], result[index])
		}
	}
}

func TestSplitJobIntoTasks(t *testing.T) {
	parser := Parser{}
	clusterWrapper, err := parser.Parse([]byte(hbaseMustache))
	if err != nil {
		t.Errorf("Parse mustache failed: %+v", err)
	}
	directive, _ := clusterWrapper.ToString()

	mockJob := &models.Job{
		JobId:     "j-1234",
		Owner:     "usr-1234",
		ClusterId: "cl-1234",
		Directive: directive,
		JobAction: constants.ActionCreateCluster,
	}

	runtime := &runtimeclient.Runtime{
		RuntimeId:         "rt-1234",
		Provider:          constants.ProviderQingCloud,
		Zone:              "testing",
		ProviderInterface: &Provider{},
	}

	frame := &vmbased.Frame{
		Job:            mockJob,
		ClusterWrapper: clusterWrapper,
		Runtime:        runtime,
	}
	testCreateCluster(t, frame)
}
