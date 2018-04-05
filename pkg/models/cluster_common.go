// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package models

import (
	"reflect"

	"openpitrix.io/openpitrix/pkg/pb"
	"openpitrix.io/openpitrix/pkg/utils"
)

const ClusterCommonTableName = "cluster_common"

type ClusterCommon struct {
	ClusterId                  string
	Role                       string
	ServerIdUpperBound         uint32
	AdvancedActions            string
	InitService                string
	StartService               string
	StopService                string
	ScaleOutService            string
	ScaleInService             string
	RestartService             string
	DestroyService             string
	UpgradeService             string
	CustomService              string
	BackupService              string
	RestoreService             string
	DeleteSnapshotService      string
	HealthCheck                string
	Monitor                    string
	Passphraseless             string
	VerticalScalingPolicy      string
	AgentInstalled             bool
	CustomMetadataScript       string
	ImageId                    string
	BackupPolicy               string
	IncrementalBackupSupported bool
	Hypervisor                 string
}

var ClusterCommonColumns = GetColumnsFromStruct(&ClusterCommon{})

func (c ClusterCommon) GetAttribute(attributeName string) interface{} {
	common := reflect.ValueOf(c)
	service := common.FieldByName(attributeName).Interface()
	return service
}

func ClusterCommonToPb(clusterCommon *ClusterCommon) *pb.ClusterCommon {
	return &pb.ClusterCommon{
		ClusterId:                  utils.ToProtoString(clusterCommon.ClusterId),
		Role:                       utils.ToProtoString(clusterCommon.Role),
		ServerIdUpperBound:         utils.ToProtoUInt32(clusterCommon.ServerIdUpperBound),
		AdvancedActions:            utils.ToProtoString(clusterCommon.AdvancedActions),
		InitService:                utils.ToProtoString(clusterCommon.InitService),
		StartService:               utils.ToProtoString(clusterCommon.StartService),
		StopService:                utils.ToProtoString(clusterCommon.StopService),
		ScaleOutService:            utils.ToProtoString(clusterCommon.ScaleOutService),
		ScaleInService:             utils.ToProtoString(clusterCommon.ScaleInService),
		RestartService:             utils.ToProtoString(clusterCommon.RestartService),
		DestroyService:             utils.ToProtoString(clusterCommon.DestroyService),
		UpgradeService:             utils.ToProtoString(clusterCommon.UpgradeService),
		CustomService:              utils.ToProtoString(clusterCommon.CustomService),
		BackupService:              utils.ToProtoString(clusterCommon.BackupService),
		RestoreService:             utils.ToProtoString(clusterCommon.RestoreService),
		DeleteSnapshotService:      utils.ToProtoString(clusterCommon.DeleteSnapshotService),
		HealthCheck:                utils.ToProtoString(clusterCommon.HealthCheck),
		Monitor:                    utils.ToProtoString(clusterCommon.Monitor),
		Passphraseless:             utils.ToProtoString(clusterCommon.Passphraseless),
		VerticalScalingPolicy:      utils.ToProtoString(clusterCommon.VerticalScalingPolicy),
		AgentInstalled:             utils.ToProtoBool(clusterCommon.AgentInstalled),
		CustomMetadataScript:       utils.ToProtoString(clusterCommon.CustomMetadataScript),
		ImageId:                    utils.ToProtoString(clusterCommon.ImageId),
		BackupPolicy:               utils.ToProtoString(clusterCommon.BackupPolicy),
		IncrementalBackupSupported: utils.ToProtoBool(clusterCommon.IncrementalBackupSupported),
		Hypervisor:                 utils.ToProtoString(clusterCommon.Hypervisor),
	}
}

func PbToClusterCommon(pbClusterCommon *pb.ClusterCommon) *ClusterCommon {
	return &ClusterCommon{
		ClusterId:                  pbClusterCommon.GetClusterId().GetValue(),
		Role:                       pbClusterCommon.GetRole().GetValue(),
		ServerIdUpperBound:         pbClusterCommon.GetServerIdUpperBound().GetValue(),
		AdvancedActions:            pbClusterCommon.GetAdvancedActions().GetValue(),
		InitService:                pbClusterCommon.GetInitService().GetValue(),
		StartService:               pbClusterCommon.GetStartService().GetValue(),
		StopService:                pbClusterCommon.GetStopService().GetValue(),
		ScaleOutService:            pbClusterCommon.GetScaleOutService().GetValue(),
		ScaleInService:             pbClusterCommon.GetScaleInService().GetValue(),
		RestartService:             pbClusterCommon.GetRestartService().GetValue(),
		DestroyService:             pbClusterCommon.GetDestroyService().GetValue(),
		UpgradeService:             pbClusterCommon.GetUpgradeService().GetValue(),
		CustomService:              pbClusterCommon.GetCustomService().GetValue(),
		BackupService:              pbClusterCommon.GetBackupService().GetValue(),
		RestoreService:             pbClusterCommon.GetRestoreService().GetValue(),
		DeleteSnapshotService:      pbClusterCommon.GetDeleteSnapshotService().GetValue(),
		HealthCheck:                pbClusterCommon.GetHealthCheck().GetValue(),
		Monitor:                    pbClusterCommon.GetMonitor().GetValue(),
		Passphraseless:             pbClusterCommon.GetPassphraseless().GetValue(),
		VerticalScalingPolicy:      pbClusterCommon.GetVerticalScalingPolicy().GetValue(),
		AgentInstalled:             pbClusterCommon.GetAgentInstalled().GetValue(),
		CustomMetadataScript:       pbClusterCommon.GetCustomService().GetValue(),
		ImageId:                    pbClusterCommon.GetImageId().GetValue(),
		BackupPolicy:               pbClusterCommon.GetBackupPolicy().GetValue(),
		IncrementalBackupSupported: pbClusterCommon.GetIncrementalBackupSupported().GetValue(),
		Hypervisor:                 pbClusterCommon.GetHypervisor().GetValue(),
	}
}

func ClusterCommonsToPbs(clusterCommons []*ClusterCommon) (pbClusterCommons []*pb.ClusterCommon) {
	for _, clusterCommon := range clusterCommons {
		pbClusterCommons = append(pbClusterCommons, ClusterCommonToPb(clusterCommon))
	}
	return
}
