// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package vmbased

import (
	"fmt"
	"strings"

	"openpitrix.io/openpitrix/pkg/constants"
	"openpitrix.io/openpitrix/pkg/logger"
	"openpitrix.io/openpitrix/pkg/models"
	"openpitrix.io/openpitrix/pkg/util/jsonutil"
)

type Metadata struct {
	ClusterWrapper *models.ClusterWrapper
}

/*
Compose cluster info into the following format,
in order to register cluster to configuration management service.
{
  "<cluster_id>": {
	 "hosts": {
		<The data from the function GetHostsCnodes below>
	 },
	 "cluster": {
		<The data from the function GetClusterMetadataCnodes below>
	 },
	 "env": { # optional
		<The data from the function GetEnvCnodes below>
	 }
   },
   "self": {
	 "192.168.100.10": {
		<The data from the function GetClusterSelfCnodes below>
	 }
   }
}
*/
func (m *Metadata) GetClusterCnodes() map[string]interface{} {
	logger.Info("Composing cluster %s", m.ClusterWrapper.Cluster.ClusterId)

	data := make(map[string]interface{})

	var nodeIds []string
	for nodeId := range m.ClusterWrapper.ClusterNodes {
		nodeIds = append(nodeIds, nodeId)
	}
	// hosts
	hosts := m.GetHostsCnodes(nodeIds)
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
	env := m.GetEnvCnodes()
	if len(env) > 0 {
		data[RegisterNodeEnv] = env
	}

	cnodes := map[string]interface{}{
		RegisterClustersRootPath: map[string]interface{}{
			m.ClusterWrapper.Cluster.ClusterId: data,
		},
	}
	logger.Info("Composed cluster %s cnodes successful: %+v", m.ClusterWrapper.Cluster.ClusterId, cnodes)

	return cnodes
}

func (m *Metadata) GetClusterNodeCnodes(nodeIds []string) map[string]interface{} {
	logger.Info("Composing cluster %s nodes", m.ClusterWrapper.Cluster.ClusterId)

	data := make(map[string]interface{})

	// hosts
	hosts := m.GetHostsCnodes(nodeIds)
	data[RegisterNodeHosts] = hosts

	cnodes := map[string]interface{}{
		RegisterClustersRootPath: map[string]interface{}{
			m.ClusterWrapper.Cluster.ClusterId: data,
		},
	}
	logger.Info("Composed cluster %s nodes cnodes successful: %+v", m.ClusterWrapper.Cluster.ClusterId, cnodes)

	return cnodes
}

func (m *Metadata) GetEmptyClusterNodeCnodes(nodeIds []string) map[string]interface{} {
	logger.Info("Composing cluster %s empty nodes", m.ClusterWrapper.Cluster.ClusterId)

	data := make(map[string]interface{})

	// hosts
	hosts := m.GetEmptyHostsCnodes(nodeIds)
	data[RegisterNodeHosts] = hosts

	cnodes := map[string]interface{}{
		RegisterClustersRootPath: map[string]interface{}{
			m.ClusterWrapper.Cluster.ClusterId: data,
		},
	}
	logger.Info("Composed cluster %s empty nodes cnodes successful: %+v", m.ClusterWrapper.Cluster.ClusterId, cnodes)

	return cnodes
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
func (m *Metadata) GetHostsCnodes(nodeIds []string) map[string]interface{} {
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
			logger.Error("No such role [%s] in cluster role [%s]. ",
				role, m.ClusterWrapper.Cluster.ClusterId)
			return nil
		}
		clusterCommon, exist := m.ClusterWrapper.ClusterCommons[role]
		if !exist {
			logger.Error("No such role [%s] in cluster common [%s]. ",
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

func (m *Metadata) GetEmptyHostsCnodes(nodeIds []string) map[string]interface{} {
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

func (m *Metadata) GetClusterMetadataCnodes() map[string]interface{} {
	clusterMetadata := map[string]interface{}{
		"cluster_id":  m.ClusterWrapper.Cluster.ClusterId,
		"app_id":      m.ClusterWrapper.Cluster.AppId,
		"subnet":      m.ClusterWrapper.Cluster.SubnetId,
		"user_id":     m.ClusterWrapper.Cluster.Owner,
		"runtime_id":  m.ClusterWrapper.Cluster.RuntimeId,
		"global_uuid": m.ClusterWrapper.Cluster.GlobalUuid,
	}
	// TODO: api_server in runtime is needed

	return clusterMetadata
}

/*
{
  "master": {
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
func (m *Metadata) GetEnvCnodes() map[string]interface{} {
	result := make(map[string]interface{})
	for _, clusterRole := range m.ClusterWrapper.ClusterRoles {
		env := clusterRole.Env
		if env != "" {
			envMap := make(map[string]interface{})
			err := jsonutil.Decode([]byte(env), &envMap)
			if err != nil {
				logger.Error("Unmarshal cluster [%s] env failed:%+v", m.ClusterWrapper.Cluster.ClusterId, err)
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

func (m *Metadata) GetScalingCnodes(nodeIds []string, path string) map[string]interface{} {
	hosts := m.GetHostsCnodes(nodeIds)
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

func (m *Metadata) GetCmdCnodes(nodeId, cmd string) map[string]interface{} {
	clusterId := m.ClusterWrapper.Cluster.ClusterId
	clusterNode := m.ClusterWrapper.ClusterNodes[nodeId]
	instanceId := clusterNode.InstanceId

	return map[string]interface{}{
		fmt.Sprintf("/%s/%s/%s/%s", RegisterClustersRootPath, clusterId, RegisterNodeCmd, instanceId): cmd,
	}
}

func (m *Metadata) GetEmptyClusterCnodes() map[string]interface{} {
	return map[string]interface{}{
		fmt.Sprintf("/%s/%s", RegisterClustersRootPath, m.ClusterWrapper.Cluster.ClusterId): "",
	}
}
