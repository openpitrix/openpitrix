// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package job

import (
	"context"

	"openpitrix.io/openpitrix/pkg/db"
	"openpitrix.io/openpitrix/pkg/gerr"
	"openpitrix.io/openpitrix/pkg/manager"
	"openpitrix.io/openpitrix/pkg/models"
	"openpitrix.io/openpitrix/pkg/pb"
	"openpitrix.io/openpitrix/pkg/util/pbutil"
	"openpitrix.io/openpitrix/pkg/util/senderutil"
)

func (p *Server) CreateJob(ctx context.Context, req *pb.CreateJobRequest) (*pb.CreateJobResponse, error) {
	s := senderutil.GetSenderFromContext(ctx)
	newJob := models.NewJob(
		"",
		req.GetClusterId().GetValue(),
		req.GetAppId().GetValue(),
		req.GetVersionId().GetValue(),
		req.GetJobAction().GetValue(),
		req.GetDirective().GetValue(),
		req.GetProvider().GetValue(),
		s.UserId,
	)

	_, err := p.Db.
		InsertInto(models.JobTableName).
		Columns(models.JobColumns...).
		Record(newJob).
		Exec()
	if err != nil {
		return nil, gerr.NewWithDetail(gerr.Internal, err, gerr.ErrorCreateResourcesFailed)
	}

	err = p.controller.queue.Enqueue(newJob.JobId)
	if err != nil {
		return nil, gerr.NewWithDetail(gerr.Internal, err, gerr.ErrorCreateResourcesFailed)
	}

	res := &pb.CreateJobResponse{
		JobId:     pbutil.ToProtoString(newJob.JobId),
		ClusterId: pbutil.ToProtoString(newJob.ClusterId),
		AppId:     pbutil.ToProtoString(newJob.AppId),
		VersionId: pbutil.ToProtoString(newJob.VersionId),
	}
	return res, nil
}

func (p *Server) DescribeJobs(ctx context.Context, req *pb.DescribeJobsRequest) (*pb.DescribeJobsResponse, error) {
	s := senderutil.GetSenderFromContext(ctx)
	var jobs []*models.Job
	offset := pbutil.GetOffsetFromRequest(req)
	limit := pbutil.GetLimitFromRequest(req)

	query := p.Db.
		Select(models.JobColumns...).
		From(models.JobTableName).
		Offset(offset).
		Limit(limit).
		Where(manager.BuildFilterConditions(req, models.JobTableName)).
		Where(db.Eq("owner", s.UserId))

	_, err := query.Load(&jobs)
	if err != nil {
		return nil, gerr.NewWithDetail(gerr.Internal, err, gerr.ErrorDescribeResourcesFailed)
	}
	count, err := query.Count()
	if err != nil {
		return nil, gerr.NewWithDetail(gerr.Internal, err, gerr.ErrorDescribeResourcesFailed)
	}

	res := &pb.DescribeJobsResponse{
		JobSet:     models.JobsToPbs(jobs),
		TotalCount: uint32(count),
	}
	return res, nil
}
