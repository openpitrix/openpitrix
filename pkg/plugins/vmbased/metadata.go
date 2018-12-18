// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package vmbased

import (
	"context"
	"fmt"
	"strings"

	"openpitrix.io/openpitrix/pkg/constants"
	"openpitrix.io/openpitrix/pkg/logger"
	"openpitrix.io/openpitrix/pkg/models"
	"openpitrix.io/openpitrix/pkg/util/jsonutil"
)

type Metadata struct {
	ClusterWrapper *models.ClusterWrapper
	RuntimeDetails *models.RuntimeDetails
}

/*
Compose cluster info into the following format,
in order to register cluster to configuration management service.
{
	"clusters": {
		"<cluster_id>": {
	 		"hosts": {
				<The data from the function GetHostsCnodes>
	 		},
	 		"cluster": {
				<The data from the function GetClusterMetadataCnodes>
	 		},
	 		"env": { # optional
				<The data from the function GetEnvCnodes>
	 		}
   		}
	}
}
*/
func (m *Metadata) GetClusterCnodes(ctx context.Context) map[string]interface{} {
	logger.Info(ctx, "Composing cluster %s", m.ClusterWrapper.Cluster.ClusterId)

	data := make(map[string]interface{})

	var nodeIds []string
	for nodeId := range m.ClusterWrapper.ClusterNodesWithKeyPairs {
		nodeIds = append(nodeIds, nodeId)
	}
	// hosts
	hosts := m.GetHostsCnodes(ctx, nodeIds)
	data[RegisterNodeHosts] = hosts

	// cluster
	clusterMetadata := m.GetClusterMetadataCnodes()
	data[RegisterNodeCluster] = clusterMetadata

	// endpoints
	endpoints, _ := m.ClusterWrapper.GetEndpoints()
	if len(endpoints) > 0 {
		data[RegisterNodeEndpoint] = endpoints
	}

	// env
	env := m.GetEnvCnodes(ctx)
	if len(env) > 0 {
		data[RegisterNodeEnv] = env
	}

	cnodes := map[string]interface{}{
		RegisterClustersRootPath: map[string]interface{}{
			m.ClusterWrapper.Cluster.ClusterId: data,
		},
	}
	logger.Info(ctx, "Composed cluster %s cnodes successful: %+v", m.ClusterWrapper.Cluster.ClusterId, cnodes)

	return cnodes
}

/*
{
	"clusters": {
		"<cluster_id>": {
	 		"env": { # optional
				<The data from the function GetEnvCnodes>
	 		}
   		}
	}
}
*/
func (m *Metadata) GetClusterEnvCnodes(ctx context.Context) map[string]interface{} {
	logger.Info(ctx, "Composing cluster %s env", m.ClusterWrapper.Cluster.ClusterId)

	data := make(map[string]interface{})

	// env
	env := m.GetEnvCnodes(ctx)
	if len(env) > 0 {
		data[RegisterNodeEnv] = env
	}

	cnodes := map[string]interface{}{
		RegisterClustersRootPath: map[string]interface{}{
			m.ClusterWrapper.Cluster.ClusterId: data,
		},
	}
	logger.Info(ctx, "Composed cluster %s env cnodes successful: %+v", m.ClusterWrapper.Cluster.ClusterId, cnodes)

	return cnodes
}

/*
{
	"self": {
		<The data from the function GetMappingCnodes below>
   	}
}
*/
func (m *Metadata) GetClusterMappingCnodes(ctx context.Context) map[string]interface{} {
	logger.Info(ctx, "Composing cluster %s mapping", m.ClusterWrapper.Cluster.ClusterId)
	var nodeIds []string
	for nodeId := range m.ClusterWrapper.ClusterNodesWithKeyPairs {
		nodeIds = append(nodeIds, nodeId)
	}
	mapping := m.GetMappingCnodes(nodeIds)
	logger.Info(ctx, "Composed cluster %s mapping successful: %+v", m.ClusterWrapper.Cluster.ClusterId, mapping)

	return mapping
}

/*
{
	"self": {
		<The data from the function GetMappingCnodes below>
   	}
}
*/
func (m *Metadata) GetClusterNodesMappingCnodes(ctx context.Context, nodeIds []string) map[string]interface{} {
	logger.Info(ctx, "Composing cluster %s mapping", m.ClusterWrapper.Cluster.ClusterId)
	mapping := m.GetMappingCnodes(nodeIds)
	logger.Info(ctx, "Composed cluster %s mapping successful: %+v", m.ClusterWrapper.Cluster.ClusterId, mapping)

	return mapping
}

