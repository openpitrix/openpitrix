// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package models

const ClusterCommonTableName = "cluster_common"

type ClusterCommon struct {
	ClusterId                  string
	Role                       string
	ServerIdUpperBound         int32
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
