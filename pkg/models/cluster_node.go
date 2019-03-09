// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package models

import (
	"time"

	"openpitrix.io/openpitrix/pkg/db"
	"openpitrix.io/openpitrix/pkg/pb"
	"openpitrix.io/openpitrix/pkg/sender"
	"openpitrix.io/openpitrix/pkg/util/idutil"
	"openpitrix.io/openpitrix/pkg/util/pbutil"
)

func NewClusterNodeId() string {
	return idutil.GetUuid("cln-")
}

type ClusterNode struct {
	NodeId           string
	ClusterId        string
	Name             string
	InstanceId       string
	VolumeId         string
	Device           string
	SubnetId         string
	PrivateIp        string
	Eip              string
	ServerId         uint32
	Role             string
	Status           string
	TransitionStatus string
	GroupId          uint32
	Owner            string
	OwnerPath        sender.OwnerPath
	GlobalServerId   uint32
	CustomMetadata   string
	PubKey           string
	HealthStatus     string
	IsBackup         bool
	AutoBackup       bool
	CreateTime       time.Time
	StatusTime       time.Time
	HostId           string
	HostIp           string
}

type ClusterNodeWithKeyPairs struct {
	*ClusterNode
	KeyPairId []string
}

var ClusterNodeColumns = db.GetColumnsFromStruct(&ClusterNode{})

func NewClusterNode() *ClusterNode {
	return &ClusterNode{
		CreateTime: time.Now(),
		StatusTime: time.Now(),
	}
}

func ClusterNodeToPb(clusterNode *ClusterNode) *pb.ClusterNode {
	return &pb.ClusterNode{
		NodeId:           pbutil.ToProtoString(clusterNode.NodeId),
		ClusterId:        pbutil.ToProtoString(clusterNode.ClusterId),
		Name:             pbutil.ToProtoString(clusterNode.Name),
		InstanceId:       pbutil.ToProtoString(clusterNode.InstanceId),
		VolumeId:         pbutil.ToProtoString(clusterNode.VolumeId),
		Device:           pbutil.ToProtoString(clusterNode.Device),
		SubnetId:         pbutil.ToProtoString(clusterNode.SubnetId),
		PrivateIp:        pbutil.ToProtoString(clusterNode.PrivateIp),
		Eip:              pbutil.ToProtoString(clusterNode.Eip),
		ServerId:         pbutil.ToProtoUInt32(clusterNode.ServerId),
		Role:             pbutil.ToProtoString(clusterNode.Role),
		Status:           pbutil.ToProtoString(clusterNode.Status),
		TransitionStatus: pbutil.ToProtoString(clusterNode.TransitionStatus),
		GroupId:          pbutil.ToProtoUInt32(clusterNode.GroupId),
		OwnerPath:        clusterNode.OwnerPath.ToProtoString(),
		Owner:            pbutil.ToProtoString(clusterNode.Owner),
		GlobalServerId:   pbutil.ToProtoUInt32(clusterNode.GlobalServerId),
		CustomMetadata:   pbutil.ToProtoString(clusterNode.CustomMetadata),
		PubKey:           pbutil.ToProtoString(clusterNode.PubKey),
		HealthStatus:     pbutil.ToProtoString(clusterNode.HealthStatus),
		IsBackup:         pbutil.ToProtoBool(clusterNode.IsBackup),
		AutoBackup:       pbutil.ToProtoBool(clusterNode.AutoBackup),
		CreateTime:       pbutil.ToProtoTimestamp(clusterNode.CreateTime),
		StatusTime:       pbutil.ToProtoTimestamp(clusterNode.StatusTime),
		HostId:           pbutil.ToProtoString(clusterNode.HostId),
		HostIp:           pbutil.ToProtoString(clusterNode.HostIp),
	}
}

