// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package plugins

import (
	"fmt"
	"time"

	"openpitrix.io/openpitrix/pkg/models"
)

var runtimePlugins map[string]RuntimeInterface

type RuntimeInterface interface {
	// Parse package and conf into cluster which clusterManager will register into db.
	ParseClusterConf(versionId, conf string) (*models.Cluster, error)
	SplitJobIntoTasks(job *models.Job) (*models.TaskLayer, error)
	HandleSubtask(task *models.Task) error
	WaitSubtask(taskId string, timeout time.Duration, waitInterval time.Duration) error
}

func RegisterRuntimePlugin(runtime string, runtimeInterface RuntimeInterface) {
	runtimePlugins[runtime] = runtimeInterface
}

func GetRuntimePlugin(runtime string) (RuntimeInterface, error) {
	runtimeInterface, exists := runtimePlugins[runtime]
	if exists {
		return runtimeInterface, nil
	} else {
		return nil, fmt.Errorf("No such runtime [%s]. ", runtime)
	}
}
