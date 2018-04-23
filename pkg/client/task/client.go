// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package task

import (
	"context"
	"fmt"
	"time"

	"openpitrix.io/openpitrix/pkg/client"
	"openpitrix.io/openpitrix/pkg/constants"
	"openpitrix.io/openpitrix/pkg/logger"
	"openpitrix.io/openpitrix/pkg/manager"
	"openpitrix.io/openpitrix/pkg/models"
	"openpitrix.io/openpitrix/pkg/pb"
	"openpitrix.io/openpitrix/pkg/utils"
)

func NewTaskManagerClient(ctx context.Context) (pb.TaskManagerClient, error) {
	conn, err := manager.NewClient(ctx, constants.TaskManagerHost, constants.TaskManagerPort)
	if err != nil {
		return nil, err
	}
	return pb.NewTaskManagerClient(conn), err
}

func CreateTask(ctx context.Context, taskRequest *pb.CreateTaskRequest) (taskId string, err error) {
	taskManagerClient, err := NewTaskManagerClient(ctx)
	if err != nil {
		return
	}
	taskResponse, err := taskManagerClient.CreateTask(ctx, taskRequest)
	if err != nil {
		return
	}
	taskId = taskResponse.GetTaskId().GetValue()
	return
}

func DescribeTasks(ctx context.Context, taskRequest *pb.DescribeTasksRequest) (*pb.DescribeTasksResponse, error) {
	taskManagerClient, err := NewTaskManagerClient(ctx)
	if err != nil {
		return nil, err
	}
	taskResponse, err := taskManagerClient.DescribeTasks(ctx, taskRequest)
	if err != nil {
		return nil, err
	}
	return taskResponse, err
}

func WaitTask(taskId string, timeout time.Duration, waitInterval time.Duration) error {
	logger.Debugf("Waiting for task [%s] finished", taskId)
	return utils.WaitForSpecificOrError(func() (bool, error) {
		taskRequest := &pb.DescribeTasksRequest{
			TaskId: []string{taskId},
		}
		taskResponse, err := DescribeTasks(client.GetSystemUserContext(), taskRequest)
		if err != nil {
			//network or api error, not considered task fail.
			return false, nil
		}
		if len(taskResponse.TaskSet) == 0 {
			return false, fmt.Errorf("Can not find task [%s]. ", taskId)
		}
		t := taskResponse.TaskSet[0]
		if t.Status == nil {
			logger.Errorf("Task [%s] status is nil", taskId)
			return false, nil
		}
		if t.Status.GetValue() == constants.StatusWorking || t.Status.GetValue() == constants.StatusPending {
			return false, nil
		}
		if t.Status.GetValue() == constants.StatusSuccessful {
			return true, nil
		}
		if t.Status.GetValue() == constants.StatusFailed {
			return false, fmt.Errorf("Task [%s] failed. ", taskId)
		}
		logger.Errorf("Unknown status [%s] for task [%s]. ", t.Status.GetValue(), taskId)
		return false, nil
	}, timeout, waitInterval)
}

func SendTask(task *models.Task) (taskId string, err error) {
	pbTask := models.TaskToPb(task)
	taskRequest := &pb.CreateTaskRequest{
		JobId:      pbTask.JobId,
		NodeId:     pbTask.NodeId,
		Target:     pbTask.Target,
		TaskAction: pbTask.TaskAction,
		Directive:  pbTask.Directive,
	}
	taskId, err = CreateTask(client.GetSystemUserContext(), taskRequest)
	if err != nil {
		logger.Errorf("Failed to create task [%s]: %+v", taskId, err)
	}
	return
}
