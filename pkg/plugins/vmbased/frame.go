// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package vmbased

import (
	"context"
	"encoding/base64"
	"fmt"
	"sort"
	"strings"

	clusterclient "openpitrix.io/openpitrix/pkg/client/cluster"
	"openpitrix.io/openpitrix/pkg/config"
	"openpitrix.io/openpitrix/pkg/constants"
	"openpitrix.io/openpitrix/pkg/devkit/opapp"
	"openpitrix.io/openpitrix/pkg/logger"
	"openpitrix.io/openpitrix/pkg/models"
	"openpitrix.io/openpitrix/pkg/pb"
	pbtypes "openpitrix.io/openpitrix/pkg/pb/metadata/types"
	"openpitrix.io/openpitrix/pkg/pi"
	"openpitrix.io/openpitrix/pkg/util/jsonutil"
	"openpitrix.io/openpitrix/pkg/util/pbutil"
	"openpitrix.io/openpitrix/pkg/util/sshutil"
)

type Frame struct {
	Ctx                   context.Context
	Job                   *models.Job
	ClusterWrapper        *models.ClusterWrapper
	Runtime               *models.RuntimeDetails
	RuntimeProviderConfig *config.RuntimeProviderConfig
}

func (f *Frame) startConfdServiceLayer(nodeIds []string, failureAllowed bool) *models.TaskLayer {
	taskLayer := new(models.TaskLayer)
	for _, nodeId := range nodeIds {
		clusterNode := f.ClusterWrapper.ClusterNodesWithKeyPairs[nodeId]
		meta := &models.Meta{
			FrontgateId: f.ClusterWrapper.Cluster.FrontgateId,
			Timeout:     TimeoutStartConfd,
			NodeId:      clusterNode.NodeId,
			DroneIp:     clusterNode.PrivateIp,
		}
		directive := jsonutil.ToString(meta)
		startConfdTask := &models.Task{
			JobId:          f.Job.JobId,
			Owner:          f.Job.Owner,
			TaskAction:     ActionStartConfd,
			Target:         constants.TargetPilot,
			NodeId:         nodeId,
			Directive:      directive,
			FailureAllowed: failureAllowed,
		}
		taskLayer.Tasks = append(taskLayer.Tasks, startConfdTask)
	}
	if len(taskLayer.Tasks) > 0 {
		return taskLayer
	} else {
		return nil
	}
}

func (f *Frame) stopConfdServiceLayer(nodeIds []string, failtureAllowed bool) *models.TaskLayer {
	taskLayer := new(models.TaskLayer)
	for _, nodeId := range nodeIds {
		clusterNode := f.ClusterWrapper.ClusterNodesWithKeyPairs[nodeId]
		meta := &models.Meta{
			FrontgateId: f.ClusterWrapper.Cluster.FrontgateId,
			Timeout:     TimeoutStopConfd,
			NodeId:      clusterNode.NodeId,
			DroneIp:     clusterNode.PrivateIp,
		}
		directive := jsonutil.ToString(meta)
		stopConfdTask := &models.Task{
			JobId:          f.Job.JobId,
			Owner:          f.Job.Owner,
			TaskAction:     ActionStopConfd,
			Target:         constants.TargetPilot,
			NodeId:         nodeId,
			Directive:      directive,
			FailureAllowed: failtureAllowed,
		}
		taskLayer.Tasks = append(taskLayer.Tasks, stopConfdTask)
	}
	if len(taskLayer.Tasks) > 0 {
		return taskLayer
	} else {
		return nil
	}
}

// Put the nodes into two groups
func (f *Frame) getPreAndPostStartGroupNodes(nodeIds []string) ([]string, []string) {
	var preGroupNodes, postGroupNodes []string
	for _, nodeId := range nodeIds {
		role := f.ClusterWrapper.ClusterNodesWithKeyPairs[nodeId].Role
		serviceStr := f.ClusterWrapper.ClusterCommons[role].InitService
		if serviceStr != "" {
			service := opapp.Service{}
			err := jsonutil.Decode([]byte(serviceStr), &service)
			if err != nil {
				logger.Error(f.Ctx, "Unmarshal cluster [%s] init service failed: %+v",
					f.ClusterWrapper.Cluster.ClusterId, err)
				return nil, nil
			}
			postStartService := false
			if service.PostStartService != nil {
				postStartService = *service.PostStartService
			}
			if postStartService {
				postGroupNodes = append(postGroupNodes, nodeId)
			} else {
				preGroupNodes = append(preGroupNodes, nodeId)
			}
		}
	}
	return preGroupNodes, postGroupNodes
}

// Put the nodes into two groups
func (f *Frame) getPreAndPostStopGroupNodes(nodeIds []string) ([]string, []string) {
	var preGroupNodes, postGroupNodes []string
	for _, nodeId := range nodeIds {
		role := f.ClusterWrapper.ClusterNodesWithKeyPairs[nodeId].Role
		serviceStr := f.ClusterWrapper.ClusterCommons[role].DestroyService
		if serviceStr != "" {
			service := opapp.Service{}
			err := jsonutil.Decode([]byte(serviceStr), &service)
			if err != nil {
				logger.Error(f.Ctx, "Unmarshal cluster [%s] init service failed: %+v",
					f.ClusterWrapper.Cluster.ClusterId, err)
				return nil, nil
			}
			postStopService := false
			if service.PostStopService != nil {
				postStopService = *service.PostStopService
			}
			if postStopService {
				postGroupNodes = append(postGroupNodes, nodeId)
			} else {
				preGroupNodes = append(preGroupNodes, nodeId)
			}
		}
	}
	return preGroupNodes, postGroupNodes
}

func (f *Frame) deregisterCmdLayer(nodeIds []string, failureAllowed bool) *models.TaskLayer {
	metadata := &Metadata{
		ClusterWrapper: f.ClusterWrapper,
		RuntimeDetails: f.Runtime,
	}
	taskLayer := new(models.TaskLayer)
	for _, nodeId := range nodeIds {
		ip := f.ClusterWrapper.ClusterNodesWithKeyPairs[nodeId].PrivateIp
		meta := &models.Meta{
			FrontgateId: f.ClusterWrapper.Cluster.FrontgateId,
			NodeId:      nodeId,
			DroneIp:     ip,
			Timeout:     TimeoutDeregister,
			Cnodes:      jsonutil.ToString(metadata.GetCmdCnodes(nodeId, nil).Format()),
		}
		directive := jsonutil.ToString(meta)
		deregisterCmdTask := &models.Task{
			JobId:          f.Job.JobId,
			Owner:          f.Job.Owner,
			TaskAction:     ActionDeregisterCmd,
			Target:         constants.TargetPilot,
			NodeId:         nodeId,
			Directive:      directive,
			FailureAllowed: failureAllowed,
		}
		taskLayer.Tasks = append(taskLayer.Tasks, deregisterCmdTask)
	}

	if len(taskLayer.Tasks) > 0 {
		return taskLayer
	} else {
		return nil
	}
}

func (f *Frame) registerCmdLayer(nodeIds []string, serviceName string, failureAllowed bool) *models.TaskLayer {
	metadata := &Metadata{
		ClusterWrapper: f.ClusterWrapper,
		RuntimeDetails: f.Runtime,
	}
	taskLayer := new(models.TaskLayer)
	for _, nodeId := range nodeIds {
		role := f.ClusterWrapper.ClusterNodesWithKeyPairs[nodeId].Role
		serviceStr := f.ClusterWrapper.GetCommonAttribute(role, serviceName)
		if serviceStr != nil {
			service := opapp.Service{}
			err := jsonutil.Decode([]byte(serviceStr.(string)), &service)
			if err != nil {
				logger.Error(f.Ctx, "Unmarshal cluster [%s] service [%s] failed: %+v",
					f.ClusterWrapper.Cluster.ClusterId, serviceName, err)
				return nil
			}
			timeout := constants.DefaultServiceTimeout
			if service.Timeout != nil {
				timeout = int(*service.Timeout)
			}
			if service.Cmd == "" {
				continue
			}
			ip := f.ClusterWrapper.ClusterNodesWithKeyPairs[nodeId].PrivateIp
			cmd := &models.Cmd{
				Cmd:     service.Cmd,
				Timeout: timeout,
			}
			meta := &models.Meta{
				FrontgateId: f.ClusterWrapper.Cluster.FrontgateId,
				NodeId:      nodeId,
				DroneIp:     ip,
				Timeout:     timeout,
				Cnodes:      jsonutil.ToString(metadata.GetCmdCnodes(nodeId, cmd)),
			}
			directive := jsonutil.ToString(meta)
			registerCmdTask := &models.Task{
				JobId:          f.Job.JobId,
				Owner:          f.Job.Owner,
				TaskAction:     ActionRegisterCmd,
				Target:         constants.TargetPilot,
				NodeId:         nodeId,
				Directive:      directive,
				FailureAllowed: failureAllowed,
			}
			taskLayer.Tasks = append(taskLayer.Tasks, registerCmdTask)
		}
	}

	if len(taskLayer.Tasks) > 0 {
		return taskLayer
	} else {
		return nil
	}
}

