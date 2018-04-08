// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package models

import (
	"openpitrix.io/openpitrix/pkg/pb"
	"openpitrix.io/openpitrix/pkg/utils"
)

const ClusterRoleTableName = "cluster_role"

type ClusterRole struct {
	ClusterId    string
	Role         string
	Cpu          uint32
	Gpu          uint32
	Memory       uint32
	InstanceSize uint32
	StorageSize  uint32
	MountPoint   string
	MountOptions string
	FileSystem   string
	Env          string
}

var ClusterRoleColumns = GetColumnsFromStruct(&ClusterRole{})

func ClusterRoleToPb(clusterRole *ClusterRole) *pb.ClusterRole {
	return &pb.ClusterRole{
		ClusterId:    utils.ToProtoString(clusterRole.ClusterId),
		Role:         utils.ToProtoString(clusterRole.Role),
		Cpu:          utils.ToProtoUInt32(clusterRole.Cpu),
		Gpu:          utils.ToProtoUInt32(clusterRole.Gpu),
		Memory:       utils.ToProtoUInt32(clusterRole.Memory),
		InstanceSize: utils.ToProtoUInt32(clusterRole.InstanceSize),
		StorageSize:  utils.ToProtoUInt32(clusterRole.StorageSize),
		MountPoint:   utils.ToProtoString(clusterRole.MountPoint),
		MountOptions: utils.ToProtoString(clusterRole.MountOptions),
		FileSystem:   utils.ToProtoString(clusterRole.FileSystem),
		Env:          utils.ToProtoString(clusterRole.Env),
	}
}

func PbToClusterRole(pbClusterRole *pb.ClusterRole) *ClusterRole {
	return &ClusterRole{
		ClusterId:    pbClusterRole.GetClusterId().GetValue(),
		Role:         pbClusterRole.GetRole().GetValue(),
		Cpu:          pbClusterRole.GetCpu().GetValue(),
		Gpu:          pbClusterRole.GetGpu().GetValue(),
		Memory:       pbClusterRole.GetMemory().GetValue(),
		InstanceSize: pbClusterRole.GetInstanceSize().GetValue(),
		StorageSize:  pbClusterRole.GetStorageSize().GetValue(),
		MountPoint:   pbClusterRole.GetMountPoint().GetValue(),
		MountOptions: pbClusterRole.GetMountOptions().GetValue(),
		FileSystem:   pbClusterRole.GetFileSystem().GetValue(),
		Env:          pbClusterRole.GetEnv().GetValue(),
	}
}

func ClusterRolesToPbs(clusterRoles []*ClusterRole) (pbClusterRoles []*pb.ClusterRole) {
	for _, clusterRole := range clusterRoles {
		pbClusterRoles = append(pbClusterRoles, ClusterRoleToPb(clusterRole))
	}
	return
}
