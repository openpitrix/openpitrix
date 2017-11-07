// Copyright 2017 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package runtime

import (
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"

	db "openpitrix.io/openpitrix/pkg/db/runtime"
	pb "openpitrix.io/openpitrix/pkg/service.pb"
)

func To_database_AppRuntime(dst *db.AppRuntime, src *pb.AppRuntime) *db.AppRuntime {
	if dst == nil {
		dst = new(db.AppRuntime)
	}

	dst.Id = src.GetId()
	dst.Name = src.GetName()
	dst.Description = src.GetDescription()
	dst.Url = src.GetUrl()

	dst.Created, _ = ptypes.Timestamp(src.Created)
	dst.LastModified, _ = ptypes.Timestamp(src.LastModified)

	return dst
}

func To_proto_AppRuntime(dst *pb.AppRuntime, src *db.AppRuntime) *pb.AppRuntime {
	if dst == nil {
		dst = new(pb.AppRuntime)
	}

	dst.Id = proto.String(src.Id)
	dst.Name = proto.String(src.Name)
	dst.Description = proto.String(src.Description)
	dst.Url = proto.String(src.Url)

	dst.Created, _ = ptypes.TimestampProto(src.Created)
	dst.LastModified, _ = ptypes.TimestampProto(src.LastModified)

	return dst
}

func To_proto_AppRuntimeList(p []db.AppRuntime, pageNumber, pageSize int) []*pb.AppRuntime {
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

	q := make([]*pb.AppRuntime, end-start)
	for i := start; i < end; i++ {
		q[i-start] = To_proto_AppRuntime(nil, &p[i])
	}
	return q
}