/*
{
	"clusters": {
		"<cluster_id>": {
	 		"hosts": {
				<The data from the function GetHostsCnodes>
	 		}
   		}
	}
}
*/
func (m *Metadata) GetClusterNodesCnodes(ctx context.Context, nodeIds []string) map[string]interface{} {
	logger.Info(ctx, "Composing cluster %s nodes", m.ClusterWrapper.Cluster.ClusterId)

	data := make(map[string]interface{})

	// hosts
	hosts := m.GetHostsCnodes(ctx, nodeIds)
	data[RegisterNodeHosts] = hosts

	cnodes := map[string]interface{}{
		RegisterClustersRootPath: map[string]interface{}{
			m.ClusterWrapper.Cluster.ClusterId: data,
		},
	}
	logger.Info(ctx, "Composed cluster %s nodes cnodes successful: %+v", m.ClusterWrapper.Cluster.ClusterId, cnodes)

	return cnodes
}

/*
{
	"clusters": {
		"<cluster_id>": {
	 		"hosts": {
				<The data from the function GetEmptyHostsCnodes>
	 		}
   		}
	}
}
*/
func (m *Metadata) GetEmptyClusterNodeCnodes(ctx context.Context, nodeIds []string) map[string]interface{} {
	logger.Info(ctx, "Composing cluster %s empty nodes", m.ClusterWrapper.Cluster.ClusterId)

	data := make(map[string]interface{})

	// hosts
	hosts := m.GetEmptyHostsCnodes(nodeIds)
	data[RegisterNodeHosts] = hosts

	cnodes := map[string]interface{}{
		RegisterClustersRootPath: map[string]interface{}{
			m.ClusterWrapper.Cluster.ClusterId: data,
		},
	}
	logger.Info(ctx, "Composed cluster %s empty nodes cnodes successful: %+v", m.ClusterWrapper.Cluster.ClusterId, cnodes)

	return cnodes
}

/*
{
    "<role>": {
    	"<instance_id>": {
			"ip":<ip>,
			"server_id":<server_id>,
			"pub_key": <pub_key>
	  	}
  	}
}
or (without role)
{
  	"<instance_id>": {
		"ip":<ip>,
	 	"server_id":<server_id>
  	}
}
*/
func (m *Metadata) GetHostsCnodes(ctx context.Context, nodeIds []string) map[string]interface{} {
	hosts := make(map[string]interface{})
	for _, nodeId := range nodeIds {
		clusterNode := m.ClusterWrapper.ClusterNodesWithKeyPairs[nodeId]
		instanceId := clusterNode.InstanceId
		role := clusterNode.Role
		if strings.HasSuffix(role, constants.ReplicaRoleSuffix) {
			role = string([]byte(role)[:len(role)-len(constants.ReplicaRoleSuffix)])
		}
		clusterRole, exist := m.ClusterWrapper.ClusterRoles[role]
		if !exist {
			logger.Error(ctx, "No such role [%s] in cluster role [%s]. ",
				role, m.ClusterWrapper.Cluster.ClusterId)
			return nil
		}
		clusterCommon, exist := m.ClusterWrapper.ClusterCommons[role]
		if !exist {
			logger.Error(ctx, "No such role [%s] in cluster common [%s]. ",
				role, m.ClusterWrapper.Cluster.ClusterId)
			return nil
		}

		host := map[string]interface{}{
			"ip":            clusterNode.PrivateIp,
			"sid":           clusterNode.ServerId,
			"gid":           clusterNode.GroupId,
			"gsid":          clusterNode.GlobalServerId,
			"node_id":       clusterNode.NodeId,
			"instance_id":   instanceId,
			"cpu":           clusterRole.Cpu,
			"gpu":           clusterRole.Gpu,
			"memory":        clusterRole.Memory,
			"volume_size":   clusterRole.StorageSize,
			"instance_size": clusterRole.InstanceSize,
		}
		if clusterCommon.Passphraseless != "" {
			host["pub_key"] = clusterNode.PubKey
		}
		if clusterNode.CustomMetadata != "" {
			host["token"] = clusterNode.CustomMetadata
		}

		if role == "" {
			hosts[instanceId] = host
		} else {
			host["role"] = role
			cnodes, exist := hosts[role]
			if exist {
				switch v := cnodes.(type) {
				case map[string]interface{}:
					v[instanceId] = host
				default:
					logger.Error(ctx, "Cnodes [%s] should be a map. ", clusterNode.NodeId)
					return nil
				}
			} else {
				hosts[role] = map[string]interface{}{instanceId: host}
			}
		}
	}
	return hosts
}

