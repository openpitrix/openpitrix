// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package task

import (
	"openpitrix.io/openpitrix/pkg/logger"
	"openpitrix.io/openpitrix/pkg/models"
)

type Processor struct {
	Task *models.Task
}

func NewProcessor(task *models.Task) *Processor {
	return &Processor{
		Task: task,
	}
}

// Post process when task is done
func (t *Processor) Post() {
	var err error
	switch t.Task.TaskAction {
	// TODO: case TaskAction
	default:
		logger.Errorf("Unknown job action [%s]", t.Task.TaskAction)
	}
	if err != nil {
		logger.Errorf("Executing task [%s] post processor failed: %+v", t.Task.TaskId, err)
	}
}
