// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package job

import (
	"context"
	"sync"
	"time"

	"openpitrix.io/openpitrix/pkg/client"
	accountclient "openpitrix.io/openpitrix/pkg/client/iam"
	taskclient "openpitrix.io/openpitrix/pkg/client/task"
	"openpitrix.io/openpitrix/pkg/constants"
	"openpitrix.io/openpitrix/pkg/db"
	"openpitrix.io/openpitrix/pkg/etcd"
	"openpitrix.io/openpitrix/pkg/logger"
	"openpitrix.io/openpitrix/pkg/models"
	"openpitrix.io/openpitrix/pkg/pi"
	"openpitrix.io/openpitrix/pkg/plugins"
	"openpitrix.io/openpitrix/pkg/util/ctxutil"
)

type Controller struct {
	runningJobs  chan string
	runningCount int32
	hostname     string
	queue        *etcd.Queue
}

func NewController(hostname string) *Controller {
	return &Controller{
		runningJobs:  make(chan string),
		runningCount: 0,
		hostname:     hostname,
		queue:        pi.Global().Etcd(nil).NewQueue("job"),
	}
}

func (c *Controller) updateJobAttributes(ctx context.Context, jobId string, attributes map[string]interface{}) error {
	_, err := pi.Global().DB(ctx).
		Update(constants.TableJob).
		SetMap(attributes).
		Where(db.Eq("job_id", jobId)).
		Exec()
	return err
}

var mutex sync.Mutex

func (c *Controller) GetJobLength() int32 {
	if pi.Global().GlobalConfig().Job.MaxWorkingJobs > 0 {
		return pi.Global().GlobalConfig().Job.MaxWorkingJobs
	} else {
		return constants.DefaultMaxWorkingJobs
	}
}

func (c *Controller) IsRunningExceed() bool {
	mutex.Lock()
	defer mutex.Unlock()
	count := c.runningCount
	return count > c.GetJobLength()
}

func (c *Controller) UpdateWorkingJobs(ctx context.Context) error {
	//TODO: retry the job
	_, err := pi.Global().DB(ctx).
		Update(constants.TableJob).
		SetMap(map[string]interface{}{"status": constants.StatusFailed}).
		Where(db.Eq("status", constants.StatusWorking)).
		Exec()
	return err
}

func (c *Controller) ExtractJobs(ctx context.Context) {
	for {
		if c.IsRunningExceed() {
			logger.Error(ctx, "Sleep 10s, running job count exceed [%d/%d]", c.runningCount, c.GetJobLength())
			time.Sleep(10 * time.Second)
			continue
		}
		jobId, err := c.queue.Dequeue()
		if err != nil {
			logger.Error(ctx, "Failed to dequeue job from etcd queue: %+v", err)
			time.Sleep(3 * time.Second)
			continue
		}
		logger.Debug(ctx, "Dequeue job [%s] from etcd queue success", jobId)
		c.runningJobs <- jobId
	}
}

func (c *Controller) HandleJob(ctx context.Context, jobId string, cb func()) error {
	ctx = ctxutil.AddMessageId(ctx, jobId)

	defer cb()

	job := &models.Job{
		JobId:  jobId,
		Status: constants.StatusWorking,
	}

	err := c.updateJobAttributes(ctx, job.JobId, map[string]interface{}{
		"status":   job.Status,
		"executor": c.hostname,
	})
	if err != nil {
		logger.Error(ctx, "Failed to update job: %+v", err)
		return err
	}

	err = func() error {
		query := pi.Global().DB(ctx).
			Select(models.JobColumns...).
			From(constants.TableJob).
			Where(db.Eq("job_id", jobId))

		err := query.LoadOne(&job)
		if err != nil {
			logger.Error(ctx, "Failed to get job: %+v", err)
			return err
		}

		accountClient, err := accountclient.NewClient()
		if err != nil {
			return err
		}
		users, err := accountClient.GetUsers(ctx, []string{job.Owner})
		if err != nil {
			return err
		}

		ctx = client.SetUserToContext(ctx, users[0])

		processor := NewProcessor(job)
		err = processor.Pre(ctx)
		if err != nil {
			return err
		}
		defer processor.Final(ctx)

		providerInterface, err := plugins.GetProviderPlugin(ctx, job.Provider)
		if err != nil {
			logger.Error(ctx, "No such provider [%s]. ", job.Provider)
			return err
		}
		module, err := providerInterface.SplitJobIntoTasks(ctx, job)
		if err != nil {
			logger.Error(ctx, "Failed to split job into tasks with provider [%s]: %+v", job.Provider, err)
			return err
		}

		taskClient, err := taskclient.NewClient()
		if err != nil {
			logger.Error(ctx, "Connect to task service failed: %+v", err)
			return err
		}

		successful := true
		module.WalkTree(func(parent *models.TaskLayer, current *models.TaskLayer) {
			if parent != nil {
				for _, parentTask := range parent.Tasks {
					err = taskClient.WaitTask(ctx, parentTask.TaskId, parentTask.GetTimeout(constants.MaxTaskTimeout), constants.WaitTaskInterval)
					if err != nil {
						logger.Error(ctx, "Failed to wait task [%s]: %+v", parentTask.TaskId, err)
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
						logger.Error(ctx, "Failed to send task [%s]: %+v", currentTask.TaskId, err)
						successful = false
					}
				}
				if current.IsLeaf() {
					for _, currentTask := range current.Tasks {
						err = taskClient.WaitTask(ctx, currentTask.TaskId, currentTask.GetTimeout(constants.MaxTaskTimeout), constants.WaitTaskInterval)
						if err != nil {
							logger.Error(ctx, "Failed to wait task [%s]: %+v", currentTask.TaskId, err)
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
		return processor.Post(ctx)
	}()

	var status = constants.StatusSuccessful
	if err != nil {
		logger.Error(ctx, "Job [%s] failed: %+v", jobId, err)
		status = constants.StatusFailed
	}

	err = c.updateJobAttributes(ctx, jobId, map[string]interface{}{
		"status":      status,
		"status_time": time.Now(),
	})
	if err != nil {
		logger.Error(ctx, "Failed to update job: %+v", err)
	}

	return err
}

func (c *Controller) HandleJobs(ctx context.Context) {
	for jobId := range c.runningJobs {
		mutex.Lock()
		c.runningCount++
		mutex.Unlock()
		go c.HandleJob(ctx, jobId, func() {
			mutex.Lock()
			c.runningCount--
			mutex.Unlock()
		})
	}
}

func (c *Controller) Serve() {
	ctx := context.Background()
	err := c.UpdateWorkingJobs(ctx)
	if err != nil {
		logger.Error(ctx, "Failed to update working jobs: %+v", err)
	}
	go c.ExtractJobs(ctx)
	go c.HandleJobs(ctx)
}
