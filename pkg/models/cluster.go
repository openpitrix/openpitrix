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

const ClusterTableName = "cluster"

func NewClusterId() string {
	return idtool.GetUuid36("cl-")
}

type Cluster struct {
	ClusterId          string
	Name               string
	Description        string
	AppId              string
	VersionId          string
	SubnetId           string
	VpcId              string
	FrontgateId        string
	ClusterType        uint32
	Endpoints          string
	Status             string
	TransitionStatus   string
	MetadataRootAccess bool
	Owner              string
	GlobalUuid         string
	UpgradeStatus      string
	UpgradeTime        *time.Time
	RuntimeId          string
	CreateTime         time.Time
	StatusTime         time.Time
}

var ClusterColumns = GetColumnsFromStruct(&Cluster{})

func NewCluster() *Cluster {
	return &Cluster{
		CreateTime: time.Now(),
		StatusTime: time.Now(),
	}
}

func ClusterToPb(cluster *Cluster) *pb.Cluster {
	c := &pb.Cluster{
		ClusterId:          utils.ToProtoString(cluster.ClusterId),
		Name:               utils.ToProtoString(cluster.Name),
		Description:        utils.ToProtoString(cluster.Description),
		AppId:              utils.ToProtoString(cluster.AppId),
		VersionId:          utils.ToProtoString(cluster.VersionId),
		SubnetId:           utils.ToProtoString(cluster.SubnetId),
		VpcId:              utils.ToProtoString(cluster.VpcId),
		FrontgateId:        utils.ToProtoString(cluster.FrontgateId),
		ClusterType:        utils.ToProtoUInt32(cluster.ClusterType),
		Endpoints:          utils.ToProtoString(cluster.Endpoints),
		Status:             utils.ToProtoString(cluster.Status),
		TransitionStatus:   utils.ToProtoString(cluster.TransitionStatus),
		MetadataRootAccess: utils.ToProtoBool(cluster.MetadataRootAccess),
		Owner:              utils.ToProtoString(cluster.Owner),
		GlobalUuid:         utils.ToProtoString(cluster.GlobalUuid),
		UpgradeStatus:      utils.ToProtoString(cluster.UpgradeStatus),
		RuntimeId:          utils.ToProtoString(cluster.RuntimeId),
		CreateTime:         utils.ToProtoTimestamp(cluster.CreateTime),
		StatusTime:         utils.ToProtoTimestamp(cluster.StatusTime),
	}
	if cluster.UpgradeTime != nil {
		c.UpgradeTime = utils.ToProtoTimestamp(*cluster.UpgradeTime)
	}
	return c
}

func PbToCluster(pbCluster *pb.Cluster) *Cluster {
	c := &Cluster{
		ClusterId:          pbCluster.GetClusterId().GetValue(),
		Name:               pbCluster.GetName().GetValue(),
		Description:        pbCluster.GetDescription().GetValue(),
		AppId:              pbCluster.GetAppId().GetValue(),
		VersionId:          pbCluster.GetVersionId().GetValue(),
		SubnetId:           pbCluster.GetSubnetId().GetValue(),
		VpcId:              pbCluster.GetVpcId().GetValue(),
		FrontgateId:        pbCluster.GetFrontgateId().GetValue(),
		ClusterType:        pbCluster.GetClusterType().GetValue(),
		Endpoints:          pbCluster.GetEndpoints().GetValue(),
		Status:             pbCluster.GetStatus().GetValue(),
		TransitionStatus:   pbCluster.GetTransitionStatus().GetValue(),
		MetadataRootAccess: pbCluster.GetMetadataRootAccess().GetValue(),
		Owner:              pbCluster.GetOwner().GetValue(),
		GlobalUuid:         pbCluster.GetGlobalUuid().GetValue(),
		UpgradeStatus:      pbCluster.GetUpgradeStatus().GetValue(),
		RuntimeId:          pbCluster.GetRuntimeId().GetValue(),
		CreateTime:         utils.FromProtoTimestamp(pbCluster.GetCreateTime()),
		StatusTime:         utils.FromProtoTimestamp(pbCluster.GetStatusTime()),
	}
	upgradeTime := utils.FromProtoTimestamp(pbCluster.GetUpgradeTime())
	c.UpgradeTime = &upgradeTime
	return c
}
