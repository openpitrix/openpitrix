// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package task

import (
	"context"
	"sync"
	"time"

	"strings"

	"openpitrix.io/openpitrix/pkg/client"
	pilotclient "openpitrix.io/openpitrix/pkg/client/pilot"
	"openpitrix.io/openpitrix/pkg/constants"
	"openpitrix.io/openpitrix/pkg/db"
	"openpitrix.io/openpitrix/pkg/etcd"
	"openpitrix.io/openpitrix/pkg/logger"
	"openpitrix.io/openpitrix/pkg/models"
	"openpitrix.io/openpitrix/pkg/pb/metadata/types"
	"openpitrix.io/openpitrix/pkg/pi"
	"openpitrix.io/openpitrix/pkg/plugins"
	"openpitrix.io/openpitrix/pkg/plugins/vmbased"
	"openpitrix.io/openpitrix/pkg/util/funcutil"
	"openpitrix.io/openpitrix/pkg/util/jsonutil"
	"openpitrix.io/openpitrix/pkg/util/retryutil"
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
	task := new(models.Task)
	query := c.Db.
		Select(models.TaskColumns...).
		From(models.TaskTableName).
		Where(db.Eq("task_id", taskId))

	err := query.LoadOne(&task)
	if err != nil {
		logger.Error("Failed to get task [%s]: %+v", task.TaskId, err)
		return err
	}

	tLogger := logger.NewLogger()
	tLogger.SetSuffix("(" + task.JobId + ")(" + taskId + ")")

	err = c.updateTaskAttributes(task.TaskId, map[string]interface{}{
		"status":   constants.StatusWorking,
		"executor": c.hostname,
	})
	if err != nil {
		tLogger.Error("Failed to update task: %+v", err)
		return err
	}

	err = func() error {
		processor := NewProcessor(task, tLogger)
		err = processor.Pre()
		if err != nil {
			tLogger.Error("Executing task pre processor failed: %+v", err)
			return err
		}

		ctx := client.GetSystemUserContext()
		pilotClient, err := pilotclient.NewClient()
		if err != nil {
			tLogger.Error("Connect to pilot service failed: %+v", err)
			return err
		}

		if task.Target == constants.TargetPilot {
			withTimeoutCtx, cancel := context.WithTimeout(ctx, constants.GrpcToPilotTimeout)
			defer cancel()
			switch task.TaskAction {
			case vmbased.ActionSetDroneConfig:
				config := new(pbtypes.SetDroneConfigRequest)
				err = jsonutil.Decode([]byte(task.Directive), config)
				if err != nil {
					tLogger.Error("Decode task directive [%s] failed: %+v", task.Directive, err)
					return err
				}
				err = retryutil.Retry(3, 0, func() error {
					_, err = pilotClient.SetDroneConfig(withTimeoutCtx, config)
					return err
				})
				if err != nil {
					tLogger.Error("Send task to pilot failed: %+v", err)
					return err
				}
			case vmbased.ActionSetFrontgateConfig:
				config := new(pbtypes.FrontgateConfig)
				err = jsonutil.Decode([]byte(task.Directive), config)
				if err != nil {
					tLogger.Error("Decode task directive [%s] failed: %+v", task.Directive, err)
					return err
				}
				err = retryutil.Retry(3, 0, func() error {
					_, err = pilotClient.SetFrontgateConfig(withTimeoutCtx, config)
					return err
				})
				if err != nil {
					tLogger.Error("Send task to pilot failed: %+v", err)
					return err
				}
			case vmbased.ActionPingDrone:
				droneEndpoint := new(pbtypes.DroneEndpoint)
				err = jsonutil.Decode([]byte(task.Directive), droneEndpoint)
				if err != nil {
					tLogger.Error("Decode task directive [%s] failed: %+v", task.Directive, err)
					return err
				}
				err = funcutil.WaitForSpecificOrError(func() (bool, error) {
					withTimeoutCtx, cancel := context.WithTimeout(ctx, constants.GrpcToPilotTimeout)
					defer cancel()
					_, err := pilotClient.PingDrone(withTimeoutCtx, droneEndpoint)
					if err != nil {
						tLogger.Warn("Send task to pilot failed, will retry: %+v", err)
						return false, nil
					} else {
						return true, nil
					}
				}, task.GetTimeout(constants.WaitDroneServiceTimeout), constants.WaitDroneServiceInterval)
				if err != nil {
					tLogger.Error("Send task to pilot failed: %+v", err)
					return err
				}

			case vmbased.ActionPingFrontgate:
				request := new(pbtypes.FrontgateId)
				err = jsonutil.Decode([]byte(task.Directive), request)
				if err != nil {
					tLogger.Error("Decode task directive [%s] failed: %+v", task.Directive, err)
					return err
				}
				err = funcutil.WaitForSpecificOrError(func() (bool, error) {
					withTimeoutCtx, cancel := context.WithTimeout(ctx, constants.GrpcToPilotTimeout)
					defer cancel()
					_, err := pilotClient.PingFrontgate(withTimeoutCtx, request)
					if err != nil {
						tLogger.Warn("Send task to pilot failed, will retry: %+v", err)
						return false, nil
					} else {
						return true, nil
					}
				}, task.GetTimeout(constants.WaitFrontgateServiceTimeout), constants.WaitFrontgateServiceInterval)
				if err != nil {
					tLogger.Error("Send task to pilot failed: %+v", err)
					return err
				}

			case vmbased.ActionStartConfd:
				pbTask := models.TaskToPb(task)
				err = retryutil.Retry(3, 0, func() error {
					_, err := pilotClient.HandleSubtask(withTimeoutCtx,
						&pbtypes.SubTaskMessage{
							TaskId:    pbTask.TaskId.GetValue(),
							Action:    pbTask.TaskAction.GetValue(),
							Directive: pbTask.Directive.GetValue(),
						})
					return err
				})
				if err != nil {
					tLogger.Error("Failed to handle task to pilot: %+v", err)
					return err
				}

				time.Sleep(1 * time.Second)

			case vmbased.ActionStopConfd:
				pbTask := models.TaskToPb(task)
				err = retryutil.Retry(3, 0, func() error {
					_, err := pilotClient.HandleSubtask(withTimeoutCtx,
						&pbtypes.SubTaskMessage{
							TaskId:    pbTask.TaskId.GetValue(),
							Action:    pbTask.TaskAction.GetValue(),
							Directive: pbTask.Directive.GetValue(),
						})
					return err
				})
				if err != nil {
					tLogger.Error("Failed to handle task to pilot: %+v", err)
					return err
				}

				tLogger.Debug("Finish subtask [%s]", task.TaskId)

				time.Sleep(1 * time.Second)

			case vmbased.ActionRemoveContainerOnDrone:
				request := new(pbtypes.RunCommandOnDroneRequest)
				err = jsonutil.Decode([]byte(task.Directive), request)
				if err != nil {
					tLogger.Error("Decode task directive [%s] failed: %+v", task.Directive, err)
					return err
				}

				err = retryutil.Retry(3, 0, func() error {
					_, err = pilotClient.RunCommandOnDrone(withTimeoutCtx, request)
					if err != nil {
						if strings.Contains(err.Error(), "transport is closing") {
							tLogger.Debug("Expected error: %+v", err)
							return nil
						} else {
							tLogger.Error("%s", err.Error())
						}
					}
					return err
				})
				if err != nil {
					tLogger.Error("Send task to pilot failed: %+v", err)
					return err
				}

			case vmbased.ActionRemoveContainerOnFrontgate:
				request := new(pbtypes.RunCommandOnFrontgateRequest)
				err = jsonutil.Decode([]byte(task.Directive), request)
				if err != nil {
					tLogger.Error("Decode task directive [%s] failed: %+v", task.Directive, err)
					return err
				}

				err = retryutil.Retry(3, 0, func() error {
					_, err = pilotClient.RunCommandOnFrontgateNode(withTimeoutCtx, request)
					if err != nil {
						if strings.Contains(err.Error(), "context canceled") {
							tLogger.Debug("Expected error: %+v", err)
							return nil
						} else {
							tLogger.Error("%s", err.Error())
						}
					}
					return err
				})
				if err != nil {
					tLogger.Error("Send task to pilot failed: %+v", err)
					return err
				}

			case vmbased.ActionRunCommandOnDrone:
				request := new(pbtypes.RunCommandOnDroneRequest)
				err = jsonutil.Decode([]byte(task.Directive), request)
				if err != nil {
					tLogger.Error("Decode task directive [%s] failed: %+v", task.Directive, err)
					return err
				}

				err = retryutil.Retry(3, 0, func() error {
					_, err = pilotClient.RunCommandOnDrone(withTimeoutCtx, request)
					return err
				})
				if err != nil {
					tLogger.Error("Send task to pilot failed: %+v", err)
					return err
				}

			case vmbased.ActionRunCommandOnFrontgateNode:
				request := new(pbtypes.RunCommandOnFrontgateRequest)
				err = jsonutil.Decode([]byte(task.Directive), request)
				if err != nil {
					tLogger.Error("Decode task directive [%s] failed: %+v", task.Directive, err)
					return err
				}
				err = retryutil.Retry(3, 0, func() error {
					_, err = pilotClient.RunCommandOnFrontgateNode(withTimeoutCtx, request)
					return err
				})
				if err != nil {
					tLogger.Error("Send task to pilot failed: %+v", err)
					return err
				}

			case vmbased.ActionRegisterMetadata, vmbased.ActionDeregisterCmd, vmbased.ActionDeregisterMetadata:
				pbTask := models.TaskToPb(task)
				err = retryutil.Retry(3, 0, func() error {
					_, err := pilotClient.HandleSubtask(withTimeoutCtx,
						&pbtypes.SubTaskMessage{
							TaskId:    pbTask.TaskId.GetValue(),
							Action:    pbTask.TaskAction.GetValue(),
							Directive: pbTask.Directive.GetValue(),
						})
					return err
				})
				if err != nil {
					tLogger.Error("Failed to handle task to pilot: %+v", err)
					return err
				}

			case vmbased.ActionRegisterCmd:
				pbTask := models.TaskToPb(task)
				err = retryutil.Retry(3, 0, func() error {
					_, err := pilotClient.HandleSubtask(withTimeoutCtx,
						&pbtypes.SubTaskMessage{
							TaskId:    pbTask.TaskId.GetValue(),
							Action:    pbTask.TaskAction.GetValue(),
							Directive: pbTask.Directive.GetValue(),
						})
					return err
				})
				if err != nil {
					tLogger.Error("Failed to handle task to pilot: %+v", err)
					return err
				}
				err = pilotClient.WaitSubtask(
					ctx, task.TaskId, task.GetTimeout(constants.WaitTaskTimeout), constants.WaitTaskInterval)
				if err != nil {
					tLogger.Error("Failed to wait task: %+v", err)
					return err
				}

			default:
				tLogger.Error("Unknown task action [%s]", task.TaskAction)
			}
		} else {
			providerInterface, err := plugins.GetProviderPlugin(task.Target, tLogger)
			if err != nil {
				tLogger.Error("No such runtime [%s]. ", task.Target)
				return err
			}
			err = providerInterface.HandleSubtask(task)
			if err != nil {
				tLogger.Error("Failed to handle subtask in runtime [%s]: %+v", task.Target, err)
				return err
			}
			err = providerInterface.WaitSubtask(
				task, task.GetTimeout(constants.WaitTaskTimeout), constants.WaitTaskInterval)
			if err != nil {
				tLogger.Error("Failed to wait subtask in runtime [%s]: %+v", task.Target, err)
				return err
			}

			tLogger.Debug("After wait subtask directive: %s", task.Directive)
		}

		if err != nil {
			return err
		}

		err = processor.Post()
		if err != nil {
			tLogger.Error("Executing task post processor failed: %+v", err)
		}
		return err
	}()
	var status = constants.StatusSuccessful
	if err != nil {
		status = constants.StatusFailed

	}
	err = c.updateTaskAttributes(task.TaskId, map[string]interface{}{
		"status":      status,
		"status_time": time.Now(),
	})
	if err != nil {
		tLogger.Error("Failed to update task: %+v", err)
	}

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
