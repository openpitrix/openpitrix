// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package plugins

import (
	"time"

	"openpitrix.io/openpitrix/pkg/models"
)

var runtimePlugins map[string]RuntimeInterface

type RuntimeInterface interface {
	SplitJobIntoTasks(job *models.Job) (*models.TaskLayer, error)
	HandleSubtask(task *models.Task) error
	WaitSubtask(taskId string, timeout time.Duration, waitInterval time.Duration) error
}

func RegisterRuntimePlugin(runtime string, runtimeInterface RuntimeInterface) {
	runtimePlugins[runtime] = runtimeInterface
}

func GetRuntimePlugin(runtime string) RuntimeInterface {
	runtimeInterface, exists := runtimePlugins[runtime]
	if exists {
		return runtimeInterface
	} else {
		return nil
	}
}
