// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package job

import (
	"sync"
	"time"

	"openpitrix.io/openpitrix/pkg/client"
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
			logger.Error("Sleep 10s, running job count exceed [%d/%d]", c.runningCount, c.GetJobLength())
			time.Sleep(10 * time.Second)
			continue
		}
		jobId, err := c.queue.Dequeue()
		if err != nil {
			logger.Error("Failed to dequeue job from etcd queue: %+v", err)
			time.Sleep(3 * time.Second)
			continue
		}
		logger.Debug("Dequeue job [%s] from etcd queue success", jobId)
		c.runningJobs <- jobId
	}
}

func (c *Controller) HandleJob(jobId string, cb func()) error {
	jLogger := logger.NewLogger()
	jLogger.SetSuffix("(" + jobId + ")")

	defer cb()

	job := &models.Job{
		JobId:  jobId,
		Status: constants.StatusWorking,
	}
	err := c.updateJobAttributes(job.JobId, map[string]interface{}{
		"status":   job.Status,
		"executor": c.hostname,
	})
	if err != nil {
		jLogger.Error("Failed to update job: %+v", err)
		return err
	}

	err = func() error {
		query := c.Db.
			Select(models.JobColumns...).
			From(models.JobTableName).
			Where(db.Eq("job_id", jobId))

		err := query.LoadOne(&job)
		if err != nil {
			jLogger.Error("Failed to get job: %+v", err)
			return err
		}

		processor := NewProcessor(job, jLogger)
		err = processor.Pre()
		if err != nil {
			return err
		}
		defer processor.Final()

		providerInterface, err := plugins.GetProviderPlugin(job.Provider, jLogger)
		if err != nil {
			jLogger.Error("No such provider [%s]. ", job.Provider)
			return err
		}
		module, err := providerInterface.SplitJobIntoTasks(job)
		if err != nil {
			jLogger.Error("Failed to split job into tasks with provider [%s]: %+v", job.Provider, err)
			return err
		}

		ctx := client.GetSystemUserContext()
		taskClient, err := taskclient.NewClient()
		if err != nil {
			jLogger.Error("Connect to task service failed: %+v", err)
			return err
		}

		successful := true
		module.WalkTree(func(parent *models.TaskLayer, current *models.TaskLayer) {
			if parent != nil {
				for _, parentTask := range parent.Tasks {
					err = taskClient.WaitTask(ctx, parentTask.TaskId, parentTask.GetTimeout(constants.MaxTaskTimeout), constants.WaitTaskInterval)
					if err != nil {
						jLogger.Error("Failed to wait task [%s]: %+v", parentTask.TaskId, err)
						if !parentTask.FailureAllowed {
							successful = false
						}
					}
				}
			}

			if current != nil {
				for _, currentTask := range current.Tasks {
					if !successful {
						currentTask.Status = constants.StatusFailed
					}
					currentTask.TaskId, err = taskClient.SendTask(ctx, currentTask)
					if err != nil {
						jLogger.Error("Failed to send task [%s]: %+v", currentTask.TaskId, err)
						successful = false
					}
				}
				if current.IsLeaf() {
					for _, currentTask := range current.Tasks {
						err = taskClient.WaitTask(ctx, currentTask.TaskId, currentTask.GetTimeout(constants.MaxTaskTimeout), constants.WaitTaskInterval)
						if err != nil {
							jLogger.Error("Failed to wait task [%s]: %+v", currentTask.TaskId, err)
							if !currentTask.FailureAllowed {
								successful = false
							}
						}
					}
				}
			}
		})
		if !successful {
			return err
		}

		processor.Job.Status = constants.StatusSuccessful
		return processor.Post()
	}()

	var status = constants.StatusSuccessful
	if err != nil {
		jLogger.Error("Job [%s] failed: %+v", jobId, err)
		status = constants.StatusFailed
	}
	err = c.updateJobAttributes(jobId, map[string]interface{}{
		"status":      status,
		"status_time": time.Now(),
	})
	if err != nil {
		jLogger.Error("Failed to update job: %+v", err)
	}

	return err
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
