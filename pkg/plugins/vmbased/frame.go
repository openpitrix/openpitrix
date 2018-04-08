// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package vmbased

import (
	"encoding/json"
	"sort"
	"strings"

	runtimeclient "openpitrix.io/openpitrix/pkg/client/runtime"
	"openpitrix.io/openpitrix/pkg/constants"
	"openpitrix.io/openpitrix/pkg/logger"
	"openpitrix.io/openpitrix/pkg/models"
)

type Frame struct {
	Job            *models.Job
	ClusterWrapper *models.ClusterWrapper
	Runtime        *runtimeclient.Runtime
}

func NewFrame(job *models.Job) (*Frame, error) {
	clusterWrapper, err := models.NewClusterWrapper(job.Directive)
	if err != nil {
		return nil, err
	}

	runtimeId := clusterWrapper.Cluster.RuntimeId
	runtime, err := runtimeclient.NewRuntime(runtimeId)
	if err != nil {
		return nil, err
	}

	return &Frame{
		Job:            job,
		ClusterWrapper: clusterWrapper,
		Runtime:        runtime,
	}, nil
}

func (f *Frame) startConfdServiceLayer() *models.TaskLayer {
	startConfdTaskLayer := new(models.TaskLayer)
	for nodeId, clusterNode := range f.ClusterWrapper.ClusterNodes {
		meta := &models.Meta{
			FrontgateId: f.ClusterWrapper.Cluster.FrontgateId,
			Timeout:     TimeoutStartConfd,
			NodeId:      clusterNode.NodeId,
			Ip:          clusterNode.PrivateIp,
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
		startConfdTaskLayer.Tasks = append(startConfdTaskLayer.Tasks, startConfdTask)
	}
	if len(startConfdTaskLayer.Tasks) > 0 {
		return startConfdTaskLayer
	} else {
		return nil
	}
}

func (f *Frame) stopConfdServiceLayer() *models.TaskLayer {
	stopConfdTaskLayer := new(models.TaskLayer)
	for nodeId, clusterNode := range f.ClusterWrapper.ClusterNodes {
		meta := &models.Meta{
			FrontgateId: f.ClusterWrapper.Cluster.FrontgateId,
			Timeout:     TimeoutStopConfd,
			NodeId:      clusterNode.NodeId,
			Ip:          clusterNode.PrivateIp,
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
		stopConfdTaskLayer.Tasks = append(stopConfdTaskLayer.Tasks, stopConfdTask)
	}
	if len(stopConfdTaskLayer.Tasks) > 0 {
		return stopConfdTaskLayer
	} else {
		return nil
	}
}

// Put the nodes into two groups
func (f *Frame) getPreAndPostInitGroupNodes(nodeIds []string) ([]string, []string) {
	var preGroupNodes, postGroupNodes []string
	for _, nodeId := range nodeIds {
		role := f.ClusterWrapper.ClusterNodes[nodeId].Role
		serviceStr := f.ClusterWrapper.ClusterCommons[role].InitService
		if serviceStr != "" {
			service := models.Service{}
			err := json.Unmarshal([]byte(serviceStr), &service)
			if err != nil {
				logger.Errorf("Unmarshal cluster [%s] init service failed: %+v",
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

func (f *Frame) deregisterCmd(nodeIds []string) *models.TaskLayer {
	taskLayer := new(models.TaskLayer)
	for _, nodeId := range nodeIds {
		meta := &models.Meta{
			FrontgateId: f.ClusterWrapper.Cluster.FrontgateId,
			NodeId:      nodeId,
			Ip:          f.ClusterWrapper.ClusterNodes[nodeId].PrivateIp,
			Timeout:     TimeoutDeregister,
			Cmd:         GetDeregisterExec("cmd"),
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

func (f *Frame) registerCmd(nodeIds []string, serviceName string) *models.TaskLayer {
	taskLayer := new(models.TaskLayer)
	for _, nodeId := range nodeIds {
		role := f.ClusterWrapper.ClusterNodes[nodeId].Role
		serviceStr := f.ClusterWrapper.GetCommonAttribute(role, serviceName)
		if serviceStr != nil {
			service := models.Service{}
			err := json.Unmarshal([]byte(serviceStr.(string)), &service)
			if err != nil {
				logger.Errorf("Unmarshal cluster [%s] service [%s] failed: %+v",
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
			meta := &models.Meta{
				FrontgateId: f.ClusterWrapper.Cluster.FrontgateId,
				NodeId:      nodeId,
				Ip:          f.ClusterWrapper.ClusterNodes[nodeId].PrivateIp,
				Timeout:     timeout,
				Cmd:         GetRegisterExec(service.Cmd),
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
			logger.Errorf("Cluster [%s] node [%s] not exist", f.ClusterWrapper.Cluster.ClusterId, nodeId)
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
	roleService := make(map[string]models.Service)
	for role, nodes := range roleNodeIds {
		serviceStr := f.ClusterWrapper.GetCommonAttribute(role, serviceName)
		if serviceStr == nil {
			return nil
		}
		service := models.Service{}
		err := json.Unmarshal([]byte(serviceStr.(string)), &service)
		if err != nil {
			logger.Errorf("Unmarshal cluster [%s] service [%s] failed: %+v",
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
		taskLayer := f.registerCmd(nodeIds, serviceName)
		headTaskLayer.Leaf().Child = taskLayer
	}
	return headTaskLayer.Child
}

func (f *Frame) initService(nodeIds []string) *models.TaskLayer {
	return f.constructServiceTasks("InitService", constants.ServiceCmdName, nodeIds, nil)
}

func (f *Frame) startService(nodeIds []string) *models.TaskLayer {
	return f.constructServiceTasks("StartService", constants.ServiceCmdName, nodeIds, nil)
}

func (f *Frame) stopService(nodeIds []string) *models.TaskLayer {
	return f.constructServiceTasks("StopService", constants.ServiceCmdName, nodeIds, nil)
}

func (f *Frame) initAndStartServiceLayer(nodeIds []string) *models.TaskLayer {
	headTaskLayer := new(models.TaskLayer)

	preInitNodes, postInitNodes := f.getPreAndPostInitGroupNodes(nodeIds)

	// Init service before start service
	headTaskLayer.Leaf().Child = f.initService(preInitNodes)

	// TODO: custom metadata
	headTaskLayer.Leaf().Child = f.startService(nodeIds)

	// Init service after start service
	headTaskLayer.Leaf().Child = f.initService(postInitNodes)

	return headTaskLayer.Child
}

func (f *Frame) createVolumesLayer() *models.TaskLayer {
	createVolumesTaskLayer := new(models.TaskLayer)
	for nodeId, clusterNode := range f.ClusterWrapper.ClusterNodes {
		role := clusterNode.Role
		if strings.HasSuffix(role, constants.ReplicaRoleSuffix) {
			role = string([]byte(role)[:len(role)-len(constants.ReplicaRoleSuffix)])
		}
		clusterRole, exist := f.ClusterWrapper.ClusterRoles[role]
		if !exist {
			logger.Errorf("No such role [%s] in cluster role [%s]. ",
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
				createVolumesTaskLayer.Tasks = append(createVolumesTaskLayer.Tasks, createVolumesTask)
			}
		}
	}
	return createVolumesTaskLayer
}

func (f *Frame) detachVolumesLayer() *models.TaskLayer {
	detachVolumesTaskLayer := new(models.TaskLayer)
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
		detachVolumesTask := &models.Task{
			JobId:      f.Job.JobId,
			Owner:      f.Job.Owner,
			TaskAction: ActionDetachVolumes,
			Target:     f.Runtime.Provider,
			NodeId:     nodeId,
			Directive:  directive,
		}
		detachVolumesTaskLayer.Tasks = append(detachVolumesTaskLayer.Tasks, detachVolumesTask)
	}
	return detachVolumesTaskLayer
}

func (f *Frame) formatAndMountVolumeLayer() *models.TaskLayer {
	formatAndMountVolumeTaskLayer := new(models.TaskLayer)

	for nodeId, clusterNode := range f.ClusterWrapper.ClusterNodes {
		// cmd will be assigned when the task is handling
		meta := &models.Meta{
			FrontgateId: f.ClusterWrapper.Cluster.FrontgateId,
			Timeout:     TimeoutFormatAndMountVolume,
			NodeId:      clusterNode.NodeId,
			Ip:          clusterNode.PrivateIp,
		}
		directive, err := meta.ToString()
		if err != nil {
			return nil
		}
		formatVolumeTask := &models.Task{
			JobId:      f.Job.JobId,
			Owner:      f.Job.Owner,
			TaskAction: ActionFormatAndMountVolume,
			Target:     f.Runtime.Provider,
			NodeId:     nodeId,
			Directive:  string(directive),
		}
		formatAndMountVolumeTaskLayer.Tasks = append(formatAndMountVolumeTaskLayer.Tasks, formatVolumeTask)
	}
	if len(formatAndMountVolumeTaskLayer.Tasks) > 0 {
		return formatAndMountVolumeTaskLayer
	} else {
		return nil
	}
}

func (f *Frame) UmountVolumeLayer() *models.TaskLayer {
	UmountVolumeTaskLayer := new(models.TaskLayer)

	for nodeId, clusterNode := range f.ClusterWrapper.ClusterNodes {
		clusterRole := f.ClusterWrapper.ClusterRoles[clusterNode.Role]
		cmd := UmountVolumeCmd(clusterRole.MountPoint)
		meta := &models.Meta{
			FrontgateId: f.ClusterWrapper.Cluster.FrontgateId,
			Timeout:     TimeoutUmountVolume,
			NodeId:      clusterNode.NodeId,
			Ip:          clusterNode.PrivateIp,
			Cmd:         cmd,
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
		UmountVolumeTaskLayer.Tasks = append(UmountVolumeTaskLayer.Tasks, umountVolumeTask)
	}
	if len(UmountVolumeTaskLayer.Tasks) > 0 {
		return UmountVolumeTaskLayer
	} else {
		return nil
	}
}

func (f *Frame) runInstancesLayer() *models.TaskLayer {
	taskLayer := new(models.TaskLayer)
	for nodeId, clusterNode := range f.ClusterWrapper.ClusterNodes {
		role := clusterNode.Role
		if strings.HasSuffix(role, constants.ReplicaRoleSuffix) {
			role = string([]byte(role)[:len(role)-len(constants.ReplicaRoleSuffix)])
		}
		clusterRole, exist := f.ClusterWrapper.ClusterRoles[role]
		if !exist {
			logger.Errorf("No such role [%s] in cluster role [%s]. ",
				role, f.ClusterWrapper.Cluster.ClusterId)
			return nil
		}

		clusterCommon, exist := f.ClusterWrapper.ClusterCommons[role]
		if !exist {
			logger.Errorf("No such role [%s] in cluster common [%s]. ",
				role, f.ClusterWrapper.Cluster.ClusterId)
			return nil
		}

		instance := &models.Instance{
			Name:      clusterNode.ClusterId + "_" + nodeId,
			NodeId:    nodeId,
			ImageId:   clusterCommon.ImageId,
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

func (f *Frame) stopInstancesLayer() *models.TaskLayer {
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
	// When the task is handled by task controller, the cmd will be filled in,
	meta := &models.Meta{
		FrontgateId: f.ClusterWrapper.Cluster.FrontgateId,
		Timeout:     TimeoutRegister,
		ClusterId:   f.ClusterWrapper.Cluster.ClusterId,
		Cmd:         "",
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

func (f *Frame) deregisterMetadataLayer() *models.TaskLayer {
	meta := &models.Meta{
		FrontgateId: f.ClusterWrapper.Cluster.FrontgateId,
		Timeout:     TimeoutDeregister,
		ClusterId:   f.ClusterWrapper.Cluster.ClusterId,
		Cmd:         GetDeregisterExec(f.ClusterWrapper.Cluster.ClusterId),
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
		Append(f.createVolumesLayer()).              // create volume
		Append(f.runInstancesLayer()).               // run instance and attach volume to instance
		Append(f.formatAndMountVolumeLayer()).       // format and mount volume to instance
		Append(f.waitFrontgateLayer()).              // wait frontgate cluster to be active
		Append(f.registerMetadataLayer()).           // register cluster metadata
		Append(f.startConfdServiceLayer()).          // start confd service
		Append(f.initAndStartServiceLayer(nodeIds)). // register init and start cmd to exec
		Append(f.deregisterCmd(nodeIds))             // deregister cmd

	return headTaskLayer.Child
}

func (f *Frame) StopClusterLayer() *models.TaskLayer {
	var nodeIds []string
	for nodeId := range f.ClusterWrapper.ClusterNodes {
		nodeIds = append(nodeIds, nodeId)
	}
	headTaskLayer := new(models.TaskLayer)

	headTaskLayer.
		Append(f.waitFrontgateLayer()).     // wait frontgate cluster to be active
		Append(f.stopService(nodeIds)).     // register stop cmd to exec
		Append(f.stopConfdServiceLayer()).  // stop confd service
		Append(f.UmountVolumeLayer()).      // umount volume from instance
		Append(f.detachVolumesLayer()).     // detach volume from instance
		Append(f.stopInstancesLayer()).     // stop instance
		Append(f.deregisterMetadataLayer()) // deregister cluster

	return headTaskLayer.Child
}
