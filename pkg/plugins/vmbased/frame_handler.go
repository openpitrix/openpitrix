// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package vmbased

import (
	"context"
	"encoding/json"
	"fmt"

	clusterclient "openpitrix.io/openpitrix/pkg/client/cluster"
	"openpitrix.io/openpitrix/pkg/constants"
	"openpitrix.io/openpitrix/pkg/logger"
	"openpitrix.io/openpitrix/pkg/models"
	"openpitrix.io/openpitrix/pkg/pb"
	"openpitrix.io/openpitrix/pkg/utils"
)

type FrameHandler struct {
}

func (f *FrameHandler) WaitFrontgateAvailable(task *models.Task) error {

	waitFrontgateDirective := make(map[string]interface{})

	if task.Directive == "" {
		logger.Warnf("Skip empty task [%s] directive", task.TaskId)
		return nil
	}
	err := json.Unmarshal([]byte(task.Directive), waitFrontgateDirective)
	if err != nil {
		logger.Errorf("Unmarshal into map failed: %+v", err)
		return err
	}

	frontgateId := waitFrontgateDirective["frontgate_id"].(string)

	ctx := context.Background()
	client, err := clusterclient.NewClusterManagerClient(ctx)
	if err != nil {
		return err
	}

	return utils.WaitForSpecificOrError(func() (bool, error) {
		response, err := client.DescribeClusters(ctx, &pb.DescribeClustersRequest{
			ClusterId: []string{frontgateId},
		})
		if err != nil {
			//network or api error, not considered task fail.
			return false, nil
		}
		if len(response.ClusterSet) == 0 {
			return false, fmt.Errorf("Can not find frontgate [%s]. ", frontgateId)
		}
		frontgate := response.ClusterSet[0]
		if frontgate.Status == nil {
			logger.Errorf("Frontgate [%s] status is nil", frontgateId)
			return false, nil
		}

		status := frontgate.Status.GetValue()
		transitionStatus := frontgate.TransitionStatus.GetValue()
		if transitionStatus != "" {
			return false, nil
		}
		if status == constants.StatusActive && transitionStatus == "" {
			return true, nil
		} else {
			return false, fmt.Errorf("Frontgate status is [%s]. ", status)
		}
	}, task.GetTimeout(constants.WaitTaskTimeout), constants.WaitTaskInterval)
}
