// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package vmbased

import (
	"fmt"
	"strings"

	pilotclient "openpitrix.io/openpitrix/pkg/client/pilot"
	"openpitrix.io/openpitrix/pkg/constants"
	"openpitrix.io/openpitrix/pkg/logger"
	"openpitrix.io/openpitrix/pkg/models"
	pbtypes "openpitrix.io/openpitrix/pkg/pb/metadata/types"
	"openpitrix.io/openpitrix/pkg/pi"
	"openpitrix.io/openpitrix/pkg/util/jsonutil"
	"openpitrix.io/openpitrix/pkg/util/retryutil"
	"openpitrix.io/openpitrix/pkg/util/sshutil"
)

type Frontgate struct {
	*Frame
}

/*
cat /opt/openpitrix/conf/frontgate.conf
IMAGE="mysql:5.7"
MOUNT_POINT="/data"
FILE_NAME="frontgate.conf"
FILE_CONF={\\"id\\":\\"cln-abcdefgh\\",\\"listen_port\\":9111,\\"pilot_host\\":192.168.0.1,\\"pilot_port\\":9110}
*/
func (f *Frontgate) getUserDataValue(nodeId string) string {
	var result string
	clusterNode := f.ClusterWrapper.ClusterNodesWithKeyPairs[nodeId]
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
	frontgateConf["id"] = f.ClusterWrapper.Cluster.ClusterId
	frontgateConf["node_id"] = nodeId
	frontgateConf["listen_port"] = constants.FrontgateServicePort
	frontgateConf["pilot_host"] = pi.Global().GlobalConfig().Pilot.Ip
	if pi.Global().GlobalConfig().Pilot.Port > 0 {
		frontgateConf["pilot_port"] = pi.Global().GlobalConfig().Pilot.Port
	} else {
		frontgateConf["pilot_port"] = constants.PilotServicePort
	}
	frontgateConfStr := strings.Replace(jsonutil.ToString(frontgateConf), "\"", "\\\\\"", -1)

	result += fmt.Sprintf("IMAGE=\"%s\"\n", imageId)
	result += fmt.Sprintf("MOUNT_POINT=\"%s\"\n", mountPoint)
	result += fmt.Sprintf("FILE_NAME=\"%s\"\n", FrontgateConfFile)
	result += fmt.Sprintf("FILE_CONF=%s\n", frontgateConfStr)

	return result
}

