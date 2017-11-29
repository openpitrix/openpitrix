// Copyright 2017 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package cluster

import (
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"

	db "openpitrix.io/openpitrix/pkg/db/cluster"
	pb "openpitrix.io/openpitrix/pkg/service.pb"
)

func To_database_Cluster(dst *db.Cluster, src *pb.Cluster) *db.Cluster {
	if dst == nil {
		dst = new(db.Cluster)
	}

	dst.Id = src.GetId()
	dst.Name = src.GetName()
	dst.Description = src.GetDescription()
	dst.AppId = src.GetAppId()
	dst.AppVersion = src.GetAppVersion()
	dst.Status = src.GetStatus()
	dst.TransitionStatus = src.GetTransitionStatus()

	dst.Created, _ = ptypes.Timestamp(src.Created)
	dst.LastModified, _ = ptypes.Timestamp(src.LastModified)

	return dst
}

func To_proto_Cluster(dst *pb.Cluster, src *db.Cluster) *pb.Cluster {
	if dst == nil {
		dst = new(pb.Cluster)
	}

	dst.Id = proto.String(src.Id)
	dst.Name = proto.String(src.Name)
	dst.Description = proto.String(src.Description)
	dst.AppId = proto.String(src.AppId)
	dst.AppVersion = proto.String(src.AppVersion)
	dst.Status = proto.String(src.Status)
	dst.TransitionStatus = proto.String(src.TransitionStatus)

	dst.Created, _ = ptypes.TimestampProto(src.Created)
	dst.LastModified, _ = ptypes.TimestampProto(src.LastModified)

	return dst
}

func To_proto_ClusterList(p []db.Cluster, pageNumber, pageSize int) []*pb.Cluster {
	if pageNumber > 0 {
		pageNumber = pageNumber - 1 // start with 1
	}

	start := pageNumber * pageSize
	end := start + pageSize

	if start >= len(p) {
		return nil
	}
	if end > len(p) {
		end = len(p)
	}

	q := make([]*pb.Cluster, end-start)
	for i := start; i < end; i++ {
		q[i-start] = To_proto_Cluster(nil, &p[i])
	}
	return q
}

func To_proto_Clusters(src []db.Cluster) *pb.Clusters {
	dst := make([]*pb.Cluster, len(src))
	for i := 0; i < len(src); i++ {
		dst[i] = To_proto_Cluster(nil, &src[i])
	}
	clusters := pb.Clusters{
		Items: dst,
	}
	return &clusters
}

func To_database_Clusters(src *pb.Clusters) []*db.Cluster {
	dst := make([]*db.Cluster, len(src.Items))
	for i := 0; i < len(src.Items); i++ {
		dst[i] = To_database_Cluster(nil, src.Items[i])
	}
	return dst
}

func To_database_ClusterNode(dst *db.ClusterNode, src *pb.ClusterNode) *db.ClusterNode {
	if dst == nil {
		dst = new(db.ClusterNode)
	}

	dst.Id = src.GetId()
	dst.InstanceId = src.GetInstanceId()
	dst.Name = src.GetName()
	dst.Description = src.GetDescription()
	dst.ClusterId = src.GetClusterId()
	dst.PrivateIp = src.GetPrivateIp()
	dst.Status = src.GetStatus()
	dst.TransitionStatus = src.GetTransitionStatus()

	dst.Created, _ = ptypes.Timestamp(src.Created)
	dst.LastModified, _ = ptypes.Timestamp(src.LastModified)

	return dst
}

func To_proto_ClusterNode(dst *pb.ClusterNode, src *db.ClusterNode) *pb.ClusterNode {
	if dst == nil {
		dst = new(pb.ClusterNode)
	}

	dst.Id = proto.String(src.Id)
	dst.InstanceId = proto.String(src.InstanceId)
	dst.Name = proto.String(src.Name)
	dst.Description = proto.String(src.Description)
	dst.ClusterId = proto.String(src.ClusterId)
	dst.PrivateIp = proto.String(src.PrivateIp)
	dst.Status = proto.String(src.Status)
	dst.TransitionStatus = proto.String(src.TransitionStatus)

	dst.Created, _ = ptypes.TimestampProto(src.Created)
	dst.LastModified, _ = ptypes.TimestampProto(src.LastModified)

	return dst
}

func To_proto_ClusterNodeList(p []db.ClusterNode, pageNumber, pageSize int) []*pb.ClusterNode {
	if pageNumber > 0 {
		pageNumber = pageNumber - 1 // start with 1
	}

	start := pageNumber * pageSize
	end := start + pageSize

	if start >= len(p) {
		return nil
	}
	if end > len(p) {
		end = len(p)
	}

	q := make([]*pb.ClusterNode, end-start)
	for i := start; i < end; i++ {
		q[i-start] = To_proto_ClusterNode(nil, &p[i])
	}
	return q
}

func To_proto_ClusterNodes(src []db.ClusterNode) *pb.ClusterNodes {
	dst := make([]*pb.ClusterNode, len(src))
	for i := 0; i < len(src); i++ {
		dst[i] = To_proto_ClusterNode(nil, &src[i])
	}
	clusterNodes := pb.ClusterNodes{
		Items: dst,
	}
	return &clusterNodes
}

func To_database_ClusterNodes(src *pb.ClusterNodes) []*db.ClusterNode {
	dst := make([]*db.ClusterNode, len(src.Items))
	for i := 0; i < len(src.Items); i++ {
		dst[i] = To_database_ClusterNode(nil, src.Items[i])
	}
	return dst
}
