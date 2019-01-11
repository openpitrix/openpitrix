// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package vmbased

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	"openpitrix.io/openpitrix/pkg/constants"
	"openpitrix.io/openpitrix/pkg/devkit/opapp"
	"openpitrix.io/openpitrix/pkg/logger"
	"openpitrix.io/openpitrix/pkg/models"
	"openpitrix.io/openpitrix/pkg/util/jsonutil"
	"openpitrix.io/openpitrix/pkg/util/reflectutil"
)

type Parser struct {
	Ctx context.Context
}

func (p *Parser) generateServerId(upperBound uint32, excludeServerIds []uint32) (uint32, error) {
	result := uint32(1)
	if len(excludeServerIds) == 0 {
		return result, nil
	}

	for _, serverId := range excludeServerIds {
		if serverId >= result {
			result = serverId + 1
		}
	}

	if upperBound > 0 && result > upperBound {
		for {
			if result == 0 {
				return 0, fmt.Errorf("Find server id failed. ")
			}
			result -= 1
			if !reflectutil.In(result, excludeServerIds) {
				break
			}
		}
	}
	return result, nil
}

func (p *Parser) getGroupAndServerIds(nodes []*models.ClusterNodeWithKeyPairs) ([]uint32, []uint32) {
	var serverIds, groupIds []uint32
	groupIdMap := make(map[uint32]interface{})
	for _, node := range nodes {
		groupIdMap[node.GroupId] = 0
		serverIds = append(serverIds, node.ServerId)
	}

	for groupId := range groupIdMap {
		groupIds = append(groupIds, groupId)
	}
	return groupIds, serverIds
}

func (p *Parser) ParseClusterRole(clusterConf opapp.ClusterConf, node opapp.Node) (*models.ClusterRole, error) {
	clusterRole := &models.ClusterRole{
		Role:         node.Role,
		Cpu:          node.CPU,
		Gpu:          node.GPU,
		Memory:       node.Memory,
		InstanceSize: node.Volume.InstanceSize,
		StorageSize:  node.Volume.Size,
		MountOptions: node.Volume.MountOptions,
		FileSystem:   node.Volume.Filesystem,
	}

	mountPoint := node.Volume.MountPoint
	switch v := mountPoint.(type) {
	case []string:
		clusterRole.MountPoint = strings.Join(v, ",")
	case string:
		if v == "" {
			clusterRole.MountPoint = DefaultMountPoint
		} else {
			clusterRole.MountPoint = v
		}
	default:
		clusterRole.MountPoint = DefaultMountPoint
	}

	if clusterRole.FileSystem == "" {
		clusterRole.FileSystem = Ext4FileSystem
	}

	if clusterRole.MountOptions == "" {
		if clusterRole.FileSystem == Ext4FileSystem {
			clusterRole.MountOptions = DefaultExt4MountOption
		} else if clusterRole.FileSystem == XfsFileSystem {
			clusterRole.MountOptions = DefaultXfsMountOption
		}
	}

	if len(node.Env) > 0 {
		clusterRole.Env = jsonutil.ToString(node.Env)
	} else if len(clusterConf.Env) > 0 {
		clusterRole.Env = jsonutil.ToString(clusterConf.Env)
	}

	if clusterRole.InstanceSize == 0 {
		clusterRole.InstanceSize = InstanceSize
	}
	return clusterRole, nil
}

