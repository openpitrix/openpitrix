// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package task

import (
	"openpitrix.io/openpitrix/pkg/client"
	clusterclient "openpitrix.io/openpitrix/pkg/client/cluster"
	"openpitrix.io/openpitrix/pkg/constants"
	"openpitrix.io/openpitrix/pkg/db"
	"openpitrix.io/openpitrix/pkg/logger"
	"openpitrix.io/openpitrix/pkg/models"
	"openpitrix.io/openpitrix/pkg/pb"
	"openpitrix.io/openpitrix/pkg/pb/metadata/types"
	"openpitrix.io/openpitrix/pkg/pi"
	"openpitrix.io/openpitrix/pkg/plugins/vmbased"
	"openpitrix.io/openpitrix/pkg/util/jsonutil"
	"openpitrix.io/openpitrix/pkg/util/pbutil"
)

type Processor struct {
	Task    *models.Task
	TLogger *logger.Logger
}

func NewProcessor(task *models.Task, tLogger *logger.Logger) *Processor {
	if tLogger == nil {
		tLogger = logger.NewLogger()
	}
	return &Processor{
		Task:    task,
		TLogger: tLogger,
	}
}

// Post process when task is start
func (p *Processor) Pre() error {
	if p.Task.Directive == "" {
		p.TLogger.Warn("Skip empty task [%s] directive", p.Task.TaskId)
		return nil
	}
	var err error
	ctx := client.GetSystemUserContext()
	clusterClient, err := clusterclient.NewClient()
	if err != nil {
		p.TLogger.Error("Executing task [%s] post processor failed: %+v", p.Task.TaskId, err)
		return err
	}

	oldDirective := p.Task.Directive
	switch p.Task.TaskAction {
	case vmbased.ActionRunInstances:
		// volume created before instance, so need to change RunInstances task directive
		instance, err := models.NewInstance(p.Task.Directive)
		if err != nil {
			return err
		}
		err = clusterClient.ModifyClusterNodeTransitionStatus(ctx, instance.NodeId, constants.StatusCreating)
		if err != nil {
			return err
		}
		clusterNodes, err := clusterClient.GetClusterNodes(ctx, []string{instance.NodeId})
		if err != nil {
			return err
		}
		instance.VolumeId = clusterNodes[0].GetVolumeId().GetValue()
		// write back
		p.Task.Directive = jsonutil.ToString(instance)

	case vmbased.ActionStartInstances:
		instance, err := models.NewInstance(p.Task.Directive)
		if err != nil {
			return err
		}
		err = clusterClient.ModifyClusterNodeTransitionStatus(ctx, instance.NodeId, constants.StatusStarting)
		if err != nil {
			return err
		}

	case vmbased.ActionStopInstances:
		instance, err := models.NewInstance(p.Task.Directive)
		if err != nil {
			return err
		}
		err = clusterClient.ModifyClusterNodeTransitionStatus(ctx, instance.NodeId, constants.StatusStopping)
		if err != nil {
			return err
		}

	case vmbased.ActionTerminateInstances:
		instance, err := models.NewInstance(p.Task.Directive)
		if err != nil {
			return err
		}
		err = clusterClient.ModifyClusterNodeTransitionStatus(ctx, instance.NodeId, constants.StatusDeleting)
		if err != nil {
			return err
		}

	case vmbased.ActionFormatAndMountVolume:
		meta, err := models.NewMeta(p.Task.Directive)
		if err != nil {
			return err
		}
		clusterNodes, err := clusterClient.GetClusterNodes(ctx, []string{meta.NodeId})
		if err != nil {
			return err
		}
		clusterNode := clusterNodes[0]
		clusterRole := clusterNode.GetClusterRole()
		cmd := vmbased.FormatAndMountVolumeCmd(
			clusterNode.GetDevice().GetValue(),
			clusterRole.GetMountPoint().GetValue(),
			clusterRole.GetFileSystem().GetValue(),
			clusterRole.GetMountOptions().GetValue(),
		)

		if meta.FrontgateId == "" {
			p.Task.TaskAction = vmbased.ActionRunCommandOnFrontgateNode
			request := &pbtypes.RunCommandOnFrontgateRequest{
				Endpoint: &pbtypes.FrontgateEndpoint{
					FrontgateId:     clusterNode.GetClusterId().GetValue(),
					FrontgateNodeId: clusterNode.GetNodeId().GetValue(),
					NodeIp:          clusterNode.GetPrivateIp().GetValue(),
					NodePort:        constants.FrontgateServicePort,
				},
				Command:        cmd,
				TimeoutSeconds: int32(meta.Timeout),
			}
			// write back
			p.Task.Directive = jsonutil.ToString(request)
		} else {
			p.Task.TaskAction = vmbased.ActionRunCommandOnDrone
			request := &pbtypes.RunCommandOnDroneRequest{
				Endpoint: &pbtypes.DroneEndpoint{
					FrontgateId: meta.FrontgateId,
					DroneIp:     clusterNode.GetPrivateIp().GetValue(),
					DronePort:   constants.DroneServicePort,
				},
				Command:        cmd,
				TimeoutSeconds: int32(meta.Timeout),
			}
			// write back
			p.Task.Directive = jsonutil.ToString(request)
		}

	case vmbased.ActionRunCommandOnDrone, vmbased.ActionRemoveContainerOnDrone:
		request := new(pbtypes.RunCommandOnDroneRequest)
		err := jsonutil.Decode([]byte(p.Task.Directive), request)
		if err != nil {
			return err
		}
		clusterNodes, err := clusterClient.GetClusterNodes(ctx, []string{p.Task.NodeId})
		if err != nil {
			return err
		}
		clusterNode := clusterNodes[0]
		request.Endpoint.DroneIp = clusterNode.GetPrivateIp().GetValue()

		// write back
		p.Task.Directive = jsonutil.ToString(request)

	case vmbased.ActionRunCommandOnFrontgateNode, vmbased.ActionRemoveContainerOnFrontgate:
		request := new(pbtypes.RunCommandOnFrontgateRequest)
		err := jsonutil.Decode([]byte(p.Task.Directive), request)
		if err != nil {
			return err
		}
		clusterNodes, err := clusterClient.GetClusterNodes(ctx, []string{p.Task.NodeId})
		if err != nil {
			return err
		}
		clusterNode := clusterNodes[0]
		request.Endpoint.NodeIp = clusterNode.GetPrivateIp().GetValue()

		// write back
		p.Task.Directive = jsonutil.ToString(request)

	case vmbased.ActionRegisterMetadata:
		meta, err := models.NewMeta(p.Task.Directive)
		if err != nil {
			return err
		}
		pbClusterWrappers, err := clusterClient.GetClusterWrappers(ctx, []string{meta.ClusterId})
		if err != nil {
			return err
		}
		metadata := &vmbased.MetadataV1{
			ClusterWrapper: pbClusterWrappers[0],
			Logger:         p.TLogger,
		}
		meta.Cnodes = jsonutil.ToString(metadata.GetClusterCnodes())

		// write back
		p.Task.Directive = jsonutil.ToString(meta)

	case vmbased.ActionRegisterCmd:
		// when CreateCluster need to reload ip
		meta, err := models.NewMeta(p.Task.Directive)
		if err != nil {
			return err
		}
		if meta.DroneIp == "" {
			clusterNodes, err := clusterClient.GetClusterNodes(ctx, []string{meta.NodeId})
			if err != nil {
				return err
			}
			meta.DroneIp = clusterNodes[0].GetPrivateIp().GetValue()
		}
		cnodes, err := models.NewCmdCnodes(meta.Cnodes)
		if err != nil {
			return err
		}
		err = cnodes.Format(meta.DroneIp, p.Task.TaskId)
		if err != nil {
			return err
		}
		meta.Cnodes = jsonutil.ToString(cnodes)
		// write back
		p.Task.Directive = jsonutil.ToString(meta)

	case vmbased.ActionDeregisterCmd:
		// when CreateCluster need to reload ip
		meta, err := models.NewMeta(p.Task.Directive)
		if err != nil {
			return err
		}
		if meta.DroneIp == "" {
			clusterNodes, err := clusterClient.GetClusterNodes(ctx, []string{meta.NodeId})
			if err != nil {
				return err
			}
			meta.DroneIp = clusterNodes[0].GetPrivateIp().GetValue()
		}
		meta.Cnodes = jsonutil.ToString(map[string]map[string]string{
			meta.DroneIp: {"cmd": ""},
		})

		// write back
		p.Task.Directive = jsonutil.ToString(meta)

	case vmbased.ActionStartConfd:
		// when CreateCluster need to reload ip
		meta, err := models.NewMeta(p.Task.Directive)
		if err != nil {
			return err
		}
		if meta.DroneIp == "" {
			clusterNodes, err := clusterClient.GetClusterNodes(ctx, []string{meta.NodeId})
			if err != nil {
				return err
			}
			meta.DroneIp = clusterNodes[0].GetPrivateIp().GetValue()

			// write back
			p.Task.Directive = jsonutil.ToString(meta)
		}

	case vmbased.ActionPingDrone:
		droneEndpoint := new(pbtypes.DroneEndpoint)
		err := jsonutil.Decode([]byte(p.Task.Directive), droneEndpoint)
		if err != nil {
			return err
		}
		if droneEndpoint.DroneIp == "" {
			clusterNodes, err := clusterClient.GetClusterNodes(ctx, []string{p.Task.NodeId})
			if err != nil {
				return err
			}
			droneEndpoint.DroneIp = clusterNodes[0].GetPrivateIp().GetValue()

			// write back
			p.Task.Directive = jsonutil.ToString(droneEndpoint)
		}

	case vmbased.ActionPingFrontgate:
		meta, err := models.NewMeta(p.Task.Directive)
		if err != nil {
			return err
		}
		directive := &pbtypes.FrontgateId{
			Id: meta.ClusterId,
		}
		p.Task.Directive = jsonutil.ToString(directive)

	case vmbased.ActionSetFrontgateConfig:
		meta, err := models.NewMeta(p.Task.Directive)
		if err != nil {
			return err
		}
		clusterId := meta.ClusterId
		pbClusterWrappers, err := clusterClient.GetClusterWrappers(ctx, []string{clusterId})
		if err != nil {
			return err
		}
		metadataConfig := &vmbased.MetadataConfig{
			ClusterWrapper: pbClusterWrappers[0],
		}
		p.Task.Directive = metadataConfig.GetFrontgateConfig(p.Task.NodeId)

	case vmbased.ActionSetDroneConfig:
		meta, err := models.NewMeta(p.Task.Directive)
		if err != nil {
			return err
		}
		clusterId := meta.ClusterId
		pbClusterWrappers, err := clusterClient.GetClusterWrappers(ctx, []string{clusterId})
		if err != nil {
			return err
		}
		metadataConfig := &vmbased.MetadataConfig{
			ClusterWrapper: pbClusterWrappers[0],
		}
		p.Task.Directive = metadataConfig.GetDroneConfig(p.Task.NodeId)

	default:
		p.TLogger.Debug("Nothing to do with task [%s] pre processor", p.Task.TaskId)
	}

	// update directive when changed
	if oldDirective != p.Task.Directive {
		attributes := map[string]interface{}{
			"directive": p.Task.Directive,
		}
		_, err := pi.Global().Db.
			Update(models.TaskTableName).
			SetMap(attributes).
			Where(db.Eq("task_id", p.Task.TaskId)).
			Exec()
		if err != nil {
			p.TLogger.Error("Failed to update task [%s]: %+v", p.Task.TaskId, err)
			return err
		}
	}
	return err
}

