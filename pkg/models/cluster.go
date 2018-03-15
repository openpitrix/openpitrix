// Copyright 2017 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package models

import (
	"time"

	"openpitrix.io/openpitrix/pkg/pb"
	"openpitrix.io/openpitrix/pkg/utils"
)

const ClusterTableName = "cluster"

func NewClusterId() string {
	return utils.GetUuid("cl-")
}

type Cluster struct {
	ClusterId          string
	Name               string
	Description        string
	AppId              string
	AppVersion         string
	FrontgateId        string
	ClusterType        int32
	Endpoints          string
	Status             string
	TransitionStatus   string
	MetadataRootAccess int32
	Owner              string
	GlobalUuid         string
	UpgradeStatus      string
	UpgradeTime        time.Time
	RuntimeEnvId       string
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
	pbCluster := pb.Cluster{}
	return &pbCluster
}

func ClustersToPbs(clusters []*Cluster) (pbClusters []*pb.Cluster) {
	for _, cluster := range clusters {
		pbClusters = append(pbClusters, ClusterToPb(cluster))
	}
	return
}
