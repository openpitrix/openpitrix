// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package vmbased

import (
	"fmt"

	"openpitrix.io/openpitrix/pkg/client"
	clusterclient "openpitrix.io/openpitrix/pkg/client/cluster"
	"openpitrix.io/openpitrix/pkg/constants"
	"openpitrix.io/openpitrix/pkg/logger"
	"openpitrix.io/openpitrix/pkg/models"
	"openpitrix.io/openpitrix/pkg/pb"
	"openpitrix.io/openpitrix/pkg/util/funcutil"
	"openpitrix.io/openpitrix/pkg/util/jsonutil"
)

type FrameHandler struct {
	Logger *logger.Logger
}

func (f *FrameHandler) WaitFrontgateAvailable(task *models.Task) error {

	waitFrontgateDirective := new(models.Meta)

	if task.Directive == "" {
		f.Logger.Warn("Skip empty task [%s] directive", task.TaskId)
		return nil
	}
	err := jsonutil.Decode([]byte(task.Directive), waitFrontgateDirective)
	if err != nil {
		f.Logger.Error("Unmarshal into map failed: %+v", err)
		return err
	}

	frontgateId := waitFrontgateDirective.FrontgateId

	ctx := client.GetSystemUserContext()
	clusterClient, err := clusterclient.NewClient()
	if err != nil {
		return err
	}

	return funcutil.WaitForSpecificOrError(func() (bool, error) {
		response, err := clusterClient.DescribeClusters(ctx, &pb.DescribeClustersRequest{
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
			f.Logger.Error("Frontgate [%s] status is nil", frontgateId)
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
