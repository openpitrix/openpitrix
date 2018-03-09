// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package task

import (
	"fmt"
	"time"

	"openpitrix.io/openpitrix/pkg/constants"
	"openpitrix.io/openpitrix/pkg/db"
	"openpitrix.io/openpitrix/pkg/etcd"
	"openpitrix.io/openpitrix/pkg/logger"
	"openpitrix.io/openpitrix/pkg/manager/pilot"
	"openpitrix.io/openpitrix/pkg/models"
	"openpitrix.io/openpitrix/pkg/pb"
	"openpitrix.io/openpitrix/pkg/pi"
	"openpitrix.io/openpitrix/pkg/plugins"
)

type Controller struct {
	runningTasks chan string
	pi           *pi.Pi
	hostname     string
	queue        *etcd.Queue
}

func NewController(pi *pi.Pi, hostname string) *Controller {
	return &Controller{
		runningTasks: make(chan string, constants.TaskLength),
		pi:           pi,
		hostname:     hostname,
		queue:        pi.Etcd.NewQueue("task"),
	}
}

func (c *Controller) updateTaskAttributes(taskId string, attributes map[string]interface{}) error {
	_, err := c.pi.Db.
		Update(models.TaskTableName).
		SetMap(attributes).
		Where(db.Eq("task_id", taskId)).
		Exec()
	if err != nil {
		logger.Errorf("Failed to update task [%s]: %+v", taskId, err)
	}
	return err
}

func (c *Controller) ExtractTask() {
	taskId, err := c.queue.Dequeue()
	if err != nil {
		logger.Errorf("Failed to dequeue task from etcd: %+v", err)
		time.Sleep(3 * time.Second)
		return
	}
	c.runningTasks <- taskId
}

func (c *Controller) HandleTask() error {
	taskId := <-c.runningTasks

	task := &models.Task{
		TaskId: taskId,
		Status: constants.StatusWorking,
	}
	c.updateTaskAttributes(task.TaskId, map[string]interface{}{
		"status":   task.Status,
		"executor": c.hostname,
	})

	task.Status = constants.StatusFailed
	defer c.updateTaskAttributes(task.TaskId, map[string]interface{}{
		"status": task.Status,
	})

	pbTask := models.TaskToPb(task)
	if task.Target == constants.PilotManagerHost {
		err := pilot.HandleSubtask(
			&pb.HandleSubtaskRequest{
				SubtaskId:     pbTask.TaskId,
				SubtaskAction: pbTask.TaskAction,
				Directive:     pbTask.Directive,
			})
		if err != nil {
			logger.Errorf("Failed to handle task [%s] to pilot: %+v", task.TaskId, err)
			return err
		}
		err = pilot.WaitSubtask(task.TaskId, constants.WaitTaskTimeout, constants.WaitTaskInterval)
		if err != nil {
			logger.Errorf("Failed to wait task [%s]: %+v", task.TaskId, err)
			return err
		}
	} else {
		runtimeInterface := plugins.GetRuntimePlugin(task.Target)
		if runtimeInterface == nil {
			logger.Errorf("No such runtime [%s]. ", task.Target)
			return fmt.Errorf("No such runtime [%s]. ", task.Target)
		}
		err := runtimeInterface.HandleSubtask(task)
		if err != nil {
			logger.Errorf("Failed to handle subtask [%s] in runtime [%s]: %+v",
				task.TaskId, task.Target, err)
			return err
		}
		err = runtimeInterface.WaitSubtask(task.TaskId, constants.WaitTaskTimeout, constants.WaitTaskInterval)
		if err != nil {
			logger.Errorf("Failed to wait subtask [%s] in runtime [%s]: %+v",
				task.TaskId, task.Target, err)
			return err
		}
	}

	task.Status = constants.StatusSuccessful
	return nil
}

func (c *Controller) Serve() {
	go c.ExtractTask()
	go c.HandleTask()
}
