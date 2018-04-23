// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package task

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

func (p *Server) CreateTask(ctx context.Context, req *pb.CreateTaskRequest) (*pb.CreateTaskResponse, error) {
	s := sender.GetSenderFromContext(ctx)
	newTask := models.NewTask(
		"",
		req.GetJobId().GetValue(),
		req.GetNodeId().GetValue(),
		req.GetTarget().GetValue(),
		req.GetTaskAction().GetValue(),
		req.GetDirective().GetValue(),
		s.UserId,
	)

	_, err := p.Db.
		InsertInto(models.TaskTableName).
		Columns(models.TaskColumns...).
		Record(newTask).
		Exec()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "CreateTask: %+v", err)
	}

	err = p.controller.queue.Enqueue(newTask.TaskId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Enqueue task [%s] failed: %+v", newTask.TaskId, err)
	}

	res := &pb.CreateTaskResponse{
		TaskId: utils.ToProtoString(newTask.TaskId),
		JobId:  utils.ToProtoString(newTask.JobId),
	}
	return res, nil
}

func (p *Server) DescribeTasks(ctx context.Context, req *pb.DescribeTasksRequest) (*pb.DescribeTasksResponse, error) {
	s := sender.GetSenderFromContext(ctx)
	var tasks []*models.Task
	offset := utils.GetOffsetFromRequest(req)
	limit := utils.GetLimitFromRequest(req)

	query := p.Db.
		Select(models.TaskColumns...).
		From(models.TaskTableName).
		Offset(offset).
		Limit(limit).
		Where(manager.BuildFilterConditions(req, models.TaskTableName)).
		Where(db.Eq("owner", s.UserId))

	_, err := query.Load(&tasks)
	if err != nil {
		// TODO: err_code should be implementation
		return nil, status.Errorf(codes.Internal, "DescribeTasks: %+v", err)
	}
	count, err := query.Count()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "DescribeTasks: %+v", err)
	}

	res := &pb.DescribeTasksResponse{
		TaskSet:    models.TasksToPbs(tasks),
		TotalCount: uint32(count),
	}
	return res, nil
}
