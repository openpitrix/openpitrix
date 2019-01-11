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

func NewClusterId() string {
	return idutil.GetUuid36("cl-")
}

type Cluster struct {
	Zone               string
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
	OwnerPath          sender.OwnerPath
	GlobalUuid         string
	UpgradeStatus      string
	UpgradeTime        *time.Time
	RuntimeId          string
	CreateTime         time.Time
	StatusTime         time.Time
	AdditionalInfo     string
	Env                string
	Debug              bool
}

var ClusterColumns = db.GetColumnsFromStruct(&Cluster{})

func NewCluster() *Cluster {
	return &Cluster{
		CreateTime: time.Now(),
		StatusTime: time.Now(),
	}
}

func ClusterToPb(cluster *Cluster) *pb.Cluster {
	if cluster == nil {
		return new(pb.Cluster)
	}
	c := &pb.Cluster{
		ClusterId:          pbutil.ToProtoString(cluster.ClusterId),
		Name:               pbutil.ToProtoString(cluster.Name),
		Description:        pbutil.ToProtoString(cluster.Description),
		AppId:              pbutil.ToProtoString(cluster.AppId),
		VersionId:          pbutil.ToProtoString(cluster.VersionId),
		SubnetId:           pbutil.ToProtoString(cluster.SubnetId),
		VpcId:              pbutil.ToProtoString(cluster.VpcId),
		FrontgateId:        pbutil.ToProtoString(cluster.FrontgateId),
		ClusterType:        pbutil.ToProtoUInt32(cluster.ClusterType),
		Endpoints:          pbutil.ToProtoString(cluster.Endpoints),
		Status:             pbutil.ToProtoString(cluster.Status),
		TransitionStatus:   pbutil.ToProtoString(cluster.TransitionStatus),
		MetadataRootAccess: pbutil.ToProtoBool(cluster.MetadataRootAccess),
		OwnerPath:          cluster.OwnerPath.ToProtoString(),
		GlobalUuid:         pbutil.ToProtoString(cluster.GlobalUuid),
		UpgradeStatus:      pbutil.ToProtoString(cluster.UpgradeStatus),
		RuntimeId:          pbutil.ToProtoString(cluster.RuntimeId),
		CreateTime:         pbutil.ToProtoTimestamp(cluster.CreateTime),
		StatusTime:         pbutil.ToProtoTimestamp(cluster.StatusTime),
		AdditionalInfo:     pbutil.ToProtoString(cluster.AdditionalInfo),
		Env:                pbutil.ToProtoString(cluster.Env),
		Debug:              pbutil.ToProtoBool(cluster.Debug),
	}
	if cluster.UpgradeTime != nil {
		c.UpgradeTime = pbutil.ToProtoTimestamp(*cluster.UpgradeTime)
	}
	return c
}

func PbToCluster(pbCluster *pb.Cluster) *Cluster {
	if pbCluster == nil {
		return new(Cluster)
	}
	ownerPath := sender.OwnerPath(pbCluster.GetOwnerPath().GetValue())
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
		OwnerPath:          ownerPath,
		Owner:              ownerPath.Owner(),
		GlobalUuid:         pbCluster.GetGlobalUuid().GetValue(),
		UpgradeStatus:      pbCluster.GetUpgradeStatus().GetValue(),
		RuntimeId:          pbCluster.GetRuntimeId().GetValue(),
		CreateTime:         pbutil.GetTime(pbCluster.GetCreateTime()),
		StatusTime:         pbutil.GetTime(pbCluster.GetStatusTime()),
		AdditionalInfo:     pbCluster.GetAdditionalInfo().GetValue(),
		Env:                pbCluster.GetEnv().GetValue(),
		Debug:              pbCluster.GetDebug().GetValue(),
	}

	upgradeTime := pbutil.GetTime(pbCluster.GetUpgradeTime())
	c.UpgradeTime = &upgradeTime

	return c
}