func (f *Frame) attachKeyPairLayer(nodeKeyPairDetail *models.NodeKeyPairDetail) *models.TaskLayer {
	taskLayer := new(models.TaskLayer)
	clusterNode := nodeKeyPairDetail.ClusterNode
	request := &pbtypes.RunCommandOnDroneRequest{
		Endpoint: &pbtypes.DroneEndpoint{
			FrontgateId: f.ClusterWrapper.Cluster.FrontgateId,
			DroneIp:     clusterNode.PrivateIp,
			DronePort:   constants.DroneServicePort,
		},
		Command:        fmt.Sprintf("%s \"%s\"", HostCmdPrefix, sshutil.DoAttachCmd(nodeKeyPairDetail.KeyPair.PubKey)),
		TimeoutSeconds: TimeoutKeyPair,
	}
	directive := jsonutil.ToString(request)
	attachKeyPairTask := &models.Task{
		JobId:          f.Job.JobId,
		Owner:          f.Job.Owner,
		TaskAction:     ActionRunCommandOnDrone,
		Target:         constants.TargetPilot,
		NodeId:         clusterNode.NodeId,
		Directive:      directive,
		FailureAllowed: false,
	}
	taskLayer.Tasks = append(taskLayer.Tasks, attachKeyPairTask)
	if len(taskLayer.Tasks) > 0 {
		return taskLayer
	} else {
		return nil
	}
}

func (f *Frame) detachKeyPairLayer(nodeKeyPairDetail *models.NodeKeyPairDetail) *models.TaskLayer {
	taskLayer := new(models.TaskLayer)
	clusterNode := nodeKeyPairDetail.ClusterNode
	request := &pbtypes.RunCommandOnDroneRequest{
		Endpoint: &pbtypes.DroneEndpoint{
			FrontgateId: f.ClusterWrapper.Cluster.FrontgateId,
			DroneIp:     clusterNode.PrivateIp,
			DronePort:   constants.DroneServicePort,
		},
		Command:        fmt.Sprintf("%s \"%s\"", HostCmdPrefix, sshutil.DoDetachCmd(nodeKeyPairDetail.KeyPair.PubKey)),
		TimeoutSeconds: TimeoutKeyPair,
	}
	directive := jsonutil.ToString(request)
	attachKeyPairTask := &models.Task{
		JobId:          f.Job.JobId,
		Owner:          f.Job.Owner,
		TaskAction:     ActionRunCommandOnDrone,
		Target:         constants.TargetPilot,
		NodeId:         clusterNode.NodeId,
		Directive:      directive,
		FailureAllowed: false,
	}
	taskLayer.Tasks = append(taskLayer.Tasks, attachKeyPairTask)
	if len(taskLayer.Tasks) > 0 {
		return taskLayer
	} else {
		return nil
	}
}

func (f *Frame) constructServiceTasks(serviceName, cmdName string, nodeIds []string,
	serviceParams map[string]interface{}, failureAllowed bool) *models.TaskLayer {
	headTaskLayer := new(models.TaskLayer)
	if len(nodeIds) == 0 {
		return nil
	}

	roleNodeIds := make(map[string][]string)
	nodeIdRole := make(map[string]string)
	for _, nodeId := range nodeIds {
		clusterNode, exist := f.ClusterWrapper.ClusterNodesWithKeyPairs[nodeId]
		if !exist {
			logger.Error(f.Ctx, "ClusterConf [%s] node [%s] not exist", f.ClusterWrapper.Cluster.ClusterId, nodeId)
			continue
		}
		role := clusterNode.Role
		service := f.ClusterWrapper.GetCommonAttribute(role, serviceName)
		if service == nil {
			continue
		}

		agentInstalled := f.ClusterWrapper.GetCommonAttribute(role, "AgentInstalled")
		if agentInstalled == nil {
			continue
		}

		if service.(string) == "" || !agentInstalled.(bool) {
			continue
		}
		roleNodeIds[role] = append(roleNodeIds[role], nodeId)
		nodeIdRole[nodeId] = role
	}

	filterNodes := make(map[string]string)
	roleService := make(map[string]opapp.Service)
	for role, nodes := range roleNodeIds {
		serviceStr := f.ClusterWrapper.GetCommonAttribute(role, serviceName)
		if serviceStr == nil {
			return nil
		}
		service := opapp.Service{}
		err := jsonutil.Decode([]byte(serviceStr.(string)), &service)
		if err != nil {
			logger.Error(f.Ctx, "Unmarshal cluster [%s] service [%s] failed: %+v",
				f.ClusterWrapper.Cluster.ClusterId, serviceName, err)
			return nil
		}
		roleService[role] = service
		execNodeNums := len(nodes)
		if service.NodesToExecuteOn != nil {
			execNodeNums = int(*service.NodesToExecuteOn)
		}
		if execNodeNums < len(nodes) && strings.HasSuffix(role, constants.ReplicaRoleSuffix) {
			// when the given nodes_to_execute_on is less than the length of the nodes, then ignore the replicas
			for _, nodeId := range nodes {
				filterNodes[nodeId] = ""
			}
			continue
		}
		num := execNodeNums
		for num < len(nodes) {
			filterNodes[nodes[num-1]] = ""
			num++
		}
	}

	orderNodeIds := make(map[int][]string)
	for nodeId, role := range nodeIdRole {
		_, exist := filterNodes[nodeId]
		if exist {
			continue
		}
		service := roleService[role]
		order := 0
		if service.Order != nil {
			order = int(*service.Order)
		}
		orderNodeIds[order] = append(orderNodeIds[order], nodeId)
	}

	var orders []int
	for order := range orderNodeIds {
		orders = append(orders, order)
	}

	sort.Ints(orders)

	for _, order := range orders {
		nodeIds := orderNodeIds[order]
		taskLayer := f.registerCmdLayer(nodeIds, serviceName, failureAllowed)
		headTaskLayer.Leaf().Child = taskLayer
	}
	return headTaskLayer.Child
}

func (f *Frame) initServiceLayer(nodeIds []string, failureAllowed bool) *models.TaskLayer {
	return f.constructServiceTasks("InitService", constants.ServiceCmdName, nodeIds, nil, failureAllowed)
}

func (f *Frame) startServiceLayer(nodeIds []string, failureAllowed bool) *models.TaskLayer {
	return f.constructServiceTasks("StartService", constants.ServiceCmdName, nodeIds, nil, failureAllowed)
}

func (f *Frame) stopServiceLayer(nodeIds []string, failureAllowed bool) *models.TaskLayer {
	return f.constructServiceTasks("StopService", constants.ServiceCmdName, nodeIds, nil, failureAllowed)
}

func (f *Frame) scaleOutPreCheckServiceLayer(nodeIds []string, failureAllowed bool) *models.TaskLayer {
	return f.constructServiceTasks("ScaleOutService", constants.ServiceCmdName, nodeIds, nil, failureAllowed)
}

func (f *Frame) scaleInPreCheckServiceLayer(nodeIds []string, failureAllowed bool) *models.TaskLayer {
	return f.constructServiceTasks("ScaleInService", constants.ServiceCmdName, nodeIds, nil, failureAllowed)
}

func (f *Frame) scaleOutServiceLayer(nodeIds []string, failureAllowed bool) *models.TaskLayer {
	return f.constructServiceTasks("ScaleOutService", constants.ServicePreCheckName, nodeIds, nil, failureAllowed)
}

func (f *Frame) scaleInServiceLayer(nodeIds []string, failureAllowed bool) *models.TaskLayer {
	return f.constructServiceTasks("ScaleInService", constants.ServicePreCheckName, nodeIds, nil, failureAllowed)
}

func (f *Frame) destroyServiceLayer(nodeIds []string, failureAllowed bool) *models.TaskLayer {
	return f.constructServiceTasks("DestroyService", constants.ServiceCmdName, nodeIds, nil, failureAllowed)
}

func (f *Frame) initAndStartServiceLayer(nodeIds []string, failureAllowed bool) *models.TaskLayer {
	headTaskLayer := new(models.TaskLayer)

	preStartNodes, postStartNodes := f.getPreAndPostStartGroupNodes(nodeIds)

	// Init service before start service
	headTaskLayer.Leaf().Child = f.initServiceLayer(preStartNodes, failureAllowed)

	// TODO: custom metadata
	headTaskLayer.Leaf().Child = f.startServiceLayer(nodeIds, failureAllowed)

	// Init service after start service
	headTaskLayer.Leaf().Child = f.initServiceLayer(postStartNodes, failureAllowed)

	return headTaskLayer.Child
}

