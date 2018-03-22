// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package models

import (
	"time"

	"openpitrix.io/openpitrix/pkg/pb"
	"openpitrix.io/openpitrix/pkg/utils"
)

const ClusterNodeTableName = "cluster_node"

func NewClusterNodeId() string {
	return utils.GetUuid("cln-")
}

type ClusterNode struct {
	NodeId           string
	ClusterId        string
	Name             string
	InstanceId       string
	VolumeId         string
	SubnetId         string
	PrivateIp        string
	ServerId         int32
	Role             string
	Status           string
	TransitionStatus string
	GroupId          int32
	Owner            string
	GlobalServerId   string
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
	pbClusterNode := pb.ClusterNode{}
	return &pbClusterNode
}

func ClusterNodesToPbs(clusterNodes []*ClusterNode) (pbClusterNodes []*pb.ClusterNode) {
	for _, clusterNode := range clusterNodes {
		pbClusterNodes = append(pbClusterNodes, ClusterNodeToPb(clusterNode))
	}
	return
}
