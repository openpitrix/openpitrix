// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package vmbased

import (
	"strings"

	"openpitrix.io/openpitrix/pkg/constants"
	"openpitrix.io/openpitrix/pkg/logger"
	"openpitrix.io/openpitrix/pkg/models"
	"openpitrix.io/openpitrix/pkg/util/jsonutil"
	"openpitrix.io/openpitrix/pkg/util/reflectutil"
)

type MetadataV1 struct {
	ClusterWrapper *models.ClusterWrapper
}

/*
Compose cluster info into the following format,
in order to register cluster to configuration management service.
{
  "<ip>": {
	 "hosts": {
		<The data from the function GetHostsCnodes below>
	 },
	 "host": {
		<The data from the function GetHostCnodes below>
	 },
	 "cluster": {
		<The data from the function GetClusterMetadataCnodes below>
	 },
	 "env": { # optional
		<The data from the function GetEnvCnodes below>
	 }
   }
}
*/
func (m *MetadataV1) GetClusterCnodes() map[string]interface{} {
	logger.Info("Composing cluster %s", m.ClusterWrapper.Cluster.ClusterId)

	data := make(map[string]interface{})

	var nodeIds []string
	for nodeId := range m.ClusterWrapper.ClusterNodes {
		nodeIds = append(nodeIds, nodeId)
	}

	hosts := m.GetHostsCnodes(nodeIds)
	clusterMetadata := m.GetClusterMetadataCnodes()

	for _, nodeId := range nodeIds {

		clusterNode := m.ClusterWrapper.ClusterNodes[nodeId]
		ip := clusterNode.PrivateIp

		selfCnodes := make(map[string]interface{})

		// hosts
		selfCnodes[RegisterNodeHosts] = hosts

		// host
		host := m.GetHostCnodes(nodeId)
		selfCnodes[RegisterNodeHost] = host

		// endpoints
		_, endpoints := m.ClusterWrapper.GetEndpoints()
		if endpoints != nil {
			clusterMetadata[RegisterNodeEndpoint] = endpoints
		}

		// cluster
		selfCnodes[RegisterNodeCluster] = clusterMetadata

		// env
		selfCnodes[RegisterNodeEnv] = m.GetSelfEnvCnodes(nodeId)

		data[ip] = selfCnodes
	}

	logger.Info("Composed cluster %s cnodes successful: %+v", m.ClusterWrapper.Cluster.ClusterId, data)

	return data
}

func (m *MetadataV1) GetClusterNodeCnodes(nodeIds []string) map[string]interface{} {
	logger.Info("Composing cluster %s nodes", m.ClusterWrapper.Cluster.ClusterId)

	data := make(map[string]interface{})

	clusterCnodes := m.GetClusterCnodes()
	// hosts
	hosts := map[string]interface{}{RegisterNodeHosts: m.GetHostsCnodes(nodeIds)}

	for nodeId, clusterNode := range m.ClusterWrapper.ClusterNodes {
		ip := clusterNode.PrivateIp
		if reflectutil.In(nodeId, nodeIds) {
			data[ip] = clusterCnodes[ip]
		} else {
			data[ip] = hosts
		}
	}

	logger.Info("Composed cluster %s nodes cnodes successful: %+v", m.ClusterWrapper.Cluster.ClusterId, data)

	return data
}

func (m *MetadataV1) GetEmptyClusterNodeCnodes(nodeIds []string) map[string]interface{} {
	logger.Info("Composing cluster %s empty nodes", m.ClusterWrapper.Cluster.ClusterId)

	data := make(map[string]interface{})

	// hosts
	hosts := map[string]interface{}{RegisterNodeHosts: m.GetEmptyHostsCnodes(nodeIds)}
	for nodeId, clusterNode := range m.ClusterWrapper.ClusterNodes {
		ip := clusterNode.PrivateIp
		if reflectutil.In(nodeId, nodeIds) {
			data[ip] = ""
		} else {
			data[ip] = hosts
		}
	}
	logger.Info("Composed cluster %s empty nodes cnodes successful: %+v", m.ClusterWrapper.Cluster.ClusterId, data)

	return data
}