func (f *Frame) destroyAndStopServiceLayer(nodeIds []string, extraLayer *models.TaskLayer, failureAllowed bool) *models.TaskLayer {
	headTaskLayer := new(models.TaskLayer)

	preStopNodes, postStopNodes := f.getPreAndPostStopGroupNodes(nodeIds)

	// Destroy service before stop service
	headTaskLayer.Leaf().Child = f.destroyServiceLayer(preStopNodes, failureAllowed)

	if extraLayer != nil {
		headTaskLayer.Leaf().Child = extraLayer
	}

	headTaskLayer.Leaf().Child = f.stopServiceLayer(nodeIds, failureAllowed)

	// Destroy service after stop service
	headTaskLayer.Leaf().Child = f.destroyServiceLayer(postStopNodes, failureAllowed)

	return headTaskLayer.Child
}

func (f *Frame) createVolumesLayer(nodeIds []string, failureAllowed bool) *models.TaskLayer {
	taskLayer := new(models.TaskLayer)
	for _, nodeId := range nodeIds {
		clusterNode := f.ClusterWrapper.ClusterNodesWithKeyPairs[nodeId]
		role := clusterNode.Role
		if strings.HasSuffix(role, constants.ReplicaRoleSuffix) {
			role = string([]byte(role)[:len(role)-len(constants.ReplicaRoleSuffix)])
		}
		clusterRole, exist := f.ClusterWrapper.ClusterRoles[role]
		if !exist {
			logger.Error(f.Ctx, "No such role [%s] in cluster role [%s]. ",
				role, f.ClusterWrapper.Cluster.ClusterId)
			return nil
		}

		size := clusterRole.StorageSize
		if size > 0 {
			mountPoints := strings.Split(clusterRole.MountPoint, ",")
			eachSize := int(size) / len(mountPoints)

			volume := &models.Volume{
				Name:      clusterNode.ClusterId + "_" + nodeId,
				Size:      eachSize,
				Zone:      f.ClusterWrapper.Cluster.Zone,
				RuntimeId: f.Runtime.RuntimeId,
			}
			directive := jsonutil.ToString(volume)
			createVolumesTask := &models.Task{
				JobId:          f.Job.JobId,
				Owner:          f.Job.Owner,
				TaskAction:     ActionCreateVolumes,
				Target:         f.Runtime.RuntimeId,
				NodeId:         nodeId,
				Directive:      directive,
				FailureAllowed: failureAllowed,
			}
			for range mountPoints {
				taskLayer.Tasks = append(taskLayer.Tasks, createVolumesTask)
			}
		}
	}
	return taskLayer
}

func (f *Frame) detachVolumesLayer(nodeIds []string, failureAllowed bool) *models.TaskLayer {
	taskLayer := new(models.TaskLayer)
	for _, nodeId := range nodeIds {
		clusterNode := f.ClusterWrapper.ClusterNodesWithKeyPairs[nodeId]
		if clusterNode.VolumeId == "" {
			continue
		}
		volume := &models.Volume{
			Name:       clusterNode.ClusterId + "_" + nodeId,
			Zone:       f.ClusterWrapper.Cluster.Zone,
			RuntimeId:  f.Runtime.RuntimeId,
			VolumeId:   clusterNode.VolumeId,
			InstanceId: clusterNode.InstanceId,
		}
		directive := jsonutil.ToString(volume)
		detachVolumesTask := &models.Task{
			JobId:          f.Job.JobId,
			Owner:          f.Job.Owner,
			TaskAction:     ActionDetachVolumes,
			Target:         f.Runtime.RuntimeId,
			NodeId:         nodeId,
			Directive:      directive,
			FailureAllowed: failureAllowed,
		}
		taskLayer.Tasks = append(taskLayer.Tasks, detachVolumesTask)
	}
	return taskLayer
}

func (f *Frame) attachVolumesLayer(nodeIds []string, failureAllowed bool) *models.TaskLayer {
	taskLayer := new(models.TaskLayer)
	for _, nodeId := range nodeIds {
		clusterNode := f.ClusterWrapper.ClusterNodesWithKeyPairs[nodeId]
		if clusterNode.VolumeId == "" {
			continue
		}
		volume := &models.Volume{
			Name:       clusterNode.ClusterId + "_" + nodeId,
			Zone:       f.ClusterWrapper.Cluster.Zone,
			RuntimeId:  f.Runtime.RuntimeId,
			VolumeId:   clusterNode.VolumeId,
			InstanceId: clusterNode.InstanceId,
		}
		directive := jsonutil.ToString(volume)
		attachVolumesTask := &models.Task{
			JobId:          f.Job.JobId,
			Owner:          f.Job.Owner,
			TaskAction:     ActionAttachVolumes,
			Target:         f.Runtime.RuntimeId,
			NodeId:         nodeId,
			Directive:      directive,
			FailureAllowed: failureAllowed,
		}
		taskLayer.Tasks = append(taskLayer.Tasks, attachVolumesTask)
	}
	return taskLayer
}

func (f *Frame) deleteVolumesLayer(nodeIds []string, failureAllowed bool) *models.TaskLayer {
	taskLayer := new(models.TaskLayer)
	for _, nodeId := range nodeIds {
		clusterNode := f.ClusterWrapper.ClusterNodesWithKeyPairs[nodeId]
		if clusterNode.VolumeId == "" {
			continue
		}
		volume := &models.Volume{
			Name:       clusterNode.ClusterId + "_" + nodeId,
			Zone:       f.ClusterWrapper.Cluster.Zone,
			RuntimeId:  f.Runtime.RuntimeId,
			VolumeId:   clusterNode.VolumeId,
			InstanceId: clusterNode.InstanceId,
		}
		directive := jsonutil.ToString(volume)
		deleteVolumesTask := &models.Task{
			JobId:          f.Job.JobId,
			Owner:          f.Job.Owner,
			TaskAction:     ActionDeleteVolumes,
			Target:         f.Runtime.RuntimeId,
			NodeId:         nodeId,
			Directive:      directive,
			FailureAllowed: failureAllowed,
		}
		taskLayer.Tasks = append(taskLayer.Tasks, deleteVolumesTask)
	}
	return taskLayer
}

func (f *Frame) resizeVolumesLayer(nodeIds []string, roleResizeResource *models.RoleResizeResource, failureAllowed bool) *models.TaskLayer {
	taskLayer := new(models.TaskLayer)
	for _, nodeId := range nodeIds {
		clusterNode := f.ClusterWrapper.ClusterNodesWithKeyPairs[nodeId]
		if clusterNode.VolumeId == "" {
			continue
		}
		clusterRole := f.ClusterWrapper.ClusterRoles[clusterNode.Role]

		if !roleResizeResource.StorageSize {
			logger.Debug(f.Ctx, "No need to resize node [%s] volume", nodeId)
			continue
		}
		volume := &models.Volume{
			Name:       clusterNode.ClusterId + "_" + nodeId,
			Zone:       f.ClusterWrapper.Cluster.Zone,
			RuntimeId:  f.Runtime.RuntimeId,
			VolumeId:   clusterNode.VolumeId,
			InstanceId: clusterNode.InstanceId,
			Size:       int(clusterRole.StorageSize),
		}
		directive := jsonutil.ToString(volume)
		deleteVolumesTask := &models.Task{
			JobId:          f.Job.JobId,
			Owner:          f.Job.Owner,
			TaskAction:     ActionResizeVolumes,
			Target:         f.Runtime.RuntimeId,
			NodeId:         nodeId,
			Directive:      directive,
			FailureAllowed: failureAllowed,
		}
		taskLayer.Tasks = append(taskLayer.Tasks, deleteVolumesTask)
	}
	return taskLayer
}

func (f *Frame) formatAndMountVolumeLayer(nodeIds []string, failureAllowed bool) *models.TaskLayer {
	taskLayer := new(models.TaskLayer)

	for _, nodeId := range nodeIds {
		clusterNode := f.ClusterWrapper.ClusterNodesWithKeyPairs[nodeId]
		role := clusterNode.Role
		clusterRole, exist := f.ClusterWrapper.ClusterRoles[role]
		if !exist {
			logger.Error(f.Ctx, "No such role [%s] in cluster role [%s]. ",
				role, f.ClusterWrapper.Cluster.ClusterId)
			return nil
		}

		size := clusterRole.StorageSize
		if size > 0 {
			// cmd will be assigned when the task is handling
			meta := &models.Meta{
				FrontgateId: f.ClusterWrapper.Cluster.FrontgateId,
				Timeout:     TimeoutFormatAndMountVolume,
				NodeId:      clusterNode.NodeId,
				DroneIp:     clusterNode.PrivateIp,
			}
			directive := jsonutil.ToString(meta)
			formatVolumeTask := &models.Task{
				JobId:          f.Job.JobId,
				Owner:          f.Job.Owner,
				TaskAction:     ActionFormatAndMountVolume,
				Target:         constants.TargetPilot,
				NodeId:         nodeId,
				Directive:      directive,
				FailureAllowed: failureAllowed,
			}
			taskLayer.Tasks = append(taskLayer.Tasks, formatVolumeTask)
		}
	}
	if len(taskLayer.Tasks) > 0 {
		return taskLayer
	} else {
		return nil
	}
}

