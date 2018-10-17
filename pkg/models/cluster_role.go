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

	ServiceType       string
	ServiceClusterIp  string
	ServiceExternalIp string
	ServicePorts      string

	ConfigMapDataCount uint32
	SecretDataCount    uint32

	PvcStatus      string
	PvcVolume      string
	PvcCapacity    string
	PvcAccessModes string

	IngressHosts   string
	IngressAddress string
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

		ServiceType:       pbutil.ToProtoString(clusterRole.ServiceType),
		ServiceClusterIp:  pbutil.ToProtoString(clusterRole.ServiceClusterIp),
		ServiceExternalIp: pbutil.ToProtoString(clusterRole.ServiceExternalIp),
		ServicePorts:      pbutil.ToProtoString(clusterRole.ServicePorts),

		ConfigMapDataCount: pbutil.ToProtoUInt32(clusterRole.ConfigMapDataCount),
		SecretDataCount:    pbutil.ToProtoUInt32(clusterRole.SecretDataCount),

		PvcStatus:      pbutil.ToProtoString(clusterRole.PvcStatus),
		PvcVolume:      pbutil.ToProtoString(clusterRole.PvcVolume),
		PvcCapacity:    pbutil.ToProtoString(clusterRole.PvcCapacity),
		PvcAccessModes: pbutil.ToProtoString(clusterRole.PvcAccessModes),

		IngressHosts:   pbutil.ToProtoString(clusterRole.IngressHosts),
		IngressAddress: pbutil.ToProtoString(clusterRole.IngressAddress),
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

		ServiceType:       pbClusterRole.GetServiceType().GetValue(),
		ServiceClusterIp:  pbClusterRole.GetServiceClusterIp().GetValue(),
		ServiceExternalIp: pbClusterRole.GetServiceExternalIp().GetValue(),
		ServicePorts:      pbClusterRole.GetServicePorts().GetValue(),

		ConfigMapDataCount: pbClusterRole.GetConfigMapDataCount().GetValue(),
		SecretDataCount:    pbClusterRole.GetSecretDataCount().GetValue(),

		PvcStatus:      pbClusterRole.GetPvcStatus().GetValue(),
		PvcVolume:      pbClusterRole.GetPvcVolume().GetValue(),
		PvcCapacity:    pbClusterRole.GetPvcCapacity().GetValue(),
		PvcAccessModes: pbClusterRole.GetPvcAccessModes().GetValue(),

		IngressHosts:   pbClusterRole.GetIngressHosts().GetValue(),
		IngressAddress: pbClusterRole.GetIngressAddress().GetValue(),
	}
}

func ClusterRolesToPbs(clusterRoles []*ClusterRole) (pbClusterRoles []*pb.ClusterRole) {
	for _, clusterRole := range clusterRoles {
		pbClusterRoles = append(pbClusterRoles, ClusterRoleToPb(clusterRole))
	}
	return
}
