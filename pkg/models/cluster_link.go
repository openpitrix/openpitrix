// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package models

import (
	"openpitrix.io/openpitrix/pkg/db"
	"openpitrix.io/openpitrix/pkg/pb"
	"openpitrix.io/openpitrix/pkg/sender"
	"openpitrix.io/openpitrix/pkg/util/pbutil"
)

type ClusterLink struct {
	ClusterId         string
	Name              string
	ExternalClusterId string
	Owner             string
	OwnerPath         sender.OwnerPath
}

var ClusterLinkColumns = db.GetColumnsFromStruct(&ClusterLink{})

func ClusterLinkToPb(clusterLink *ClusterLink) *pb.ClusterLink {
	return &pb.ClusterLink{
		ClusterId:         pbutil.ToProtoString(clusterLink.ClusterId),
		Name:              pbutil.ToProtoString(clusterLink.Name),
		ExternalClusterId: pbutil.ToProtoString(clusterLink.ExternalClusterId),
		OwnerPath:         clusterLink.OwnerPath.ToProtoString(),
		Owner:             pbutil.ToProtoString(clusterLink.Owner),
	}
}

func PbToClusterLink(pbClusterLink *pb.ClusterLink) *ClusterLink {
	ownerPath := sender.OwnerPath(pbClusterLink.GetOwnerPath().GetValue())
	return &ClusterLink{
		ClusterId:         pbClusterLink.GetClusterId().GetValue(),
		Name:              pbClusterLink.GetName().GetValue(),
		ExternalClusterId: pbClusterLink.GetExternalClusterId().GetValue(),
		Owner:             ownerPath.Owner(),
		OwnerPath:         ownerPath,
	}
}

func ClusterLinksToPbs(clusterLinks []*ClusterLink) (pbClusterLinks []*pb.ClusterLink) {
	for _, clusterLink := range clusterLinks {
		pbClusterLinks = append(pbClusterLinks, ClusterLinkToPb(clusterLink))
	}
	return
}
