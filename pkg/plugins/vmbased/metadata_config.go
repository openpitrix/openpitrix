// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package vmbased

import (
	"fmt"

	"openpitrix.io/openpitrix/pkg/constants"
	"openpitrix.io/openpitrix/pkg/models"
	pbtypes "openpitrix.io/openpitrix/pkg/pb/metadata/types"
	"openpitrix.io/openpitrix/pkg/pi"
	"openpitrix.io/openpitrix/pkg/util/jsonutil"
)

type MetadataConfig struct {
	ClusterWrapper *models.ClusterWrapper
}

func (m *MetadataConfig) GetDroneConfig(nodeId string) string {
	clusterNode := m.ClusterWrapper.ClusterNodesWithKeyPairs[nodeId]

	droneEndpoint := &pbtypes.DroneEndpoint{
		FrontgateId: m.ClusterWrapper.Cluster.FrontgateId,
		DroneIp:     clusterNode.PrivateIp,
		DronePort:   constants.DroneServicePort,
	}
	droneConfig := &pbtypes.DroneConfig{
		Id:             nodeId,
		Host:           clusterNode.PrivateIp,
		ListenPort:     constants.DroneServicePort,
		CmdInfoLogPath: ConfdCmdLogPath,
		ConfdSelfHost:  clusterNode.PrivateIp,
		LogLevel:       MetadataLogLevel,
	}
	config := &pbtypes.SetDroneConfigRequest{
		Endpoint: droneEndpoint,
		Config:   droneConfig,
	}
	return jsonutil.ToString(config)
}

func (m *MetadataConfig) GetFrontgateConfig(nodeId string) string {
	clusterNode := m.ClusterWrapper.ClusterNodesWithKeyPairs[nodeId]

	var frontgateEndpoints []*pbtypes.FrontgateEndpoint
	var etcdEndpoints []*pbtypes.EtcdEndpoint
	var backendHosts []string
	for _, node := range m.ClusterWrapper.ClusterNodesWithKeyPairs {
		frontgateNode := &pbtypes.FrontgateEndpoint{
			FrontgateId: m.ClusterWrapper.Cluster.ClusterId,
			NodeIp:      node.PrivateIp,
			NodePort:    constants.FrontgateServicePort,
		}
		frontgateEndpoints = append(frontgateEndpoints, frontgateNode)

		etcdNode := &pbtypes.EtcdEndpoint{
			Host: node.PrivateIp,
			Port: EtcdPort,
		}
		etcdEndpoints = append(etcdEndpoints, etcdNode)

		backendHosts = append(backendHosts, fmt.Sprintf("%s:%d", etcdNode.Host, MetadPort))
	}

	etcdConfig := &pbtypes.EtcdConfig{
		NodeList: etcdEndpoints,
	}

	confdConfig := &pbtypes.ConfdConfig{
		ProcessorConfig: &pbtypes.ConfdProcessorConfig{
			Confdir:       ConfdPath,
			Interval:      10,
			Noop:          false,
			Prefix:        "/self",
			SyncOnly:      false,
			LogLevel:      MetadataLogLevel,
			Onetime:       false,
			Watch:         true,
			KeepStageFile: false,
		},
		BackendConfig: &pbtypes.ConfdBackendConfig{
			Type: ConfdBackendType,
			Host: backendHosts,
		},
	}

	config := &pbtypes.FrontgateConfig{
		Id:          m.ClusterWrapper.Cluster.ClusterId,
		NodeId:      nodeId,
		Host:        clusterNode.PrivateIp,
		ListenPort:  constants.FrontgateServicePort,
		PilotHost:   pi.Global().GlobalConfig().Pilot.Ip,
		NodeList:    frontgateEndpoints,
		EtcdConfig:  etcdConfig,
		ConfdConfig: confdConfig,
		LogLevel:    MetadataLogLevel,
		AutoUpdate:  pi.Global().GlobalConfig().Cluster.FrontgateAutoUpdate,
	}
	if pi.Global().GlobalConfig().Pilot.Port > 0 {
		config.PilotPort = pi.Global().GlobalConfig().Pilot.Port
	} else {
		config.PilotPort = constants.PilotServicePort
	}

	return jsonutil.ToString(config)
}
