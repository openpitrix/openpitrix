// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package task

import (
	"context"
	"fmt"
	"time"

	"openpitrix.io/openpitrix/pkg/constants"
	"openpitrix.io/openpitrix/pkg/logger"
	"openpitrix.io/openpitrix/pkg/manager"
	"openpitrix.io/openpitrix/pkg/models"
	"openpitrix.io/openpitrix/pkg/pb"
	"openpitrix.io/openpitrix/pkg/util/funcutil"
)

type Client struct {
	pb.TaskManagerClient
}

func NewClient() (*Client, error) {
	conn, err := manager.NewClient(constants.TaskManagerHost, constants.TaskManagerPort)
	if err != nil {
		return nil, err
	}
	return &Client{
		TaskManagerClient: pb.NewTaskManagerClient(conn),
	}, nil
}

func (c *Client) WaitTask(ctx context.Context, taskId string, timeout time.Duration, waitInterval time.Duration) error {
	logger.Debug("Waiting for task [%s] finished", taskId)
	return funcutil.WaitForSpecificOrError(func() (bool, error) {
		taskRequest := &pb.DescribeTasksRequest{
			TaskId: []string{taskId},
		}
		taskResponse, err := c.DescribeTasks(ctx, taskRequest)
		if err != nil {
			//network or api error, not considered task fail.
			return false, nil
		}
		if len(taskResponse.TaskSet) == 0 {
			return false, fmt.Errorf("Can not find task [%s]. ", taskId)
		}
		t := taskResponse.TaskSet[0]
		if t.Status == nil {
			logger.Error("Task [%s] status is nil", taskId)
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
		logger.Error("Unknown status [%s] for task [%s]. ", t.Status.GetValue(), taskId)
		return false, nil
	}, timeout, waitInterval)
}

func (c *Client) SendTask(ctx context.Context, task *models.Task) (string, error) {
	pbTask := models.TaskToPb(task)
	taskRequest := &pb.CreateTaskRequest{
		JobId:          pbTask.JobId,
		NodeId:         pbTask.NodeId,
		Target:         pbTask.Target,
		TaskAction:     pbTask.TaskAction,
		Directive:      pbTask.Directive,
		FailureAllowed: pbTask.FailureAllowed,
		Status:         pbTask.Status,
	}
	response, err := c.CreateTask(ctx, taskRequest)
	taskId := response.GetTaskId().GetValue()
	if err != nil {
		logger.Error("Failed to create task [%s]: %+v", taskId, err)
		return "", err
	}
	return taskId, nil
}
