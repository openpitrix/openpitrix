// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package models

import (
	"openpitrix.io/openpitrix/pkg/db"
	"openpitrix.io/openpitrix/pkg/pb"
	"openpitrix.io/openpitrix/pkg/util/pbutil"
)

type ClusterLoadbalancer struct {
	ClusterId              string
	Role                   string
	LoadbalancerListenerId string
	LoadbalancerPort       uint32
	LoadbalancerPolicyId   string
}

var ClusterLoadbalancerColumns = db.GetColumnsFromStruct(&ClusterLoadbalancer{})

func ClusterLoadbalancerToPb(clusterLoadbalancer *ClusterLoadbalancer) *pb.ClusterLoadbalancer {
	return &pb.ClusterLoadbalancer{
		ClusterId:              pbutil.ToProtoString(clusterLoadbalancer.ClusterId),
		Role:                   pbutil.ToProtoString(clusterLoadbalancer.Role),
		LoadbalancerListenerId: pbutil.ToProtoString(clusterLoadbalancer.LoadbalancerListenerId),
		LoadbalancerPort:       pbutil.ToProtoUInt32(clusterLoadbalancer.LoadbalancerPort),
		LoadbalancerPolicyId:   pbutil.ToProtoString(clusterLoadbalancer.LoadbalancerPolicyId),
	}
}

func PbToClusterLoadbalancer(pbClusterLoadbalancer *pb.ClusterLoadbalancer) *ClusterLoadbalancer {
	return &ClusterLoadbalancer{
		ClusterId:              pbClusterLoadbalancer.GetClusterId().GetValue(),
		Role:                   pbClusterLoadbalancer.GetRole().GetValue(),
		LoadbalancerListenerId: pbClusterLoadbalancer.GetLoadbalancerListenerId().GetValue(),
		LoadbalancerPort:       pbClusterLoadbalancer.GetLoadbalancerPort().GetValue(),
		LoadbalancerPolicyId:   pbClusterLoadbalancer.GetLoadbalancerPolicyId().GetValue(),
	}
}

func ClusterLoadbalancersToPbs(clusterLoadbalancers []*ClusterLoadbalancer) (pbClusterLoadbalancers []*pb.ClusterLoadbalancer) {
	for _, clusterLoadbalancer := range clusterLoadbalancers {
		pbClusterLoadbalancers = append(pbClusterLoadbalancers, ClusterLoadbalancerToPb(clusterLoadbalancer))
	}
	return
}
