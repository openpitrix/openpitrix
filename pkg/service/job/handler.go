// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package job

import (
	"context"

	"openpitrix.io/openpitrix/pkg/constants"
	"openpitrix.io/openpitrix/pkg/gerr"
	"openpitrix.io/openpitrix/pkg/manager"
	"openpitrix.io/openpitrix/pkg/models"
	"openpitrix.io/openpitrix/pkg/pb"
	"openpitrix.io/openpitrix/pkg/pi"
	"openpitrix.io/openpitrix/pkg/util/ctxutil"
	"openpitrix.io/openpitrix/pkg/util/pbutil"
)

func (p *Server) CreateJob(ctx context.Context, req *pb.CreateJobRequest) (*pb.CreateJobResponse, error) {
	s := ctxutil.GetSender(ctx)
	newJob := models.NewJob(
		"",
		req.GetClusterId().GetValue(),
		req.GetAppId().GetValue(),
		req.GetVersionId().GetValue(),
		req.GetJobAction().GetValue(),
		req.GetDirective().GetValue(),
		req.GetProvider().GetValue(),
		s.GetOwnerPath(),
		req.GetRuntimeId().GetValue(),
	)

	_, err := pi.Global().DB(ctx).
		InsertInto(constants.TableJob).
		Record(newJob).
		Exec()
	if err != nil {
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorCreateResourcesFailed)
	}

	err = p.controller.queue.Enqueue(newJob.JobId)
	if err != nil {
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorCreateResourcesFailed)
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
	var jobs []*models.Job
	offset := pbutil.GetOffsetFromRequest(req)
	limit := pbutil.GetLimitFromRequest(req)

	displayColumns := manager.GetDisplayColumns(req.GetDisplayColumns(), models.JobColumns)
	query := pi.Global().DB(ctx).
		Select(displayColumns...).
		From(constants.TableJob).
		Offset(offset).
		Limit(limit).
		Where(manager.BuildPermissionFilter(ctx)).
		Where(manager.BuildFilterConditions(req, constants.TableJob))
	query = manager.AddQueryOrderDir(query, req, constants.ColumnCreateTime)
	if len(displayColumns) > 0 {
		_, err := query.Load(&jobs)
		if err != nil {
			return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorDescribeResourcesFailed)
		}
	}
	count, err := query.Count()
	if err != nil {
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorDescribeResourcesFailed)
	}

	res := &pb.DescribeJobsResponse{
		JobSet:     models.JobsToPbs(jobs),
		TotalCount: uint32(count),
	}
	return res, nil
}
