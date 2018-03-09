// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package job

import (
	"fmt"
	"time"

	"openpitrix.io/openpitrix/pkg/constants"
	"openpitrix.io/openpitrix/pkg/db"
	"openpitrix.io/openpitrix/pkg/etcd"
	"openpitrix.io/openpitrix/pkg/logger"
	"openpitrix.io/openpitrix/pkg/manager/task"
	"openpitrix.io/openpitrix/pkg/models"
	"openpitrix.io/openpitrix/pkg/pi"
	"openpitrix.io/openpitrix/pkg/plugins"
)

type Controller struct {
	runningJobs chan string
	pi          *pi.Pi
	hostname    string
	queue       *etcd.Queue
}

func NewController(pi *pi.Pi, hostname string) *Controller {
	return &Controller{
		runningJobs: make(chan string, constants.JobLength),
		pi:          pi,
		hostname:    hostname,
		queue:       pi.Etcd.NewQueue("job"),
	}
}

func (c *Controller) updateJobAttributes(jobId string, attributes map[string]interface{}) error {
	_, err := c.pi.Db.
		Update(models.JobTableName).
		SetMap(attributes).
		Where(db.Eq("job_id", jobId)).
		Exec()
	if err != nil {
		logger.Errorf("Failed to update job [%s]: %+v", jobId, err)
	}
	return err
}

func (c *Controller) ExtractJob() {
	jobId, err := c.queue.Dequeue()
	if err != nil {
		logger.Errorf("Failed to dequeue job from etcd: %+v", err)
		time.Sleep(3 * time.Second)
		return
	}
	c.runningJobs <- jobId
}

func (c *Controller) HandleJob() error {
	jobId := <-c.runningJobs
	job := &models.Job{
		JobId:  jobId,
		Status: constants.StatusWorking,
	}
	c.updateJobAttributes(job.JobId, map[string]interface{}{
		"status":   job.Status,
		"executor": c.hostname,
	})

	job.Status = constants.StatusFailed
	defer c.updateJobAttributes(job.JobId, map[string]interface{}{
		"status": job.Status,
	})

	query := c.pi.Db.
		Select(models.JobColumns...).
		From(models.JobTableName).
		Where(db.Eq("job_id", jobId))

	err := query.LoadOne(&job)
	if err != nil {
		return err
	}

	runtimeInterface := plugins.GetRuntimePlugin(job.Runtime)
	if runtimeInterface == nil {
		logger.Errorf("No such runtime [%s]. ", job.Runtime)
		return fmt.Errorf("No such runtime [%s]. ", job.Runtime)
	}
	module, err := runtimeInterface.SplitJobIntoTasks(job)
	if err != nil {
		logger.Errorf("Failed to split job [%s] into tasks with runtime [%s]: %+v",
			job.JobId, job.Runtime, err)
		return err
	}

	err = module.WalkTree(func(parent *models.Module, current *models.Module) error {
		if parent != nil {
			err = task.WaitTask(parent.Task.TaskId, constants.WaitTaskTimeout, constants.WaitTaskInterval)
			if err != nil {
				logger.Errorf("Failed to wait task [%s]: %+v", parent.Task.TaskId, err)
				return err
			}
		}
		if current != nil {
			err = task.SendTask(current.Task)
			if err != nil {
				logger.Errorf("Failed to send task [%s]: %+v", current.Task.TaskId, err)
				return err
			}
		}
		return nil
	})

	if err == nil {
		job.Status = constants.StatusSuccessful
	}

	return err
}

func (c *Controller) Serve() {
	go c.ExtractJob()
	go c.HandleJob()
}