/*
{
  "master": {
	 "i-abcdefg": {
		"ip":<ip>,
		"server_id":<server id>,
		"pub_key": <pub_key>
	  },
	  "i-xuzabcd": {
		 "ip":<ip>,
		 "server_id":<server id>,
		 "pub_key": <pub_key>
	  }
  }
}
or (without role)
{
  "i-abcdefg": {
	 "ip":<ip>,
	 "server_id":<server id>
  },
  "i-xuzabcd": {
	 "ip":<ip>,
	 "server_id":<server id>
  }
}
*/
func (m *MetadataV1) GetHostsCnodes(nodeIds []string) map[string]interface{} {
	hosts := make(map[string]interface{})
	for _, nodeId := range nodeIds {
		clusterNode := m.ClusterWrapper.ClusterNodes[nodeId]
		instanceId := clusterNode.InstanceId
		role := clusterNode.Role
		if strings.HasSuffix(role, constants.ReplicaRoleSuffix) {
			role = string([]byte(role)[:len(role)-len(constants.ReplicaRoleSuffix)])
		}
		clusterRole, exist := m.ClusterWrapper.ClusterRoles[role]
		if !exist {
			logger.Error("No such role [%s] in cluster role [%s]. ", role, m.ClusterWrapper.Cluster.ClusterId)
			return nil
		}
		clusterCommon, exist := m.ClusterWrapper.ClusterCommons[role]
		if !exist {
			logger.Error("No such role [%s] in cluster common [%s]. ", role, m.ClusterWrapper.Cluster.ClusterId)
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
					logger.Error("Cnodes [%s] should be a map. ", clusterNode.NodeId)
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
  "ip":<ip>,
  "server_id":<server id>
}
*/
func (m *MetadataV1) GetHostCnodes(nodeId string) map[string]interface{} {
	clusterNode := m.ClusterWrapper.ClusterNodes[nodeId]
	instanceId := clusterNode.InstanceId
	role := clusterNode.Role
	if strings.HasSuffix(role, constants.ReplicaRoleSuffix) {
		role = string([]byte(role)[:len(role)-len(constants.ReplicaRoleSuffix)])
	}
	clusterRole, exist := m.ClusterWrapper.ClusterRoles[role]
	if !exist {
		logger.Error("No such role [%s] in cluster role [%s]. ", role, m.ClusterWrapper.Cluster.ClusterId)
		return nil
	}
	clusterCommon, exist := m.ClusterWrapper.ClusterCommons[role]
	if !exist {
		logger.Error("No such role [%s] in cluster common [%s]. ", role, m.ClusterWrapper.Cluster.ClusterId)
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

	return host
}

func (m *MetadataV1) GetEmptyHostsCnodes(nodeIds []string) map[string]interface{} {
	hosts := make(map[string]interface{})
	for _, nodeId := range nodeIds {
		clusterNode := m.ClusterWrapper.ClusterNodes[nodeId]
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

func (m *MetadataV1) GetClusterMetadataCnodes() map[string]interface{} {
	clusterMetadata := map[string]interface{}{
		"cluster_id":  m.ClusterWrapper.Cluster.ClusterId,
		"app_id":      m.ClusterWrapper.Cluster.AppId,
		"vxnet":       m.ClusterWrapper.Cluster.SubnetId,
		"user_id":     m.ClusterWrapper.Cluster.Owner,
		"runtime_id":  m.ClusterWrapper.Cluster.RuntimeId,
		"global_uuid": m.ClusterWrapper.Cluster.GlobalUuid,
	}
	// TODO: api_server in runtime is needed

	return clusterMetadata
}

/*
{
  "p1": "v1",
  "p2": "v2"
}
*/
func (m *MetadataV1) GetSelfEnvCnodes(nodeId string) map[string]interface{} {
	result := make(map[string]interface{})
	clusterNode := m.ClusterWrapper.ClusterNodes[nodeId]
	role := clusterNode.Role
	if strings.HasSuffix(role, constants.ReplicaRoleSuffix) {
		role = string([]byte(role)[:len(role)-len(constants.ReplicaRoleSuffix)])
	}
	clusterRole := m.ClusterWrapper.ClusterRoles[role]
	env := clusterRole.Env
	if env != "" {
		err := jsonutil.Decode([]byte(env), &result)
		if err != nil {
			logger.Error("Unmarshal cluster [%s] env failed: %+v", m.ClusterWrapper.Cluster.ClusterId, err)
			return nil
		}
	}
	return result
}

func (m *MetadataV1) GetScalingCnodes(nodeIds []string, path string) map[string]interface{} {
	hosts := m.GetHostsCnodes(nodeIds)
	if len(hosts) == 0 {
		return nil
	}

	data := make(map[string]interface{})
	for _, clusterNode := range m.ClusterWrapper.ClusterNodes {
		ip := clusterNode.PrivateIp
		data[ip] = map[string]interface{}{path: hosts}
	}
	return data
}

func (m *MetadataV1) GetEmptyClusterCnodes() map[string]interface{} {
	data := make(map[string]interface{})
	for _, clusterNode := range m.ClusterWrapper.ClusterNodes {
		ip := clusterNode.PrivateIp
		data[ip] = ""
	}
	return data
}

func GetCmdCnodes(ip string, cmd *models.Cmd) map[string]interface{} {
	if cmd == nil {
		// deregister
		return map[string]interface{}{
			ip: map[string]interface{}{
				RegisterNodeCmd: "",
			},
		}
	} else {
		return map[string]interface{}{
			ip: map[string]interface{}{
				RegisterNodeCmd: cmd,
			},
		}
	}
}
