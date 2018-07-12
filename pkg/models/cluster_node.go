// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package models

import (
	"time"

	"openpitrix.io/openpitrix/pkg/pb"
	"openpitrix.io/openpitrix/pkg/util/idutil"
	"openpitrix.io/openpitrix/pkg/util/pbutil"
)

const ClusterNodeTableName = "cluster_node"

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
	GlobalServerId   uint32
	CustomMetadata   string
	PubKey           string
	HealthStatus     string
	IsBackup         bool
	AutoBackup       bool
	CreateTime       time.Time
	StatusTime       time.Time
}

var ClusterNodeColumns = GetColumnsFromStruct(&ClusterNode{})

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
		ServerId:         pbutil.ToProtoUInt32(clusterNode.ServerId),
		Role:             pbutil.ToProtoString(clusterNode.Role),
		Status:           pbutil.ToProtoString(clusterNode.Status),
		TransitionStatus: pbutil.ToProtoString(clusterNode.TransitionStatus),
		GroupId:          pbutil.ToProtoUInt32(clusterNode.GroupId),
		Owner:            pbutil.ToProtoString(clusterNode.Owner),
		GlobalServerId:   pbutil.ToProtoUInt32(clusterNode.GlobalServerId),
		CustomMetadata:   pbutil.ToProtoString(clusterNode.CustomMetadata),
		PubKey:           pbutil.ToProtoString(clusterNode.PubKey),
		HealthStatus:     pbutil.ToProtoString(clusterNode.HealthStatus),
		IsBackup:         pbutil.ToProtoBool(clusterNode.IsBackup),
		AutoBackup:       pbutil.ToProtoBool(clusterNode.AutoBackup),
		CreateTime:       pbutil.ToProtoTimestamp(clusterNode.CreateTime),
		StatusTime:       pbutil.ToProtoTimestamp(clusterNode.StatusTime),
	}
}

func PbToClusterNode(pbClusterNode *pb.ClusterNode) *ClusterNode {
	return &ClusterNode{
		NodeId:           pbClusterNode.GetNodeId().GetValue(),
		ClusterId:        pbClusterNode.GetClusterId().GetValue(),
		Name:             pbClusterNode.GetName().GetValue(),
		InstanceId:       pbClusterNode.GetInstanceId().GetValue(),
		VolumeId:         pbClusterNode.GetVolumeId().GetValue(),
		Device:           pbClusterNode.GetDevice().GetValue(),
		SubnetId:         pbClusterNode.GetSubnetId().GetValue(),
		PrivateIp:        pbClusterNode.GetPrivateIp().GetValue(),
		ServerId:         pbClusterNode.GetServerId().GetValue(),
		Role:             pbClusterNode.GetRole().GetValue(),
		Status:           pbClusterNode.GetStatus().GetValue(),
		TransitionStatus: pbClusterNode.GetTransitionStatus().GetValue(),
		GroupId:          pbClusterNode.GetGroupId().GetValue(),
		Owner:            pbClusterNode.GetOwner().GetValue(),
		GlobalServerId:   pbClusterNode.GetGlobalServerId().GetValue(),
		CustomMetadata:   pbClusterNode.GetCustomMetadata().GetValue(),
		PubKey:           pbClusterNode.GetPubKey().GetValue(),
		HealthStatus:     pbClusterNode.GetHealthStatus().GetValue(),
		IsBackup:         pbClusterNode.GetIsBackup().GetValue(),
		AutoBackup:       pbClusterNode.GetAutoBackup().GetValue(),
		CreateTime:       pbutil.FromProtoTimestamp(pbClusterNode.GetCreateTime()),
		StatusTime:       pbutil.FromProtoTimestamp(pbClusterNode.GetCreateTime()),
	}
}

func ClusterNodesToPbs(clusterNodes []*ClusterNode) (pbClusterNodes []*pb.ClusterNode) {
	for _, clusterNode := range clusterNodes {
		pbClusterNodes = append(pbClusterNodes, ClusterNodeToPb(clusterNode))
	}
	return
}