func ClusterNodeWithKeyPairsToPb(clusterNodeKeyPairs *ClusterNodeWithKeyPairs) *pb.ClusterNode {
	return &pb.ClusterNode{
		NodeId:           pbutil.ToProtoString(clusterNodeKeyPairs.NodeId),
		ClusterId:        pbutil.ToProtoString(clusterNodeKeyPairs.ClusterId),
		Name:             pbutil.ToProtoString(clusterNodeKeyPairs.Name),
		InstanceId:       pbutil.ToProtoString(clusterNodeKeyPairs.InstanceId),
		VolumeId:         pbutil.ToProtoString(clusterNodeKeyPairs.VolumeId),
		Device:           pbutil.ToProtoString(clusterNodeKeyPairs.Device),
		SubnetId:         pbutil.ToProtoString(clusterNodeKeyPairs.SubnetId),
		PrivateIp:        pbutil.ToProtoString(clusterNodeKeyPairs.PrivateIp),
		Eip:              pbutil.ToProtoString(clusterNodeKeyPairs.Eip),
		ServerId:         pbutil.ToProtoUInt32(clusterNodeKeyPairs.ServerId),
		Role:             pbutil.ToProtoString(clusterNodeKeyPairs.Role),
		Status:           pbutil.ToProtoString(clusterNodeKeyPairs.Status),
		TransitionStatus: pbutil.ToProtoString(clusterNodeKeyPairs.TransitionStatus),
		GroupId:          pbutil.ToProtoUInt32(clusterNodeKeyPairs.GroupId),
		OwnerPath:        clusterNodeKeyPairs.OwnerPath.ToProtoString(),
		Owner:            pbutil.ToProtoString(clusterNodeKeyPairs.Owner),
		GlobalServerId:   pbutil.ToProtoUInt32(clusterNodeKeyPairs.GlobalServerId),
		CustomMetadata:   pbutil.ToProtoString(clusterNodeKeyPairs.CustomMetadata),
		PubKey:           pbutil.ToProtoString(clusterNodeKeyPairs.PubKey),
		HealthStatus:     pbutil.ToProtoString(clusterNodeKeyPairs.HealthStatus),
		IsBackup:         pbutil.ToProtoBool(clusterNodeKeyPairs.IsBackup),
		AutoBackup:       pbutil.ToProtoBool(clusterNodeKeyPairs.AutoBackup),
		CreateTime:       pbutil.ToProtoTimestamp(clusterNodeKeyPairs.CreateTime),
		StatusTime:       pbutil.ToProtoTimestamp(clusterNodeKeyPairs.StatusTime),
		HostId:           pbutil.ToProtoString(clusterNodeKeyPairs.HostId),
		HostIp:           pbutil.ToProtoString(clusterNodeKeyPairs.HostIp),
		KeyPairId:        clusterNodeKeyPairs.KeyPairId,
	}
}

func PbToClusterNode(pbClusterNode *pb.ClusterNode) *ClusterNodeWithKeyPairs {
	ownerPath := sender.OwnerPath(pbClusterNode.GetOwnerPath().GetValue())
	clusterNodeKeyPairs := &ClusterNodeWithKeyPairs{
		ClusterNode: &ClusterNode{
			NodeId:           pbClusterNode.GetNodeId().GetValue(),
			ClusterId:        pbClusterNode.GetClusterId().GetValue(),
			Name:             pbClusterNode.GetName().GetValue(),
			InstanceId:       pbClusterNode.GetInstanceId().GetValue(),
			VolumeId:         pbClusterNode.GetVolumeId().GetValue(),
			Device:           pbClusterNode.GetDevice().GetValue(),
			SubnetId:         pbClusterNode.GetSubnetId().GetValue(),
			PrivateIp:        pbClusterNode.GetPrivateIp().GetValue(),
			Eip:              pbClusterNode.GetEip().GetValue(),
			ServerId:         pbClusterNode.GetServerId().GetValue(),
			Role:             pbClusterNode.GetRole().GetValue(),
			Status:           pbClusterNode.GetStatus().GetValue(),
			TransitionStatus: pbClusterNode.GetTransitionStatus().GetValue(),
			GroupId:          pbClusterNode.GetGroupId().GetValue(),
			OwnerPath:        ownerPath,
			Owner:            ownerPath.Owner(),
			GlobalServerId:   pbClusterNode.GetGlobalServerId().GetValue(),
			CustomMetadata:   pbClusterNode.GetCustomMetadata().GetValue(),
			PubKey:           pbClusterNode.GetPubKey().GetValue(),
			HealthStatus:     pbClusterNode.GetHealthStatus().GetValue(),
			IsBackup:         pbClusterNode.GetIsBackup().GetValue(),
			AutoBackup:       pbClusterNode.GetAutoBackup().GetValue(),
			CreateTime:       pbutil.GetTime(pbClusterNode.GetCreateTime()),
			StatusTime:       pbutil.GetTime(pbClusterNode.GetCreateTime()),
			HostId:           pbClusterNode.HostId.GetValue(),
			HostIp:           pbClusterNode.HostIp.GetValue(),
		},
	}
	clusterNodeKeyPairs.KeyPairId = pbClusterNode.KeyPairId
	return clusterNodeKeyPairs
}

func ClusterNodesWithKeyPairsToPbs(clusterNodeKeyPairs []*ClusterNodeWithKeyPairs) (pbClusterNodes []*pb.ClusterNode) {
	for _, clusterNodeKeyPairsItem := range clusterNodeKeyPairs {
		pbClusterNodes = append(pbClusterNodes, ClusterNodeWithKeyPairsToPb(clusterNodeKeyPairsItem))
	}
	return
}
