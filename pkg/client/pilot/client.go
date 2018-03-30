// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package pilot

import (
	"context"
	"fmt"
	"time"

	"openpitrix.io/openpitrix/pkg/constants"
	"openpitrix.io/openpitrix/pkg/logger"
	"openpitrix.io/openpitrix/pkg/manager"
	"openpitrix.io/openpitrix/pkg/pb"
	"openpitrix.io/openpitrix/pkg/utils"
)

func NewPilotManagerClient(ctx context.Context) (pb.PilotServiceClient, error) {
	conn, err := manager.NewClient(ctx, constants.PilotManagerHost, constants.PilotManagerPort)
	if err != nil {
		return nil, err
	}
	return pb.NewPilotServiceClient(conn), err
}

func HandleSubtask(subtaskRequest *pb.HandleSubtaskRequest) error {
	ctx := context.Background()
	client, err := NewPilotManagerClient(ctx)
	if err != nil {
		return err
	}
	_, err = client.HandleSubtask(ctx, subtaskRequest)
	if err != nil {
		return err
	}
	return nil
}

func GetSubtaskStatus(subtaskStatusRequest *pb.GetSubtaskStatusRequest) (*pb.GetSubtaskStatusResponse, error) {
	ctx := context.Background()
	client, err := NewPilotManagerClient(ctx)
	if err != nil {
		return nil, err
	}
	subtaskStatusResponse, err := client.GetSubtaskStatus(ctx, subtaskStatusRequest)
	if err != nil {
		return nil, err
	}
	return subtaskStatusResponse, err
}

func WaitSubtask(taskId string, timeout time.Duration, waitInterval time.Duration) error {
	logger.Debug("Waiting for task [%s] finished", taskId)
	return utils.WaitForSpecificOrError(func() (bool, error) {
		taskStatusRequest := &pb.GetSubtaskStatusRequest{
			SubtaskId: []string{taskId},
		}
		taskStatusResponse, err := GetSubtaskStatus(taskStatusRequest)
		if err != nil {
			//network or api error, not considered task fail.
			return false, nil
		}
		if len(taskStatusResponse.SubtaskStatusSet) == 0 {
			return false, fmt.Errorf("Can not find task [%s]. ", taskId)
		}
		t := taskStatusResponse.SubtaskStatusSet[0]
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
