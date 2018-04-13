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
	"openpitrix.io/openpitrix/pkg/pb"
	"openpitrix.io/openpitrix/pkg/utils"
	"openpitrix.io/openpitrix/pkg/utils/sender"
)

func (p *Server) IndexRepo(ctx context.Context, req *pb.IndexRepoRequest) (*pb.IndexRepoResponse, error) {
	s := sender.GetSenderFromContext(ctx)
	repoId := req.GetRepoId().GetValue()
	if repoId == "" {
		// TODO: api gateway params validate
		return nil, status.Errorf(codes.InvalidArgument, "Invalid argument: [repo_id]")
	}
	repoEvent, err := p.controller.NewRepoEvent(repoId, s.UserId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "IndexRepo error: %+v", err)
	}
	ret := pb.IndexRepoResponse{
		RepoEvent: models.RepoEventToPb(repoEvent),
	}
	return &ret, nil
}

func (p *Server) DescribeRepoEvents(ctx context.Context, req *pb.DescribeRepoEventsRequest) (*pb.DescribeRepoEventsResponse, error) {
	var repoEvents []*models.RepoEvent
	offset := utils.GetOffsetFromRequest(req)
	limit := utils.GetLimitFromRequest(req)

	query := p.Db.
		Select(models.RepoEventColumns...).
		From(models.RepoEventTableName).
		Offset(offset).
		Limit(limit).
		Where(manager.BuildFilterConditions(req, models.RepoEventTableName))
	_, err := query.Load(&repoEvents)
	if err != nil {
		logger.Errorf("DescribeRepoEvents error: %+v", err)
		return nil, status.Errorf(codes.Internal, "DescribeRepoEvents: %+v", err)
	}
	count, err := query.Count()
	if err != nil {
		logger.Errorf("DescribeRepoEvents error: %+v", err)
		return nil, status.Errorf(codes.Internal, "DescribeRepoEvents: %+v", err)
	}
	res := &pb.DescribeRepoEventsResponse{
		RepoEventSet: models.RepoEventsToPbs(repoEvents),
		TotalCount:   count,
	}
	return res, nil
}
