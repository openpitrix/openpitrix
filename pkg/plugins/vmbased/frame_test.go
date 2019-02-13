// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package vmbased

import (
	"context"
	"testing"

	"openpitrix.io/openpitrix/pkg/config"
	"openpitrix.io/openpitrix/pkg/constants"
	"openpitrix.io/openpitrix/pkg/models"
	"openpitrix.io/openpitrix/pkg/util/jsonutil"
)

type ActionNum struct {
	Action string
	Num    int
}

func testCreateCluster(t *testing.T, frame *Frame) {
	rootTaskLayer := frame.CreateClusterLayer()

	expectResult := []ActionNum{
		{ActionCreateVolumes, 5},
		{ActionWaitFrontgateAvailable, 1},
		{ActionRunInstances, 5},
		{ActionPingDrone, 5},
		{ActionSetDroneConfig, 5},
		{ActionFormatAndMountVolume, 5},
		{ActionRemoveContainerOnDrone, 5},
		{ActionPingDrone, 5},
		{ActionSetDroneConfig, 5},
		{ActionDeregisterMetadata, 1},
		{ActionDeregisterMetadataMapping, 1},
		{ActionRegisterMetadata, 1},
		{ActionRegisterMetadataMapping, 1},
		{ActionStartConfd, 5},
		{ActionRegisterCmd, 1}, // hbase-hdfs-master init
		{ActionRegisterCmd, 1}, // hbase-hdfs-master start
		{ActionRegisterCmd, 1}, // hbase-master start
		{ActionRegisterCmd, 3}, // hbase-slave start
		{ActionDeregisterCmd, 5},
	}

	var result []ActionNum
	for rootTaskLayer != nil {
		result = append(result, ActionNum{rootTaskLayer.Tasks[0].TaskAction, len(rootTaskLayer.Tasks)})
		rootTaskLayer = rootTaskLayer.Child
	}

	if len(result) != len(expectResult) {
		t.Errorf("Expect [%d] task layer, while get [%d] task layer", len(expectResult), len(result))
	}

	for index := range result {
		if result[index] != expectResult[index] {
			t.Errorf("Index [%d] expect [%+v], while get [%+v]", index, expectResult[index], result[index])
		}
	}
}

func TestSplitJobIntoTasks(t *testing.T) {
	clusterWrapper := getTestClusterWrapper(t)
	directive := jsonutil.ToString(clusterWrapper)

	mockJob := &models.Job{
		JobId:     "j-1234",
		Owner:     "usr-1234",
		ClusterId: "cl-1234",
		Directive: directive,
		JobAction: constants.ActionCreateCluster,
	}

	runtime := new(models.RuntimeDetails)
	runtime.RuntimeId = "rt-1234"
	runtime.Runtime.Provider = constants.ProviderQingCloud
	runtime.Zone = "testing"

	frame := &Frame{
		Job:            mockJob,
		ClusterWrapper: clusterWrapper,
		Runtime:        runtime,
		Ctx:            context.Background(),
		RuntimeProviderConfig: &config.RuntimeProviderConfig{
			ImageId: "img:abcd",
		},
	}
	testCreateCluster(t, frame)
}
