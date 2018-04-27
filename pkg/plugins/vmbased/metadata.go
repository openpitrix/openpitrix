// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package vmbased

import (
	"encoding/json"
	"strings"

	"openpitrix.io/openpitrix/pkg/constants"
	"openpitrix.io/openpitrix/pkg/logger"
	"openpitrix.io/openpitrix/pkg/models"
)

type Metadata struct {
	ClusterWrapper *models.ClusterWrapper
}

func (m *Metadata) GetClusterCnodesString() string {
	cnodes, err := json.Marshal(m.GetClusterCnodes())
	if err != nil {
		logger.Errorf("Marshal cluster [%s] metadata cnodes failed: %+v",
			m.ClusterWrapper.Cluster.ClusterId, err)
	}
	return string(cnodes)
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
	logger.Infof("Composing cluster %s", m.ClusterWrapper.Cluster.ClusterId)

	data := make(map[string]interface{})

	// hosts
	hosts := m.GetHostsCnodes()
	data[RegisterNodeHosts] = hosts

	// cluster
	clusterMetadata := m.GetClusterMetadataCnodes()
	data[RegisterNodeCluster] = clusterMetadata

	// endpoints
	_, endpoints := m.ClusterWrapper.GetEndpoints()
	if endpoints != nil {
		data[RegisterNodeEndpoint] = endpoints
	}

	cnodes := map[string]interface{}{
		RegisterClustersRootPath: map[string]interface{}{
			m.ClusterWrapper.Cluster.ClusterId: data,
		},
	}
	logger.Infof("Composed cluster %s cnodes %s successful", m.ClusterWrapper.Cluster.ClusterId)

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
func (m *Metadata) GetHostsCnodes() map[string]interface{} {
	hosts := make(map[string]interface{})
	for _, clusterNode := range m.ClusterWrapper.ClusterNodes {
		instanceId := clusterNode.InstanceId
		role := clusterNode.Role
		if strings.HasSuffix(role, constants.ReplicaRoleSuffix) {
			role = string([]byte(role)[:len(role)-len(constants.ReplicaRoleSuffix)])
		}
		clusterRole, exist := m.ClusterWrapper.ClusterRoles[role]
		if !exist {
			logger.Errorf("No such role [%s] in cluster role [%s]. ",
				role, m.ClusterWrapper.Cluster.ClusterId)
			return nil
		}
		clusterCommon, exist := m.ClusterWrapper.ClusterCommons[role]
		if !exist {
			logger.Errorf("No such role [%s] in cluster common [%s]. ",
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
					logger.Errorf("Cnodes [%s] should be a map. ", clusterNode.NodeId)
					return nil
				}
			} else {
				hosts[role] = map[string]interface{}{instanceId: host}
			}
		}
	}
	return hosts
}

func (m *Metadata) GetClusterMetadataCnodes() map[string]interface{} {
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
			err := json.Unmarshal([]byte(env), &envMap)
			if err != nil {
				logger.Errorf("Unmarshal cluster [%s] env failed:%+v", m.ClusterWrapper.Cluster.ClusterId, err)
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