func (f *Frontgate) getCertificateExec() string {
	var pilotClientTLSConfig *pbtypes.PilotClientTLSConfig

	err := retryutil.Retry(3, constants.RetryInterval, func() error {
		pilotClient, err := pilotclient.NewClient()
		if err != nil {
			return err
		}
		pilotClientTLSConfig, err = pilotClient.GetPilotClientTLSConfig(f.Ctx, &pbtypes.Empty{})
		if err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		logger.Critical(f.Ctx, "Get pilot client tls config failed, %+v", err)
		return ""
	}

	exec := fmt.Sprintf(`
echo '%s' > /opt/openpitrix/conf/openpitrix-ca.crt
echo '%s' > /opt/openpitrix/conf/pilot-client.crt
echo '%s' > /opt/openpitrix/conf/pilot-client.key
`, pilotClientTLSConfig.CaCrtData, pilotClientTLSConfig.ClientCrtData, pilotClientTLSConfig.ClientKeyData)

	return exec
}

func (f *Frontgate) pingFrontgateLayer(failureAllowed bool) *models.TaskLayer {
	taskLayer := new(models.TaskLayer)

	directive := jsonutil.ToString(&models.Meta{
		ClusterId: f.ClusterWrapper.Cluster.ClusterId,
	})

	task := &models.Task{
		JobId:          f.Job.JobId,
		Owner:          f.Job.Owner,
		TaskAction:     ActionPingFrontgate,
		Target:         constants.TargetPilot,
		NodeId:         f.ClusterWrapper.Cluster.ClusterId,
		Directive:      directive,
		FailureAllowed: failureAllowed,
	}
	taskLayer.Tasks = append(taskLayer.Tasks, task)
	if len(taskLayer.Tasks) > 0 {
		return taskLayer
	} else {
		return nil
	}
}

func (f *Frontgate) setFrontgateConfigLayer(nodeIds []string, failureAllowed bool) *models.TaskLayer {
	var tasks []*models.Task
	directive := jsonutil.ToString(&models.Meta{
		ClusterId: f.ClusterWrapper.Cluster.ClusterId,
	})

	for _, nodeId := range nodeIds {
		// get frontgate config when pre task
		task := &models.Task{
			JobId:          f.Job.JobId,
			Owner:          f.Job.Owner,
			TaskAction:     ActionSetFrontgateConfig,
			Target:         constants.TargetPilot,
			NodeId:         nodeId,
			Directive:      directive,
			FailureAllowed: failureAllowed,
		}
		tasks = append(tasks, task)
	}
	return &models.TaskLayer{
		Tasks: tasks,
	}
}

func (f *Frontgate) pingMetadataBackendLayer(failureAllowed bool) *models.TaskLayer {
	taskLayer := new(models.TaskLayer)

	directive := jsonutil.ToString(&models.Meta{
		ClusterId: f.ClusterWrapper.Cluster.ClusterId,
	})

	task := &models.Task{
		JobId:          f.Job.JobId,
		Owner:          f.Job.Owner,
		TaskAction:     PingMetadataBackend,
		Target:         constants.TargetPilot,
		NodeId:         f.ClusterWrapper.Cluster.ClusterId,
		Directive:      directive,
		FailureAllowed: failureAllowed,
	}
	taskLayer.Tasks = append(taskLayer.Tasks, task)
	if len(taskLayer.Tasks) > 0 {
		return taskLayer
	} else {
		return nil
	}
}

func (f *Frontgate) removeContainerLayer(nodeIds []string, failureAllowed bool) *models.TaskLayer {
	taskLayer := new(models.TaskLayer)

	for nodeId, clusterNode := range f.ClusterWrapper.ClusterNodesWithKeyPairs {
		ip := clusterNode.PrivateIp
		cmd := fmt.Sprintf("%s \"docker rm -f default\"", HostCmdPrefix)
		request := &pbtypes.RunCommandOnFrontgateRequest{
			Endpoint: &pbtypes.FrontgateEndpoint{
				FrontgateId:     f.ClusterWrapper.Cluster.ClusterId,
				FrontgateNodeId: nodeId,
				NodeIp:          ip,
				NodePort:        constants.FrontgateServicePort,
			},
			Command:        cmd,
			TimeoutSeconds: TimeoutRemoveContainer,
		}
		directive := jsonutil.ToString(request)
		formatVolumeTask := &models.Task{
			JobId:          f.Job.JobId,
			Owner:          f.Job.Owner,
			TaskAction:     ActionRemoveContainerOnFrontgate,
			Target:         constants.TargetPilot,
			NodeId:         nodeId,
			Directive:      directive,
			FailureAllowed: failureAllowed,
		}
		taskLayer.Tasks = append(taskLayer.Tasks, formatVolumeTask)
	}
	if len(taskLayer.Tasks) > 0 {
		return taskLayer
	} else {
		return nil
	}
}

func (f *Frontgate) attachKeyPairLayer(nodeKeyPairDetail *models.NodeKeyPairDetail) *models.TaskLayer {
	taskLayer := new(models.TaskLayer)
	clusterNode := nodeKeyPairDetail.ClusterNode
	request := &pbtypes.RunCommandOnFrontgateRequest{
		Endpoint: &pbtypes.FrontgateEndpoint{
			FrontgateId:     f.ClusterWrapper.Cluster.ClusterId,
			FrontgateNodeId: clusterNode.NodeId,
			NodeIp:          clusterNode.PrivateIp,
			NodePort:        constants.FrontgateServicePort,
		},
		Command:        fmt.Sprintf("%s \"%s\"", HostCmdPrefix, sshutil.DoAttachCmd(nodeKeyPairDetail.KeyPair.PubKey)),
		TimeoutSeconds: TimeoutKeyPair,
	}
	directive := jsonutil.ToString(request)
	attachKeyPairTask := &models.Task{
		JobId:          f.Job.JobId,
		Owner:          f.Job.Owner,
		TaskAction:     ActionRunCommandOnFrontgateNode,
		Target:         constants.TargetPilot,
		NodeId:         clusterNode.NodeId,
		Directive:      directive,
		FailureAllowed: false,
	}
	taskLayer.Tasks = append(taskLayer.Tasks, attachKeyPairTask)
	if len(taskLayer.Tasks) > 0 {
		return taskLayer
	} else {
		return nil
	}
}

func (f *Frontgate) detachKeyPairLayer(nodeKeyPairDetail *models.NodeKeyPairDetail) *models.TaskLayer {
	taskLayer := new(models.TaskLayer)
	clusterNode := nodeKeyPairDetail.ClusterNode
	request := &pbtypes.RunCommandOnFrontgateRequest{
		Endpoint: &pbtypes.FrontgateEndpoint{
			FrontgateId:     f.ClusterWrapper.Cluster.ClusterId,
			FrontgateNodeId: clusterNode.NodeId,
			NodeIp:          clusterNode.PrivateIp,
			NodePort:        constants.FrontgateServicePort,
		},
		Command:        fmt.Sprintf("%s \"%s\"", HostCmdPrefix, sshutil.DoDetachCmd(nodeKeyPairDetail.KeyPair.PubKey)),
		TimeoutSeconds: TimeoutKeyPair,
	}
	directive := jsonutil.ToString(request)
	attachKeyPairTask := &models.Task{
		JobId:          f.Job.JobId,
		Owner:          f.Job.Owner,
		TaskAction:     ActionRunCommandOnFrontgateNode,
		Target:         constants.TargetPilot,
		NodeId:         clusterNode.NodeId,
		Directive:      directive,
		FailureAllowed: false,
	}
	taskLayer.Tasks = append(taskLayer.Tasks, attachKeyPairTask)
	if len(taskLayer.Tasks) > 0 {
		return taskLayer
	} else {
		return nil
	}
}

func (f *Frontgate) CreateClusterLayer() *models.TaskLayer {
	var nodeIds []string
	for nodeId := range f.ClusterWrapper.ClusterNodesWithKeyPairs {
		nodeIds = append(nodeIds, nodeId)
	}
	headTaskLayer := new(models.TaskLayer)

	headTaskLayer.
		Append(f.createVolumesLayer(nodeIds, false)).        // create volume
		Append(f.runInstancesLayer(nodeIds, false)).         // run instance and attach volume to instance
		Append(f.pingFrontgateLayer(false)).                 // ping frontgate
		Append(f.setFrontgateConfigLayer(nodeIds, false)).   // set frontgate config
		Append(f.formatAndMountVolumeLayer(nodeIds, false)). // format and mount volume to instance
		Append(f.removeContainerLayer(nodeIds, false)).      // remove default container
		Append(f.pingFrontgateLayer(false)).                 // ping frontgate
		Append(f.pingMetadataBackendLayer(false))            // ping metadata backend

	return headTaskLayer.Child
}

func (f *Frontgate) DeleteClusterLayer() *models.TaskLayer {
	var nodeIds []string
	for nodeId := range f.ClusterWrapper.ClusterNodesWithKeyPairs {
		nodeIds = append(nodeIds, nodeId)
	}
	headTaskLayer := new(models.TaskLayer)

	if f.ClusterWrapper.Cluster.Status == constants.StatusActive {
		headTaskLayer.
			Append(f.umountVolumeLayer(nodeIds, true)).  // umount volume from instance
			Append(f.stopInstancesLayer(nodeIds, true)). // stop instance
			Append(f.detachVolumesLayer(nodeIds, false)) // detach volume from instance
	}

	headTaskLayer.
		Append(f.deleteInstancesLayer(nodeIds, false)). // delete instance
		Append(f.deleteVolumesLayer(nodeIds, false))    // delete volume
	return headTaskLayer.Child
}

func (f *Frontgate) StartClusterLayer() *models.TaskLayer {
	var nodeIds []string
	for nodeId := range f.ClusterWrapper.ClusterNodesWithKeyPairs {
		nodeIds = append(nodeIds, nodeId)
	}
	headTaskLayer := new(models.TaskLayer)

	headTaskLayer.
		Append(f.attachVolumesLayer(nodeIds, false)).      // attach volume to instance, will auto mount
		Append(f.startInstancesLayer(nodeIds, false)).     // run instance and attach volume to instance
		Append(f.pingFrontgateLayer(false)).               // ping frontgate
		Append(f.setFrontgateConfigLayer(nodeIds, false)). // set frontgate config
		Append(f.pingMetadataBackendLayer(false))          // ping metadata backend

	return headTaskLayer.Child
}

func (f *Frontgate) StopClusterLayer() *models.TaskLayer {
	var nodeIds []string
	for nodeId := range f.ClusterWrapper.ClusterNodesWithKeyPairs {
		nodeIds = append(nodeIds, nodeId)
	}
	headTaskLayer := new(models.TaskLayer)

	headTaskLayer.
		Append(f.umountVolumeLayer(nodeIds, true)).   // umount volume from instance
		Append(f.stopInstancesLayer(nodeIds, false)). // delete instance
		Append(f.detachVolumesLayer(nodeIds, false))  // detach volume from instance

	return headTaskLayer.Child
}

func (f *Frontgate) AttachKeyPairsLayer(nodeKeyPairDetails models.NodeKeyPairDetails) *models.TaskLayer {
	headTaskLayer := new(models.TaskLayer)

	for _, nodeKeyPairDetail := range nodeKeyPairDetails {
		headTaskLayer.Append(f.attachKeyPairLayer(&nodeKeyPairDetail))
	}

	return headTaskLayer.Child
}

func (f *Frontgate) DetachKeyPairsLayer(nodeKeyPairDetails models.NodeKeyPairDetails) *models.TaskLayer {
	headTaskLayer := new(models.TaskLayer)

	for _, nodeKeyPairDetail := range nodeKeyPairDetails {
		headTaskLayer.Append(f.detachKeyPairLayer(&nodeKeyPairDetail))
	}

	return headTaskLayer.Child
}