/*
{
    "<role>": {
    	"<instance_id>": ""
  	}
}
or (without role)
{
  	"<instance_id>": ""
}
*/
func (m *Metadata) GetEmptyHostsCnodes(nodeIds []string) map[string]interface{} {
	hosts := make(map[string]interface{})
	for _, nodeId := range nodeIds {
		clusterNode := m.ClusterWrapper.ClusterNodesWithKeyPairs[nodeId]
		instanceId := clusterNode.InstanceId
		role := clusterNode.Role
		if strings.HasSuffix(role, constants.ReplicaRoleSuffix) {
			role = string([]byte(role)[:len(role)-len(constants.ReplicaRoleSuffix)])
		}
		if role == "" {
			hosts[instanceId] = ""
		} else {
			hosts[role] = map[string]interface{}{instanceId: ""}
		}
	}
	return hosts
}

/*
{
	"cluster_id":  <cluster_id>,
	"app_id":      <app_id>,
	"subnet":      <subnet>,
	"user_id":     <user_id>,
	"global_uuid": <global_uuid>,
	"zone":        <zone>,
	"provider":    <provider>,
	"runtime_url": <runtime_url>,
}
*/
func (m *Metadata) GetClusterMetadataCnodes() map[string]interface{} {
	clusterMetadata := map[string]interface{}{
		"cluster_id":  m.ClusterWrapper.Cluster.ClusterId,
		"app_id":      m.ClusterWrapper.Cluster.AppId,
		"subnet":      m.ClusterWrapper.Cluster.SubnetId,
		"user_id":     m.ClusterWrapper.Cluster.Owner,
		"global_uuid": m.ClusterWrapper.Cluster.GlobalUuid,
		"zone":        m.RuntimeDetails.Zone,
		"provider":    m.RuntimeDetails.Runtime.Provider,
		"runtime_url": m.RuntimeDetails.RuntimeUrl,
	}

	return clusterMetadata
}

/*
{
  	"<role>": {
	 	"p1": "v1",
	 	"p2": "v2"
  	}
}

or (without role)

{
  	"p1": "v1",
  	"p2": "v2"
}
*/
func (m *Metadata) GetEnvCnodes(ctx context.Context) map[string]interface{} {
	result := make(map[string]interface{})
	for _, clusterRole := range m.ClusterWrapper.ClusterRoles {
		env := clusterRole.Env
		if env != "" {
			envMap := make(map[string]interface{})
			err := jsonutil.Decode([]byte(env), &envMap)
			if err != nil {
				logger.Error(ctx, "Unmarshal cluster [%s] env failed:%+v", m.ClusterWrapper.Cluster.ClusterId, err)
				return nil
			}
			if clusterRole.Role == "" {
				result = envMap
			} else {
				result[clusterRole.Role] = envMap
			}
		}
	}
	return result
}

/*
{
  	clusters: {
		<cluster_id>: {
			<RegisterNodeAdding/RegisterNodeDeleting>: {
				<The data from the function GetHostsCnodes>
			}
		}
	}
}
*/
func (m *Metadata) GetScalingCnodes(ctx context.Context, nodeIds []string, path string) map[string]interface{} {
	hosts := m.GetHostsCnodes(ctx, nodeIds)
	if len(hosts) == 0 {
		return nil
	}
	return map[string]interface{}{
		RegisterClustersRootPath: map[string]interface{}{
			m.ClusterWrapper.Cluster.ClusterId: map[string]interface{}{
				path: hosts,
			},
		},
	}
}

func (m *Metadata) GetCmdCnodes(nodeId string, cmd *models.Cmd) *models.CmdCnodes {
	clusterId := m.ClusterWrapper.Cluster.ClusterId
	clusterNode := m.ClusterWrapper.ClusterNodesWithKeyPairs[nodeId]
	instanceId := clusterNode.InstanceId

	cmdCnodes := &models.CmdCnodes{
		RootPath:   RegisterClustersRootPath,
		ClusterId:  clusterId,
		CmdKey:     RegisterNodeCmd,
		InstanceId: instanceId,
		Cmd:        cmd,
	}
	return cmdCnodes
}

/*
{
  	clusters: {
		<cluster_id>: ""
	}
}
*/
func (m *Metadata) GetEmptyClusterCnodes() map[string]interface{} {
	return map[string]interface{}{
		RegisterClustersRootPath: map[string]interface{}{
			m.ClusterWrapper.Cluster.ClusterId: "",
		},
	}
}

