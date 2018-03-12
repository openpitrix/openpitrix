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

func (c *Controller) ExtractJobs() {
	for {
		jobId, err := c.queue.Dequeue()
		if err != nil {
			logger.Errorf("Failed to dequeue job from etcd queue: %+v", err)
			time.Sleep(3 * time.Second)
			continue
		}
		logger.Debugf("Dequeue job [%s] from etcd queue success", jobId)
		c.runningJobs <- jobId
	}
}

func (c *Controller) HandleJob(jobId string) error {
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
		logger.Errorf("Failed to get job [%s]: %+v", job.JobId, err)
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

	err = module.WalkTree(func(parent *models.TaskLayer, current *models.TaskLayer) error {
		if parent != nil {
			for _, parentTask := range parent.Tasks {
				err = task.WaitTask(parentTask.TaskId, constants.WaitTaskTimeout, constants.WaitTaskInterval)
				if err != nil {
					logger.Errorf("Failed to wait task [%s]: %+v", parentTask.TaskId, err)
					return err
				}
			}
		}

		for _, currentTask := range current.Tasks {
			err = task.SendTask(currentTask)
			if err != nil {
				logger.Errorf("Failed to send task [%s]: %+v", currentTask.TaskId, err)
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

func (c *Controller) HandleJobs() {
	for {
		jobId, ok := <-c.runningJobs
		if !ok {
			logger.Errorf("Channel controller runningJobs is closed")
			return
		}
		go c.HandleJob(jobId)
	}
}

func (c *Controller) Serve() {
	go c.ExtractJobs()
	go c.HandleJobs()
}
