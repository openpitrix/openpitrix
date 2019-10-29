// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package task

import (
	"context"

	clusterclient "openpitrix.io/openpitrix/pkg/client/cluster"
	runtimeclient "openpitrix.io/openpitrix/pkg/client/runtime"
	"openpitrix.io/openpitrix/pkg/constants"
	"openpitrix.io/openpitrix/pkg/db"
	"openpitrix.io/openpitrix/pkg/logger"
	"openpitrix.io/openpitrix/pkg/models"
	"openpitrix.io/openpitrix/pkg/pb"
	pbtypes "openpitrix.io/openpitrix/pkg/pb/metadata/types"
	"openpitrix.io/openpitrix/pkg/pi"
	"openpitrix.io/openpitrix/pkg/plugins/vmbased"
	"openpitrix.io/openpitrix/pkg/util/jsonutil"
	"openpitrix.io/openpitrix/pkg/util/pbutil"
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
func (p *Processor) Pre(ctx context.Context) error {
	if p.Task.Directive == "" {
		logger.Warn(ctx, "Skip empty task [%s] directive", p.Task.TaskId)
		return nil
	}
	var err error
	clusterClient, err := clusterclient.NewClient()
	if err != nil {
		logger.Error(ctx, "Executing task [%s] post processor failed: %+v", p.Task.TaskId, err)
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
		clusterNode := clusterNodes[0]
		instance.VolumeId = clusterNode.GetVolumeId().GetValue()

		frontgateIp := ""
		// Get frontgate node ip
		clusters, err := clusterClient.GetClusters(ctx, []string{clusterNode.GetClusterId().GetValue()})
		if err != nil {
			return err
		}
		cluster := clusters[0]
		if cluster.GetClusterType().GetValue() == constants.NormalClusterType {
			frontgates, err := clusterClient.GetClusters(ctx, []string{cluster.GetFrontgateId().GetValue()})
			if err != nil {
				return err
			}
			for _, frontgateNode := range frontgates[0].ClusterNodeSet {
				frontgateIp = frontgateNode.GetPrivateIp().GetValue()
			}
		}
		instance.UserDataValue = vmbased.FormatUserData(instance.UserDataValue, frontgateIp)
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
		runtimeId := pbClusterWrappers[0].Cluster.RuntimeId
		runtime, err := runtimeclient.NewRuntime(ctx, runtimeId)
		if err != nil {
			logger.Error(ctx, "Get runtime [%s] failed: %+v", runtimeId, err)
			return err
		}
		metadata := &vmbased.Metadata{
			ClusterWrapper: pbClusterWrappers[0],
			RuntimeDetails: runtime,
		}
		meta.Cnodes = jsonutil.ToString(metadata.GetClusterCnodes(ctx))

		// write back
		p.Task.Directive = jsonutil.ToString(meta)

	case vmbased.ActionRegisterMetadataMapping:
		meta, err := models.NewMeta(p.Task.Directive)
		if err != nil {
			return err
		}
		pbClusterWrappers, err := clusterClient.GetClusterWrappers(ctx, []string{meta.ClusterId})
		if err != nil {
			return err
		}
		runtimeId := pbClusterWrappers[0].Cluster.RuntimeId
		runtime, err := runtimeclient.NewRuntime(ctx, runtimeId)
		if err != nil {
			logger.Error(ctx, "Get runtime [%s] failed: %+v", runtimeId, err)
			return err
		}
		metadata := &vmbased.Metadata{
			ClusterWrapper: pbClusterWrappers[0],
			RuntimeDetails: runtime,
		}
		meta.Cnodes = jsonutil.ToString(metadata.GetClusterMappingCnodes(ctx))

		// write back
		p.Task.Directive = jsonutil.ToString(meta)

	case vmbased.ActionDeregisterMetadataMapping:
		meta, err := models.NewMeta(p.Task.Directive)
		if err != nil {
			return err
		}
		pbClusterWrappers, err := clusterClient.GetClusterWrappers(ctx, []string{meta.ClusterId})
		if err != nil {
			return err
		}
		runtimeId := pbClusterWrappers[0].Cluster.RuntimeId
		runtime, err := runtimeclient.NewRuntime(ctx, runtimeId)
		if err != nil {
			logger.Error(ctx, "Get runtime [%s] failed: %+v", runtimeId, err)
			return err
		}
		metadata := &vmbased.Metadata{
			ClusterWrapper: pbClusterWrappers[0],
			RuntimeDetails: runtime,
		}
		meta.Cnodes = jsonutil.ToString(metadata.GetEmptyClusterMappingCnodes())

		// write back
		p.Task.Directive = jsonutil.ToString(meta)

	case vmbased.ActionRegisterNodesMetadata, vmbased.ActionRegisterEnvMetadata:
		p.Task.TaskAction = vmbased.ActionRegisterMetadata

	case vmbased.ActionRegisterNodesMetadataMapping:
		p.Task.TaskAction = vmbased.ActionRegisterMetadataMapping

	case vmbased.ActionRegisterCmd:
		meta, err := models.NewMeta(p.Task.Directive)
		if err != nil {
			return err
		}
		err = clusterClient.ModifyClusterNodeTransitionStatus(ctx, meta.NodeId, constants.StatusUpdating)
		if err != nil {
			return err
		}

		cmdCnodes, err := models.NewCmdCnodes(meta.Cnodes)
		if err != nil {
			return err
		}

		if meta.DroneIp == "" || cmdCnodes.InstanceId == "" {
			clusterNodes, err := clusterClient.GetClusterNodes(ctx, []string{meta.NodeId})
			if err != nil {
				return err
			}
			meta.DroneIp = clusterNodes[0].GetPrivateIp().GetValue()
			cmdCnodes.InstanceId = clusterNodes[0].GetInstanceId().GetValue()
		}

		cmdCnodes.Cmd.Id = p.Task.TaskId

		meta.Cnodes = jsonutil.ToString(cmdCnodes.Format())
		// write back
		p.Task.Directive = jsonutil.ToString(meta)

	case vmbased.ActionStartConfd:
		// when CreateCluster need to reload ip
		meta, err := models.NewMeta(p.Task.Directive)
		if err != nil {
			return err
		}
		err = clusterClient.ModifyClusterNodeTransitionStatus(ctx, meta.NodeId, constants.StatusUpdating)
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
		err = clusterClient.ModifyClusterNodeTransitionStatus(ctx, p.Task.NodeId, constants.StatusUpdating)
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

	case vmbased.ActionPingFrontgate, vmbased.PingMetadataBackend:
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
		logger.Debug(ctx, "Nothing to do with task [%s] pre processor", p.Task.TaskId)
	}

	// update directive when changed
	if oldDirective != p.Task.Directive {
		attributes := map[string]interface{}{
			constants.ColumnDirective: p.Task.Directive,
		}
		_, err := pi.Global().DB(ctx).
			Update(constants.TableTask).
			SetMap(attributes).
			Where(db.Eq(constants.ColumnTaskId, p.Task.TaskId)).
			Exec()
		if err != nil {
			logger.Error(ctx, "Failed to update task [%s]: %+v", p.Task.TaskId, err)
			return err
		}
	}
	return err
}

// Post process when task is done
func (p *Processor) Post(ctx context.Context) error {
	var err error
	clusterClient, err := clusterclient.NewClient()
	if err != nil {
		logger.Error(ctx, "Executing task [%s] post processor failed: %+v", p.Task.TaskId, err)
		return err
	}
	switch p.Task.TaskAction {
	case vmbased.ActionRunInstances:
		if p.Task.Directive == "" {
			logger.Warn(ctx, "Skip empty task [%s] directive", p.Task.TaskId)
		}
		instance, err := models.NewInstance(p.Task.Directive)
		if err != nil {
			return err
		}
		_, err = clusterClient.ModifyClusterNode(ctx, &pb.ModifyClusterNodeRequest{
			ClusterNode: &pb.ClusterNode{
				NodeId:     pbutil.ToProtoString(instance.NodeId),
				InstanceId: pbutil.ToProtoString(instance.InstanceId),
				Device:     pbutil.ToProtoString(instance.Device),
				PrivateIp:  pbutil.ToProtoString(instance.PrivateIp),
				Eip:        pbutil.ToProtoString(instance.Eip),
			},
		})
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

	case vmbased.ActionStartInstances:
		instance, err := models.NewInstance(p.Task.Directive)
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
		instance, err := models.NewInstance(p.Task.Directive)
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
		instance, err := models.NewInstance(p.Task.Directive)
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
		if p.Task.Directive == "" {
			logger.Warn(ctx, "Skip empty task [%s] directive", p.Task.TaskId)
		}
		volume, err := models.NewVolume(p.Task.Directive)
		if err != nil {
			return err
		}
		_, err = clusterClient.ModifyClusterNode(ctx, &pb.ModifyClusterNodeRequest{
			ClusterNode: &pb.ClusterNode{
				NodeId:   pbutil.ToProtoString(p.Task.NodeId),
				VolumeId: pbutil.ToProtoString(volume.VolumeId),
			},
		})
		if err != nil {
			return err
		}

	case vmbased.ActionPingDrone, vmbased.ActionRegisterCmd, vmbased.ActionStartConfd:
		err = clusterClient.ModifyClusterNodeTransitionStatus(ctx, p.Task.NodeId, "")
		if err != nil {
			return err
		}
	default:
		logger.Debug(ctx, "Nothing to do with task [%s] post processor", p.Task.TaskId)
	}
	return err
}
