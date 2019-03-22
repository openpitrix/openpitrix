// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package task

import (
	"context"
	"strings"

	"openpitrix.io/openpitrix/pkg/constants"
	"openpitrix.io/openpitrix/pkg/gerr"
	"openpitrix.io/openpitrix/pkg/manager"
	"openpitrix.io/openpitrix/pkg/models"
	"openpitrix.io/openpitrix/pkg/pb"
	"openpitrix.io/openpitrix/pkg/pi"
	"openpitrix.io/openpitrix/pkg/util/ctxutil"
	"openpitrix.io/openpitrix/pkg/util/pbutil"
)

func (p *Server) CreateTask(ctx context.Context, req *pb.CreateTaskRequest) (*pb.CreateTaskResponse, error) {
	s := ctxutil.GetSender(ctx)
	newTask := models.NewTask(
		"",
		req.GetJobId().GetValue(),
		req.GetNodeId().GetValue(),
		req.GetTarget().GetValue(),
		req.GetTaskAction().GetValue(),
		req.GetDirective().GetValue(),
		s.GetOwnerPath(),
		req.GetFailureAllowed().GetValue(),
	)

	if req.GetStatus().GetValue() == constants.StatusFailed {
		newTask.Status = req.GetStatus().GetValue()
	}

	_, err := pi.Global().DB(ctx).
		InsertInto(constants.TableTask).
		Record(newTask).
		Exec()
	if err != nil {
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorCreateResourcesFailed)
	}

	if newTask.Status != constants.StatusFailed {
		err = p.controller.queue.Enqueue(newTask.TaskId)
		if err != nil {
			return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorCreateResourcesFailed)
		}
	}

	res := &pb.CreateTaskResponse{
		TaskId: pbutil.ToProtoString(newTask.TaskId),
		JobId:  pbutil.ToProtoString(newTask.JobId),
	}
	return res, nil
}

func (p *Server) DescribeTasks(ctx context.Context, req *pb.DescribeTasksRequest) (*pb.DescribeTasksResponse, error) {
	var tasks []*models.Task
	offset := pbutil.GetOffsetFromRequest(req)
	limit := pbutil.GetLimitFromRequest(req)

	displayColumns := manager.GetDisplayColumns(req.GetDisplayColumns(), models.TaskColumns)
	query := pi.Global().DB(ctx).
		Select(displayColumns...).
		From(constants.TableTask).
		Offset(offset).
		Limit(limit).
		Where(manager.BuildPermissionFilter(ctx)).
		Where(manager.BuildFilterConditions(req, constants.TableTask)).
		OrderDir(constants.ColumnCreateTime, true)

	if len(displayColumns) > 0 {
		_, err := query.Load(&tasks)
		if err != nil {
			return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorDescribeResourcesFailed)
		}
	}
	count, err := query.Count()
	if err != nil {
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorDescribeResourcesFailed)
	}

	res := &pb.DescribeTasksResponse{
		TaskSet:    models.TasksToPbs(tasks),
		TotalCount: uint32(count),
	}
	return res, nil
}

func (p *Server) RetryTasks(ctx context.Context, req *pb.RetryTasksRequest) (*pb.RetryTasksResponse, error) {
	taskIds := req.GetTaskId()
	tasks, err := CheckTasksPermission(ctx, taskIds)
	if err != nil {
		return nil, err
	}

	for _, taskId := range taskIds {
		err = p.controller.queue.Enqueue(taskId)
		if err != nil {
			return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorRetryTaskFailed, strings.Join(taskIds, ","))
		}
	}

	res := &pb.RetryTasksResponse{
		TaskSet: models.TasksToPbs(tasks),
	}
	return res, nil
}
