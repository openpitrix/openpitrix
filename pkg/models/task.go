// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package models

import (
	"time"

	"openpitrix.io/openpitrix/pkg/constants"
	"openpitrix.io/openpitrix/pkg/pb"
	"openpitrix.io/openpitrix/pkg/utils"
)

const TaskTableName = "task"

func NewTaskId() string {
	return utils.GetUuid("t-")
}

type Task struct {
	TaskId     string
	JobId      string
	TaskAction string
	Directive  string
	Owner      string
	Status     string
	ErrorCode  uint32
	Executor   string
	Target     string
	NodeId     string
	CreateTime time.Time
	StatusTime time.Time
}

var TaskColumns = GetColumnsFromStruct(&Task{})

func NewTask(taskId, jobId, nodeId, target, taskAction, directive, userId string) *Task {
	if taskId == "" {
		taskId = NewTaskId()
	} else if taskId == constants.PlaceHolder {
		taskId = ""
	}
	return &Task{
		TaskId:     taskId,
		JobId:      jobId,
		NodeId:     nodeId,
		Target:     target,
		TaskAction: taskAction,
		Directive:  directive,
		Owner:      userId,
		Status:     constants.StatusPending,
		CreateTime: time.Now(),
		StatusTime: time.Now(),
	}
}

func TaskToPb(task *Task) *pb.Task {
	pbTask := pb.Task{}
	pbTask.TaskId = utils.ToProtoString(task.TaskId)
	pbTask.JobId = utils.ToProtoString(task.JobId)
	pbTask.TaskAction = utils.ToProtoString(task.TaskAction)
	pbTask.Directive = utils.ToProtoString(task.Directive)
	pbTask.Owner = utils.ToProtoString(task.Owner)
	pbTask.Status = utils.ToProtoString(task.Status)
	pbTask.ErrorCode = utils.ToProtoUInt32(task.ErrorCode)
	pbTask.Executor = utils.ToProtoString(task.Executor)
	pbTask.Target = utils.ToProtoString(task.Target)
	pbTask.NodeId = utils.ToProtoString(task.NodeId)
	pbTask.CreateTime = utils.ToProtoTimestamp(task.CreateTime)
	pbTask.StatusTime = utils.ToProtoTimestamp(task.StatusTime)
	return &pbTask
}

func TasksToPbs(tasks []*Task) (pbTasks []*pb.Task) {
	for _, task := range tasks {
		pbTasks = append(pbTasks, TaskToPb(task))
	}
	return
}

type TaskLayer struct {
	Tasks []*Task
	Child *TaskLayer
}

// WalkFunc is a callback type for use with TaskLayer.WalkTree
type WalkFunc func(parent *TaskLayer, current *TaskLayer) error

func (t *TaskLayer) WalkTree(cb WalkFunc) error {
	return walkTaskLayerTree(nil, t, cb)
}

func (t *TaskLayer) Leaf() *TaskLayer {
	current := t
	for {
		if current.Child == nil {
			return current
		} else {
			current = current.Child
		}
	}
}

func walkTaskLayerTree(parent *TaskLayer, current *TaskLayer, cb WalkFunc) error {
	err := cb(parent, current)
	if err != nil {
		return err
	}

	if current.Child == nil {
		return nil
	} else {
		err = walkTaskLayerTree(current, current.Child, cb)
		return err
	}
}