func (f *Frame) removeContainerLayer(nodeIds []string, failureAllowed bool) *models.TaskLayer {
	taskLayer := new(models.TaskLayer)

	for nodeId, clusterNode := range f.ClusterWrapper.ClusterNodesWithKeyPairs {
		role := clusterNode.Role
		clusterRole, exist := f.ClusterWrapper.ClusterRoles[role]
		if !exist {
			logger.Error(f.Ctx, "No such role [%s] in cluster role [%s]. ",
				role, f.ClusterWrapper.Cluster.ClusterId)
			return nil
		}

		size := clusterRole.StorageSize
		if size > 0 {
			ip := clusterNode.PrivateIp
			cmd := fmt.Sprintf("%s \"docker rm -f default\"", HostCmdPrefix)
			request := &pbtypes.RunCommandOnDroneRequest{
				Endpoint: &pbtypes.DroneEndpoint{
					FrontgateId: f.ClusterWrapper.Cluster.FrontgateId,
					DroneIp:     ip,
					DronePort:   constants.DroneServicePort,
				},
				Command:        cmd,
				TimeoutSeconds: TimeoutRemoveContainer,
			}
			directive := jsonutil.ToString(request)
			formatVolumeTask := &models.Task{
				JobId:          f.Job.JobId,
				Owner:          f.Job.Owner,
				TaskAction:     ActionRemoveContainerOnDrone,
				Target:         constants.TargetPilot,
				NodeId:         nodeId,
				Directive:      directive,
				FailureAllowed: failureAllowed,
			}
			taskLayer.Tasks = append(taskLayer.Tasks, formatVolumeTask)
		}
	}
	if len(taskLayer.Tasks) > 0 {
		return taskLayer
	} else {
		return nil
	}
}

func (f *Frame) sshKeygenLayer(failureAllowed bool) *models.TaskLayer {
	taskLayer := new(models.TaskLayer)
	clusterClient, err := clusterclient.NewClient()
	if err != nil {
		logger.Error(f.Ctx, "New ssh key gen task layer failed: %+v", err)
		return nil
	}

	for nodeId, clusterNode := range f.ClusterWrapper.ClusterNodesWithKeyPairs {
		role := clusterNode.Role
		clusterCommon := f.ClusterWrapper.ClusterCommons[role]
		keyType := clusterCommon.Passphraseless
		if keyType != "" {
			private, public, err := sshutil.MakeSSHKeyPair(keyType)
			if err != nil {
				logger.Error(f.Ctx, "Generate ssh key [%s] in cluster node [%s] failed",
					clusterCommon.Passphraseless, nodeId)
				return nil
			}
			_, err = clusterClient.ModifyClusterNode(f.Ctx, &pb.ModifyClusterNodeRequest{
				ClusterNode: &pb.ClusterNode{
					NodeId: pbutil.ToProtoString(nodeId),
					PubKey: pbutil.ToProtoString(public),
				},
			})
			cmd := fmt.Sprintf("mkdir -p /root/.ssh/ && chmod 700 /root/.ssh/ && "+
				"echo \"%s\" > /root/.ssh/id_%s && echo \"%s\" > /root/.ssh/id_%s.pub && "+
				"chown 600 /root/.ssh/id_%s && chown 644 /root/.ssh/id_%s.pub",
				private, keyType, public, keyType, keyType, keyType)
			ip := clusterNode.PrivateIp

			request := &pbtypes.RunCommandOnDroneRequest{
				Endpoint: &pbtypes.DroneEndpoint{
					FrontgateId: f.ClusterWrapper.Cluster.FrontgateId,
					DroneIp:     ip,
					DronePort:   constants.DroneServicePort,
				},
				Command:        cmd,
				TimeoutSeconds: TimeoutSshKeygen,
			}
			directive := jsonutil.ToString(request)
			formatVolumeTask := &models.Task{
				JobId:          f.Job.JobId,
				Owner:          f.Job.Owner,
				TaskAction:     ActionRunCommandOnDrone,
				Target:         constants.TargetPilot,
				NodeId:         nodeId,
				Directive:      directive,
				FailureAllowed: failureAllowed,
			}
			taskLayer.Tasks = append(taskLayer.Tasks, formatVolumeTask)
		}
	}
	if len(taskLayer.Tasks) > 0 {
		return taskLayer
	} else {
		return nil
	}
}

func (f *Frame) umountVolumeLayer(nodeIds []string, failureAllowed bool) *models.TaskLayer {
	taskLayer := new(models.TaskLayer)

	for _, nodeId := range nodeIds {
		clusterNode := f.ClusterWrapper.ClusterNodesWithKeyPairs[nodeId]
		if clusterNode.VolumeId == "" {
			continue
		}
		clusterRole := f.ClusterWrapper.ClusterRoles[clusterNode.Role]
		cmd := UmountVolumeCmd(clusterRole.MountPoint)
		ip := clusterNode.PrivateIp

		umountVolumeTask := &models.Task{
			JobId:          f.Job.JobId,
			Owner:          f.Job.Owner,
			Target:         constants.TargetPilot,
			NodeId:         nodeId,
			FailureAllowed: failureAllowed,
		}

		if f.ClusterWrapper.Cluster.ClusterType == constants.FrontgateClusterType {
			request := &pbtypes.RunCommandOnFrontgateRequest{
				Endpoint: &pbtypes.FrontgateEndpoint{
					FrontgateId:     f.ClusterWrapper.Cluster.ClusterId,
					FrontgateNodeId: nodeId,
					NodeIp:          ip,
					NodePort:        constants.FrontgateServicePort,
				},
				Command:        cmd,
				TimeoutSeconds: TimeoutUmountVolume,
			}
			umountVolumeTask.Directive = jsonutil.ToString(request)
			umountVolumeTask.TaskAction = ActionRunCommandOnFrontgateNode
		} else {
			request := &pbtypes.RunCommandOnDroneRequest{
				Endpoint: &pbtypes.DroneEndpoint{
					FrontgateId: f.ClusterWrapper.Cluster.FrontgateId,
					DroneIp:     ip,
					DronePort:   constants.DroneServicePort,
				},
				Command:        cmd,
				TimeoutSeconds: TimeoutUmountVolume,
			}
			umountVolumeTask.Directive = jsonutil.ToString(request)
			umountVolumeTask.TaskAction = ActionRunCommandOnDrone
		}
		taskLayer.Tasks = append(taskLayer.Tasks, umountVolumeTask)
	}
	if len(taskLayer.Tasks) > 0 {
		return taskLayer
	} else {
		return nil
	}
}

func (f *Frame) getUserDataExec(filename, contents, imageUrl, certificateExec string) string {
	if pi.Global() == nil {
		logger.Error(f.Ctx, "Pi global should be init.")
		return ""
	}

	exec := fmt.Sprintf(`
mkdir -p /opt/openpitrix/image/ /opt/openpitrix/conf/
test -f /opt/openpitrix/conf/init && exit 0
%s
echo '%s' > %s
for i in $(seq 1 100); do cd /opt/openpitrix/image/ && rm -rf * && wget %s && tar -xzvf * && break || sleep 3; done
/opt/openpitrix/image/install_service.sh %s
touch /opt/openpitrix/conf/init
`, certificateExec, contents, f.getConfFile(), imageUrl, filename)
	return exec
}

func FormatUserData(userData, frontgateIp string) string {
	var mirror string
	if len(frontgateIp) > 0 {
		mirror = fmt.Sprintf("http://%s:5000", frontgateIp)
	} else {
		mirror = pi.Global().GlobalConfig().Cluster.RegistryMirror
	}

	data := ""
	if len(mirror) > 0 {
		data = fmt.Sprintf(`#!/bin/bash -e

mkdir -p /etc/docker/
echo '{
  "registry-mirrors": ["%s"]
}' > /etc/docker/daemon.json
%s
`, mirror, userData)
	} else {
		data = fmt.Sprintf(`#!/bin/bash -e

%s
`, userData)
	}

	return base64.StdEncoding.EncodeToString([]byte(data))
}

/*
cat /opt/openpitrix/conf/drone.conf
IMAGE="mysql:5.7"
MOUNT_POINT="/data"
FILE_NAME="drone.conf"
FILE_CONF={\\"id\\":\\"cln-abcdefgh\\",\\"listen_port\\":9112}
*/
func (f *Frame) getUserDataValue(nodeId string) string {
	var result string
	clusterNode := f.ClusterWrapper.ClusterNodesWithKeyPairs[nodeId]
	role := clusterNode.Role
	if strings.HasSuffix(role, constants.ReplicaRoleSuffix) {
		role = string([]byte(role)[:len(role)-len(constants.ReplicaRoleSuffix)])
	}
	clusterRole, _ := f.ClusterWrapper.ClusterRoles[role]
	clusterCommon, _ := f.ClusterWrapper.ClusterCommons[role]
	mountPoint := clusterRole.MountPoint
	// Empty string can not be a parameter
	if len(mountPoint) == 0 {
		mountPoint = "#"
	}
	imageId := clusterCommon.ImageId

	droneConf := make(map[string]interface{})
	droneConf["id"] = nodeId
	droneConf["listen_port"] = constants.DroneServicePort
	droneConfStr := strings.Replace(jsonutil.ToString(droneConf), "\"", "\\\\\"", -1)

	result += fmt.Sprintf("IMAGE=\"%s\"\n", imageId)
	result += fmt.Sprintf("MOUNT_POINT=\"%s\"\n", mountPoint)
	result += fmt.Sprintf("FILE_NAME=\"%s\"\n", DroneConfFile)
	result += fmt.Sprintf("FILE_CONF=%s\n", droneConfStr)

	return result
}

