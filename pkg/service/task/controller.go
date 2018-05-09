// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package task

import (
	"sync"
	"time"

	"openpitrix.io/openpitrix/pkg/client"
	pilotclient "openpitrix.io/openpitrix/pkg/client/pilot"
	"openpitrix.io/openpitrix/pkg/constants"
	"openpitrix.io/openpitrix/pkg/db"
	"openpitrix.io/openpitrix/pkg/etcd"
	"openpitrix.io/openpitrix/pkg/logger"
	"openpitrix.io/openpitrix/pkg/models"
	"openpitrix.io/openpitrix/pkg/pb/types"
	"openpitrix.io/openpitrix/pkg/pi"
	"openpitrix.io/openpitrix/pkg/plugins"
	"openpitrix.io/openpitrix/pkg/plugins/vmbased"
	"openpitrix.io/openpitrix/pkg/util/jsonutil"
)

type Controller struct {
	*pi.Pi
	runningTasks chan string
	runningCount uint32
	hostname     string
	queue        *etcd.Queue
}

func NewController(pi *pi.Pi, hostname string) *Controller {
	return &Controller{
		Pi:           pi,
		runningTasks: make(chan string),
		runningCount: 0,
		hostname:     hostname,
		queue:        pi.Etcd.NewQueue("task"),
	}
}

func (c *Controller) updateTaskAttributes(taskId string, attributes map[string]interface{}) error {
	_, err := c.Db.
		Update(models.TaskTableName).
		SetMap(attributes).
		Where(db.Eq("task_id", taskId)).
		Exec()
	if err != nil {
		logger.Error("Failed to update task [%s]: %+v", taskId, err)
	}
	return err
}

var mutex sync.Mutex

func (c *Controller) GetTaskLength() uint32 {
	// TODO: from global config
	return constants.TaskLength
}

func (c *Controller) IsRunningExceed() bool {
	mutex.Lock()
	defer mutex.Unlock()
	count := c.runningCount
	return count > c.GetTaskLength()
}

func (c *Controller) ExtractTasks() {
	for {
		if c.IsRunningExceed() {
			logger.Error("Sleep 10s, running task count exceed [%d/%d]", c.runningCount, c.GetTaskLength())
			time.Sleep(10 * time.Second)
			continue
		}
		taskId, err := c.queue.Dequeue()
		if err != nil {
			logger.Error("Failed to dequeue task from etcd queue: %+v", err)
			time.Sleep(3 * time.Second)
			continue
		}
		logger.Debug("Dequeue task [%s] from etcd queue success", taskId)
		c.runningTasks <- taskId
	}
}

func (c *Controller) HandleTask(taskId string, cb func()) error {
	defer cb()
	task := &models.Task{
		TaskId: taskId,
		Status: constants.StatusWorking,
	}
	c.updateTaskAttributes(task.TaskId, map[string]interface{}{
		"status":   task.Status,
		"executor": c.hostname,
	})
	err := func() error {
		query := c.Db.
			Select(models.TaskColumns...).
			From(models.TaskTableName).
			Where(db.Eq("task_id", taskId))

		err := query.LoadOne(&task)
		if err != nil {
			logger.Error("Failed to get task [%s]: %+v", task.TaskId, err)
			return err
		}

		processor := NewProcessor(task)
		err = processor.Pre()
		if err != nil {
			logger.Error("Executing task [%s] pre processor failed: %+v", task.TaskId, err)
			return err
		}

		ctx := client.GetSystemUserContext()
		pilotClient, err := pilotclient.NewClient(ctx)
		if err != nil {
			logger.Error("Connect to pilot service failed: %+v", err)
			return err
		}

		if task.Target == constants.TargetPilot {
			switch task.TaskAction {
			case vmbased.ActionSetDroneConfig:
				config := new(pbtypes.SetDroneConfigRequest)
				err = jsonutil.Decode([]byte(task.Directive), config)
				if err != nil {
					logger.Error("Decode task [%s] directive [%s] failed: %+v", taskId, task.Directive, err)
					return err
				}
				_, err = pilotClient.SetDroneConfig(ctx, config)
				if err != nil {
					logger.Error("Send task [%s] to pilot failed: %+v", taskId, err)
					return err
				}
			case vmbased.ActionSetFrontgateConfig:
				config := new(pbtypes.FrontgateConfig)
				err = jsonutil.Decode([]byte(task.Directive), config)
				if err != nil {
					logger.Error("Decode task [%s] directive [%s] failed: %+v", taskId, task.Directive, err)
					return err
				}
				_, err = pilotClient.SetFrontgateConfig(ctx, config)
				if err != nil {
					logger.Error("Send task [%s] to pilot failed: %+v", taskId, err)
					return err
				}
			default:
				pbTask := models.TaskToPb(task)
				_, err := pilotClient.HandleSubtask(ctx,
					&pbtypes.SubTaskMessage{
						TaskId:    pbTask.TaskId.GetValue(),
						Action:    pbTask.TaskAction.GetValue(),
						Directive: pbTask.Directive.GetValue(),
					})
				if err != nil {
					logger.Error("Failed to handle task [%s] to pilot: %+v", task.TaskId, err)
					return err
				}
				err = pilotClient.WaitSubtask(
					ctx, task.TaskId, task.GetTimeout(constants.WaitTaskTimeout), constants.WaitTaskInterval)
				if err != nil {
					logger.Error("Failed to wait task [%s]: %+v", task.TaskId, err)
					return err
				}
			}
		} else {
			providerInterface, err := plugins.GetProviderPlugin(task.Target)
			if err != nil {
				logger.Error("No such runtime [%s]. ", task.Target)
				return err
			}
			err = providerInterface.HandleSubtask(task)
			if err != nil {
				logger.Error("Failed to handle subtask [%s] in runtime [%s]: %+v",
					task.TaskId, task.Target, err)
				return err
			}
			err = providerInterface.WaitSubtask(
				task, task.GetTimeout(constants.WaitTaskTimeout), constants.WaitTaskInterval)
			if err != nil {
				logger.Error("Failed to wait subtask [%s] in runtime [%s]: %+v",
					task.TaskId, task.Target, err)
				return err
			}

			logger.Debug("After wait subtask [%s] directive: %s", task.TaskId, task.Directive)
		}

		if err != nil {
			return err
		}

		err = processor.Post()
		if err != nil {
			logger.Error("Executing task [%s] post processor failed: %+v", task.TaskId, err)
		}
		return err
	}()
	var status = constants.StatusSuccessful
	if err != nil {
		status = constants.StatusFailed

	}
	c.updateTaskAttributes(task.TaskId, map[string]interface{}{
		"status": status,
	})
	return err
}

func (c *Controller) HandleTasks() {
	for taskId := range c.runningTasks {
		mutex.Lock()
		c.runningCount++
		mutex.Unlock()

		go c.HandleTask(taskId, func() {
			mutex.Lock()
			c.runningCount--
			mutex.Unlock()
		})
	}
}

func (c *Controller) Serve() {
	go c.ExtractTasks()
	go c.HandleTasks()
}
