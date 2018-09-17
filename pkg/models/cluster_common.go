// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package models

import (
	"reflect"

	"openpitrix.io/openpitrix/pkg/db"
	"openpitrix.io/openpitrix/pkg/pb"
	"openpitrix.io/openpitrix/pkg/util/pbutil"
)

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

var ClusterCommonColumns = db.GetColumnsFromStruct(&ClusterCommon{})

func (c ClusterCommon) GetAttribute(attributeName string) interface{} {
	common := reflect.ValueOf(c)
	service := common.FieldByName(attributeName).Interface()
	return service
}

func ClusterCommonToPb(clusterCommon *ClusterCommon) *pb.ClusterCommon {
	return &pb.ClusterCommon{
		ClusterId:                  pbutil.ToProtoString(clusterCommon.ClusterId),
		Role:                       pbutil.ToProtoString(clusterCommon.Role),
		ServerIdUpperBound:         pbutil.ToProtoUInt32(clusterCommon.ServerIdUpperBound),
		AdvancedActions:            pbutil.ToProtoString(clusterCommon.AdvancedActions),
		InitService:                pbutil.ToProtoString(clusterCommon.InitService),
		StartService:               pbutil.ToProtoString(clusterCommon.StartService),
		StopService:                pbutil.ToProtoString(clusterCommon.StopService),
		ScaleOutService:            pbutil.ToProtoString(clusterCommon.ScaleOutService),
		ScaleInService:             pbutil.ToProtoString(clusterCommon.ScaleInService),
		RestartService:             pbutil.ToProtoString(clusterCommon.RestartService),
		DestroyService:             pbutil.ToProtoString(clusterCommon.DestroyService),
		UpgradeService:             pbutil.ToProtoString(clusterCommon.UpgradeService),
		CustomService:              pbutil.ToProtoString(clusterCommon.CustomService),
		BackupService:              pbutil.ToProtoString(clusterCommon.BackupService),
		RestoreService:             pbutil.ToProtoString(clusterCommon.RestoreService),
		DeleteSnapshotService:      pbutil.ToProtoString(clusterCommon.DeleteSnapshotService),
		HealthCheck:                pbutil.ToProtoString(clusterCommon.HealthCheck),
		Monitor:                    pbutil.ToProtoString(clusterCommon.Monitor),
		Passphraseless:             pbutil.ToProtoString(clusterCommon.Passphraseless),
		VerticalScalingPolicy:      pbutil.ToProtoString(clusterCommon.VerticalScalingPolicy),
		AgentInstalled:             pbutil.ToProtoBool(clusterCommon.AgentInstalled),
		CustomMetadataScript:       pbutil.ToProtoString(clusterCommon.CustomMetadataScript),
		ImageId:                    pbutil.ToProtoString(clusterCommon.ImageId),
		BackupPolicy:               pbutil.ToProtoString(clusterCommon.BackupPolicy),
		IncrementalBackupSupported: pbutil.ToProtoBool(clusterCommon.IncrementalBackupSupported),
		Hypervisor:                 pbutil.ToProtoString(clusterCommon.Hypervisor),
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