func (p *Parser) ParseAddClusterNode(clusterConf opapp.ClusterConf, clusterWrapper *models.ClusterWrapper) error {
	existRoleNodes := make(map[string][]*models.ClusterNodeWithKeyPairs)
	addRoleNodes := make(map[string][]*models.ClusterNodeWithKeyPairs)
	clusterNodes := make(map[string]*models.ClusterNodeWithKeyPairs)
	for _, clusterNode := range clusterWrapper.ClusterNodesWithKeyPairs {
		role := clusterNode.Role
		if strings.HasSuffix(role, constants.ReplicaRoleSuffix) {
			role = string([]byte(role)[:len(role)-len(constants.ReplicaRoleSuffix)])
		}
		if clusterNode.Status != constants.StatusPending {
			nodes, isExist := existRoleNodes[role]
			if isExist {
				nodes = append(nodes, clusterNode)
			} else {
				existRoleNodes[role] = []*models.ClusterNodeWithKeyPairs{clusterNode}
			}
		} else {
			nodes, isExist := addRoleNodes[role]
			if isExist {
				nodes = append(nodes, clusterNode)
			} else {
				addRoleNodes[role] = []*models.ClusterNodeWithKeyPairs{clusterNode}
			}
		}
	}

	for role, addNodes := range addRoleNodes {
		// add replica
		addNodeRole := addNodes[0].Role
		if strings.HasSuffix(addNodeRole, constants.ReplicaRoleSuffix) {
			nodes, isExist := existRoleNodes[role]
			if !isExist {
				err := fmt.Errorf("role [%s] not exist", role)
				return err
			}

			count := len(addNodes)
			serverIdUpperBound := clusterWrapper.ClusterCommons[role].ServerIdUpperBound
			groupIds, serverIds := p.getGroupAndServerIds(nodes)
			for _, groupId := range groupIds {
				for i := 1; i <= count; i++ {
					serverId, err := p.generateServerId(serverIdUpperBound, serverIds)
					if err != nil {
						logger.Error(p.Ctx, "Generate server id failed: %v", err)
						return err
					}
					serverIds = append(serverIds, serverId)
					clusterNode := &models.ClusterNode{
						Name:      "",
						ClusterId: nodes[0].ClusterId,
						Role:      addNodeRole,
						SubnetId:  nodes[0].SubnetId,
						ServerId:  serverId,
						Status:    constants.StatusPending,
						GroupId:   groupId,
					}
					clusterNodeWIthKeyPairs := &models.ClusterNodeWithKeyPairs{
						ClusterNode: clusterNode,
					}
					// NodeId has not been generated yet.
					clusterNodes[clusterNode.Role+fmt.Sprintf("%d", serverId)] = clusterNodeWIthKeyPairs
				}
			}
		} else {
			var replicaCount int
			var serverIds, groupIds []uint32
			var subnetId string
			count := len(addNodes)
			serverIdUpperBound := clusterWrapper.ClusterCommons[role].ServerIdUpperBound
			nodes, isExist := existRoleNodes[role]
			if !isExist {
				var nodeConf *opapp.Node
				if clusterConf.Nodes != nil {
					for _, node := range clusterConf.Nodes {
						if node.Role == role {
							nodeConf = &node
						}
					}
				}
				if nodeConf == nil {
					err := fmt.Errorf("role [%s] not exist in package", role)
					return err
				}

				replicaCount = int(nodeConf.Replica)
				subnetId = clusterWrapper.Cluster.SubnetId
			} else {
				groupIds, serverIds = p.getGroupAndServerIds(nodes)
				replicaCount = (len(nodes) - len(groupIds)) / len(groupIds)
				subnetId = nodes[0].SubnetId
			}

			for i := 1; i <= count; i++ {
				serverId, err := p.generateServerId(serverIdUpperBound, serverIds)
				if err != nil {
					logger.Error(p.Ctx, "Generate server id failed: %v", err)
					return err
				}
				serverIds = append(serverIds, serverId)

				groupId, err := p.generateServerId(serverIdUpperBound, groupIds)
				if err != nil {
					logger.Error(p.Ctx, "Generate group id failed: %v", err)
					return err
				}
				groupIds = append(groupIds, groupId)

				clusterNode := &models.ClusterNode{
					Name:      "",
					Role:      addNodeRole,
					ClusterId: clusterWrapper.Cluster.ClusterId,
					SubnetId:  subnetId,
					ServerId:  serverId,
					Status:    constants.StatusPending,
					GroupId:   groupId,
				}
				clusterNodeWIthKeyPairs := &models.ClusterNodeWithKeyPairs{
					ClusterNode: clusterNode,
				}
				// NodeId has not been generated yet.
				clusterNodes[clusterNode.Role+fmt.Sprintf("%d", serverId)] = clusterNodeWIthKeyPairs

				for j := 1; j <= replicaCount; j++ {
					serverId, err = p.generateServerId(serverIdUpperBound, serverIds)
					if err != nil {
						logger.Error(p.Ctx, "Generate server id failed: %v", err)
						return err
					}
					serverIds = append(serverIds, serverId)

					replicaRole := addNodeRole + constants.ReplicaRoleSuffix

					clusterNode := &models.ClusterNode{
						Name:      "",
						Role:      replicaRole,
						ClusterId: clusterWrapper.Cluster.ClusterId,
						SubnetId:  subnetId,
						ServerId:  serverId,
						Status:    constants.StatusPending,
						GroupId:   groupId,
					}
					clusterNodeWIthKeyPairs := &models.ClusterNodeWithKeyPairs{
						ClusterNode: clusterNode,
					}
					clusterNodes[clusterNode.Role+fmt.Sprintf("%d", serverId)] = clusterNodeWIthKeyPairs
				}
			}

		}
	}
	clusterWrapper.ClusterNodesWithKeyPairs = clusterNodes

	return nil
}

