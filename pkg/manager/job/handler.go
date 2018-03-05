// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package job

import (
	"context"

	"github.com/golang/protobuf/ptypes/wrappers"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"openpitrix.io/openpitrix/pkg/db"
	"openpitrix.io/openpitrix/pkg/logger"
	"openpitrix.io/openpitrix/pkg/manager"
	"openpitrix.io/openpitrix/pkg/models"
	"openpitrix.io/openpitrix/pkg/pb"
	"openpitrix.io/openpitrix/pkg/utils"
	"openpitrix.io/openpitrix/pkg/utils/sender"
)

func (p *Server) CreateJob(ctx context.Context, req *pb.CreateJobRequest) (*pb.CreateJobResponse, error) {
	logger.Debugf("Got req: %+v", req)

	s := sender.GetSenderFromContext(ctx)
	newJob := models.NewJob(
		req.GetClusterId().GetValue(),
		req.GetAppId().GetValue(),
		req.GetAppVersion().GetValue(),
		req.GetJobAction().GetValue(),
		req.GetDirective().GetValue(),
		s.UserId,
	)

	_, err := p.db.
		InsertInto(models.JobTableName).
		Columns(models.JobColumns...).
		Record(newJob).
		Exec()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "CreateJob: %+v", err)
	}

	//TODO: push job into jobQueue

	res := &pb.CreateJobResponse{
		JobId:      &wrappers.StringValue{Value: newJob.JobId},
		ClusterId:  &wrappers.StringValue{Value: newJob.ClusterId},
		AppId:      &wrappers.StringValue{Value: newJob.AppId},
		AppVersion: &wrappers.StringValue{Value: newJob.AppVersion},
	}
	return res, nil
}

func (p *Server) DescribeJobs(ctx context.Context, req *pb.DescribeJobsRequest) (*pb.DescribeJobsResponse, error) {
	logger.Debugf("Got req: %+v", req)

	s := sender.GetSenderFromContext(ctx)
	var jobs []*models.Job
	offset := utils.GetOffsetFromRequest(req)
	limit := utils.GetLimitFromRequest(req)

	query := p.db.
		Select(models.JobColumns...).
		From(models.JobTableName).
		Offset(offset).
		Limit(limit).
		Where(manager.BuildFilterConditions(req, models.JobColumns...))
	query = query.Where(db.Eq("owner", s.UserId))

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