func (f *Frame) getUserDataFile() string {
	return OpenPitrixExecFile
}

func (f *Frame) getConfFile() string {
	return OpenPitrixConfPath + OpenPitrixConfFile
}

func (f *Frame) setDroneConfigLayer(nodeIds []string, failureAllowed bool) *models.TaskLayer {
	var tasks []*models.Task
	directive := jsonutil.ToString(&models.Meta{
		ClusterId: f.ClusterWrapper.Cluster.ClusterId,
	})

	for _, nodeId := range nodeIds {
		// get drone config when pre task
		task := &models.Task{
			JobId:          f.Job.JobId,
			Owner:          f.Job.Owner,
			TaskAction:     ActionSetDroneConfig,
			Target:         constants.TargetPilot,
			NodeId:         nodeId,
			Directive:      directive,
			FailureAllowed: failureAllowed,
		}
		tasks = append(tasks, task)
	}
	return &models.TaskLayer{
		Tasks: tasks,
	}
}

func (f *Frame) runInstancesLayer(nodeIds []string, failureAllowed bool) *models.TaskLayer {
	taskLayer := new(models.TaskLayer)

	for _, nodeId := range nodeIds {
		clusterNode := f.ClusterWrapper.ClusterNodesWithKeyPairs[nodeId]
		role := clusterNode.Role
		if strings.HasSuffix(role, constants.ReplicaRoleSuffix) {
			role = string([]byte(role)[:len(role)-len(constants.ReplicaRoleSuffix)])
		}
		clusterRole, exist := f.ClusterWrapper.ClusterRoles[role]
		if !exist {
			logger.Error(f.Ctx, "No such role [%s] in cluster role [%s]. ",
				role, f.ClusterWrapper.Cluster.ClusterId)
			return nil
		}

		showName := role
		if f.ClusterWrapper.Cluster.ClusterType == constants.FrontgateClusterType {
			showName = MetadataNodeName
		} else if role == "" {
			showName = DefaultNodeName
		}
		instance := &models.Instance{
			Name:         clusterNode.ClusterId + "_" + nodeId + "_" + showName,
			Hostname:     nodeId,
			NodeId:       nodeId,
			ImageId:      f.RuntimeProviderConfig.ImageId,
			Cpu:          int(clusterRole.Cpu),
			Memory:       int(clusterRole.Memory),
			Gpu:          int(clusterRole.Gpu),
			Subnet:       clusterNode.SubnetId,
			RuntimeId:    f.Runtime.RuntimeId,
			Zone:         f.ClusterWrapper.Cluster.Zone,
			NeedUserData: 1,
			UserdataFile: f.getUserDataFile(),
		}
		if f.ClusterWrapper.Cluster.ClusterType == constants.FrontgateClusterType {
			frontgate := &Frontgate{f}
			instance.UserDataValue = f.getUserDataExec(FrontgateConfFile, frontgate.getUserDataValue(nodeId), f.RuntimeProviderConfig.ImageUrl, frontgate.getCertificateExec())
		} else {
			instance.UserDataValue = f.getUserDataExec(DroneConfFile, f.getUserDataValue(nodeId), f.RuntimeProviderConfig.ImageUrl, "")
		}
		directive := jsonutil.ToString(instance)
		runInstanceTask := &models.Task{
			JobId:          f.Job.JobId,
			Owner:          f.Job.Owner,
			TaskAction:     ActionRunInstances,
			Target:         f.Runtime.RuntimeId,
			NodeId:         nodeId,
			Directive:      directive,
			FailureAllowed: failureAllowed,
		}
		taskLayer.Tasks = append(taskLayer.Tasks, runInstanceTask)
	}

	if len(taskLayer.Tasks) > 0 {
		return taskLayer
	} else {
		return nil
	}
}

func (f *Frame) stopInstancesLayer(nodeIds []string, failureAllowed bool) *models.TaskLayer {
	taskLayer := new(models.TaskLayer)
	for _, nodeId := range nodeIds {
		clusterNode := f.ClusterWrapper.ClusterNodesWithKeyPairs[nodeId]
		instance := &models.Instance{
			Name:       clusterNode.ClusterId + "_" + nodeId,
			NodeId:     nodeId,
			InstanceId: clusterNode.InstanceId,
			RuntimeId:  f.Runtime.RuntimeId,
			Zone:       f.ClusterWrapper.Cluster.Zone,
		}
		directive := jsonutil.ToString(instance)
		stopInstanceTask := &models.Task{
			JobId:          f.Job.JobId,
			Owner:          f.Job.Owner,
			TaskAction:     ActionStopInstances,
			Target:         f.Runtime.RuntimeId,
			NodeId:         nodeId,
			Directive:      directive,
			FailureAllowed: failureAllowed,
		}
		taskLayer.Tasks = append(taskLayer.Tasks, stopInstanceTask)
	}

	if len(taskLayer.Tasks) > 0 {
		return taskLayer
	} else {
		return nil
	}
}

func (f *Frame) deleteInstancesLayer(nodeIds []string, failureAllowed bool) *models.TaskLayer {
	taskLayer := new(models.TaskLayer)
	for _, nodeId := range nodeIds {
		clusterNode := f.ClusterWrapper.ClusterNodesWithKeyPairs[nodeId]
		instance := &models.Instance{
			Name:       clusterNode.ClusterId + "_" + nodeId,
			NodeId:     nodeId,
			InstanceId: clusterNode.InstanceId,
			RuntimeId:  f.Runtime.RuntimeId,
			Zone:       f.ClusterWrapper.Cluster.Zone,
		}
		directive := jsonutil.ToString(instance)
		deleteInstanceTask := &models.Task{
			JobId:          f.Job.JobId,
			Owner:          f.Job.Owner,
			TaskAction:     ActionTerminateInstances,
			Target:         f.Runtime.RuntimeId,
			NodeId:         nodeId,
			Directive:      directive,
			FailureAllowed: false,
		}
		taskLayer.Tasks = append(taskLayer.Tasks, deleteInstanceTask)
	}

	if len(taskLayer.Tasks) > 0 {
		return taskLayer
	} else {
		return nil
	}
}

func (f *Frame) startInstancesLayer(nodeIds []string, failureAllowed bool) *models.TaskLayer {
	taskLayer := new(models.TaskLayer)
	for _, nodeId := range nodeIds {
		clusterNode := f.ClusterWrapper.ClusterNodesWithKeyPairs[nodeId]
		instance := &models.Instance{
			Name:       clusterNode.ClusterId + "_" + nodeId,
			NodeId:     nodeId,
			InstanceId: clusterNode.InstanceId,
			RuntimeId:  f.Runtime.RuntimeId,
			Zone:       f.ClusterWrapper.Cluster.Zone,
		}
		directive := jsonutil.ToString(instance)
		startInstanceTask := &models.Task{
			JobId:          f.Job.JobId,
			Owner:          f.Job.Owner,
			TaskAction:     ActionStartInstances,
			Target:         f.Runtime.RuntimeId,
			NodeId:         nodeId,
			Directive:      directive,
			FailureAllowed: failureAllowed,
		}
		taskLayer.Tasks = append(taskLayer.Tasks, startInstanceTask)
	}

	if len(taskLayer.Tasks) > 0 {
		return taskLayer
	} else {
		return nil
	}
}

func (f *Frame) resizeInstancesLayer(nodeIds []string, roleResizeResource *models.RoleResizeResource, failureAllowed bool) *models.TaskLayer {
	taskLayer := new(models.TaskLayer)
	for _, nodeId := range nodeIds {
		clusterNode := f.ClusterWrapper.ClusterNodesWithKeyPairs[nodeId]
		clusterRole := f.ClusterWrapper.ClusterRoles[clusterNode.Role]
		if !roleResizeResource.Cpu &&
			!roleResizeResource.Gpu &&
			!roleResizeResource.Memory &&
			!roleResizeResource.InstanceSize {
			logger.Debug(f.Ctx, "No need to resize node [%s] instance", nodeId)
			continue
		}
		instance := &models.Instance{
			Name:       clusterNode.ClusterId + "_" + nodeId,
			NodeId:     nodeId,
			InstanceId: clusterNode.InstanceId,
			RuntimeId:  f.Runtime.RuntimeId,
			Zone:       f.ClusterWrapper.Cluster.Zone,
			Cpu:        int(clusterRole.Cpu),
			Memory:     int(clusterRole.Memory),
			Gpu:        int(clusterRole.Gpu),
		}
		directive := jsonutil.ToString(instance)
		resizeInstanceTask := &models.Task{
			JobId:          f.Job.JobId,
			Owner:          f.Job.Owner,
			TaskAction:     ActionResizeInstances,
			Target:         f.Runtime.RuntimeId,
			NodeId:         nodeId,
			Directive:      directive,
			FailureAllowed: failureAllowed,
		}
		taskLayer.Tasks = append(taskLayer.Tasks, resizeInstanceTask)
	}

	if len(taskLayer.Tasks) > 0 {
		return taskLayer
	} else {
		return nil
	}
}

