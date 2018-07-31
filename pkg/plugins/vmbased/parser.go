// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package vmbased

import (
	"fmt"
	"reflect"
	"strings"

	"openpitrix.io/openpitrix/pkg/constants"
	"openpitrix.io/openpitrix/pkg/devkit/app"
	"openpitrix.io/openpitrix/pkg/logger"
	"openpitrix.io/openpitrix/pkg/models"
	"openpitrix.io/openpitrix/pkg/util/jsonutil"
	"openpitrix.io/openpitrix/pkg/util/reflectutil"
)

type Parser struct {
	Logger *logger.Logger
}

func (p *Parser) generateServerId(upperBound int, excludeServerIds []int) (int, error) {
	result := 1
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

func (p *Parser) ParseClusterRole(clusterConf app.ClusterConf, node app.Node) (*models.ClusterRole, error) {
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
			clusterRole.MountPoint = constants.DefaultMountPoint
		} else {
			clusterRole.MountPoint = v
		}
	default:
		clusterRole.MountPoint = constants.DefaultMountPoint
	}

	if clusterRole.FileSystem == "" {
		clusterRole.FileSystem = constants.Ext4FileSystem
	}

	if clusterRole.MountOptions == "" {
		if clusterRole.FileSystem == constants.Ext4FileSystem {
			clusterRole.MountOptions = constants.DefaultExt4MountOption
		} else if clusterRole.FileSystem == constants.XfsFileSystem {
			clusterRole.MountOptions = constants.DefaultXfsMountOption
		}
	}

	if len(node.Env) > 0 {
		clusterRole.Env = jsonutil.ToString(node.Env)
	} else if len(clusterConf.Env) > 0 {
		clusterRole.Env = jsonutil.ToString(clusterConf.Env)
	}

	if clusterRole.InstanceSize == 0 {
		clusterRole.InstanceSize = constants.InstanceSize
	}
	return clusterRole, nil
}

func (p *Parser) ParseClusterNode(node app.Node, subnetId string) (map[string]*models.ClusterNodeWithKeyPairs, error) {
	count := int(node.Count)
	serverIdUpperBound := node.ServerIDUpperBound
	replicaRole := node.Role + constants.ReplicaRoleSuffix
	var serverIds, groupIds []int
	clusterNodes := make(map[string]*models.ClusterNodeWithKeyPairs)
	for i := 1; i <= count; i++ {
		serverId, err := p.generateServerId(int(serverIdUpperBound), serverIds)
		if err != nil {
			p.Logger.Error("Generate server id failed: %v", err)
			return nil, err
		}
		serverIds = append(serverIds, serverId)

		groupId, err := p.generateServerId(int(serverIdUpperBound), groupIds)
		if err != nil {
			p.Logger.Error("Generate group id failed: %v", err)
			return nil, err
		}
		groupIds = append(groupIds, groupId)

		clusterNode := &models.ClusterNode{
			Name:     "",
			Role:     node.Role,
			SubnetId: subnetId,
			ServerId: uint32(serverId),
			Status:   constants.StatusPending,
			GroupId:  uint32(groupId),
		}
		clusterNodeWIthKeyPairs := &models.ClusterNodeWithKeyPairs{
			ClusterNode: clusterNode,
		}
		// NodeId has not been generated yet.
		clusterNodes[clusterNode.Role+fmt.Sprintf("%d", serverId)] = clusterNodeWIthKeyPairs

		replica := int(node.Replica)
		for j := 1; j <= replica; j++ {
			serverId, err = p.generateServerId(int(serverIdUpperBound), serverIds)
			if err != nil {
				p.Logger.Error("Generate server id failed: %v", err)
				return nil, err
			}
			serverIds = append(serverIds, serverId)

			clusterNode := &models.ClusterNode{
				Name:     "",
				Role:     replicaRole,
				SubnetId: subnetId,
				ServerId: uint32(serverId),
				Status:   constants.StatusPending,
				GroupId:  uint32(groupId),
			}
			clusterNodeWIthKeyPairs := &models.ClusterNodeWithKeyPairs{
				ClusterNode: clusterNode,
			}
			clusterNodes[clusterNode.Role+fmt.Sprintf("%d", serverId)] = clusterNodeWIthKeyPairs
		}
	}
	return clusterNodes, nil
}

func (p *Parser) ParseClusterLoadbalancer(node app.Node) []*models.ClusterLoadbalancer {
	var clusterLoadbalancers []*models.ClusterLoadbalancer
	for _, loadbalancer := range node.Loadbalancer {
		clusterLoadbalancer := &models.ClusterLoadbalancer{
			Role: node.Role,
			LoadbalancerListenerId: loadbalancer.Listener,
			LoadbalancerPolicyId:   loadbalancer.Policy,
			LoadbalancerPort:       loadbalancer.Port,
		}
		clusterLoadbalancers = append(clusterLoadbalancers, clusterLoadbalancer)
	}

	return clusterLoadbalancers
}

func (p *Parser) ParseClusterLinks(clusterConf app.ClusterConf) map[string]*models.ClusterLink {
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

func (p *Parser) ParseCluster(clusterConf app.ClusterConf) (*models.Cluster, error) {
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
	}
	return cluster, nil
}

func (p *Parser) ParseClusterCommon(clusterConf app.ClusterConf, node app.Node) (*models.ClusterCommon, error) {

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
			p.Logger.Error("Unknown type of service [%s] ", serviceName)
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

func (p *Parser) Parse(clusterConf app.ClusterConf) (*models.ClusterWrapper, error) {
	var cluster *models.Cluster
	clusterNodes := make(map[string]*models.ClusterNodeWithKeyPairs)
	clusterCommons := make(map[string]*models.ClusterCommon)
	clusterLinks := make(map[string]*models.ClusterLink)
	clusterRoles := make(map[string]*models.ClusterRole)
	clusterLoadbalancers := make(map[string][]*models.ClusterLoadbalancer)

	// Parse cluster
	cluster, err := p.ParseCluster(clusterConf)
	if err != nil {
		return nil, err
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
				return nil, err
			}
			clusterRoles[clusterRole.Role] = clusterRole

			// Parse cluster node
			addClusterNodes, err := p.ParseClusterNode(node, clusterConf.Subnet)
			if err != nil {
				return nil, err
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

	clusterWrapper := &models.ClusterWrapper{
		Cluster:                  cluster,
		ClusterNodesWithKeyPairs: clusterNodes,
		ClusterCommons:           clusterCommons,
		ClusterLinks:             clusterLinks,
		ClusterRoles:             clusterRoles,
		ClusterLoadbalancers:     clusterLoadbalancers,
	}
	return clusterWrapper, nil
}
