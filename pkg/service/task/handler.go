// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package task

import (
	"context"
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"openpitrix.io/openpitrix/pkg/constants"
	"openpitrix.io/openpitrix/pkg/db"
	"openpitrix.io/openpitrix/pkg/logger"
	"openpitrix.io/openpitrix/pkg/manager"
	"openpitrix.io/openpitrix/pkg/models"
	"openpitrix.io/openpitrix/pkg/pb"
	"openpitrix.io/openpitrix/pkg/util/pbutil"
	"openpitrix.io/openpitrix/pkg/util/senderutil"
)

func (p *Server) CreateTask(ctx context.Context, req *pb.CreateTaskRequest) (*pb.CreateTaskResponse, error) {
	s := senderutil.GetSenderFromContext(ctx)
	newTask := models.NewTask(
		"",
		req.GetJobId().GetValue(),
		req.GetNodeId().GetValue(),
		req.GetTarget().GetValue(),
		req.GetTaskAction().GetValue(),
		req.GetDirective().GetValue(),
		s.UserId,
		req.GetFailureAllowed().GetValue(),
	)

	if req.GetStatus().GetValue() == constants.StatusFailed {
		newTask.Status = req.GetStatus().GetValue()
	}

	_, err := p.Db.
		InsertInto(models.TaskTableName).
		Columns(models.TaskColumns...).
		Record(newTask).
		Exec()
	if err != nil {
		logger.Error("CreateTask failed: %+v", err)
		return nil, status.Errorf(codes.Internal, "CreateTask: %+v", err)
	}

	if newTask.Status != constants.StatusFailed {
		err = p.controller.queue.Enqueue(newTask.TaskId)
		if err != nil {
			logger.Error("CreateTask [%s] failed: %+v", newTask.TaskId, err)
			return nil, status.Errorf(codes.Internal, "Enqueue task [%s] failed: %+v", newTask.TaskId, err)
		}
	}

	res := &pb.CreateTaskResponse{
		TaskId: pbutil.ToProtoString(newTask.TaskId),
		JobId:  pbutil.ToProtoString(newTask.JobId),
	}
	return res, nil
}

func (p *Server) DescribeTasks(ctx context.Context, req *pb.DescribeTasksRequest) (*pb.DescribeTasksResponse, error) {
	s := senderutil.GetSenderFromContext(ctx)
	var tasks []*models.Task
	offset := pbutil.GetOffsetFromRequest(req)
	limit := pbutil.GetLimitFromRequest(req)

	query := p.Db.
		Select(models.TaskColumns...).
		From(models.TaskTableName).
		Offset(offset).
		Limit(limit).
		Where(manager.BuildFilterConditions(req, models.TaskTableName)).
		Where(db.Eq("owner", s.UserId)).
		OrderDir("create_time", true)

	_, err := query.Load(&tasks)
	if err != nil {
		// TODO: err_code should be implementation
		logger.Error("DescribeTasks failed: %+v", err)
		return nil, status.Errorf(codes.Internal, "DescribeTasks: %+v", err)
	}
	count, err := query.Count()
	if err != nil {
		logger.Error("DescribeTasks failed: %+v", err)
		return nil, status.Errorf(codes.Internal, "DescribeTasks: %+v", err)
	}

	res := &pb.DescribeTasksResponse{
		TaskSet:    models.TasksToPbs(tasks),
		TotalCount: uint32(count),
	}
	return res, nil
}

func (p *Server) RetryTasks(ctx context.Context, req *pb.RetryTasksRequest) (*pb.RetryTasksResponse, error) {
	s := senderutil.GetSenderFromContext(ctx)

	taskIds := req.GetTaskId()
	var tasks []*models.Task
	query := p.Db.
		Select(models.TaskColumns...).
		From(models.TaskTableName).
		Where(db.Eq("task_id", taskIds)).
		Where(db.Eq("owner", s.UserId))

	_, err := query.Load(&tasks)
	if err != nil {
		// TODO: err_code should be implementation
		logger.Error("RetryTasks %s failed: %+v", taskIds, err)
		return nil, status.Errorf(codes.Internal, "RetryTasks: %+v", err)
	}

	if len(tasks) != len(taskIds) {
		logger.Error("RetryTasks %s with count [%d]", taskIds, len(tasks))
		return nil, fmt.Errorf("retryTasks %s with count [%d]", taskIds, len(tasks))
	}

	for _, taskId := range taskIds {
		err = p.controller.queue.Enqueue(taskId)
		if err != nil {
			logger.Error("Enqueue [%s] failed: %+v", taskId, err)
			return nil, status.Errorf(codes.Internal, "Enqueue task [%s] failed: %+v", taskId, err)
		}
	}

	res := &pb.RetryTasksResponse{
		TaskSet: models.TasksToPbs(tasks),
	}
	return res, nil
}
