// Copyright 2017 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package app

import (
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"

	db "openpitrix.io/openpitrix/pkg/db/app"
	pb "openpitrix.io/openpitrix/pkg/service.pb"
)

func To_database_App(dst *db.App, src *pb.App) *db.App {
	if dst == nil {
		dst = new(db.App)
	}

	dst.Id = src.GetId()
	dst.Name = src.GetName()
	dst.Description = src.GetDescription()
	dst.RepoId = src.GetRepoId()

	dst.Created, _ = ptypes.Timestamp(src.Created)
	dst.LastModified, _ = ptypes.Timestamp(src.LastModified)

	return dst
}

func To_proto_App(dst *pb.App, src *db.App) *pb.App {
	if dst == nil {
		dst = new(pb.App)
	}

	dst.Id = proto.String(src.Id)
	dst.Name = proto.String(src.Name)
	dst.Description = proto.String(src.Description)
	dst.RepoId = proto.String(src.RepoId)

	dst.Created, _ = ptypes.TimestampProto(src.Created)
	dst.LastModified, _ = ptypes.TimestampProto(src.LastModified)

	return dst
}

func To_proto_AppList(p []db.App, pageNumber, pageSize int) []*pb.App {
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

	q := make([]*pb.App, end-start)
	for i := start; i < end; i++ {
		q[i-start] = To_proto_App(nil, &p[i])
	}
	return q
}