/*
{
  	self: {
    	<ip>: ""
  	}
}
*/
func (m *Metadata) GetEmptyClusterMappingCnodes() map[string]interface{} {
	cnodes := make(map[string]interface{})
	for _, clusterNode := range m.ClusterWrapper.ClusterNodesWithKeyPairs {
		cnodes[clusterNode.PrivateIp] = ""
	}

	return cnodes
}

/*
{
  	self: {
    	<ip>: ""
  	}
}
*/
func (m *Metadata) GetEmptyClusterNodeMappingCnodes(ctx context.Context, nodeIds []string) map[string]interface{} {
	cnodes := make(map[string]interface{})
	for _, nodeId := range nodeIds {
		cnodes[m.ClusterWrapper.ClusterNodesWithKeyPairs[nodeId].PrivateIp] = ""
	}
	return cnodes
}

/*
{
	"<ip>":
	{
		"host":"/clusters/</cluster_id>/hosts/master/<instance_id>",
		"hosts":"/clusters/</cluster_id>/hosts",
		"cluster":"/clusters/</cluster_id>/cluster",
		"env":"/clusters/</cluster_id>/env/<role>",
		"cmd":"/clusters/<cluster_id>/cmd/<instance_id>,
		"links":"/clusters/<cluster_id>/links"
	}
}
*/
func (m *Metadata) GetMappingCnodes(nodeIds []string) map[string]interface{} {
	cnodes := make(map[string]interface{})
	for _, nodeId := range nodeIds {
		clusterNode := m.ClusterWrapper.ClusterNodesWithKeyPairs[nodeId]
		clusterId := clusterNode.ClusterId
		instanceId := clusterNode.InstanceId
		role := clusterNode.Role

		var hostTarget, envTarget string
		if len(role) > 0 {
			hostTarget = fmt.Sprintf("/%s/%s/%s/%s/%s", RegisterClustersRootPath, clusterId, RegisterNodeHosts, role, instanceId)
			envTarget = fmt.Sprintf("/%s/%s/%s/%s", RegisterClustersRootPath, clusterId, RegisterNodeEnv, role)
		} else {
			hostTarget = fmt.Sprintf("/%s/%s/%s/%s", RegisterClustersRootPath, clusterId, RegisterNodeHosts, instanceId)
			envTarget = fmt.Sprintf("/%s/%s/%s", RegisterClustersRootPath, clusterId, RegisterNodeEnv)
		}

		mapping := map[string]interface{}{
			RegisterNodeHost:     hostTarget,
			RegisterNodeEnv:      envTarget,
			RegisterNodeHosts:    fmt.Sprintf("/%s/%s/%s", RegisterClustersRootPath, clusterId, RegisterNodeHosts),
			RegisterNodeCluster:  fmt.Sprintf("/%s/%s/%s", RegisterClustersRootPath, clusterId, RegisterNodeCluster),
			RegisterNodeCmd:      fmt.Sprintf("/%s/%s/%s/%s", RegisterClustersRootPath, clusterId, RegisterNodeCmd, instanceId),
			RegisterNodeAdding:   fmt.Sprintf("/%s/%s/%s", RegisterClustersRootPath, clusterId, RegisterNodeAdding),
			RegisterNodeDeleting: fmt.Sprintf("/%s/%s/%s", RegisterClustersRootPath, clusterId, RegisterNodeDeleting),
			RegisterNodeScaling:  fmt.Sprintf("/%s/%s/%s", RegisterClustersRootPath, clusterId, RegisterNodeScaling),
			RegisterNodeStopping: fmt.Sprintf("/%s/%s/%s", RegisterClustersRootPath, clusterId, RegisterNodeStopping),
			RegisterNodeStarting: fmt.Sprintf("/%s/%s/%s", RegisterClustersRootPath, clusterId, RegisterNodeStarting),
		}

		for name, clusterLink := range m.ClusterWrapper.ClusterLinks {
			if len(clusterLink.ExternalClusterId) == 0 {
				continue
			}
			_, isExist := mapping[RegisterNodeLinks]
			if !isExist {
				mapping[RegisterNodeLinks] = make(map[string]interface{})
			}
			mapping[RegisterNodeLinks].(map[string]interface{})[name] = fmt.Sprintf("/%s/%s", RegisterClustersRootPath, clusterLink.ExternalClusterId)
		}
		cnodes[clusterNode.PrivateIp] = mapping
	}

	return cnodes
}