func (f *Frame) waitFrontgateLayer(failureAllowed bool) *models.TaskLayer {
	meta := &models.Meta{
		FrontgateId: f.ClusterWrapper.Cluster.FrontgateId,
	}
	directive := jsonutil.ToString(meta)
	// Wait frontgate available
	waitFrontgateTask := &models.Task{
		JobId:          f.Job.JobId,
		Owner:          f.Job.Owner,
		TaskAction:     ActionWaitFrontgateAvailable,
		Target:         f.Runtime.RuntimeId,
		NodeId:         f.ClusterWrapper.Cluster.ClusterId,
		Directive:      directive,
		FailureAllowed: failureAllowed,
	}
	return &models.TaskLayer{
		Tasks: []*models.Task{waitFrontgateTask},
	}
}

func (f *Frame) pingDroneLayer(nodeIds []string, failureAllowed bool) *models.TaskLayer {
	taskLayer := new(models.TaskLayer)
	for _, nodeId := range nodeIds {
		clusterNode := f.ClusterWrapper.ClusterNodesWithKeyPairs[nodeId]
		droneEndpoint := &pbtypes.DroneEndpoint{
			FrontgateId: f.ClusterWrapper.Cluster.FrontgateId,
			DroneIp:     clusterNode.PrivateIp,
			DronePort:   constants.DroneServicePort,
		}
		task := &models.Task{
			JobId:          f.Job.JobId,
			Owner:          f.Job.Owner,
			TaskAction:     ActionPingDrone,
			Target:         constants.TargetPilot,
			NodeId:         nodeId,
			Directive:      jsonutil.ToString(droneEndpoint),
			FailureAllowed: failureAllowed,
		}
		taskLayer.Tasks = append(taskLayer.Tasks, task)
	}
	if len(taskLayer.Tasks) > 0 {
		return taskLayer
	} else {
		return nil
	}
}

func (f *Frame) registerMetadataLayer(failureAllowed bool) *models.TaskLayer {
	// When the task is handled by task controller, the cnodes will be filled in,
	meta := &models.Meta{
		FrontgateId: f.ClusterWrapper.Cluster.FrontgateId,
		Timeout:     TimeoutRegister,
		ClusterId:   f.ClusterWrapper.Cluster.ClusterId,
	}
	directive := jsonutil.ToString(meta)
	registerMetadataTask := &models.Task{
		JobId:          f.Job.JobId,
		Owner:          f.Job.Owner,
		TaskAction:     ActionRegisterMetadata,
		Target:         constants.TargetPilot,
		NodeId:         f.ClusterWrapper.Cluster.ClusterId,
		Directive:      directive,
		FailureAllowed: failureAllowed,
	}
	return &models.TaskLayer{
		Tasks: []*models.Task{registerMetadataTask},
	}
}

func (f *Frame) registerMetadataMappingLayer(failureAllowed bool) *models.TaskLayer {
	// When the task is handled by task controller, the cnodes will be filled in,
	meta := &models.Meta{
		FrontgateId: f.ClusterWrapper.Cluster.FrontgateId,
		Timeout:     TimeoutRegister,
		ClusterId:   f.ClusterWrapper.Cluster.ClusterId,
	}
	directive := jsonutil.ToString(meta)
	registerMetadataTask := &models.Task{
		JobId:          f.Job.JobId,
		Owner:          f.Job.Owner,
		TaskAction:     ActionRegisterMetadataMapping,
		Target:         constants.TargetPilot,
		NodeId:         f.ClusterWrapper.Cluster.ClusterId,
		Directive:      directive,
		FailureAllowed: failureAllowed,
	}
	return &models.TaskLayer{
		Tasks: []*models.Task{registerMetadataTask},
	}
}

func (f *Frame) registerNodesMetadataLayer(nodeIds []string, failureAllowed bool) *models.TaskLayer {
	metadata := &Metadata{
		ClusterWrapper: f.ClusterWrapper,
		RuntimeDetails: f.Runtime,
	}
	cnodes := jsonutil.ToString(metadata.GetClusterNodesCnodes(f.Ctx, nodeIds))
	meta := &models.Meta{
		FrontgateId: f.ClusterWrapper.Cluster.FrontgateId,
		Timeout:     TimeoutRegister,
		ClusterId:   f.ClusterWrapper.Cluster.ClusterId,
		Cnodes:      cnodes,
	}
	directive := jsonutil.ToString(meta)
	task := &models.Task{
		JobId:          f.Job.JobId,
		Owner:          f.Job.Owner,
		TaskAction:     ActionRegisterNodesMetadata,
		Target:         constants.TargetPilot,
		NodeId:         f.ClusterWrapper.Cluster.ClusterId,
		Directive:      directive,
		FailureAllowed: failureAllowed,
	}
	return &models.TaskLayer{
		Tasks: []*models.Task{task},
	}
}

func (f *Frame) registerNodesMetadataMappingLayer(nodeIds []string, failureAllowed bool) *models.TaskLayer {
	metadata := &Metadata{
		ClusterWrapper: f.ClusterWrapper,
		RuntimeDetails: f.Runtime,
	}
	cnodes := jsonutil.ToString(metadata.GetClusterNodesMappingCnodes(f.Ctx, nodeIds))
	meta := &models.Meta{
		FrontgateId: f.ClusterWrapper.Cluster.FrontgateId,
		Timeout:     TimeoutRegister,
		ClusterId:   f.ClusterWrapper.Cluster.ClusterId,
		Cnodes:      cnodes,
	}
	directive := jsonutil.ToString(meta)
	task := &models.Task{
		JobId:          f.Job.JobId,
		Owner:          f.Job.Owner,
		TaskAction:     ActionRegisterNodesMetadataMapping,
		Target:         constants.TargetPilot,
		NodeId:         f.ClusterWrapper.Cluster.ClusterId,
		Directive:      directive,
		FailureAllowed: failureAllowed,
	}
	return &models.TaskLayer{
		Tasks: []*models.Task{task},
	}
}

func (f *Frame) registerScalingNodesMetadataLayer(nodeIds []string, path string, failureAllowed bool) *models.TaskLayer {
	clusterId := f.ClusterWrapper.Cluster.ClusterId
	metadata := &Metadata{
		ClusterWrapper: f.ClusterWrapper,
		RuntimeDetails: f.Runtime,
	}
	scalingCnodes := metadata.GetScalingCnodes(f.Ctx, nodeIds, path)
	if scalingCnodes == nil {
		logger.Info(f.Ctx, "No new nodes for cluster [%s] is registered", clusterId)
		return nil
	}
	meta := &models.Meta{
		FrontgateId: f.ClusterWrapper.Cluster.FrontgateId,
		Timeout:     TimeoutRegister,
		ClusterId:   f.ClusterWrapper.Cluster.ClusterId,
		Cnodes:      jsonutil.ToString(scalingCnodes),
	}
	directive := jsonutil.ToString(meta)
	task := &models.Task{
		JobId:          f.Job.JobId,
		Owner:          f.Job.Owner,
		TaskAction:     ActionRegisterNodesMetadata,
		Target:         constants.TargetPilot,
		NodeId:         f.ClusterWrapper.Cluster.ClusterId,
		Directive:      directive,
		FailureAllowed: failureAllowed,
	}
	return &models.TaskLayer{
		Tasks: []*models.Task{task},
	}
}

func (f *Frame) registerEnvMetadataLayer(failureAllowed bool) *models.TaskLayer {
	metadata := &Metadata{
		ClusterWrapper: f.ClusterWrapper,
		RuntimeDetails: f.Runtime,
	}
	cnodes := jsonutil.ToString(metadata.GetClusterEnvCnodes(f.Ctx))
	meta := &models.Meta{
		FrontgateId: f.ClusterWrapper.Cluster.FrontgateId,
		Timeout:     TimeoutRegister,
		ClusterId:   f.ClusterWrapper.Cluster.ClusterId,
		Cnodes:      cnodes,
	}
	directive := jsonutil.ToString(meta)
	task := &models.Task{
		JobId:          f.Job.JobId,
		Owner:          f.Job.Owner,
		TaskAction:     ActionRegisterEnvMetadata,
		Target:         constants.TargetPilot,
		NodeId:         f.ClusterWrapper.Cluster.ClusterId,
		Directive:      directive,
		FailureAllowed: failureAllowed,
	}
	return &models.TaskLayer{
		Tasks: []*models.Task{task},
	}
}

