// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package job

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"openpitrix.io/openpitrix/pkg/logger"
	"openpitrix.io/openpitrix/pkg/pb"
	"openpitrix.io/openpitrix/pkg/utils/sender"
	"openpitrix.io/openpitrix/pkg/models"
	"openpitrix.io/openpitrix/pkg/utils"
	"openpitrix.io/openpitrix/pkg/db"
)

func (p *Server) CreateJob(ctx context.Context, req *pb.CreateJobRequest) (*pb.CreateJobResponse, error) {
	s := sender.GetSenderFromContext(ctx)
	logger.Debugf("Got req: %+v", req)
	newJob := models.NewJob(req.GetClusterId(), req.GetAppId(), req.GetAppVersion(), req.GetJobAction(), req.GetDirective(), s.UserId)

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
		JobId:      newJob.JobId,
		ClusterId:  newJob.ClusterId,
		AppId:      newJob.AppId,
		AppVersion: newJob.AppVersion,
	}
	return res, nil
}

func (p *Server) DescribeJobs(ctx context.Context, req *pb.DescribeJobsRequest) (*pb.DescribeJobsResponse, error) {
	s := sender.GetSenderFromContext(ctx)
	logger.Debugf("Got req: %+v", req)
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
	if len(req.GetClusterId()) > 0 {
		query = query.Where(db.Eq("cluster_id", req.GetClusterId()))
	}
	if len(req.GetAppId()) > 0 {
		query = query.Where(db.Eq("app_id", req.GetAppId()))
	}
	if len(req.GetAppVersion()) > 0 {
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
	}
	return res, nil
}
