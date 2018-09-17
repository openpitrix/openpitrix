// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package models

import (
	"openpitrix.io/openpitrix/pkg/db"
	"openpitrix.io/openpitrix/pkg/pb"
	"openpitrix.io/openpitrix/pkg/util/pbutil"
)

type ClusterRole struct {
	ClusterId     string
	Role          string
	Cpu           uint32
	Gpu           uint32
	Memory        uint32
	InstanceSize  uint32
	StorageSize   uint32
	MountPoint    string
	MountOptions  string
	FileSystem    string
	Env           string
	Replicas      uint32
	ReadyReplicas uint32
	ApiVersion    string
}

var ClusterRoleColumns = db.GetColumnsFromStruct(&ClusterRole{})

func ClusterRoleToPb(clusterRole *ClusterRole) *pb.ClusterRole {
	return &pb.ClusterRole{
		ClusterId:     pbutil.ToProtoString(clusterRole.ClusterId),
		Role:          pbutil.ToProtoString(clusterRole.Role),
		Cpu:           pbutil.ToProtoUInt32(clusterRole.Cpu),
		Gpu:           pbutil.ToProtoUInt32(clusterRole.Gpu),
		Memory:        pbutil.ToProtoUInt32(clusterRole.Memory),
		InstanceSize:  pbutil.ToProtoUInt32(clusterRole.InstanceSize),
		StorageSize:   pbutil.ToProtoUInt32(clusterRole.StorageSize),
		MountPoint:    pbutil.ToProtoString(clusterRole.MountPoint),
		MountOptions:  pbutil.ToProtoString(clusterRole.MountOptions),
		FileSystem:    pbutil.ToProtoString(clusterRole.FileSystem),
		Env:           pbutil.ToProtoString(clusterRole.Env),
		Replicas:      pbutil.ToProtoUInt32(clusterRole.Replicas),
		ReadyReplicas: pbutil.ToProtoUInt32(clusterRole.ReadyReplicas),
		ApiVersion:    pbutil.ToProtoString(clusterRole.ApiVersion),
	}
}

func PbToClusterRole(pbClusterRole *pb.ClusterRole) *ClusterRole {
	return &ClusterRole{
		ClusterId:     pbClusterRole.GetClusterId().GetValue(),
		Role:          pbClusterRole.GetRole().GetValue(),
		Cpu:           pbClusterRole.GetCpu().GetValue(),
		Gpu:           pbClusterRole.GetGpu().GetValue(),
		Memory:        pbClusterRole.GetMemory().GetValue(),
		InstanceSize:  pbClusterRole.GetInstanceSize().GetValue(),
		StorageSize:   pbClusterRole.GetStorageSize().GetValue(),
		MountPoint:    pbClusterRole.GetMountPoint().GetValue(),
		MountOptions:  pbClusterRole.GetMountOptions().GetValue(),
		FileSystem:    pbClusterRole.GetFileSystem().GetValue(),
		Env:           pbClusterRole.GetEnv().GetValue(),
		Replicas:      pbClusterRole.GetReplicas().GetValue(),
		ReadyReplicas: pbClusterRole.GetReadyReplicas().GetValue(),
		ApiVersion:    pbClusterRole.GetApiVersion().GetValue(),
	}
}

func ClusterRolesToPbs(clusterRoles []*ClusterRole) (pbClusterRoles []*pb.ClusterRole) {
	for _, clusterRole := range clusterRoles {
		pbClusterRoles = append(pbClusterRoles, ClusterRoleToPb(clusterRole))
	}
	return
}
