// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package job

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/golang/protobuf/ptypes/wrappers"

	"openpitrix.io/openpitrix/pkg/db"
	"openpitrix.io/openpitrix/pkg/logger"
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
		Limit(limit)
	query = query.Where(db.Eq("owner", s.UserId))

	// TODO: filter condition
	if len(req.GetJobId()) > 0 {
		query = query.Where(db.Eq("job_id", req.GetJobId()))
	}
	if len(req.GetClusterId().GetValue()) > 0 {
		query = query.Where(db.Eq("cluster_id", req.GetClusterId()))
	}
	if len(req.GetAppId().GetValue()) > 0 {
		query = query.Where(db.Eq("app_id", req.GetAppId()))
	}
	if len(req.GetAppVersion().GetValue()) > 0 {
		query = query.Where(db.Eq("app_versoin", req.GetAppVersion()))
	}
	if len(req.GetStatus()) > 0 {
		query = query.Where(db.Eq("status", req.GetStatus()))
	}

	count, err := query.Load(&jobs)
	if err != nil {
		// TODO: err_code should be implementation
		return nil, status.Errorf(codes.Internal, "DescribeJobs: %+v", err)
	}

	res := &pb.DescribeJobsResponse{
		JobSet:     models.JobsToPbs(jobs),
		TotalCount: uint32(count),
		Limit:      uint32(limit),
	}
	return res, nil
}