func (p *Parser) ParseClusterNode(node opapp.Node, subnetId string) (map[string]*models.ClusterNodeWithKeyPairs, error) {
	count := int(node.Count)
	serverIdUpperBound := node.ServerIDUpperBound
	replicaRole := node.Role + constants.ReplicaRoleSuffix
	var serverIds, groupIds []uint32
	clusterNodes := make(map[string]*models.ClusterNodeWithKeyPairs)
	for i := 1; i <= count; i++ {
		serverId, err := p.generateServerId(serverIdUpperBound, serverIds)
		if err != nil {
			logger.Error(p.Ctx, "Generate server id failed: %v", err)
			return nil, err
		}
		serverIds = append(serverIds, serverId)

		groupId, err := p.generateServerId(serverIdUpperBound, groupIds)
		if err != nil {
			logger.Error(p.Ctx, "Generate group id failed: %v", err)
			return nil, err
		}
		groupIds = append(groupIds, groupId)

		clusterNode := &models.ClusterNode{
			NodeId:   node.Role + fmt.Sprintf("%d", serverId),
			Name:     "",
			Role:     node.Role,
			SubnetId: subnetId,
			ServerId: serverId,
			Status:   constants.StatusPending,
			GroupId:  groupId,
		}
		clusterNodeWIthKeyPairs := &models.ClusterNodeWithKeyPairs{
			ClusterNode: clusterNode,
		}
		// NodeId has not been generated yet.
		clusterNodes[clusterNode.NodeId] = clusterNodeWIthKeyPairs

		replica := int(node.Replica)
		for j := 1; j <= replica; j++ {
			serverId, err = p.generateServerId(serverIdUpperBound, serverIds)
			if err != nil {
				logger.Error(p.Ctx, "Generate server id failed: %v", err)
				return nil, err
			}
			serverIds = append(serverIds, serverId)

			clusterNode := &models.ClusterNode{
				NodeId:   clusterNode.Role + fmt.Sprintf("%d", serverId),
				Name:     "",
				Role:     replicaRole,
				SubnetId: subnetId,
				ServerId: serverId,
				Status:   constants.StatusPending,
				GroupId:  groupId,
			}
			clusterNodeWIthKeyPairs := &models.ClusterNodeWithKeyPairs{
				ClusterNode: clusterNode,
			}
			clusterNodes[clusterNode.NodeId] = clusterNodeWIthKeyPairs
		}
	}
	return clusterNodes, nil
}

func (p *Parser) ParseClusterLoadbalancer(node opapp.Node) []*models.ClusterLoadbalancer {
	var clusterLoadbalancers []*models.ClusterLoadbalancer
	for _, loadbalancer := range node.Loadbalancer {
		clusterLoadbalancer := &models.ClusterLoadbalancer{
			Role:                   node.Role,
			LoadbalancerListenerId: loadbalancer.Listener,
			LoadbalancerPolicyId:   loadbalancer.Policy,
			LoadbalancerPort:       loadbalancer.Port,
		}
		clusterLoadbalancers = append(clusterLoadbalancers, clusterLoadbalancer)
	}

	return clusterLoadbalancers
}

