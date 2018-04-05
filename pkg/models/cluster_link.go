// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package models

import (
	"openpitrix.io/openpitrix/pkg/pb"
	"openpitrix.io/openpitrix/pkg/utils"
)

const ClusterLinkTableName = "cluster_link"

type ClusterLink struct {
	ClusterId         string
	Name              string
	ExternalClusterId string
	Owner             string
}

var ClusterLinkColumns = GetColumnsFromStruct(&ClusterLink{})

func ClusterLinkToPb(clusterLink *ClusterLink) *pb.ClusterLink {
	return &pb.ClusterLink{
		ClusterId:         utils.ToProtoString(clusterLink.ClusterId),
		Name:              utils.ToProtoString(clusterLink.Name),
		ExternalClusterId: utils.ToProtoString(clusterLink.ExternalClusterId),
		Owner:             utils.ToProtoString(clusterLink.Owner),
	}
}

func PbToClusterLink(pbClusterLink *pb.ClusterLink) *ClusterLink {
	return &ClusterLink{
		ClusterId:         pbClusterLink.GetClusterId().GetValue(),
		Name:              pbClusterLink.GetName().GetValue(),
		ExternalClusterId: pbClusterLink.GetExternalClusterId().GetValue(),
		Owner:             pbClusterLink.GetOwner().GetValue(),
	}
}

func ClusterLinksToPbs(clusterLinks []*ClusterLink) (pbClusterLinks []*pb.ClusterLink) {
	for _, clusterLink := range clusterLinks {
		pbClusterLinks = append(pbClusterLinks, ClusterLinkToPb(clusterLink))
	}
	return
}
