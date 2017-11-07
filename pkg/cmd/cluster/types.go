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
