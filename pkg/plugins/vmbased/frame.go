// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package vmbased

import (
	"encoding/json"
	"sort"
	"strings"

	runtimeenvclient "openpitrix.io/openpitrix/pkg/client/runtimeenv"
	"openpitrix.io/openpitrix/pkg/constants"
	"openpitrix.io/openpitrix/pkg/logger"
	"openpitrix.io/openpitrix/pkg/models"
)

type Frame struct {
	Job            *models.Job
	ClusterWrapper *models.ClusterWrapper
	Runtime        *runtimeenvclient.Runtime
}

func NewFrame(job *models.Job) (*Frame, error) {
	clusterWrapper, err := models.NewClusterWrapper(job.Directive)
	if err != nil {
		return nil, err
	}

	runtimeEnvId := clusterWrapper.Cluster.RuntimeEnvId
	runtime, err := runtimeenvclient.NewRuntime(runtimeEnvId)
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
		startConfdDirective := map[string]interface{}{
			"frontgate_id": f.ClusterWrapper.Cluster.FrontgateId,
			"ip":           clusterNode.PrivateIp,
			"timeout":      TimeoutStartConfd,
		}
		directive, _ := json.Marshal(startConfdDirective)
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
		// TODO: construct cmd
		deregisterCmdDirective := map[string]interface{}{
			"frontgate_id": f.ClusterWrapper.Cluster.FrontgateId,
			"ip":           f.ClusterWrapper.ClusterNodes[nodeId].PrivateIp,
			"timeout":      TimeoutDeregisterCmd,
			"cmd":          "",
		}
		directive, _ := json.Marshal(deregisterCmdDirective)
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
				timeout = *service.Timeout
			}
			if service.Cmd == "" {
				continue
			}
			// TODO: construct cmd
			registerCmdDirective := map[string]interface{}{
				"frontgate_id": f.ClusterWrapper.Cluster.FrontgateId,
				"ip":           f.ClusterWrapper.ClusterNodes[nodeId].PrivateIp,
				"timeout":      timeout,
				"cmd":          service.Cmd,
			}
			directive, _ := json.Marshal(registerCmdDirective)
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
			execNodeNums = *service.NodesToExecuteOn
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
			order = *service.Order
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

func (f *Frame) initAndStartServiceLayer(nodeIds []string) *models.TaskLayer {
	headTaskLayer := new(models.TaskLayer)
	// Start confd service
	headTaskLayer.Leaf().Child = f.startConfdServiceLayer()

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
				RuntimeId: f.Runtime.RuntimeEnvId,
			}
			volumeTaskDirective, err := volume.ToString()
			if err != nil {
				return nil
			}

			createVolumesTask := &models.Task{
				JobId:      f.Job.JobId,
				Owner:      f.Job.Owner,
				TaskAction: ActionCreateVolumes,
				Target:     f.Runtime.Runtime,
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
			RuntimeId: f.Runtime.RuntimeEnvId,
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
			Target:     f.Runtime.Runtime,
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

func (f *Frame) waitFrontgateLayer() *models.TaskLayer {
	waitFrontgateDirective := map[string]interface{}{
		"frontgate_id": f.ClusterWrapper.Cluster.FrontgateId,
	}
	directive, _ := json.Marshal(waitFrontgateDirective)
	// Wait frontgate available
	waitFrontgateTask := &models.Task{
		JobId:      f.Job.JobId,
		Owner:      f.Job.Owner,
		TaskAction: ActionWaitFrontgateAvailable,
		Target:     f.Runtime.Runtime,
		NodeId:     f.ClusterWrapper.Cluster.ClusterId,
		Directive:  string(directive),
	}
	return &models.TaskLayer{
		Tasks: []*models.Task{waitFrontgateTask},
	}
}

func (f *Frame) registerMetadataLayer() *models.TaskLayer {
	metadata := &Metadata{
		ClusterWrapper: f.ClusterWrapper,
		Runtime:        f.Runtime,
	}
	cnodes := metadata.GetClusterCnodes()
	directive, err := json.Marshal(cnodes)
	if err != nil {
		logger.Errorf("Marshal cluster [%s] metadata cnodes failed: %+v",
			f.ClusterWrapper.Cluster.ClusterId, err)
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

func (f *Frame) CreateClusterLayer() *models.TaskLayer {
	headTaskLayer := new(models.TaskLayer)
	headTaskLayer.Leaf().Child = f.createVolumesLayer()
	headTaskLayer.Leaf().Child = f.runInstancesLayer()
	headTaskLayer.Leaf().Child = f.waitFrontgateLayer()
	headTaskLayer.Leaf().Child = f.registerMetadataLayer()
	var nodeIds []string
	for nodeId := range f.ClusterWrapper.ClusterNodes {
		nodeIds = append(nodeIds, nodeId)
	}
	headTaskLayer.Leaf().Child = f.initAndStartServiceLayer(nodeIds)
	headTaskLayer.Leaf().Child = f.deregisterCmd(nodeIds)
	return headTaskLayer.Child
}
