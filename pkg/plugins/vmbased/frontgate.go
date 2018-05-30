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
	frontgateConf["id"] = f.ClusterWrapper.Cluster.ClusterId
	frontgateConf["node_id"] = nodeId
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

func (f *Frame) pingFrontgateLayer(failureAllowed bool) *models.TaskLayer {
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

func (f *Frontgate) CreateClusterLayer() *models.TaskLayer {
	var nodeIds []string
	for nodeId := range f.ClusterWrapper.ClusterNodes {
		nodeIds = append(nodeIds, nodeId)
	}
	headTaskLayer := new(models.TaskLayer)

	headTaskLayer.
		Append(f.runInstancesLayer(nodeIds, false)).      // run instance and attach volume to instance
		Append(f.pingFrontgateLayer(false)).              // ping frontgate
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

func (f *Frontgate) StartClusterLayer() *models.TaskLayer {
	var nodeIds []string
	for nodeId := range f.ClusterWrapper.ClusterNodes {
		nodeIds = append(nodeIds, nodeId)
	}
	headTaskLayer := new(models.TaskLayer)

	headTaskLayer.
		Append(f.startInstancesLayer(false)).             // run instance and attach volume to instance
		Append(f.pingFrontgateLayer(false)).              // ping frontgate
		Append(f.setFrontgateConfigLayer(nodeIds, false)) // set frontgate config

	return headTaskLayer.Child
}

func (f *Frontgate) StopClusterLayer() *models.TaskLayer {
	var nodeIds []string
	for nodeId := range f.ClusterWrapper.ClusterNodes {
		nodeIds = append(nodeIds, nodeId)
	}
	headTaskLayer := new(models.TaskLayer)

	headTaskLayer.
		Append(f.stopInstancesLayer(nodeIds, false)) // delete instance

	return headTaskLayer.Child
}
