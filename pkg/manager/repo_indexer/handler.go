// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package repo_indexer

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"openpitrix.io/openpitrix/pkg/logger"

	"openpitrix.io/openpitrix/pkg/manager"
	"openpitrix.io/openpitrix/pkg/models"
	"openpitrix.io/openpitrix/pkg/utils"

	"openpitrix.io/openpitrix/pkg/pb"
)

func (p *Server) IndexRepo(ctx context.Context, req *pb.IndexRepoRequest) (*pb.IndexRepoResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "")
}

func (p *Server) DescribeRepoTasks(ctx context.Context, req *pb.DescribeRepoTasksRequest) (*pb.DescribeRepoTasksResponse, error) {
	var repoTasks []*models.RepoTask
	offset := utils.GetOffsetFromRequest(req)
	limit := utils.GetLimitFromRequest(req)

	query := p.Db.
		Select(models.RepoTaskColumns...).
		From(models.RepoTaskTableName).
		Offset(offset).
		Limit(limit).
		Where(manager.BuildFilterConditions(req, models.RepoTaskTableName))
	_, err := query.Load(&repoTasks)
	if err != nil {
		logger.Errorf("DescribeRepoTasks error: %+v", err)
		return nil, status.Errorf(codes.Internal, "DescribeRepoTasks: %+v", err)
	}
	count, err := query.Count()
	if err != nil {
		logger.Errorf("DescribeRepoTasks error: %+v", err)
		return nil, status.Errorf(codes.Internal, "DescribeRepoTasks: %+v", err)
	}
	res := &pb.DescribeRepoTasksResponse{
		RepoTaskSet: models.RepoTasksToPbs(repoTasks),
		TotalCount:  count,
	}
	return res, nil
}
