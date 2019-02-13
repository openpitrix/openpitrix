// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package vmbased

import (
	"context"
	"fmt"

	clusterclient "openpitrix.io/openpitrix/pkg/client/cluster"
	"openpitrix.io/openpitrix/pkg/constants"
	"openpitrix.io/openpitrix/pkg/logger"
	"openpitrix.io/openpitrix/pkg/models"
	"openpitrix.io/openpitrix/pkg/util/funcutil"
	"openpitrix.io/openpitrix/pkg/util/jsonutil"
)

type FrameHandler struct {
}

func (f *FrameHandler) WaitFrontgateAvailable(ctx context.Context, task *models.Task) (*models.Task, error) {

	waitFrontgateDirective := new(models.Meta)

	if task.Directive == "" {
		logger.Warn(ctx, "Skip empty task [%s] directive", task.TaskId)
		return task, nil
	}
	err := jsonutil.Decode([]byte(task.Directive), waitFrontgateDirective)
	if err != nil {
		logger.Error(ctx, "Unmarshal into map failed: %+v", err)
		return task, err
	}

	frontgateId := waitFrontgateDirective.FrontgateId

	clusterClient, err := clusterclient.NewClient()
	if err != nil {
		return task, err
	}

	return task, funcutil.WaitForSpecificOrError(func() (bool, error) {
		response, err := clusterClient.GetClusters(ctx, []string{frontgateId})
		if err != nil {
			//network or api error, not considered task fail.
			return false, nil
		}
		frontgate := response[0]
		if frontgate.Status == nil {
			logger.Error(ctx, "Frontgate [%s] status is nil", frontgateId)
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
	}, task.GetTimeout(constants.WaitFrontgateServiceTimeout), constants.WaitFrontgateServiceInterval)
}
