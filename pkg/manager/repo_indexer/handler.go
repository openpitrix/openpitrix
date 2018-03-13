// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package repo_indexer

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"openpitrix.io/openpitrix/pkg/pb"
)

func (p *Server) IndexRepo(ctx context.Context, req *pb.IndexRepoRequest) (*pb.IndexRepoResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "")
}

func (p *Server) DescribeRepoTasks(ctx context.Context, req *pb.DescribeRepoTasksRequest) (*pb.DescribeRepoTasksResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "")
}
