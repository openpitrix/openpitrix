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

func (p *Parser) ParseClusterRole(mustache *models.ClusterJsonMustache, node *models.Node) (*models.ClusterRole, error) {
	clusterRole := &models.ClusterRole{
		Role:         node.Role,
		Cpu:          node.CPU,
		Gpu:          node.GPU,
		Memory:       node.Memory,
		InstanceSize: node.Volume.InstanceSize,
		StorageSize:  node.Volume.Size,
	}
	if len(node.Env) > 0 {
		env, err := json.Marshal(node.Env)
		if err != nil {
			logger.Errorf("Encode env of nodes to json failed: %v", err)
			return nil, err
		}
		clusterRole.Env = string(env)
	} else if len(mustache.Env) > 0 {
		env, err := json.Marshal(mustache.Env)
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

func (p *Parser) ParseClusterNode(node *models.Node, vxnet string) ([]*models.ClusterNode, error) {
	count := int(node.Count)
	serverIdUpperBound := node.ServerIDUpperBound
	replicaRole := node.Role + constants.ReplicaRoleSuffix
	var serverIds, groupIds []int
	var clusterNodes []*models.ClusterNode
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
			VxnetId:  vxnet,
			ServerId: int32(serverId),
			Status:   constants.StatusPending,
			GroupId:  int32(groupId),
		}
		clusterNodes = append(clusterNodes, clusterNode)

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
				VxnetId:  vxnet,
				ServerId: int32(serverId),
				Status:   constants.StatusPending,
				GroupId:  int32(groupId),
			}
			clusterNodes = append(clusterNodes, clusterNode)
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

func (p *Parser) ParseClusterLinks(mustache *models.ClusterJsonMustache) []*models.ClusterLink {
	var clusterLinks []*models.ClusterLink
	for name, link := range mustache.Links {
		clusterLink := &models.ClusterLink{
			Name:              name,
			ExternalClusterId: link,
		}
		clusterLinks = append(clusterLinks, clusterLink)
	}

	return clusterLinks
}

func (p *Parser) ParseCluster(mustache *models.ClusterJsonMustache) (*models.Cluster, error) {
	endpoints, err := json.Marshal(mustache.Endpoints)
	if err != nil {
		logger.Errorf("Encode endpoint to json failed: %v", err)
		return nil, err
	}

	metadataRootAccess := false
	if mustache.MetadataRootAccess != nil {
		metadataRootAccess = *mustache.MetadataRootAccess
	}

	cluster := &models.Cluster{
		Name:               mustache.Name,
		Description:        mustache.Description,
		AppId:              mustache.AppId,
		VersionId:          mustache.VersionId,
		VxnetId:            mustache.Vxnet,
		Endpoints:          string(endpoints),
		Status:             constants.StatusPending,
		MetadataRootAccess: metadataRootAccess,
		GlobalUuid:         mustache.GlobalUuid,
	}
	return cluster, nil
}

func (p *Parser) ParseClusterCommon(mustache *models.ClusterJsonMustache,
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
	if mustache.IncrementalBackupSupported != nil {
		incrementalBackupSupported = *mustache.IncrementalBackupSupported
	}

	agentInstalled := true
	if node.AgentInstalled != nil {
		agentInstalled = *node.AgentInstalled
	}

	clusterCommon := &models.ClusterCommon{
		Role:                       node.Role,
		ServerIdUpperBound:         node.ServerIDUpperBound,
		AdvancedActions:            strings.Join(node.AdvancedActions, ","),
		BackupPolicy:               mustache.BackupPolicy,
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
	} else if mustache.HealthCheck != nil {
		healthCheck, err := json.Marshal(*mustache.HealthCheck)
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
	} else if mustache.Monitor != nil {
		monitor, err := json.Marshal(*mustache.Monitor)
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
	var clusterNodes []*models.ClusterNode
	var clusterCommons []*models.ClusterCommon
	var clusterLinks []*models.ClusterLink
	var clusterRoles []*models.ClusterRole
	var clusterLoadbalancers []*models.ClusterLoadbalancer

	var mustache models.ClusterJsonMustache
	if err := yaml.Decode(conf, &mustache); err != nil {
		logger.Errorf("Decode conf to mustache struct failed: %v", err)
		return nil, err
	}

	// Parse cluster
	cluster, err := p.ParseCluster(&mustache)
	if err != nil {
		return nil, err
	}

	// Parse cluster link
	clusterLinks = p.ParseClusterLinks(&mustache)

	if mustache.Nodes != nil {
		for _, node := range mustache.Nodes {
			// Parse cluster common
			clusterCommon, err := p.ParseClusterCommon(&mustache, &node)
			clusterCommons = append(clusterCommons, clusterCommon)

			// Parse cluster role
			clusterRole, err := p.ParseClusterRole(&mustache, &node)
			if err != nil {
				return nil, err
			}
			clusterRoles = append(clusterRoles, clusterRole)

			// Parse cluster node
			addClusterNodes, err := p.ParseClusterNode(&node, mustache.Vxnet)
			if err != nil {
				return nil, err
			}
			clusterNodes = append(clusterNodes, addClusterNodes...)

			// Parse cluster loadbalancer
			addClusterLoadblancers := p.ParseClusterLoadbalancer(&node)
			clusterLoadbalancers = append(clusterLoadbalancers, addClusterLoadblancers...)
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
