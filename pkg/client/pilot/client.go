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
	"openpitrix.io/openpitrix/pkg/pb/metadata/pilot"
	"openpitrix.io/openpitrix/pkg/pb/metadata/types"
	"openpitrix.io/openpitrix/pkg/util/funcutil"
)

type Client struct {
	pbpilot.PilotServiceClient
}

func NewClient() (*Client, error) {
	conn, err := manager.NewClient(constants.PilotServiceHost, constants.PilotServicePort)
	if err != nil {
		return nil, err
	}
	return &Client{
		PilotServiceClient: pbpilot.NewPilotServiceClient(conn),
	}, nil
}

func (c *Client) WaitSubtask(ctx context.Context, taskId string, timeout time.Duration, waitInterval time.Duration) error {
	logger.Debug("Waiting for task [%s] finished", taskId)
	return funcutil.WaitForSpecificOrError(func() (bool, error) {
		taskStatusRequest := &pbtypes.SubTaskId{
			TaskId: taskId,
		}
		withTimeoutCtx, cancel := context.WithTimeout(ctx, constants.GrpcToPilotTimeout)
		defer cancel()
		taskStatusResponse, err := c.GetSubtaskStatus(withTimeoutCtx, taskStatusRequest)
		if err != nil {
			//network or api error, not considered task fail.
			return false, nil
		}

		t := taskStatusResponse
		if t.Status == constants.StatusWorking || t.Status == constants.StatusPending {
			return false, nil
		}
		if t.Status == constants.StatusSuccessful {
			return true, nil
		}
		if t.Status == constants.StatusFailed {
			return false, fmt.Errorf("Task [%s] failed. ", taskId)
		}
		logger.Error("Unknown status [%s] for task [%s]. ", t.Status, taskId)
		return false, nil
	}, timeout, waitInterval)
}
