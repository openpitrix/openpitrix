// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package vmbased

import (
	"encoding/base64"
	"fmt"
	"strings"

	"openpitrix.io/openpitrix/pkg/constants"
	"openpitrix.io/openpitrix/pkg/models"
	"openpitrix.io/openpitrix/pkg/pb/types"
	"openpitrix.io/openpitrix/pkg/pi"
	"openpitrix.io/openpitrix/pkg/util/jsonutil"
)

type Frontgate struct {
	*Frame
}

/*
cat /opt/openpitrix/conf/drone.conf
IMAGE="mysql:5.7"
MOUNT_POINT="/data"
FILE_NAME="frontgate.conf"
FILE_CONF={\\"id\\":\\"cln-abcdefgh\\",\\"listen_port\\":9111,\\"pilot_host\\":192.168.0.1,\\"pilot_port\\":9110}
*/
func (f *Frontgate) getUserDataValue(nodeId string) string {
	var result string
	clusterNode := f.ClusterWrapper.ClusterNodes[nodeId]
	role := clusterNode.Role
	if strings.HasSuffix(role, constants.ReplicaRoleSuffix) {
		role = string([]byte(role)[:len(role)-len(constants.ReplicaRoleSuffix)])
	}
	clusterRole, _ := f.ClusterWrapper.ClusterRoles[role]
	clusterCommon, _ := f.ClusterWrapper.ClusterCommons[role]
	mountPoint := clusterRole.MountPoint
	// Empty string can not be a parameter
	if len(mountPoint) == 0 {
		mountPoint = "#"
	}
	imageId := clusterCommon.ImageId

	frontgateConf := make(map[string]interface{})
	frontgateConf["id"] = nodeId
	frontgateConf["listen_port"] = constants.FrontgateServicePort
	frontgateConf["pilot_host"] = pi.Global().GlobalConfig().Pilot.Ip
	frontgateConf["pilot_port"] = constants.PilotServicePort
	frontgateConfStr := strings.Replace(jsonutil.ToString(frontgateConf), "\"", "\\\\\"", -1)

	result += fmt.Sprintf("IMAGE=\"%s\"\n", imageId)
	result += fmt.Sprintf("MOUNT_POINT=\"%s\"\n", mountPoint)
	result += fmt.Sprintf("FILE_NAME=\"%s\"\n", FrontgateConfFile)
	result += fmt.Sprintf("FILE_CONF=%s\n", frontgateConfStr)

	return base64.StdEncoding.EncodeToString([]byte(result))
}

func (f *Frontgate) getConfig(nodeId string) string {
	clusterNode := f.ClusterWrapper.ClusterNodes[nodeId]

	var frontgateEndpoints []*pbtypes.FrontgateEndpoint
	var etcdEndpoints []*pbtypes.EtcdEndpoint
	var backendHosts []string
	for _, node := range f.ClusterWrapper.ClusterNodes {
		frontgateNode := &pbtypes.FrontgateEndpoint{
			FrontgateId: f.ClusterWrapper.Cluster.ClusterId,
			NodeIp:      node.PrivateIp,
			NodePort:    constants.FrontgateServicePort,
		}
		frontgateEndpoints = append(frontgateEndpoints, frontgateNode)

		etcdNode := &pbtypes.EtcdEndpoint{
			Host: node.PrivateIp,
			Port: constants.EtcdServicePort,
		}
		etcdEndpoints = append(etcdEndpoints, etcdNode)

		backendHosts = append(backendHosts, fmt.Sprintf("%s:%d", etcdNode.Host, etcdNode.Port))
	}

	etcdConfig := &pbtypes.EtcdConfig{
		NodeList: etcdEndpoints,
	}

	confdConfig := &pbtypes.ConfdConfig{
		ProcessorConfig: &pbtypes.ConfdProcessorConfig{
			Confdir:       ConfdPath,
			Interval:      10,
			Noop:          false,
			Prefix:        "",
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
		Id:          f.ClusterWrapper.Cluster.ClusterId,
		NodeId:      nodeId,
		Host:        clusterNode.PrivateIp,
		ListenPort:  constants.FrontgateServicePort,
		PilotHost:   pi.Global().GlobalConfig().Pilot.Ip,
		PilotPort:   constants.PilotServicePort,
		NodeList:    frontgateEndpoints,
		EtcdConfig:  etcdConfig,
		ConfdConfig: confdConfig,
		LogLevel:    MetadataLogLevel,
	}

	return jsonutil.ToString(config)
}

func (f *Frontgate) setFrontgateConfigLayer(nodeIds []string, failureAllowed bool) *models.TaskLayer {
	var tasks []*models.Task
	for _, nodeId := range nodeIds {
		task := &models.Task{
			JobId:          f.Job.JobId,
			Owner:          f.Job.Owner,
			TaskAction:     ActionSetFrontgateConfig,
			Target:         constants.TargetPilot,
			NodeId:         nodeId,
			Directive:      f.getConfig(nodeId),
			FailureAllowed: failureAllowed,
		}
		tasks = append(tasks, task)
	}
	return &models.TaskLayer{
		Tasks: tasks,
	}
}

func (f *Frontgate) CreateClusterLayer() *models.TaskLayer {
	var nodeIds []string
	for nodeId := range f.ClusterWrapper.ClusterNodes {
		nodeIds = append(nodeIds, nodeId)
	}
	headTaskLayer := new(models.TaskLayer)

	headTaskLayer.
		Append(f.runInstancesLayer(nodeIds, false)).      // run instance and attach volume to instance
		Append(f.setFrontgateConfigLayer(nodeIds, false)) // set frontgate config

	return headTaskLayer.Child
}

func (f *Frontgate) DeleteClusterLayer() *models.TaskLayer {
	var nodeIds []string
	for nodeId := range f.ClusterWrapper.ClusterNodes {
		nodeIds = append(nodeIds, nodeId)
	}
	headTaskLayer := new(models.TaskLayer)

	headTaskLayer.
		Append(f.deleteInstancesLayer(nodeIds, false)) // delete instance

	return headTaskLayer.Child
}