// Post process when task is done
func (t *Processor) Post() error {
	t.TLogger.Debug("Post task [%s] directive: %s", t.Task.TaskId, t.Task.Directive)
	var err error
	ctx := client.GetSystemUserContext()
	clusterClient, err := clusterclient.NewClient()
	if err != nil {
		t.TLogger.Error("Executing task [%s] post processor failed: %+v", t.Task.TaskId, err)
		return err
	}
	switch t.Task.TaskAction {
	case vmbased.ActionRunInstances:
		if t.Task.Directive == "" {
			t.TLogger.Warn("Skip empty task [%s] directive", t.Task.TaskId)
		}
		instance, err := models.NewInstance(t.Task.Directive)
		if err != nil {
			return err
		}
		_, err = clusterClient.ModifyClusterNode(ctx, &pb.ModifyClusterNodeRequest{
			ClusterNode: &pb.ClusterNode{
				NodeId:           pbutil.ToProtoString(instance.NodeId),
				InstanceId:       pbutil.ToProtoString(instance.InstanceId),
				Device:           pbutil.ToProtoString(instance.Device),
				PrivateIp:        pbutil.ToProtoString(instance.PrivateIp),
				TransitionStatus: pbutil.ToProtoString(""),
				Status:           pbutil.ToProtoString(constants.StatusActive),
			},
		})
		if err != nil {
			return err
		}

	case vmbased.ActionStartInstances:
		instance, err := models.NewInstance(t.Task.Directive)
		if err != nil {
			return err
		}
		err = clusterClient.ModifyClusterNodeTransitionStatus(ctx, instance.NodeId, "")
		if err != nil {
			return err
		}
		err = clusterClient.ModifyClusterNodeStatus(ctx, instance.NodeId, constants.StatusActive)
		if err != nil {
			return err
		}

	case vmbased.ActionStopInstances:
		instance, err := models.NewInstance(t.Task.Directive)
		if err != nil {
			return err
		}
		err = clusterClient.ModifyClusterNodeTransitionStatus(ctx, instance.NodeId, "")
		if err != nil {
			return err
		}
		err = clusterClient.ModifyClusterNodeStatus(ctx, instance.NodeId, constants.StatusStopped)
		if err != nil {
			return err
		}

	case vmbased.ActionTerminateInstances:
		instance, err := models.NewInstance(t.Task.Directive)
		if err != nil {
			return err
		}
		err = clusterClient.ModifyClusterNodeTransitionStatus(ctx, instance.NodeId, "")
		if err != nil {
			return err
		}
		err = clusterClient.ModifyClusterNodeStatus(ctx, instance.NodeId, constants.StatusDeleted)
		if err != nil {
			return err
		}

	case vmbased.ActionCreateVolumes:
		if t.Task.Directive == "" {
			t.TLogger.Warn("Skip empty task [%s] directive", t.Task.TaskId)
		}
		volume, err := models.NewVolume(t.Task.Directive)
		if err != nil {
			return err
		}
		_, err = clusterClient.ModifyClusterNode(ctx, &pb.ModifyClusterNodeRequest{
			ClusterNode: &pb.ClusterNode{
				NodeId:   pbutil.ToProtoString(t.Task.NodeId),
				VolumeId: pbutil.ToProtoString(volume.VolumeId),
			},
		})
		if err != nil {
			return err
		}

	default:
		t.TLogger.Debug("Nothing to do with task [%s] post processor", t.Task.TaskId)
	}
	return err
}
