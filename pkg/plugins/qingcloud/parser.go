// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package qingcloud

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"openpitrix.io/openpitrix/pkg/constants"
	"openpitrix.io/openpitrix/pkg/logger"
	"openpitrix.io/openpitrix/pkg/models"
	"openpitrix.io/openpitrix/pkg/utils"
	"openpitrix.io/openpitrix/pkg/utils/yaml"
)

type Parser struct {
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
			if !utils.In(result, excludeServerIds) {
				break
			}
		}
	}
	return result, nil
}

func (p *Parser) ParseClusterRole(tmpl *models.ClusterJsonTmpl, node *models.Node) (*models.ClusterRole, error) {
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
		env, err := json.Marshal(node.Env)
		if err != nil {
			logger.Errorf("Encode env of nodes to json failed: %v", err)
			return nil, err
		}
		clusterRole.Env = string(env)
	} else if len(tmpl.Env) > 0 {
		env, err := json.Marshal(tmpl.Env)
		if err != nil {
			logger.Errorf("Encode env of cluster to json failed: %v", err)
			return nil, err
		}
		clusterRole.Env = string(env)
	}

	if clusterRole.InstanceSize == 0 {
		clusterRole.InstanceSize = constants.InstanceSize
	}
	return clusterRole, nil
}

func (p *Parser) ParseClusterNode(node *models.Node, subnetId string) (map[string]*models.ClusterNode, error) {
	count := int(node.Count)
	serverIdUpperBound := node.ServerIDUpperBound
	replicaRole := node.Role + constants.ReplicaRoleSuffix
	var serverIds, groupIds []int
	clusterNodes := make(map[string]*models.ClusterNode)
	for i := 1; i <= count; i++ {
		serverId, err := p.generateServerId(int(serverIdUpperBound), serverIds)
		if err != nil {
			logger.Errorf("Generate server id failed: %v", err)
			return nil, err
		}
		serverIds = append(serverIds, serverId)

		groupId, err := p.generateServerId(int(serverIdUpperBound), groupIds)
		if err != nil {
			logger.Errorf("Generate group id failed: %v", err)
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
		// NodeId has not been generated yet.
		clusterNodes[clusterNode.Role+fmt.Sprintf("%d", serverId)] = clusterNode

		replica := int(node.Replica)
		for j := 1; j <= replica; j++ {
			serverId, err = p.generateServerId(int(serverIdUpperBound), serverIds)
			if err != nil {
				logger.Errorf("Generate server id failed: %v", err)
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
			clusterNodes[clusterNode.Role+fmt.Sprintf("%d", serverId)] = clusterNode
		}
	}
	return clusterNodes, nil
}

func (p *Parser) ParseClusterLoadbalancer(node *models.Node) []*models.ClusterLoadbalancer {
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

func (p *Parser) ParseClusterLinks(tmpl *models.ClusterJsonTmpl) map[string]*models.ClusterLink {
	clusterLinks := make(map[string]*models.ClusterLink)
	for name, link := range tmpl.Links {
		clusterLink := &models.ClusterLink{
			Name:              name,
			ExternalClusterId: link,
		}
		clusterLinks[name] = clusterLink
	}

	return clusterLinks
}

func (p *Parser) ParseCluster(tmpl *models.ClusterJsonTmpl) (*models.Cluster, error) {
	endpoints, err := json.Marshal(tmpl.Endpoints)
	if err != nil {
		logger.Errorf("Encode endpoint to json failed: %v", err)
		return nil, err
	}

	metadataRootAccess := false
	if tmpl.MetadataRootAccess != nil {
		metadataRootAccess = *tmpl.MetadataRootAccess
	}

	cluster := &models.Cluster{
		Name:               tmpl.Name,
		Description:        tmpl.Description,
		AppId:              tmpl.AppId,
		VersionId:          tmpl.VersionId,
		SubnetId:           tmpl.Subnet,
		Endpoints:          string(endpoints),
		Status:             constants.StatusPending,
		MetadataRootAccess: metadataRootAccess,
		GlobalUuid:         tmpl.GlobalUuid,
	}
	return cluster, nil
}

func (p *Parser) ParseClusterCommon(tmpl *models.ClusterJsonTmpl,
	node *models.Node) (*models.ClusterCommon, error) {

	customMetadata := ""
	if len(node.CustomMetadata) != 0 {
		customMetadataByte, err := json.Marshal(node.CustomMetadata)
		if err != nil {
			logger.Errorf("Encode custom metadata to json failed: %v", err)
			return nil, err
		}
		customMetadata = string(customMetadataByte)
	}

	incrementalBackupSupported := false
	if tmpl.IncrementalBackupSupported != nil {
		incrementalBackupSupported = *tmpl.IncrementalBackupSupported
	}

	agentInstalled := true
	if node.AgentInstalled != nil {
		agentInstalled = *node.AgentInstalled
	}

	clusterCommon := &models.ClusterCommon{
		Role:                       node.Role,
		ServerIdUpperBound:         node.ServerIDUpperBound,
		AdvancedActions:            strings.Join(node.AdvancedActions, ","),
		BackupPolicy:               tmpl.BackupPolicy,
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
		healthCheck, err := json.Marshal(*node.HealthCheck)
		if err != nil {
			logger.Errorf("Encode node health check to json failed: %v", err)
			return nil, err
		}
		clusterCommon.HealthCheck = string(healthCheck)
	} else if tmpl.HealthCheck != nil {
		healthCheck, err := json.Marshal(*tmpl.HealthCheck)
		if err != nil {
			logger.Errorf("Encode cluster health check to json failed: %v", err)
			return nil, err
		}
		clusterCommon.HealthCheck = string(healthCheck)
	} else {
		clusterCommon.HealthCheck = ""
	}

	if node.Monitor != nil {
		monitor, err := json.Marshal(*node.Monitor)
		if err != nil {
			logger.Errorf("Encode node monitor to json failed: %v", err)
			return nil, err
		}
		clusterCommon.Monitor = string(monitor)
	} else if tmpl.Monitor != nil {
		monitor, err := json.Marshal(*tmpl.Monitor)
		if err != nil {
			logger.Errorf("Encode cluster Monitor to json failed: %v", err)
			return nil, err
		}
		clusterCommon.Monitor = string(monitor)
	} else {
		clusterCommon.Monitor = ""
	}

	for serviceName, service := range node.Services {
		var serviceValue map[string]interface{}
		switch reflect.TypeOf(service).Kind() {
		case reflect.Map:
			serviceValue = service.(map[string]interface{})
			if utils.In(serviceName, constants.ServiceNames) {
				_, exist := serviceValue["order"]
				if !exist {
					serviceValue["order"] = 0
				}
			}
		default:
			logger.Errorf("Unknown type of service [%s] ", serviceName)
			return nil, fmt.Errorf("Unknown type of service [%s] ", serviceName)
		}
		serviceByte, err := json.Marshal(serviceValue)
		if err != nil {
			logger.Errorf("Encode service [%s] to json failed: %v", serviceName, err)
			return nil, err
		}
		switch serviceName {
		case constants.ServiceInit:
			clusterCommon.InitService = string(serviceByte)
		case constants.ServiceStart:
			clusterCommon.StartService = string(serviceByte)
		case constants.ServiceStop:
			clusterCommon.StopService = string(serviceByte)
		case constants.ServiceScaleIn:
			clusterCommon.ScaleInService = string(serviceByte)
		case constants.ServiceScaleOut:
			clusterCommon.ScaleOutService = string(serviceByte)
		case constants.ServiceRestart:
			clusterCommon.RestartService = string(serviceByte)
		case constants.ServiceDestroy:
			clusterCommon.DestroyService = string(serviceByte)
		case constants.ServiceBackup:
			clusterCommon.BackupService = string(serviceByte)
		case constants.ServiceRestore:
			clusterCommon.RestoreService = string(serviceByte)
		case constants.ServiceDeleteSnapshot:
			clusterCommon.DeleteSnapshotService = string(serviceByte)
		case constants.ServiceUpgrade:
			clusterCommon.UpgradeService = string(serviceByte)
		default:
			customService := map[string]interface{}{constants.ServiceCustom: service}
			customServiceByte, err := json.Marshal(customService)
			if err != nil {
				logger.Errorf("Encode custom service [%s] to json failed: %v", serviceName, err)
				return nil, err
			}
			clusterCommon.CustomService = string(customServiceByte)
		}
	}

	return clusterCommon, nil
}

func (p *Parser) Parse(conf []byte) (*models.ClusterWrapper, error) {
	var cluster *models.Cluster
	clusterNodes := make(map[string]*models.ClusterNode)
	clusterCommons := make(map[string]*models.ClusterCommon)
	clusterLinks := make(map[string]*models.ClusterLink)
	clusterRoles := make(map[string]*models.ClusterRole)
	clusterLoadbalancers := make(map[string][]*models.ClusterLoadbalancer)

	var tmpl models.ClusterJsonTmpl
	if err := yaml.Decode(conf, &tmpl); err != nil {
		logger.Errorf("Decode conf to tmpl struct failed: %v", err)
		return nil, err
	}

	// Parse cluster
	cluster, err := p.ParseCluster(&tmpl)
	if err != nil {
		return nil, err
	}

	// Parse cluster link
	clusterLinks = p.ParseClusterLinks(&tmpl)

	if tmpl.Nodes != nil {
		for _, node := range tmpl.Nodes {
			// Parse cluster common
			clusterCommon, err := p.ParseClusterCommon(&tmpl, &node)
			clusterCommons[clusterCommon.Role] = clusterCommon

			// Parse cluster role
			clusterRole, err := p.ParseClusterRole(&tmpl, &node)
			if err != nil {
				return nil, err
			}
			clusterRoles[clusterRole.Role] = clusterRole

			// Parse cluster node
			addClusterNodes, err := p.ParseClusterNode(&node, tmpl.Subnet)
			if err != nil {
				return nil, err
			}
			for key, value := range addClusterNodes {
				clusterNodes[key] = value
			}

			// Parse cluster loadbalancer
			addClusterLoadblancers := p.ParseClusterLoadbalancer(&node)
			if len(addClusterLoadblancers) > 0 {
				clusterLoadbalancers[addClusterLoadblancers[0].Role] = addClusterLoadblancers
			}
		}
	}

	clusterWrapper := &models.ClusterWrapper{
		Cluster:              cluster,
		ClusterNodes:         clusterNodes,
		ClusterCommons:       clusterCommons,
		ClusterLinks:         clusterLinks,
		ClusterRoles:         clusterRoles,
		ClusterLoadbalancers: clusterLoadbalancers,
	}
	return clusterWrapper, nil
}