func (p *Parser) ParseClusterLinks(clusterConf opapp.ClusterConf) map[string]*models.ClusterLink {
	clusterLinks := make(map[string]*models.ClusterLink)
	for name, link := range clusterConf.Links {
		clusterLink := &models.ClusterLink{
			Name:              name,
			ExternalClusterId: link,
		}
		clusterLinks[name] = clusterLink
	}

	return clusterLinks
}

func (p *Parser) ParseCluster(clusterConf opapp.ClusterConf, clusterEnv string) (*models.Cluster, error) {
	endpoints := jsonutil.ToString(clusterConf.Endpoints)

	metadataRootAccess := false
	if clusterConf.MetadataRootAccess != nil {
		metadataRootAccess = *clusterConf.MetadataRootAccess
	}

	cluster := &models.Cluster{
		Name:               clusterConf.Name,
		Description:        clusterConf.Description,
		AppId:              clusterConf.AppId,
		VersionId:          clusterConf.VersionId,
		SubnetId:           clusterConf.Subnet,
		Endpoints:          endpoints,
		Status:             constants.StatusPending,
		MetadataRootAccess: metadataRootAccess,
		GlobalUuid:         clusterConf.GlobalUuid,
		Env:                clusterEnv,
	}
	return cluster, nil
}

func (p *Parser) ParseClusterCommon(clusterConf opapp.ClusterConf, node opapp.Node) (*models.ClusterCommon, error) {

	customMetadata := ""
	if len(node.CustomMetadata) != 0 {
		customMetadata = jsonutil.ToString(node.CustomMetadata)
	}

	incrementalBackupSupported := false
	if clusterConf.IncrementalBackupSupported != nil {
		incrementalBackupSupported = *clusterConf.IncrementalBackupSupported
	}

	agentInstalled := true
	if node.AgentInstalled != nil {
		agentInstalled = *node.AgentInstalled
	}

	clusterCommon := &models.ClusterCommon{
		Role:                       node.Role,
		ServerIdUpperBound:         node.ServerIDUpperBound,
		AdvancedActions:            strings.Join(node.AdvancedActions, ","),
		BackupPolicy:               clusterConf.BackupPolicy,
		IncrementalBackupSupported: incrementalBackupSupported,
		Passphraseless:             node.Passphraseless,
		CustomMetadataScript:       customMetadata,
		VerticalScalingPolicy:      node.VerticalScalingPolicy,
		AgentInstalled:             agentInstalled,
		ImageId:                    node.Container.Image,
		Hypervisor:                 node.Container.Type,
	}

	if clusterCommon.VerticalScalingPolicy == "" {
		clusterCommon.VerticalScalingPolicy = constants.ScalingPolicyParallel
	}

	if node.HealthCheck != nil {
		clusterCommon.HealthCheck = jsonutil.ToString(node.HealthCheck)
	} else if clusterConf.HealthCheck != nil {
		clusterCommon.HealthCheck = jsonutil.ToString(clusterConf.HealthCheck)
	} else {
		clusterCommon.HealthCheck = ""
	}

	if node.Monitor != nil {
		clusterCommon.Monitor = jsonutil.ToString(node.Monitor)
	} else if clusterConf.Monitor != nil {
		clusterCommon.Monitor = jsonutil.ToString(clusterConf.Monitor)
	} else {
		clusterCommon.Monitor = ""
	}

	for serviceName, service := range node.Services {
		var serviceValue map[string]interface{}
		switch reflect.TypeOf(service).Kind() {
		case reflect.Map:
			serviceValue = service.(map[string]interface{})
			if reflectutil.In(serviceName, constants.ServiceNames) {
				_, exist := serviceValue["order"]
				if !exist {
					serviceValue["order"] = 0
				}
			}
		default:
			logger.Error(p.Ctx, "Unknown type of service [%s] ", serviceName)
			return nil, fmt.Errorf("Unknown type of service [%s] ", serviceName)
		}
		serviceStr := jsonutil.ToString(serviceValue)
		switch serviceName {
		case constants.ServiceInit:
			clusterCommon.InitService = serviceStr
		case constants.ServiceStart:
			clusterCommon.StartService = serviceStr
		case constants.ServiceStop:
			clusterCommon.StopService = serviceStr
		case constants.ServiceScaleIn:
			clusterCommon.ScaleInService = serviceStr
		case constants.ServiceScaleOut:
			clusterCommon.ScaleOutService = serviceStr
		case constants.ServiceRestart:
			clusterCommon.RestartService = serviceStr
		case constants.ServiceDestroy:
			clusterCommon.DestroyService = serviceStr
		case constants.ServiceBackup:
			clusterCommon.BackupService = serviceStr
		case constants.ServiceRestore:
			clusterCommon.RestoreService = serviceStr
		case constants.ServiceDeleteSnapshot:
			clusterCommon.DeleteSnapshotService = serviceStr
		case constants.ServiceUpgrade:
			clusterCommon.UpgradeService = serviceStr
		default:
			customService := map[string]interface{}{constants.ServiceCustom: service}
			clusterCommon.CustomService = jsonutil.ToString(customService)
		}
	}

	return clusterCommon, nil
}

