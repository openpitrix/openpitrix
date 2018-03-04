// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package task

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

func (p *Server) CreateTask(ctx context.Context, req *pb.CreateTaskRequest) (*pb.CreateTaskResponse, error) {
	logger.Debugf("Got req: %+v", req)

	s := sender.GetSenderFromContext(ctx)
	newTask := models.NewTask(
		req.GetJobId().GetValue(),
		req.GetTaskAction().GetValue(),
		req.GetDirective().GetValue(),
		s.UserId,
	)

	_, err := p.db.
		InsertInto(models.TaskTableName).
		Columns(models.TaskColumns...).
		Record(newTask).
		Exec()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "CreateTask: %+v", err)
	}

	//TODO: push task into taskQueue

	res := &pb.CreateTaskResponse{
		TaskId: &wrappers.StringValue{Value: newTask.TaskId},
		JobId:  &wrappers.StringValue{Value: newTask.JobId},
	}
	return res, nil
}

func (p *Server) DescribeTasks(ctx context.Context, req *pb.DescribeTasksRequest) (*pb.DescribeTasksResponse, error) {
	logger.Debugf("Got req: %+v", req)

	s := sender.GetSenderFromContext(ctx)
	var tasks []*models.Task
	offset := utils.GetOffsetFromRequest(req)
	limit := utils.GetLimitFromRequest(req)

	query := p.db.
		Select(models.TaskColumns...).
		From(models.TaskTableName).
		Offset(offset).
		Limit(limit)
	query = query.Where(db.Eq("owner", s.UserId))

	// TODO: filter condition
	if len(req.GetTaskId()) > 0 {
		query = query.Where(db.Eq("task_id", req.GetTaskId()))
	}
	if len(req.GetJobId()) > 0 {
		query = query.Where(db.Eq("job_id", req.GetJobId()))
	}
	if len(req.GetStatus()) > 0 {
		query = query.Where(db.Eq("status", req.GetStatus()))
	}

	count, err := query.Load(&tasks)
	if err != nil {
		// TODO: err_code should be implementation
		return nil, status.Errorf(codes.Internal, "DescribeTasks: %+v", err)
	}

	res := &pb.DescribeTasksResponse{
		TaskSet:    models.TasksToPbs(tasks),
		TotalCount: uint32(count),
		Limit:      uint32(limit),
	}
	return res, nil
}
