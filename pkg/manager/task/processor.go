// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package task

import (
	"context"

	clusterclient "openpitrix.io/openpitrix/pkg/client/cluster"
	"openpitrix.io/openpitrix/pkg/logger"
	"openpitrix.io/openpitrix/pkg/models"
	"openpitrix.io/openpitrix/pkg/pb"
	"openpitrix.io/openpitrix/pkg/plugins/vmbased"
	"openpitrix.io/openpitrix/pkg/utils"
)

type Processor struct {
	Task *models.Task
}

func NewProcessor(task *models.Task) *Processor {
	return &Processor{
		Task: task,
	}
}

// Post process when task is start
func (t *Processor) Pre() error {
	var err error
	ctx := context.Background()
	client, err := clusterclient.NewClusterManagerClient(ctx)
	if err != nil {
		logger.Errorf("Executing task [%s] post processor failed: %+v", t.Task.TaskId, err)
		return err
	}
	switch t.Task.TaskAction {
	case vmbased.ActionRunInstances:
		// volume created before instance, so need to change RunInstances task directive
		if t.Task.Directive == "" {
			logger.Warnf("Skip empty task [%p] directive", t.Task.TaskId)
		}
		instance, err := models.NewInstance(t.Task.Directive)
		if err == nil {
			clusterNodes, err := clusterclient.GetClusterNodes(ctx, client, []string{instance.NodeId})
			if err == nil {
				instance.VolumeId = clusterNodes[0].GetVolumeId().GetValue()
				t.Task.Directive, err = instance.ToString()
			}
		}
	case vmbased.ActionCreateVolumes:

	case vmbased.ActionWaitFrontgateAvailable:
	default:
		logger.Infof("Nothing to do with task [%s] pre processor", t.Task.TaskId)
	}
	if err != nil {
		logger.Errorf("Executing task [%s] pre processor failed: %+v", t.Task.TaskId, err)
	}
	return err
}

// Post process when task is done
func (t *Processor) Post() error {
	var err error
	ctx := context.Background()
	client, err := clusterclient.NewClusterManagerClient(ctx)
	if err != nil {
		logger.Errorf("Executing task [%s] post processor failed: %+v", t.Task.TaskId, err)
		return err
	}
	switch t.Task.TaskAction {
	case vmbased.ActionRunInstances:
		if t.Task.Directive == "" {
			logger.Warnf("Skip empty task [%p] directive", t.Task.TaskId)
		}
		instance, err := models.NewInstance(t.Task.Directive)
		if err != nil {
			_, err = client.ModifyClusterNode(ctx, &pb.ModifyClusterNodeRequest{
				NodeId:    utils.ToProtoString(instance.NodeId),
				PrivateIp: utils.ToProtoString(instance.PrivateIp),
			})
		}
	case vmbased.ActionCreateVolumes:
		if t.Task.Directive == "" {
			logger.Warnf("Skip empty task [%p] directive", t.Task.TaskId)
		}
		volume, err := models.NewVolume(t.Task.Directive)
		if err != nil {
			_, err = client.ModifyClusterNode(ctx, &pb.ModifyClusterNodeRequest{
				NodeId:   utils.ToProtoString(volume.NodeId),
				VolumeId: utils.ToProtoString(volume.VolumeId),
			})
		}
	default:
		logger.Infof("Nothing to do with task [%s] post processor", t.Task.TaskId)
	}
	if err != nil {
		logger.Errorf("Executing task [%s] post processor failed: %+v", t.Task.TaskId, err)
	}
	return err
}
