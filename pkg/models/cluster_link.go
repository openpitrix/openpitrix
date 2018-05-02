// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package models

import (
	"openpitrix.io/openpitrix/pkg/pb"
	"openpitrix.io/openpitrix/pkg/util/pbutil"
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
		ClusterId:         pbutil.ToProtoString(clusterLink.ClusterId),
		Name:              pbutil.ToProtoString(clusterLink.Name),
		ExternalClusterId: pbutil.ToProtoString(clusterLink.ExternalClusterId),
		Owner:             pbutil.ToProtoString(clusterLink.Owner),
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
