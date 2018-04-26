// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package job

import (
	"sync"
	"time"

	taskclient "openpitrix.io/openpitrix/pkg/client/task"
	"openpitrix.io/openpitrix/pkg/constants"
	"openpitrix.io/openpitrix/pkg/db"
	"openpitrix.io/openpitrix/pkg/etcd"
	"openpitrix.io/openpitrix/pkg/logger"
	"openpitrix.io/openpitrix/pkg/models"
	"openpitrix.io/openpitrix/pkg/pi"
	"openpitrix.io/openpitrix/pkg/plugins"
)

type Controller struct {
	*pi.Pi
	runningJobs  chan string
	runningCount uint32
	hostname     string
	queue        *etcd.Queue
}

func NewController(pi *pi.Pi, hostname string) *Controller {
	return &Controller{
		Pi:           pi,
		runningJobs:  make(chan string),
		runningCount: 0,
		hostname:     hostname,
		queue:        pi.Etcd.NewQueue("job"),
	}
}

func (c *Controller) updateJobAttributes(jobId string, attributes map[string]interface{}) error {
	_, err := c.Db.
		Update(models.JobTableName).
		SetMap(attributes).
		Where(db.Eq("job_id", jobId)).
		Exec()
	if err != nil {
		logger.Errorf("Failed to update job [%s]: %+v", jobId, err)
	}
	return err
}

var mutex sync.Mutex

func (c *Controller) GetJobLength() uint32 {
	// TODO: from global config
	return constants.JobLength
}

func (c *Controller) IsRunningExceed() bool {
	mutex.Lock()
	defer mutex.Unlock()
	count := c.runningCount
	return count > c.GetJobLength()
}

func (c *Controller) ExtractJobs() {
	for {
		if c.IsRunningExceed() {
			logger.Errorf("Sleep 10s, running job count exceed [%d/%d]", c.runningCount, c.GetJobLength())
			time.Sleep(10 * time.Second)
			continue
		}
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

func (c *Controller) HandleJob(jobId string, cb func()) error {
	defer cb()

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

	query := c.Db.
		Select(models.JobColumns...).
		From(models.JobTableName).
		Where(db.Eq("job_id", jobId))

	err := query.LoadOne(&job)
	if err != nil {
		logger.Errorf("Failed to get job [%s]: %+v", job.JobId, err)
		return err
	}

	processor := NewProcessor(job)
	err = processor.Pre()
	if err != nil {
		return err
	}
	defer processor.Final()

	providerInterface, err := plugins.GetProviderPlugin(job.Provider)
	if err != nil {
		logger.Errorf("No such provider [%s]. ", job.Provider)
		return err
	}
	module, err := providerInterface.SplitJobIntoTasks(job)
	if err != nil {
		logger.Errorf("Failed to split job [%s] into tasks with provider [%s]: %+v",
			job.JobId, job.Provider, err)
		return err
	}

	err = module.WalkTree(func(parent *models.TaskLayer, current *models.TaskLayer) error {
		if parent != nil {
			for _, parentTask := range parent.Tasks {
				err = taskclient.WaitTask(parentTask.TaskId, parentTask.GetTimeout(constants.WaitTaskTimeout), constants.WaitTaskInterval)
				if err != nil {
					logger.Errorf("Failed to wait task [%s]: %+v", parentTask.TaskId, err)
					return err
				}
			}
		}

		if current != nil {
			for _, currentTask := range current.Tasks {
				taskId, err := taskclient.SendTask(currentTask)
				if err != nil {
					logger.Errorf("Failed to send task [%s]: %+v", currentTask.TaskId, err)
					return err
				}
				currentTask.TaskId = taskId
			}
			if current.IsLeaf() {
				for _, currentTask := range current.Tasks {
					err = taskclient.WaitTask(currentTask.TaskId, currentTask.GetTimeout(constants.WaitTaskTimeout), constants.WaitTaskInterval)
					if err != nil {
						logger.Errorf("Failed to wait task [%s]: %+v", currentTask.TaskId, err)
						return err
					}
				}
			}
		}

		return nil
	})

	if err != nil {
		return err
	}

	job.Status = constants.StatusSuccessful
	return processor.Post()
}

func (c *Controller) HandleJobs() {
	for jobId := range c.runningJobs {
		mutex.Lock()
		c.runningCount++
		mutex.Unlock()
		go c.HandleJob(jobId, func() {
			mutex.Lock()
			c.runningCount--
			mutex.Unlock()
		})
	}
}

func (c *Controller) Serve() {
	go c.ExtractJobs()
	go c.HandleJobs()
}