func (f *Frame) deregisterNodesMetadataLayer(nodeIds []string, failureAllowed bool) *models.TaskLayer {
	metadata := &Metadata{
		ClusterWrapper: f.ClusterWrapper,
		RuntimeDetails: f.Runtime,
	}
	cnodes := jsonutil.ToString(metadata.GetEmptyClusterNodeCnodes(f.Ctx, nodeIds))
	meta := &models.Meta{
		FrontgateId: f.ClusterWrapper.Cluster.FrontgateId,
		Timeout:     TimeoutDeregister,
		ClusterId:   f.ClusterWrapper.Cluster.ClusterId,
		Cnodes:      cnodes,
	}
	directive := jsonutil.ToString(meta)
	task := &models.Task{
		JobId:          f.Job.JobId,
		Owner:          f.Job.Owner,
		TaskAction:     ActionDeregisterMetadata,
		Target:         constants.TargetPilot,
		NodeId:         f.ClusterWrapper.Cluster.ClusterId,
		Directive:      directive,
		FailureAllowed: failureAllowed,
	}
	return &models.TaskLayer{
		Tasks: []*models.Task{task},
	}
}

func (f *Frame) deregisterNodesMetadataMappingLayer(nodeIds []string, failureAllowed bool) *models.TaskLayer {
	metadata := &Metadata{
		ClusterWrapper: f.ClusterWrapper,
		RuntimeDetails: f.Runtime,
	}
	cnodes := jsonutil.ToString(metadata.GetEmptyClusterNodeMappingCnodes(f.Ctx, nodeIds))
	meta := &models.Meta{
		FrontgateId: f.ClusterWrapper.Cluster.FrontgateId,
		Timeout:     TimeoutDeregister,
		ClusterId:   f.ClusterWrapper.Cluster.ClusterId,
		Cnodes:      cnodes,
	}
	directive := jsonutil.ToString(meta)
	task := &models.Task{
		JobId:          f.Job.JobId,
		Owner:          f.Job.Owner,
		TaskAction:     ActionDeregisterMetadataMapping,
		Target:         constants.TargetPilot,
		NodeId:         f.ClusterWrapper.Cluster.ClusterId,
		Directive:      directive,
		FailureAllowed: failureAllowed,
	}
	return &models.TaskLayer{
		Tasks: []*models.Task{task},
	}
}

func (f *Frame) deregisterScalingNodesMetadataLayer(path string, failureAllowed bool) *models.TaskLayer {
	clusterId := f.ClusterWrapper.Cluster.ClusterId
	cnodes := map[string]interface{}{
		RegisterClustersRootPath: map[string]interface{}{
			clusterId: map[string]interface{}{
				path: "",
			},
		},
	}
	meta := &models.Meta{
		FrontgateId: f.ClusterWrapper.Cluster.FrontgateId,
		Timeout:     TimeoutDeregister,
		ClusterId:   f.ClusterWrapper.Cluster.ClusterId,
		Cnodes:      jsonutil.ToString(cnodes),
	}
	directive := jsonutil.ToString(meta)
	task := &models.Task{
		JobId:          f.Job.JobId,
		Owner:          f.Job.Owner,
		TaskAction:     ActionDeregisterMetadata,
		Target:         constants.TargetPilot,
		NodeId:         f.ClusterWrapper.Cluster.ClusterId,
		Directive:      directive,
		FailureAllowed: failureAllowed,
	}
	return &models.TaskLayer{
		Tasks: []*models.Task{task},
	}
}

func (f *Frame) deregisterMetadataLayer(failureAllowed bool) *models.TaskLayer {
	metadata := &Metadata{
		ClusterWrapper: f.ClusterWrapper,
		RuntimeDetails: f.Runtime,
	}
	cnodes := metadata.GetEmptyClusterCnodes()
	meta := &models.Meta{
		FrontgateId: f.ClusterWrapper.Cluster.FrontgateId,
		Timeout:     TimeoutDeregister,
		ClusterId:   f.ClusterWrapper.Cluster.ClusterId,
		Cnodes:      jsonutil.ToString(cnodes),
	}
	directive := jsonutil.ToString(meta)
	deregisterMetadataTask := &models.Task{
		JobId:          f.Job.JobId,
		Owner:          f.Job.Owner,
		TaskAction:     ActionDeregisterMetadata,
		Target:         constants.TargetPilot,
		NodeId:         f.ClusterWrapper.Cluster.ClusterId,
		Directive:      directive,
		FailureAllowed: failureAllowed,
	}
	return &models.TaskLayer{
		Tasks: []*models.Task{deregisterMetadataTask},
	}
}

func (f *Frame) deregisterMetadataMappingLayer(failureAllowed bool) *models.TaskLayer {
	metadata := &Metadata{
		ClusterWrapper: f.ClusterWrapper,
		RuntimeDetails: f.Runtime,
	}
	cnodes := metadata.GetEmptyClusterMappingCnodes()
	meta := &models.Meta{
		FrontgateId: f.ClusterWrapper.Cluster.FrontgateId,
		Timeout:     TimeoutDeregister,
		ClusterId:   f.ClusterWrapper.Cluster.ClusterId,
		Cnodes:      jsonutil.ToString(cnodes),
	}
	directive := jsonutil.ToString(meta)
	deregisterMetadataTask := &models.Task{
		JobId:          f.Job.JobId,
		Owner:          f.Job.Owner,
		TaskAction:     ActionDeregisterMetadataMapping,
		Target:         constants.TargetPilot,
		NodeId:         f.ClusterWrapper.Cluster.ClusterId,
		Directive:      directive,
		FailureAllowed: failureAllowed,
	}
	return &models.TaskLayer{
		Tasks: []*models.Task{deregisterMetadataTask},
	}
}

func (f *Frame) CreateClusterLayer() *models.TaskLayer {
	var nodeIds []string
	for nodeId := range f.ClusterWrapper.ClusterNodesWithKeyPairs {
		nodeIds = append(nodeIds, nodeId)
	}
	headTaskLayer := new(models.TaskLayer)

	headTaskLayer.
		Append(f.createVolumesLayer(nodeIds, false)).        // create volume
		Append(f.waitFrontgateLayer(false)).                 // wait frontgate cluster to be active
		Append(f.runInstancesLayer(nodeIds, false)).         // run instance and attach volume to instance
		Append(f.pingDroneLayer(nodeIds, false)).            // ping drone
		Append(f.setDroneConfigLayer(nodeIds, false)).       // set drone config
		Append(f.formatAndMountVolumeLayer(nodeIds, false)). // format and mount volume to instance
		Append(f.removeContainerLayer(nodeIds, false)).      // remove default container
		Append(f.pingDroneLayer(nodeIds, false)).            // ping drone
		Append(f.setDroneConfigLayer(nodeIds, false)).       // set drone config
		Append(f.sshKeygenLayer(false)).                     // generate ssh key
		Append(f.deregisterMetadataLayer(true)).             // deregister cluster metadata
		Append(f.deregisterMetadataMappingLayer(true)).      // deregister cluster metadata mapping
		Append(f.registerMetadataLayer(false)).              // register cluster metadata
		Append(f.registerMetadataMappingLayer(false)).       // register cluster metadata mapping
		Append(f.startConfdServiceLayer(nodeIds, false)).    // start confd service
		Append(f.initAndStartServiceLayer(nodeIds, false)).  // register init and start cmd to exec
		Append(f.deregisterCmdLayer(nodeIds, true))          // deregister cmd

	return headTaskLayer.Child
}

func (f *Frame) StopClusterLayer() *models.TaskLayer {
	var nodeIds []string
	for nodeId := range f.ClusterWrapper.ClusterNodesWithKeyPairs {
		nodeIds = append(nodeIds, nodeId)
	}
	headTaskLayer := new(models.TaskLayer)

	headTaskLayer.
		Append(f.stopServiceLayer(nodeIds, true)).      // register stop cmd to exec
		Append(f.stopConfdServiceLayer(nodeIds, true)). // stop confd service
		Append(f.umountVolumeLayer(nodeIds, true)).     // umount volume from instance
		Append(f.stopInstancesLayer(nodeIds, false)).   // stop instance
		Append(f.detachVolumesLayer(nodeIds, false)).   // detach volume from instance
		Append(f.deregisterMetadataLayer(true)).        // deregister cluster metadata
		Append(f.deregisterMetadataMappingLayer(true))  // deregister cluster metadata mapping

	return headTaskLayer.Child
}

