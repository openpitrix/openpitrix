// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package vmbased

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	clusterclient "openpitrix.io/openpitrix/pkg/client/cluster"
	runtimeclient "openpitrix.io/openpitrix/pkg/client/runtime"
	"openpitrix.io/openpitrix/pkg/constants"
	"openpitrix.io/openpitrix/pkg/devkit/app"
	"openpitrix.io/openpitrix/pkg/logger"
	"openpitrix.io/openpitrix/pkg/models"
	"openpitrix.io/openpitrix/pkg/pb"
	"openpitrix.io/openpitrix/pkg/pi"
	"openpitrix.io/openpitrix/pkg/util/jsonutil"
	"openpitrix.io/openpitrix/pkg/util/pbutil"
	"openpitrix.io/openpitrix/pkg/util/sshutil"
)

type Frame struct {
	Job            *models.Job
	ClusterWrapper *models.ClusterWrapper
	Runtime        *runtimeclient.Runtime
}

func (f *Frame) startConfdServiceLayer(nodeIds []string) *models.TaskLayer {
	taskLayer := new(models.TaskLayer)
	for _, nodeId := range nodeIds {
		clusterNode := f.ClusterWrapper.ClusterNodes[nodeId]
		meta := &models.Meta{
			FrontgateId: f.ClusterWrapper.Cluster.FrontgateId,
			Timeout:     TimeoutStartConfd,
			NodeId:      clusterNode.NodeId,
			DroneIp:     clusterNode.PrivateIp,
		}
		directive, err := meta.ToString()
		if err != nil {
			return nil
		}
		startConfdTask := &models.Task{
			JobId:      f.Job.JobId,
			Owner:      f.Job.Owner,
			TaskAction: ActionStartConfd,
			Target:     constants.TargetPilot,
			NodeId:     nodeId,
			Directive:  string(directive),
		}
		taskLayer.Tasks = append(taskLayer.Tasks, startConfdTask)
	}
	if len(taskLayer.Tasks) > 0 {
		return taskLayer
	} else {
		return nil
	}
}

