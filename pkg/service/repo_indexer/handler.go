// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package repo_indexer

import (
	"context"

	"openpitrix.io/openpitrix/pkg/gerr"
	"openpitrix.io/openpitrix/pkg/manager"
	"openpitrix.io/openpitrix/pkg/models"
	"openpitrix.io/openpitrix/pkg/pb"
	"openpitrix.io/openpitrix/pkg/util/pbutil"
	"openpitrix.io/openpitrix/pkg/util/senderutil"
)

func (p *Server) IndexRepo(ctx context.Context, req *pb.IndexRepoRequest) (*pb.IndexRepoResponse, error) {
	s := senderutil.GetSenderFromContext(ctx)
	repoId := req.GetRepoId().GetValue()
	repoEvent, err := p.controller.NewRepoEvent(repoId, s.UserId)
	if err != nil {
		return nil, gerr.NewWithDetail(gerr.Internal, err, gerr.ErrorCreateResourcesFailed)
	}
	ret := pb.IndexRepoResponse{
		RepoEvent: models.RepoEventToPb(repoEvent),
		RepoId:    req.GetRepoId(),
	}
	return &ret, nil
}

func (p *Server) DescribeRepoEvents(ctx context.Context, req *pb.DescribeRepoEventsRequest) (*pb.DescribeRepoEventsResponse, error) {
	var repoEvents []*models.RepoEvent
	offset := pbutil.GetOffsetFromRequest(req)
	limit := pbutil.GetLimitFromRequest(req)

	query := p.Db.
		Select(models.RepoEventColumns...).
		From(models.RepoEventTableName).
		Offset(offset).
		Limit(limit).
		Where(manager.BuildFilterConditions(req, models.RepoEventTableName))
	query = manager.AddQueryOrderDir(query, req, models.ColumnCreateTime)
	_, err := query.Load(&repoEvents)
	if err != nil {
		return nil, gerr.NewWithDetail(gerr.Internal, err, gerr.ErrorDescribeResourcesFailed)
	}
	count, err := query.Count()
	if err != nil {
		return nil, gerr.NewWithDetail(gerr.Internal, err, gerr.ErrorDescribeResourcesFailed)
	}
	res := &pb.DescribeRepoEventsResponse{
		RepoEventSet: models.RepoEventsToPbs(repoEvents),
		TotalCount:   count,
	}
	return res, nil
}
