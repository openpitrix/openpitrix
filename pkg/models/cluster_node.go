// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package models

import (
	"time"

	"openpitrix.io/openpitrix/pkg/pb"
	"openpitrix.io/openpitrix/pkg/utils"
	"openpitrix.io/openpitrix/pkg/utils/idtool"
)

const ClusterNodeTableName = "cluster_node"

func NewClusterNodeId() string {
	return idtool.GetUuid("cln-")
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
		NodeId:           utils.ToProtoString(clusterNode.NodeId),
		ClusterId:        utils.ToProtoString(clusterNode.ClusterId),
		Name:             utils.ToProtoString(clusterNode.Name),
		InstanceId:       utils.ToProtoString(clusterNode.InstanceId),
		VolumeId:         utils.ToProtoString(clusterNode.VolumeId),
		Device:           utils.ToProtoString(clusterNode.Device),
		SubnetId:         utils.ToProtoString(clusterNode.SubnetId),
		PrivateIp:        utils.ToProtoString(clusterNode.PrivateIp),
		ServerId:         utils.ToProtoUInt32(clusterNode.ServerId),
		Role:             utils.ToProtoString(clusterNode.Role),
		Status:           utils.ToProtoString(clusterNode.Status),
		TransitionStatus: utils.ToProtoString(clusterNode.TransitionStatus),
		GroupId:          utils.ToProtoUInt32(clusterNode.GroupId),
		Owner:            utils.ToProtoString(clusterNode.Owner),
		GlobalServerId:   utils.ToProtoUInt32(clusterNode.GlobalServerId),
		CustomMetadata:   utils.ToProtoString(clusterNode.CustomMetadata),
		PubKey:           utils.ToProtoString(clusterNode.PubKey),
		HealthStatus:     utils.ToProtoString(clusterNode.HealthStatus),
		IsBackup:         utils.ToProtoBool(clusterNode.IsBackup),
		AutoBackup:       utils.ToProtoBool(clusterNode.AutoBackup),
		CreateTime:       utils.ToProtoTimestamp(clusterNode.CreateTime),
		StatusTime:       utils.ToProtoTimestamp(clusterNode.StatusTime),
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
		CreateTime:       utils.FromProtoTimestamp(pbClusterNode.GetCreateTime()),
		StatusTime:       utils.FromProtoTimestamp(pbClusterNode.GetCreateTime()),
	}
}

func ClusterNodesToPbs(clusterNodes []*ClusterNode) (pbClusterNodes []*pb.ClusterNode) {
	for _, clusterNode := range clusterNodes {
		pbClusterNodes = append(pbClusterNodes, ClusterNodeToPb(clusterNode))
	}
	return
}
