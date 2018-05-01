// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package job

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"openpitrix.io/openpitrix/pkg/db"
	"openpitrix.io/openpitrix/pkg/manager"
	"openpitrix.io/openpitrix/pkg/models"
	"openpitrix.io/openpitrix/pkg/pb"
	"openpitrix.io/openpitrix/pkg/utils"
	"openpitrix.io/openpitrix/pkg/utils/sender"
)

func (p *Server) CreateJob(ctx context.Context, req *pb.CreateJobRequest) (*pb.CreateJobResponse, error) {
	s := sender.GetSenderFromContext(ctx)
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
		return nil, status.Errorf(codes.Internal, "CreateJob [%s] failed: %+v", newJob.JobId, err)
	}

	err = p.controller.queue.Enqueue(newJob.JobId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Enqueue job [%s] failed: %+v", newJob.JobId, err)
	}

	res := &pb.CreateJobResponse{
		JobId:     utils.ToProtoString(newJob.JobId),
		ClusterId: utils.ToProtoString(newJob.ClusterId),
		AppId:     utils.ToProtoString(newJob.AppId),
		VersionId: utils.ToProtoString(newJob.VersionId),
	}
	return res, nil
}

func (p *Server) DescribeJobs(ctx context.Context, req *pb.DescribeJobsRequest) (*pb.DescribeJobsResponse, error) {
	s := sender.GetSenderFromContext(ctx)
	var jobs []*models.Job
	offset := utils.GetOffsetFromRequest(req)
	limit := utils.GetLimitFromRequest(req)

	query := p.Db.
		Select(models.JobColumns...).
		From(models.JobTableName).
		Offset(offset).
		Limit(limit).
		Where(manager.BuildFilterConditions(req, models.JobTableName)).
		Where(db.Eq("owner", s.UserId))

	_, err := query.Load(&jobs)
	if err != nil {
		// TODO: err_code should be implementation
		return nil, status.Errorf(codes.Internal, "DescribeJobs: %+v", err)
	}
	count, err := query.Count()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "DescribeJobs: %+v", err)
	}

	res := &pb.DescribeJobsResponse{
		JobSet:     models.JobsToPbs(jobs),
		TotalCount: uint32(count),
	}
	return res, nil
}