func (p *Parser) Parse(clusterConf opapp.ClusterConf, clusterWrapper *models.ClusterWrapper, clusterEnv string) error {
	var cluster *models.Cluster
	clusterNodes := make(map[string]*models.ClusterNodeWithKeyPairs)
	clusterCommons := make(map[string]*models.ClusterCommon)
	clusterLinks := make(map[string]*models.ClusterLink)
	clusterRoles := make(map[string]*models.ClusterRole)
	clusterLoadbalancers := make(map[string][]*models.ClusterLoadbalancer)

	// Parse cluster
	cluster, err := p.ParseCluster(clusterConf, clusterEnv)
	if err != nil {
		return err
	}

	// Parse cluster link
	clusterLinks = p.ParseClusterLinks(clusterConf)

	if clusterConf.Nodes != nil {
		for _, node := range clusterConf.Nodes {
			// Parse cluster common
			clusterCommon, err := p.ParseClusterCommon(clusterConf, node)
			clusterCommons[clusterCommon.Role] = clusterCommon

			// Parse cluster role
			clusterRole, err := p.ParseClusterRole(clusterConf, node)
			if err != nil {
				return err
			}
			clusterRoles[clusterRole.Role] = clusterRole

			// Parse cluster node
			addClusterNodes, err := p.ParseClusterNode(node, clusterConf.Subnet)
			if err != nil {
				return err
			}
			for key, value := range addClusterNodes {
				clusterNodes[key] = value
			}

			// Parse cluster loadbalancer
			addClusterLoadblancers := p.ParseClusterLoadbalancer(node)
			if len(addClusterLoadblancers) > 0 {
				clusterLoadbalancers[addClusterLoadblancers[0].Role] = addClusterLoadblancers
			}
		}
	}

	// add cluster nodes
	if clusterWrapper.Cluster != nil && len(clusterWrapper.Cluster.ClusterId) > 0 {
		logger.Debug(p.Ctx, "Add cluster [%s] node.", clusterWrapper.Cluster.ClusterId)
		err = p.ParseAddClusterNode(clusterConf, clusterWrapper)
		if err != nil {
			return err
		}
		clusterWrapper.Cluster.Env = clusterEnv
	} else {
		clusterWrapper.ClusterNodesWithKeyPairs = clusterNodes
		clusterWrapper.Cluster = cluster
		clusterWrapper.ClusterLinks = clusterLinks
		clusterWrapper.ClusterLoadbalancers = clusterLoadbalancers
	}
	clusterWrapper.ClusterCommons = clusterCommons
	clusterWrapper.ClusterRoles = clusterRoles

	return nil
}