func (f *Frame) StartClusterLayer() *models.TaskLayer {
	var nodeIds []string
	for nodeId := range f.ClusterWrapper.ClusterNodesWithKeyPairs {
		nodeIds = append(nodeIds, nodeId)
	}
	headTaskLayer := new(models.TaskLayer)

	headTaskLayer.
		Append(f.attachVolumesLayer(nodeIds, false)).     // attach volume to instance, will auto mount
		Append(f.startInstancesLayer(nodeIds, false)).    // start instance
		Append(f.waitFrontgateLayer(false)).              // wait frontgate cluster to be active
		Append(f.registerMetadataLayer(false)).           // register cluster metadata
		Append(f.registerMetadataMappingLayer(false)).    // register cluster metadata mapping
		Append(f.pingDroneLayer(nodeIds, false)).         // ping drone
		Append(f.setDroneConfigLayer(nodeIds, false)).    // set drone config
		Append(f.startConfdServiceLayer(nodeIds, false)). // start confd service
		Append(f.startServiceLayer(nodeIds, false)).      // register start cmd to exec
		Append(f.deregisterCmdLayer(nodeIds, true))       // deregister cmd

	return headTaskLayer.Child
}

func (f *Frame) DeleteClusterLayer() *models.TaskLayer {
	var nodeIds []string
	for nodeId := range f.ClusterWrapper.ClusterNodesWithKeyPairs {
		nodeIds = append(nodeIds, nodeId)
	}
	headTaskLayer := new(models.TaskLayer)

	if f.ClusterWrapper.Cluster.Status == constants.StatusActive {
		headTaskLayer.
			Append(f.destroyAndStopServiceLayer(nodeIds, nil, true)). // register destroy and stop cmd to exec
			Append(f.stopConfdServiceLayer(nodeIds, true)).           // stop confd service
			Append(f.umountVolumeLayer(nodeIds, true)).               // umount volume from instance
			Append(f.stopInstancesLayer(nodeIds, true)).              // stop instance
			Append(f.detachVolumesLayer(nodeIds, false))              // detach volume from instance
	}

	headTaskLayer.
		Append(f.deleteInstancesLayer(nodeIds, false)). // delete instance
		Append(f.deleteVolumesLayer(nodeIds, false)).   // delete volume
		Append(f.deregisterMetadataLayer(true)).        // deregister cluster metadata
		Append(f.deregisterMetadataMappingLayer(true))  // deregister cluster metadata mapping

	return headTaskLayer.Child
}

func (f *Frame) AddClusterNodesLayer() *models.TaskLayer {
	var addNodeIds, nonAddNodeIds []string
	for nodeId, node := range f.ClusterWrapper.ClusterNodesWithKeyPairs {
		if node.Status == constants.StatusPending {
			addNodeIds = append(addNodeIds, nodeId)
		} else {
			nonAddNodeIds = append(nonAddNodeIds, nodeId)
		}
	}
	headTaskLayer := new(models.TaskLayer)

	headTaskLayer.
		Append(f.scaleOutPreCheckServiceLayer(nonAddNodeIds, false)).                       // register scale out pre check to exec
		Append(f.createVolumesLayer(addNodeIds, false)).                                    // create volume
		Append(f.runInstancesLayer(addNodeIds, false)).                                     // run instance and attach volume to instance
		Append(f.pingDroneLayer(addNodeIds, false)).                                        // ping drone
		Append(f.setDroneConfigLayer(addNodeIds, false)).                                   // set drone config
		Append(f.formatAndMountVolumeLayer(addNodeIds, false)).                             // format and mount volume to instance
		Append(f.registerNodesMetadataLayer(addNodeIds, false)).                            // register cluster nodes metadata
		Append(f.registerScalingNodesMetadataLayer(addNodeIds, RegisterNodeAdding, false)). // register adding hosts metadata
		Append(f.startConfdServiceLayer(addNodeIds, false)).                                // start confd service
		Append(f.initAndStartServiceLayer(addNodeIds, false)).                              // register init and start cmd to exec
		Append(f.scaleOutServiceLayer(nonAddNodeIds, false)).                               // register scale out cmd to exec
		Append(f.deregisterScalingNodesMetadataLayer(RegisterNodeAdding, true))             // deregister adding host metadata
	return headTaskLayer.Child
}

func (f *Frame) DeleteClusterNodesLayer() *models.TaskLayer {
	var deleteNodeIds, nonDeleteNodeIds []string
	for nodeId, node := range f.ClusterWrapper.ClusterNodesWithKeyPairs {
		if node.Status == constants.StatusDeleting {
			deleteNodeIds = append(deleteNodeIds, nodeId)
		} else {
			nonDeleteNodeIds = append(nonDeleteNodeIds, nodeId)
		}
	}
	headTaskLayer := new(models.TaskLayer)

	headTaskLayer.
		Append(f.registerScalingNodesMetadataLayer(deleteNodeIds, RegisterNodeDeleting, false)).                    // register scale in node metadata
		Append(f.scaleInPreCheckServiceLayer(nonDeleteNodeIds, false)).                                             // register scale in pre check to exec
		Append(f.destroyAndStopServiceLayer(deleteNodeIds, f.scaleInServiceLayer(nonDeleteNodeIds, false), false)). // register destroy, scale in and stop cmd to exec
		Append(f.stopConfdServiceLayer(deleteNodeIds, false)).                                                      // stop confd service
		Append(f.umountVolumeLayer(deleteNodeIds, false)).                                                          // umount volume from instance
		Append(f.stopInstancesLayer(deleteNodeIds, false)).
		Append(f.detachVolumesLayer(deleteNodeIds, false)).                        // detach volume from instance
		Append(f.deleteInstancesLayer(deleteNodeIds, false)).                      // delete instance
		Append(f.deleteVolumesLayer(deleteNodeIds, false)).                        // delete volume
		Append(f.deregisterNodesMetadataLayer(deleteNodeIds, false)).              // deregister deleting cluster nodes metadata
		Append(f.deregisterScalingNodesMetadataLayer(RegisterNodeDeleting, false)) // deregister deleting nodes metadata
	return headTaskLayer.Child
}

func (f *Frame) UpdateClusterEnvLayer() *models.TaskLayer {
	headTaskLayer := new(models.TaskLayer)
	headTaskLayer.Append(f.registerMetadataLayer(false)) // register cluster metadata
	return headTaskLayer.Child
}

func (f *Frame) resizeClusterLayer(headTaskLayer *models.TaskLayer, nodeIds []string, roleResizeResource *models.RoleResizeResource) {
	headTaskLayer.
		Append(f.stopServiceLayer(nodeIds, true)).      // register stop cmd to exec
		Append(f.stopConfdServiceLayer(nodeIds, true)). // stop confd service
		Append(f.umountVolumeLayer(nodeIds, true)).     // umount volume from instance
		Append(f.stopInstancesLayer(nodeIds, false)).   // stop instance
		Append(f.detachVolumesLayer(nodeIds, false)).   // detach volume from instance
		Append(f.resizeInstancesLayer(nodeIds, roleResizeResource, false)).
		Append(f.resizeVolumesLayer(nodeIds, roleResizeResource, false)).
		Append(f.attachVolumesLayer(nodeIds, false)).         // attach volume to instance, will auto mount
		Append(f.startInstancesLayer(nodeIds, false)).        // start instance
		Append(f.registerNodesMetadataLayer(nodeIds, false)). // register cluster metadata
		Append(f.pingDroneLayer(nodeIds, false)).             // ping drone
		Append(f.setDroneConfigLayer(nodeIds, false)).        // set drone config
		Append(f.startConfdServiceLayer(nodeIds, false)).     // start confd service
		Append(f.startServiceLayer(nodeIds, false)).          // register start cmd to exec
		Append(f.deregisterCmdLayer(nodeIds, true))           // deregister cmd
}

func (f *Frame) ResizeClusterLayer(roleResizeResources models.RoleResizeResources) *models.TaskLayer {
	headTaskLayer := new(models.TaskLayer)
	for _, roleResizeResource := range roleResizeResources {
		var nodeIds []string
		for nodeId, node := range f.ClusterWrapper.ClusterNodesWithKeyPairs {
			if node.Role == roleResizeResource.Role {
				nodeIds = append(nodeIds, nodeId)
			}
		}
		if f.ClusterWrapper.ClusterCommons[roleResizeResource.Role].VerticalScalingPolicy == constants.ScalingPolicySequential {
			for _, nodeId := range nodeIds {
				f.resizeClusterLayer(headTaskLayer, []string{nodeId}, roleResizeResource)
			}
		} else {
			f.resizeClusterLayer(headTaskLayer, nodeIds, roleResizeResource)
		}
	}

	return headTaskLayer.Child
}

func (f *Frame) AttachKeyPairsLayer(nodeKeyPairDetails models.NodeKeyPairDetails) *models.TaskLayer {
	headTaskLayer := new(models.TaskLayer)

	for _, nodeKeyPairDetail := range nodeKeyPairDetails {
		headTaskLayer.Append(f.attachKeyPairLayer(&nodeKeyPairDetail))
	}

	return headTaskLayer.Child
}

func (f *Frame) DetachKeyPairsLayer(nodeKeyPairDetails models.NodeKeyPairDetails) *models.TaskLayer {
	headTaskLayer := new(models.TaskLayer)

	for _, nodeKeyPairDetail := range nodeKeyPairDetails {
		headTaskLayer.Append(f.detachKeyPairLayer(&nodeKeyPairDetail))
	}

	return headTaskLayer.Child
}
