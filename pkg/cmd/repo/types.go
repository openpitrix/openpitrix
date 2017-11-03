// Copyright 2017 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package repo

import (
	"github.com/golang/protobuf/ptypes"

	db "openpitrix.io/openpitrix/pkg/db/repo"
	pb "openpitrix.io/openpitrix/pkg/service.pb"
)

func To_database_Repo(dst *db.Repo, src *pb.Repo) *db.Repo {
	if dst == nil {
		dst = new(db.Repo)
	}

	dst.Id = src.Id
	dst.Name = src.Name
	dst.Description = src.Description
	dst.Url = src.Url

	dst.Created, _ = ptypes.Timestamp(src.Created)
	dst.LastModified, _ = ptypes.Timestamp(src.LastModified)

	return dst
}

func To_proto_Repo(dst *pb.Repo, src *db.Repo) *pb.Repo {
	if dst == nil {
		dst = new(pb.Repo)
	}

	dst.Id = src.Id
	dst.Name = src.Name
	dst.Description = src.Description
	dst.Url = src.Url

	dst.Created, _ = ptypes.TimestampProto(src.Created)
	dst.LastModified, _ = ptypes.TimestampProto(src.LastModified)

	return dst
}

func To_proto_RepoList(p []db.Repo, pageNumber, pageSize int) []*pb.Repo {
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

	q := make([]*pb.Repo, end-start)
	for i := start; i < end; i++ {
		q[i-start] = To_proto_Repo(nil, &p[i])
	}
	return q
}
