// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package models

import (
	"time"

	"openpitrix.io/openpitrix/pkg/constants"
	"openpitrix.io/openpitrix/pkg/db"
	"openpitrix.io/openpitrix/pkg/logger"
	"openpitrix.io/openpitrix/pkg/pb"
	"openpitrix.io/openpitrix/pkg/sender"
	"openpitrix.io/openpitrix/pkg/util/idutil"
	"openpitrix.io/openpitrix/pkg/util/jsonutil"
	"openpitrix.io/openpitrix/pkg/util/pbutil"
)

func NewTaskId() string {
	return idutil.GetUuid("t-")
}

type Task struct {
	TaskId         string
	JobId          string
	TaskAction     string
	Directive      string
	Owner          string
	OwnerPath      sender.OwnerPath
	Status         string
	ErrorCode      uint32
	Executor       string
	Target         string
	NodeId         string
	FailureAllowed bool
	CreateTime     time.Time
	StatusTime     time.Time
}

var TaskColumns = db.GetColumnsFromStruct(&Task{})

func NewTask(taskId, jobId, nodeId, target, taskAction, directive string, ownerPath sender.OwnerPath, failureAllowed bool) *Task {
	if taskId == "" {
		taskId = NewTaskId()
	} else if taskId == constants.PlaceHolder {
		taskId = ""
	}
	return &Task{
		TaskId:         taskId,
		JobId:          jobId,
		NodeId:         nodeId,
		Target:         target,
		TaskAction:     taskAction,
		Directive:      directive,
		Owner:          ownerPath.Owner(),
		OwnerPath:      ownerPath,
		Status:         constants.StatusPending,
		CreateTime:     time.Now(),
		StatusTime:     time.Now(),
		FailureAllowed: failureAllowed,
	}
}

func TaskToPb(task *Task) *pb.Task {
	pbTask := pb.Task{}
	pbTask.TaskId = pbutil.ToProtoString(task.TaskId)
	pbTask.JobId = pbutil.ToProtoString(task.JobId)
	pbTask.TaskAction = pbutil.ToProtoString(task.TaskAction)
	pbTask.Directive = pbutil.ToProtoString(task.Directive)
	pbTask.OwnerPath = task.OwnerPath.ToProtoString()
	pbTask.Status = pbutil.ToProtoString(task.Status)
	pbTask.ErrorCode = pbutil.ToProtoUInt32(task.ErrorCode)
	pbTask.Executor = pbutil.ToProtoString(task.Executor)
	pbTask.Target = pbutil.ToProtoString(task.Target)
	pbTask.NodeId = pbutil.ToProtoString(task.NodeId)
	pbTask.CreateTime = pbutil.ToProtoTimestamp(task.CreateTime)
	pbTask.StatusTime = pbutil.ToProtoTimestamp(task.StatusTime)
	pbTask.FailureAllowed = pbutil.ToProtoBool(task.FailureAllowed)
	return &pbTask
}

func TasksToPbs(tasks []*Task) (pbTasks []*pb.Task) {
	for _, task := range tasks {
		pbTasks = append(pbTasks, TaskToPb(task))
	}
	return
}

func PbToTask(pbTask *pb.Task) *Task {
	ownerPath := sender.OwnerPath(pbTask.GetOwnerPath().GetValue())
	return &Task{
		TaskId:         pbTask.GetTaskId().GetValue(),
		JobId:          pbTask.GetJobId().GetValue(),
		TaskAction:     pbTask.GetTaskAction().GetValue(),
		Directive:      pbTask.GetDirective().GetValue(),
		OwnerPath:      ownerPath,
		Owner:          ownerPath.Owner(),
		Status:         pbTask.GetStatus().GetValue(),
		ErrorCode:      pbTask.GetErrorCode().GetValue(),
		Executor:       pbTask.GetExecutor().GetValue(),
		Target:         pbTask.GetTarget().GetValue(),
		NodeId:         pbTask.GetNodeId().GetValue(),
		FailureAllowed: pbTask.GetFailureAllowed().GetValue(),
		CreateTime:     pbutil.GetTime(pbTask.GetCreateTime()),
		StatusTime:     pbutil.GetTime(pbTask.GetStatusTime()),
	}
}

func PbsToTasks(pbTasks []*pb.Task) (tasks []*Task) {
	for _, pbTask := range pbTasks {
		tasks = append(tasks, PbToTask(pbTask))
	}
	return
}

func (t *Task) GetTimeout(defaultTimeout time.Duration) time.Duration {
	if t.Directive == "" {
		return defaultTimeout
	}

	directive := make(map[string]interface{})
	err := jsonutil.Decode([]byte(t.Directive), &directive)
	if err != nil {
		logger.Error(nil, "Decode task [%s] directive [%s] failed: %+v.", t.TaskId, t.Directive, err)
		return defaultTimeout
	}

	timeout, exist := directive[constants.TimeoutName]
	if !exist {
		return defaultTimeout
	}
	tm := timeout.(float64)
	if tm <= 0 {
		return defaultTimeout
	}
	return time.Duration(tm) * time.Second
}