func (f *Frame) stopConfdServiceLayer(nodeIds []string) *models.TaskLayer {
	taskLayer := new(models.TaskLayer)
	for _, nodeId := range nodeIds {
		clusterNode := f.ClusterWrapper.ClusterNodes[nodeId]
		meta := &models.Meta{
			FrontgateId: f.ClusterWrapper.Cluster.FrontgateId,
			Timeout:     TimeoutStopConfd,
			NodeId:      clusterNode.NodeId,
			DroneIp:     clusterNode.PrivateIp,
		}
		directive, err := meta.ToString()
		if err != nil {
			return nil
		}
		stopConfdTask := &models.Task{
			JobId:      f.Job.JobId,
			Owner:      f.Job.Owner,
			TaskAction: ActionStopConfd,
			Target:     constants.TargetPilot,
			NodeId:     nodeId,
			Directive:  string(directive),
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
		role := f.ClusterWrapper.ClusterNodes[nodeId].Role
		serviceStr := f.ClusterWrapper.ClusterCommons[role].InitService
		if serviceStr != "" {
			service := app.Service{}
			err := json.Unmarshal([]byte(serviceStr), &service)
			if err != nil {
				logger.Error("Unmarshal cluster [%s] init service failed: %+v",
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
		role := f.ClusterWrapper.ClusterNodes[nodeId].Role
		serviceStr := f.ClusterWrapper.ClusterCommons[role].DestroyService
		if serviceStr != "" {
			service := app.Service{}
			err := json.Unmarshal([]byte(serviceStr), &service)
			if err != nil {
				logger.Error("Unmarshal cluster [%s] init service failed: %+v",
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

func (f *Frame) deregisterCmdLayer(nodeIds []string) *models.TaskLayer {
	taskLayer := new(models.TaskLayer)
	for _, nodeId := range nodeIds {
		cnodes := GetCmdCnodes(
			f.ClusterWrapper.Cluster.ClusterId,
			f.ClusterWrapper.ClusterNodes[nodeId].InstanceId,
			"",
		)
		meta := &models.Meta{
			FrontgateId: f.ClusterWrapper.Cluster.FrontgateId,
			NodeId:      nodeId,
			DroneIp:     f.ClusterWrapper.ClusterNodes[nodeId].PrivateIp,
			Timeout:     TimeoutDeregister,
			Cnodes:      cnodes,
		}
		directive, err := meta.ToString()
		if err != nil {
			return nil
		}
		deregisterCmdTask := &models.Task{
			JobId:      f.Job.JobId,
			Owner:      f.Job.Owner,
			TaskAction: ActionDeregesterCmd,
			Target:     constants.TargetPilot,
			NodeId:     nodeId,
			Directive:  string(directive),
		}
		taskLayer.Tasks = append(taskLayer.Tasks, deregisterCmdTask)
	}

	if len(taskLayer.Tasks) > 0 {
		return taskLayer
	} else {
		return nil
	}
}

func (f *Frame) registerCmdLayer(nodeIds []string, serviceName string) *models.TaskLayer {
	taskLayer := new(models.TaskLayer)
	for _, nodeId := range nodeIds {
		role := f.ClusterWrapper.ClusterNodes[nodeId].Role
		serviceStr := f.ClusterWrapper.GetCommonAttribute(role, serviceName)
		if serviceStr != nil {
			service := app.Service{}
			err := json.Unmarshal([]byte(serviceStr.(string)), &service)
			if err != nil {
				logger.Error("Unmarshal cluster [%s] service [%s] failed: %+v",
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
			cnodes := GetCmdCnodes(
				f.ClusterWrapper.Cluster.ClusterId,
				f.ClusterWrapper.ClusterNodes[nodeId].InstanceId,
				service.Cmd,
			)
			meta := &models.Meta{
				FrontgateId: f.ClusterWrapper.Cluster.FrontgateId,
				NodeId:      nodeId,
				DroneIp:     f.ClusterWrapper.ClusterNodes[nodeId].PrivateIp,
				Timeout:     timeout,
				Cnodes:      cnodes,
			}
			directive, err := meta.ToString()
			if err != nil {
				return nil
			}
			registerCmdTask := &models.Task{
				JobId:      f.Job.JobId,
				Owner:      f.Job.Owner,
				TaskAction: ActionRegisterCmd,
				Target:     constants.TargetPilot,
				NodeId:     nodeId,
				Directive:  string(directive),
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

func (f *Frame) constructServiceTasks(serviceName, cmdName string, nodeIds []string,
	serviceParams map[string]interface{}) *models.TaskLayer {
	headTaskLayer := new(models.TaskLayer)
	if len(nodeIds) == 0 {
		return nil
	}

	roleNodeIds := make(map[string][]string)
	nodeIdRole := make(map[string]string)
	for _, nodeId := range nodeIds {
		clusterNode, exist := f.ClusterWrapper.ClusterNodes[nodeId]
		if !exist {
			logger.Error("ClusterConf [%s] node [%s] not exist", f.ClusterWrapper.Cluster.ClusterId, nodeId)
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
	roleService := make(map[string]app.Service)
	for role, nodes := range roleNodeIds {
		serviceStr := f.ClusterWrapper.GetCommonAttribute(role, serviceName)
		if serviceStr == nil {
			return nil
		}
		service := app.Service{}
		err := json.Unmarshal([]byte(serviceStr.(string)), &service)
		if err != nil {
			logger.Error("Unmarshal cluster [%s] service [%s] failed: %+v",
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
		taskLayer := f.registerCmdLayer(nodeIds, serviceName)
		headTaskLayer.Leaf().Child = taskLayer
	}
	return headTaskLayer.Child
}

func (f *Frame) initServiceLayer(nodeIds []string) *models.TaskLayer {
	return f.constructServiceTasks("InitService", constants.ServiceCmdName, nodeIds, nil)
}

func (f *Frame) startServiceLayer(nodeIds []string) *models.TaskLayer {
	return f.constructServiceTasks("StartService", constants.ServiceCmdName, nodeIds, nil)
}

func (f *Frame) stopServiceLayer(nodeIds []string) *models.TaskLayer {
	return f.constructServiceTasks("StopService", constants.ServiceCmdName, nodeIds, nil)
}

func (f *Frame) scaleOutPreCheckServiceLayer(nodeIds []string) *models.TaskLayer {
	return f.constructServiceTasks("ScaleOutService", constants.ServiceCmdName, nodeIds, nil)
}

func (f *Frame) scaleInPreCheckServiceLayer(nodeIds []string) *models.TaskLayer {
	return f.constructServiceTasks("ScaleInService", constants.ServiceCmdName, nodeIds, nil)
}

func (f *Frame) scaleOutServiceLayer(nodeIds []string) *models.TaskLayer {
	return f.constructServiceTasks("ScaleOutService", constants.ServicePreCheckName, nodeIds, nil)
}

func (f *Frame) scaleInServiceLayer(nodeIds []string) *models.TaskLayer {
	return f.constructServiceTasks("ScaleInService", constants.ServicePreCheckName, nodeIds, nil)
}

func (f *Frame) destroyServiceLayer(nodeIds []string) *models.TaskLayer {
	return f.constructServiceTasks("DestroyService", constants.ServiceCmdName, nodeIds, nil)
}

func (f *Frame) initAndStartServiceLayer(nodeIds []string) *models.TaskLayer {
	headTaskLayer := new(models.TaskLayer)

	preStartNodes, postStartNodes := f.getPreAndPostStartGroupNodes(nodeIds)

	// Init service before start service
	headTaskLayer.Leaf().Child = f.initServiceLayer(preStartNodes)

	// TODO: custom metadata
	headTaskLayer.Leaf().Child = f.startServiceLayer(nodeIds)

	// Init service after start service
	headTaskLayer.Leaf().Child = f.initServiceLayer(postStartNodes)

	return headTaskLayer.Child
}

func (f *Frame) destroyAndStopServiceLayer(nodeIds []string, extraLayer *models.TaskLayer) *models.TaskLayer {
	headTaskLayer := new(models.TaskLayer)

	preStopNodes, postStopNodes := f.getPreAndPostStopGroupNodes(nodeIds)

	// Destroy service before stop service
	headTaskLayer.Leaf().Child = f.destroyServiceLayer(preStopNodes)

	if extraLayer != nil {
		headTaskLayer.Leaf().Child = extraLayer
	}

	headTaskLayer.Leaf().Child = f.stopServiceLayer(nodeIds)

	// Destroy service after stop service
	headTaskLayer.Leaf().Child = f.destroyServiceLayer(postStopNodes)

	return headTaskLayer.Child
}

func (f *Frame) createVolumesLayer(nodeIds []string) *models.TaskLayer {
	taskLayer := new(models.TaskLayer)
	for _, nodeId := range nodeIds {
		clusterNode := f.ClusterWrapper.ClusterNodes[nodeId]
		role := clusterNode.Role
		if strings.HasSuffix(role, constants.ReplicaRoleSuffix) {
			role = string([]byte(role)[:len(role)-len(constants.ReplicaRoleSuffix)])
		}
		clusterRole, exist := f.ClusterWrapper.ClusterRoles[role]
		if !exist {
			logger.Error("No such role [%s] in cluster role [%s]. ",
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
				Zone:      f.Runtime.Zone,
				RuntimeId: f.Runtime.RuntimeId,
			}
			volumeTaskDirective, err := volume.ToString()
			if err != nil {
				return nil
			}

			createVolumesTask := &models.Task{
				JobId:      f.Job.JobId,
				Owner:      f.Job.Owner,
				TaskAction: ActionCreateVolumes,
				Target:     f.Runtime.Provider,
				NodeId:     nodeId,
				Directive:  volumeTaskDirective,
			}
			for range mountPoints {
				taskLayer.Tasks = append(taskLayer.Tasks, createVolumesTask)
			}
		}
	}
	return taskLayer
}

func (f *Frame) detachVolumesLayer(nodeIds []string) *models.TaskLayer {
	taskLayer := new(models.TaskLayer)
	for _, nodeId := range nodeIds {
		clusterNode := f.ClusterWrapper.ClusterNodes[nodeId]
		volume := &models.Volume{
			Name:       clusterNode.ClusterId + "_" + nodeId,
			Zone:       f.Runtime.Zone,
			RuntimeId:  f.Runtime.RuntimeId,
			VolumeId:   clusterNode.VolumeId,
			InstanceId: clusterNode.InstanceId,
		}
		directive, err := volume.ToString()
		if err != nil {
			return nil
		}
		detachVolumesTask := &models.Task{
			JobId:      f.Job.JobId,
			Owner:      f.Job.Owner,
			TaskAction: ActionDetachVolumes,
			Target:     f.Runtime.Provider,
			NodeId:     nodeId,
			Directive:  directive,
		}
		taskLayer.Tasks = append(taskLayer.Tasks, detachVolumesTask)
	}
	return taskLayer
}

func (f *Frame) attachVolumesLayer() *models.TaskLayer {
	taskLayer := new(models.TaskLayer)
	for nodeId, clusterNode := range f.ClusterWrapper.ClusterNodes {
		volume := &models.Volume{
			Name:       clusterNode.ClusterId + "_" + nodeId,
			Zone:       f.Runtime.Zone,
			RuntimeId:  f.Runtime.RuntimeId,
			VolumeId:   clusterNode.VolumeId,
			InstanceId: clusterNode.InstanceId,
		}
		directive, err := volume.ToString()
		if err != nil {
			return nil
		}
		attachVolumesTask := &models.Task{
			JobId:      f.Job.JobId,
			Owner:      f.Job.Owner,
			TaskAction: ActionAttachVolumes,
			Target:     f.Runtime.Provider,
			NodeId:     nodeId,
			Directive:  directive,
		}
		taskLayer.Tasks = append(taskLayer.Tasks, attachVolumesTask)
	}
	return taskLayer
}

func (f *Frame) deleteVolumesLayer(nodeIds []string) *models.TaskLayer {
	taskLayer := new(models.TaskLayer)
	for _, nodeId := range nodeIds {
		clusterNode := f.ClusterWrapper.ClusterNodes[nodeId]
		volume := &models.Volume{
			Name:       clusterNode.ClusterId + "_" + nodeId,
			Zone:       f.Runtime.Zone,
			RuntimeId:  f.Runtime.RuntimeId,
			VolumeId:   clusterNode.VolumeId,
			InstanceId: clusterNode.InstanceId,
		}
		directive, err := volume.ToString()
		if err != nil {
			return nil
		}
		deleteVolumesTask := &models.Task{
			JobId:      f.Job.JobId,
			Owner:      f.Job.Owner,
			TaskAction: ActionDeleteVolumes,
			Target:     f.Runtime.Provider,
			NodeId:     nodeId,
			Directive:  directive,
		}
		taskLayer.Tasks = append(taskLayer.Tasks, deleteVolumesTask)
	}
	return taskLayer
}

func (f *Frame) formatAndMountVolumeLayer(nodeIds []string) *models.TaskLayer {
	taskLayer := new(models.TaskLayer)

	for _, nodeId := range nodeIds {
		clusterNode := f.ClusterWrapper.ClusterNodes[nodeId]
		// cmd will be assigned when the task is handling
		meta := &models.Meta{
			FrontgateId: f.ClusterWrapper.Cluster.FrontgateId,
			Timeout:     TimeoutFormatAndMountVolume,
			NodeId:      clusterNode.NodeId,
			DroneIp:     clusterNode.PrivateIp,
		}
		directive, err := meta.ToString()
		if err != nil {
			return nil
		}
		formatVolumeTask := &models.Task{
			JobId:      f.Job.JobId,
			Owner:      f.Job.Owner,
			TaskAction: ActionFormatAndMountVolume,
			Target:     constants.TargetPilot,
			NodeId:     nodeId,
			Directive:  string(directive),
		}
		taskLayer.Tasks = append(taskLayer.Tasks, formatVolumeTask)
	}
	if len(taskLayer.Tasks) > 0 {
		return taskLayer
	} else {
		return nil
	}
}

func (f *Frame) sshKeygenLayer() *models.TaskLayer {
	taskLayer := new(models.TaskLayer)
	ctx := context.Background()
	client, err := clusterclient.NewClusterManagerClient(ctx)
	if err != nil {
		logger.Error("New ssh key gen task layer failed: %+v", err)
		return nil
	}

	for nodeId, clusterNode := range f.ClusterWrapper.ClusterNodes {
		role := clusterNode.Role
		clusterCommon := f.ClusterWrapper.ClusterCommons[role]
		keyType := clusterCommon.Passphraseless
		if keyType != "" {
			private, public, err := sshutil.MakeSSHKeyPair(keyType)
			if err != nil {
				logger.Error("Generate ssh key [%s] in cluster node [%s] failed",
					clusterCommon.Passphraseless, nodeId)
				return nil
			}
			_, err = client.ModifyClusterNode(ctx, &pb.ModifyClusterNodeRequest{
				ClusterNode: &pb.ClusterNode{
					NodeId: pbutil.ToProtoString(nodeId),
					PubKey: pbutil.ToProtoString(public),
				},
			})
			cmd := fmt.Sprintf("mkdir -p /root/.ssh/;chmod 700 /root/.ssh/;"+
				"echo \"%s\" > /root/.ssh/id_%s;echo \"%s\" > /root/.ssh/id_%s.pub;"+
				"chown 600 /root/.ssh/id_%s;chown 644 /root/.ssh/id_%s.pub",
				private, keyType, public, keyType, keyType, keyType)
			cnodes := GetCmdCnodes(
				f.ClusterWrapper.Cluster.ClusterId,
				f.ClusterWrapper.ClusterNodes[nodeId].InstanceId,
				cmd,
			)
			meta := &models.Meta{
				FrontgateId: f.ClusterWrapper.Cluster.FrontgateId,
				Timeout:     TimeoutSshKeygen,
				NodeId:      clusterNode.NodeId,
				DroneIp:     clusterNode.PrivateIp,
				Cnodes:      cnodes,
			}
			directive, err := meta.ToString()
			if err != nil {
				return nil
			}
			formatVolumeTask := &models.Task{
				JobId:      f.Job.JobId,
				Owner:      f.Job.Owner,
				TaskAction: ActionRegisterCmd,
				Target:     constants.TargetPilot,
				NodeId:     nodeId,
				Directive:  string(directive),
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

func (f *Frame) umountVolumeLayer(nodeIds []string) *models.TaskLayer {
	taskLayer := new(models.TaskLayer)

	for _, nodeId := range nodeIds {
		clusterNode := f.ClusterWrapper.ClusterNodes[nodeId]
		clusterRole := f.ClusterWrapper.ClusterRoles[clusterNode.Role]
		cmd := UmountVolumeCmd(clusterRole.MountPoint)
		cnodes := GetCmdCnodes(
			f.ClusterWrapper.Cluster.ClusterId,
			f.ClusterWrapper.ClusterNodes[nodeId].InstanceId,
			cmd,
		)
		meta := &models.Meta{
			FrontgateId: f.ClusterWrapper.Cluster.FrontgateId,
			Timeout:     TimeoutUmountVolume,
			NodeId:      clusterNode.NodeId,
			DroneIp:     clusterNode.PrivateIp,
			Cnodes:      cnodes,
		}
		directive, err := meta.ToString()
		if err != nil {
			return nil
		}
		umountVolumeTask := &models.Task{
			JobId:      f.Job.JobId,
			Owner:      f.Job.Owner,
			TaskAction: ActionRegisterCmd,
			Target:     constants.TargetPilot,
			NodeId:     nodeId,
			Directive:  string(directive),
		}
		taskLayer.Tasks = append(taskLayer.Tasks, umountVolumeTask)
	}
	if len(taskLayer.Tasks) > 0 {
		return taskLayer
	} else {
		return nil
	}
}

func (f *Frame) runInstancesLayer(nodeIds []string) *models.TaskLayer {
	taskLayer := new(models.TaskLayer)
	apiServer := f.Runtime.RuntimeUrl
	zone := f.Runtime.Zone
	globalPi := pi.Global()
	imageId := ""
	var err error
	if globalPi != nil {
		imageId, err = globalPi.GlobalConfig().GetRuntimeImageId(apiServer, zone)
		if err != nil {
			return nil
		}
	}
	for _, nodeId := range nodeIds {
		clusterNode := f.ClusterWrapper.ClusterNodes[nodeId]
		role := clusterNode.Role
		if strings.HasSuffix(role, constants.ReplicaRoleSuffix) {
			role = string([]byte(role)[:len(role)-len(constants.ReplicaRoleSuffix)])
		}
		clusterRole, exist := f.ClusterWrapper.ClusterRoles[role]
		if !exist {
			logger.Error("No such role [%s] in cluster role [%s]. ",
				role, f.ClusterWrapper.Cluster.ClusterId)
			return nil
		}
		instance := &models.Instance{
			Name:      clusterNode.ClusterId + "_" + nodeId,
			NodeId:    nodeId,
			ImageId:   imageId,
			Cpu:       int(clusterRole.Cpu),
			Memory:    int(clusterRole.Memory),
			Gpu:       int(clusterRole.Gpu),
			Subnet:    clusterNode.SubnetId,
			RuntimeId: f.Runtime.RuntimeId,
			Zone:      f.Runtime.Zone,
		}
		instanceTaskDirective, err := instance.ToString()
		if err != nil {
			return nil
		}
		runInstanceTask := &models.Task{
			JobId:      f.Job.JobId,
			Owner:      f.Job.Owner,
			TaskAction: ActionRunInstances,
			Target:     f.Runtime.Provider,
			NodeId:     nodeId,
			Directive:  instanceTaskDirective,
		}
		taskLayer.Tasks = append(taskLayer.Tasks, runInstanceTask)
	}

	if len(taskLayer.Tasks) > 0 {
		return taskLayer
	} else {
		return nil
	}
}

func (f *Frame) stopInstancesLayer(nodeIds []string) *models.TaskLayer {
	taskLayer := new(models.TaskLayer)
	for _, nodeId := range nodeIds {
		clusterNode := f.ClusterWrapper.ClusterNodes[nodeId]
		instance := &models.Instance{
			Name:       clusterNode.ClusterId + "_" + nodeId,
			NodeId:     nodeId,
			InstanceId: clusterNode.InstanceId,
			RuntimeId:  f.Runtime.RuntimeId,
			Zone:       f.Runtime.Zone,
		}
		instanceTaskDirective, err := instance.ToString()
		if err != nil {
			return nil
		}
		stopInstanceTask := &models.Task{
			JobId:      f.Job.JobId,
			Owner:      f.Job.Owner,
			TaskAction: ActionStopInstances,
			Target:     f.Runtime.Provider,
			NodeId:     nodeId,
			Directive:  instanceTaskDirective,
		}
		taskLayer.Tasks = append(taskLayer.Tasks, stopInstanceTask)
	}

	if len(taskLayer.Tasks) > 0 {
		return taskLayer
	} else {
		return nil
	}
}

func (f *Frame) deleteInstancesLayer(nodeIds []string) *models.TaskLayer {
	taskLayer := new(models.TaskLayer)
	for _, nodeId := range nodeIds {
		clusterNode := f.ClusterWrapper.ClusterNodes[nodeId]
		instance := &models.Instance{
			Name:       clusterNode.ClusterId + "_" + nodeId,
			NodeId:     nodeId,
			InstanceId: clusterNode.InstanceId,
			RuntimeId:  f.Runtime.RuntimeId,
			Zone:       f.Runtime.Zone,
		}
		instanceTaskDirective, err := instance.ToString()
		if err != nil {
			return nil
		}
		deleteInstanceTask := &models.Task{
			JobId:      f.Job.JobId,
			Owner:      f.Job.Owner,
			TaskAction: ActionTerminateInstances,
			Target:     f.Runtime.Provider,
			NodeId:     nodeId,
			Directive:  instanceTaskDirective,
		}
		taskLayer.Tasks = append(taskLayer.Tasks, deleteInstanceTask)
	}

	if len(taskLayer.Tasks) > 0 {
		return taskLayer
	} else {
		return nil
	}
}

func (f *Frame) startInstancesLayer() *models.TaskLayer {
	taskLayer := new(models.TaskLayer)
	for nodeId, clusterNode := range f.ClusterWrapper.ClusterNodes {
		instance := &models.Instance{
			Name:       clusterNode.ClusterId + "_" + nodeId,
			NodeId:     nodeId,
			InstanceId: clusterNode.InstanceId,
			RuntimeId:  f.Runtime.RuntimeId,
			Zone:       f.Runtime.Zone,
		}
		instanceTaskDirective, err := instance.ToString()
		if err != nil {
			return nil
		}
		startInstanceTask := &models.Task{
			JobId:      f.Job.JobId,
			Owner:      f.Job.Owner,
			TaskAction: ActionStartInstances,
			Target:     f.Runtime.Provider,
			NodeId:     nodeId,
			Directive:  instanceTaskDirective,
		}
		taskLayer.Tasks = append(taskLayer.Tasks, startInstanceTask)
	}

	if len(taskLayer.Tasks) > 0 {
		return taskLayer
	} else {
		return nil
	}
}

func (f *Frame) waitFrontgateLayer() *models.TaskLayer {
	meta := &models.Meta{
		FrontgateId: f.ClusterWrapper.Cluster.FrontgateId,
	}
	directive, err := meta.ToString()
	if err != nil {
		return nil
	}
	// Wait frontgate available
	waitFrontgateTask := &models.Task{
		JobId:      f.Job.JobId,
		Owner:      f.Job.Owner,
		TaskAction: ActionWaitFrontgateAvailable,
		Target:     f.Runtime.Provider,
		NodeId:     f.ClusterWrapper.Cluster.ClusterId,
		Directive:  string(directive),
	}
	return &models.TaskLayer{
		Tasks: []*models.Task{waitFrontgateTask},
	}
}

func (f *Frame) registerMetadataLayer() *models.TaskLayer {
	// When the task is handled by task controller, the cnodes will be filled in,
	meta := &models.Meta{
		FrontgateId: f.ClusterWrapper.Cluster.FrontgateId,
		Timeout:     TimeoutRegister,
		ClusterId:   f.ClusterWrapper.Cluster.ClusterId,
		Cnodes:      "",
	}
	directive, err := meta.ToString()
	if err != nil {
		return nil
	}
	registerMetadataTask := &models.Task{
		JobId:      f.Job.JobId,
		Owner:      f.Job.Owner,
		TaskAction: ActionRegisterMetadata,
		Target:     constants.TargetPilot,
		NodeId:     f.ClusterWrapper.Cluster.ClusterId,
		Directive:  string(directive),
	}
	return &models.TaskLayer{
		Tasks: []*models.Task{registerMetadataTask},
	}
}

func (f *Frame) registerNodesMetadataLayer(nodeIds []string) *models.TaskLayer {
	metadata := &Metadata{
		ClusterWrapper: f.ClusterWrapper,
	}
	cnodes := jsonutil.ToString(metadata.GetClusterNodeCnodes(nodeIds))
	meta := &models.Meta{
		FrontgateId: f.ClusterWrapper.Cluster.FrontgateId,
		Timeout:     TimeoutRegister,
		ClusterId:   f.ClusterWrapper.Cluster.ClusterId,
		Cnodes:      cnodes,
	}
	directive, err := meta.ToString()
	if err != nil {
		return nil
	}
	task := &models.Task{
		JobId:      f.Job.JobId,
		Owner:      f.Job.Owner,
		TaskAction: ActionRegisterMetadata,
		Target:     constants.TargetPilot,
		NodeId:     f.ClusterWrapper.Cluster.ClusterId,
		Directive:  string(directive),
	}
	return &models.TaskLayer{
		Tasks: []*models.Task{task},
	}
}

func (f *Frame) registerScalingNodesMetadataLayer(nodeIds []string, path string) *models.TaskLayer {
	clusterId := f.ClusterWrapper.Cluster.ClusterId
	metadata := &Metadata{
		ClusterWrapper: f.ClusterWrapper,
	}
	hosts := metadata.GetHostsCnodes(nodeIds)
	if len(hosts) == 0 {
		logger.Info("No new nodes for cluster [%s] is registered", clusterId)
		return nil
	}
	cnodes := map[string]interface{}{
		RegisterClustersRootPath: map[string]interface{}{
			clusterId: map[string]interface{}{
				path: hosts,
			},
		},
	}
	meta := &models.Meta{
		FrontgateId: f.ClusterWrapper.Cluster.FrontgateId,
		Timeout:     TimeoutRegister,
		ClusterId:   f.ClusterWrapper.Cluster.ClusterId,
		Cnodes:      jsonutil.ToString(cnodes),
	}
	directive, err := meta.ToString()
	if err != nil {
		return nil
	}
	task := &models.Task{
		JobId:      f.Job.JobId,
		Owner:      f.Job.Owner,
		TaskAction: ActionRegisterMetadata,
		Target:     constants.TargetPilot,
		NodeId:     f.ClusterWrapper.Cluster.ClusterId,
		Directive:  string(directive),
	}
	return &models.TaskLayer{
		Tasks: []*models.Task{task},
	}
}

func (f *Frame) deregisterNodesMetadataLayer(nodeIds []string) *models.TaskLayer {
	metadata := &Metadata{
		ClusterWrapper: f.ClusterWrapper,
	}
	cnodes := jsonutil.ToString(metadata.GetEmptyClusterNodeCnodes(nodeIds))
	meta := &models.Meta{
		FrontgateId: f.ClusterWrapper.Cluster.FrontgateId,
		Timeout:     TimeoutDeregister,
		ClusterId:   f.ClusterWrapper.Cluster.ClusterId,
		Cnodes:      cnodes,
	}
	directive, err := meta.ToString()
	if err != nil {
		return nil
	}
	task := &models.Task{
		JobId:      f.Job.JobId,
		Owner:      f.Job.Owner,
		TaskAction: ActionDeregisterMetadata,
		Target:     constants.TargetPilot,
		NodeId:     f.ClusterWrapper.Cluster.ClusterId,
		Directive:  string(directive),
	}
	return &models.TaskLayer{
		Tasks: []*models.Task{task},
	}
}

func (f *Frame) deregisterScalingNodesMetadataLayer(path string) *models.TaskLayer {
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
	directive, err := meta.ToString()
	if err != nil {
		return nil
	}
	task := &models.Task{
		JobId:      f.Job.JobId,
		Owner:      f.Job.Owner,
		TaskAction: ActionDeregisterMetadata,
		Target:     constants.TargetPilot,
		NodeId:     f.ClusterWrapper.Cluster.ClusterId,
		Directive:  string(directive),
	}
	return &models.TaskLayer{
		Tasks: []*models.Task{task},
	}
}

func (f *Frame) deregisterMetadataLayer() *models.TaskLayer {
	cnodes := GetEmptyClusterCnodes(f.ClusterWrapper.Cluster.ClusterId)
	meta := &models.Meta{
		FrontgateId: f.ClusterWrapper.Cluster.FrontgateId,
		Timeout:     TimeoutDeregister,
		ClusterId:   f.ClusterWrapper.Cluster.ClusterId,
		Cnodes:      cnodes,
	}
	directive, err := meta.ToString()
	if err != nil {
		return nil
	}
	deregisterMetadataTask := &models.Task{
		JobId:      f.Job.JobId,
		Owner:      f.Job.Owner,
		TaskAction: ActionDeregisterMetadata,
		Target:     constants.TargetPilot,
		NodeId:     f.ClusterWrapper.Cluster.ClusterId,
		Directive:  directive,
	}
	return &models.TaskLayer{
		Tasks: []*models.Task{deregisterMetadataTask},
	}
}

func (f *Frame) CreateClusterLayer() *models.TaskLayer {
	var nodeIds []string
	for nodeId := range f.ClusterWrapper.ClusterNodes {
		nodeIds = append(nodeIds, nodeId)
	}
	headTaskLayer := new(models.TaskLayer)

	headTaskLayer.
		Append(f.createVolumesLayer(nodeIds)).        // create volume
		Append(f.runInstancesLayer(nodeIds)).         // run instance and attach volume to instance
		Append(f.formatAndMountVolumeLayer(nodeIds)). // format and mount volume to instance
		Append(f.waitFrontgateLayer()).               // wait frontgate cluster to be active
		Append(f.sshKeygenLayer()).                   // generate ssh key
		Append(f.registerMetadataLayer()).            // register cluster metadata
		Append(f.startConfdServiceLayer(nodeIds)).    // start confd service
		Append(f.initAndStartServiceLayer(nodeIds)).  // register init and start cmd to exec
		Append(f.deregisterCmdLayer(nodeIds))         // deregister cmd

	return headTaskLayer.Child
}

func (f *Frame) StopClusterLayer() *models.TaskLayer {
	var nodeIds []string
	for nodeId := range f.ClusterWrapper.ClusterNodes {
		nodeIds = append(nodeIds, nodeId)
	}
	headTaskLayer := new(models.TaskLayer)

	headTaskLayer.
		Append(f.waitFrontgateLayer()).           // wait frontgate cluster to be active
		Append(f.stopServiceLayer(nodeIds)).      // register stop cmd to exec
		Append(f.stopConfdServiceLayer(nodeIds)). // stop confd service
		Append(f.umountVolumeLayer(nodeIds)).     // umount volume from instance
		Append(f.detachVolumesLayer(nodeIds)).    // detach volume from instance
		Append(f.stopInstancesLayer(nodeIds)).    // stop instance
		Append(f.deregisterMetadataLayer())       // deregister cluster

	return headTaskLayer.Child
}

func (f *Frame) StartClusterLayer() *models.TaskLayer {
	var nodeIds []string
	for nodeId := range f.ClusterWrapper.ClusterNodes {
		nodeIds = append(nodeIds, nodeId)
	}
	headTaskLayer := new(models.TaskLayer)

	headTaskLayer.
		Append(f.startInstancesLayer()).           // start instance
		Append(f.attachVolumesLayer()).            // attach volume to instance, will auto mount
		Append(f.waitFrontgateLayer()).            // wait frontgate cluster to be active
		Append(f.registerMetadataLayer()).         // register cluster metadata
		Append(f.startConfdServiceLayer(nodeIds)). // start confd service
		Append(f.startServiceLayer(nodeIds)).      // register start cmd to exec
		Append(f.deregisterCmdLayer(nodeIds))      // deregister cmd

	return headTaskLayer.Child
}

func (f *Frame) DeleteClusterLayer() *models.TaskLayer {
	var nodeIds []string
	for nodeId := range f.ClusterWrapper.ClusterNodes {
		nodeIds = append(nodeIds, nodeId)
	}
	headTaskLayer := new(models.TaskLayer)

	headTaskLayer.
		Append(f.destroyAndStopServiceLayer(nodeIds, nil)). // register destroy and stop cmd to exec
		Append(f.stopConfdServiceLayer(nodeIds)).           // stop confd service
		Append(f.umountVolumeLayer(nodeIds)).               // umount volume from instance
		Append(f.detachVolumesLayer(nodeIds)).              // detach volume from instance
		Append(f.deleteInstancesLayer(nodeIds)).            // delete instance
		Append(f.deleteVolumesLayer(nodeIds)).              // delete volume
		Append(f.deregisterMetadataLayer())                 // deregister cluster

	return headTaskLayer.Child
}

func (f *Frame) AddClusterNodesLayer() *models.TaskLayer {
	var addNodeIds, nonAddNodeIds []string
	for nodeId, node := range f.ClusterWrapper.ClusterNodes {
		if node.Status == constants.StatusPending {
			addNodeIds = append(addNodeIds, nodeId)
		} else {
			nonAddNodeIds = append(nonAddNodeIds, nodeId)
		}
	}
	headTaskLayer := new(models.TaskLayer)

	headTaskLayer.
		Append(f.scaleOutPreCheckServiceLayer(nonAddNodeIds)).                       // register scale out pre check to exec
		Append(f.createVolumesLayer(addNodeIds)).                                    // create volume
		Append(f.runInstancesLayer(addNodeIds)).                                     // run instance and attach volume to instance
		Append(f.formatAndMountVolumeLayer(addNodeIds)).                             // format and mount volume to instance
		Append(f.registerNodesMetadataLayer(addNodeIds)).                            // register cluster nodes metadata
		Append(f.registerScalingNodesMetadataLayer(addNodeIds, RegisterNodeAdding)). // register adding hosts metadata
		Append(f.startConfdServiceLayer(addNodeIds)).                                // start confd service
		Append(f.initAndStartServiceLayer(addNodeIds)).                              // register init and start cmd to exec
		Append(f.scaleOutServiceLayer(nonAddNodeIds)).                               // register scale out cmd to exec
		Append(f.deregisterScalingNodesMetadataLayer(RegisterNodeAdding))            // deregister adding host metadata
	return headTaskLayer.Child
}

func (f *Frame) DeleteClusterNodesLayer() *models.TaskLayer {
	var deleteNodeIds, nonDeleteNodeIds []string
	for nodeId, node := range f.ClusterWrapper.ClusterNodes {
		if node.Status == constants.StatusDeleting {
			deleteNodeIds = append(deleteNodeIds, nodeId)
		} else {
			nonDeleteNodeIds = append(nonDeleteNodeIds, nodeId)
		}
	}
	headTaskLayer := new(models.TaskLayer)

	headTaskLayer.
		Append(f.registerScalingNodesMetadataLayer(deleteNodeIds, RegisterNodeDeleting)).             // register scale in node metadata
		Append(f.scaleInPreCheckServiceLayer(nonDeleteNodeIds)).                                      // register scale in pre check to exec
		Append(f.destroyAndStopServiceLayer(deleteNodeIds, f.scaleInServiceLayer(nonDeleteNodeIds))). // register destroy, scale in and stop cmd to exec
		Append(f.stopConfdServiceLayer(deleteNodeIds)).                                               // stop confd service
		Append(f.umountVolumeLayer(deleteNodeIds)).                                                   // umount volume from instance
		Append(f.detachVolumesLayer(deleteNodeIds)).                                                  // detach volume from instance
		Append(f.deleteInstancesLayer(deleteNodeIds)).                                                // delete instance
		Append(f.deleteVolumesLayer(deleteNodeIds)).                                                  // delete volume
		Append(f.deregisterNodesMetadataLayer(deleteNodeIds)).                                        // deregister deleting cluster nodes metadata
		Append(f.deregisterScalingNodesMetadataLayer(RegisterNodeDeleting))                           // deregister deleting nodes metadata
	return headTaskLayer.Child
}
