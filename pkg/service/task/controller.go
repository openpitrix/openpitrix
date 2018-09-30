// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package task

import (
	"context"
	"strings"
	"sync"
	"time"

	"openpitrix.io/openpitrix/pkg/client"
	accountclient "openpitrix.io/openpitrix/pkg/client/iam"
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
	"openpitrix.io/openpitrix/pkg/util/ctxutil"
	"openpitrix.io/openpitrix/pkg/util/funcutil"
	"openpitrix.io/openpitrix/pkg/util/jsonutil"
	"openpitrix.io/openpitrix/pkg/util/retryutil"
)

type Controller struct {
	runningTasks chan string
	runningCount int32
	hostname     string
	queue        *etcd.Queue
}

func NewController(hostname string) *Controller {
	return &Controller{
		runningTasks: make(chan string),
		runningCount: 0,
		hostname:     hostname,
		queue:        pi.Global().Etcd(context.Background()).NewQueue("task"),
	}
}

func (c *Controller) updateTaskAttributes(ctx context.Context, taskId string, attributes map[string]interface{}) error {
	_, err := pi.Global().DB(ctx).
		Update(constants.TableTask).
		SetMap(attributes).
		Where(db.Eq("task_id", taskId)).
		Exec()

	return err
}

var mutex sync.Mutex

func (c *Controller) GetTaskLength() int32 {
	if pi.Global().GlobalConfig().Task.MaxWorkingTasks > 0 {
		return pi.Global().GlobalConfig().Task.MaxWorkingTasks
	} else {
		return constants.DefaultMaxWorkingTasks
	}
}

func (c *Controller) IsRunningExceed() bool {
	mutex.Lock()
	defer mutex.Unlock()
	count := c.runningCount
	return count > c.GetTaskLength()
}

func (c *Controller) UpdateWorkingTasks(ctx context.Context) error {
	//TODO: retry the tasks
	_, err := pi.Global().DB(ctx).
		Update(constants.TableTask).
		SetMap(map[string]interface{}{"status": constants.StatusFailed}).
		Where(db.Eq("status", constants.StatusWorking)).
		Exec()
	return err
}

func (c *Controller) ExtractTasks(ctx context.Context) {
	for {
		if c.IsRunningExceed() {
			logger.Error(ctx, "Sleep 10s, running task count exceed [%d/%d]", c.runningCount, c.GetTaskLength())
			time.Sleep(10 * time.Second)
			continue
		}
		taskId, err := c.queue.Dequeue()
		if err != nil {
			logger.Error(ctx, "Failed to dequeue task from etcd queue: %+v", err)
			time.Sleep(3 * time.Second)
			continue
		}
		logger.Debug(ctx, "Dequeue task [%s] from etcd queue success", taskId)
		c.runningTasks <- taskId
	}
}

func (c *Controller) HandleTask(ctx context.Context, taskId string, cb func()) error {
	ctx = ctxutil.AddMessageId(ctx, taskId)
	defer cb()
	task := new(models.Task)
	query := pi.Global().DB(ctx).
		Select(models.TaskColumns...).
		From(constants.TableTask).
		Where(db.Eq("task_id", taskId))

	err := query.LoadOne(&task)
	if err != nil {
		logger.Error(ctx, "Failed to get task [%s]: %+v", task.TaskId, err)
		return err
	}
	ctx = ctxutil.AddMessageId(ctx, task.JobId)

	accountClient, err := accountclient.NewClient()
	if err != nil {
		return err
	}
	users, err := accountClient.GetUsers(ctx, []string{task.Owner})
	if err != nil {
		return err
	}

	ctx = client.SetUserToContext(ctx, users[0])

	err = c.updateTaskAttributes(ctx, task.TaskId, map[string]interface{}{
		"status":   constants.StatusWorking,
		"executor": c.hostname,
	})
	if err != nil {
		logger.Error(ctx, "Failed to update task: %+v", err)
		return err
	}

	err = func() error {
		processor := NewProcessor(task)
		err = processor.Pre(ctx)
		if err != nil {
			logger.Error(ctx, "Executing task pre processor failed: %+v", err)
			return err
		}

		pilotClient, err := pilotclient.NewClient()
		if err != nil {
			logger.Error(ctx, "Connect to pilot service failed: %+v", err)
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
					logger.Error(ctx, "Decode task directive [%s] failed: %+v", task.Directive, err)
					return err
				}
				err = retryutil.Retry(3, 0, func() error {
					_, err = pilotClient.SetDroneConfig(withTimeoutCtx, config)
					return err
				})
				if err != nil {
					logger.Error(ctx, "Send task to pilot failed: %+v", err)
					return err
				}
			case vmbased.ActionSetFrontgateConfig:
				config := new(pbtypes.FrontgateConfig)
				err = jsonutil.Decode([]byte(task.Directive), config)
				if err != nil {
					logger.Error(ctx, "Decode task directive [%s] failed: %+v", task.Directive, err)
					return err
				}
				err = retryutil.Retry(3, 0, func() error {
					_, err = pilotClient.SetFrontgateConfig(withTimeoutCtx, config)
					return err
				})
				if err != nil {
					logger.Error(ctx, "Send task to pilot failed: %+v", err)
					return err
				}
			case vmbased.ActionPingDrone:
				droneEndpoint := new(pbtypes.DroneEndpoint)
				err = jsonutil.Decode([]byte(task.Directive), droneEndpoint)
				if err != nil {
					logger.Error(ctx, "Decode task directive [%s] failed: %+v", task.Directive, err)
					return err
				}
				err = funcutil.WaitForSpecificOrError(func() (bool, error) {
					withTimeoutCtx, cancel := context.WithTimeout(ctx, constants.GrpcToPilotTimeout)
					defer cancel()
					_, err := pilotClient.PingDrone(withTimeoutCtx, droneEndpoint)
					if err != nil {
						logger.Warn(ctx, "Send task to pilot failed, will retry: %+v", err)
						return false, nil
					} else {
						return true, nil
					}
				}, task.GetTimeout(constants.WaitDroneServiceTimeout), constants.WaitDroneServiceInterval)
				if err != nil {
					logger.Error(ctx, "Send task to pilot failed: %+v", err)
					return err
				}

			case vmbased.ActionPingFrontgate:
				request := new(pbtypes.FrontgateId)
				err = jsonutil.Decode([]byte(task.Directive), request)
				if err != nil {
					logger.Error(ctx, "Decode task directive [%s] failed: %+v", task.Directive, err)
					return err
				}
				err = funcutil.WaitForSpecificOrError(func() (bool, error) {
					withTimeoutCtx, cancel := context.WithTimeout(ctx, constants.GrpcToPilotTimeout)
					defer cancel()
					_, err := pilotClient.PingFrontgate(withTimeoutCtx, request)
					if err != nil {
						logger.Warn(ctx, "Send task to pilot failed, will retry: %+v", err)
						return false, nil
					} else {
						return true, nil
					}
				}, task.GetTimeout(constants.WaitFrontgateServiceTimeout), constants.WaitFrontgateServiceInterval)
				if err != nil {
					logger.Error(ctx, "Send task to pilot failed: %+v", err)
					return err
				}

			case vmbased.PingMetadataBackend:
				request := new(pbtypes.FrontgateId)
				err = jsonutil.Decode([]byte(task.Directive), request)
				if err != nil {
					logger.Error(ctx, "Decode task directive [%s] failed: %+v", task.Directive, err)
					return err
				}
				err = funcutil.WaitForSpecificOrError(func() (bool, error) {
					withTimeoutCtx, cancel := context.WithTimeout(ctx, constants.GrpcToPilotTimeout)
					defer cancel()
					_, err := pilotClient.PingMetadataBackend(withTimeoutCtx, request)
					if err != nil {
						logger.Warn(ctx, "Send task to pilot failed, will retry: %+v", err)
						return false, nil
					} else {
						return true, nil
					}
				}, task.GetTimeout(constants.WaitFrontgateServiceTimeout), constants.WaitFrontgateServiceInterval)
				if err != nil {
					logger.Error(ctx, "Send task to pilot failed: %+v", err)
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
					logger.Error(ctx, "Failed to handle task to pilot: %+v", err)
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
					logger.Error(ctx, "Failed to handle task to pilot: %+v", err)
					return err
				}

				logger.Debug(ctx, "Finish subtask [%s]", task.TaskId)

				time.Sleep(1 * time.Second)

			case vmbased.ActionRemoveContainerOnDrone:
				request := new(pbtypes.RunCommandOnDroneRequest)
				err = jsonutil.Decode([]byte(task.Directive), request)
				if err != nil {
					logger.Error(ctx, "Decode task directive [%s] failed: %+v", task.Directive, err)
					return err
				}

				err = retryutil.Retry(3, 0, func() error {
					_, err = pilotClient.RunCommandOnDrone(withTimeoutCtx, request)
					if err != nil {
						if strings.Contains(err.Error(), "transport is closing") {
							logger.Debug(ctx, "Expected error: %+v", err)
							return nil
						} else {
							logger.Error(ctx, "%s", err.Error())
						}
					}
					return err
				})
				if err != nil {
					logger.Error(ctx, "Send task to pilot failed: %+v", err)
					return err
				}

			case vmbased.ActionRemoveContainerOnFrontgate:
				request := new(pbtypes.RunCommandOnFrontgateRequest)
				err = jsonutil.Decode([]byte(task.Directive), request)
				if err != nil {
					logger.Error(ctx, "Decode task directive [%s] failed: %+v", task.Directive, err)
					return err
				}

				err = retryutil.Retry(3, 0, func() error {
					_, err = pilotClient.RunCommandOnFrontgateNode(withTimeoutCtx, request)
					if err != nil {
						if strings.Contains(err.Error(), "context canceled") {
							logger.Debug(ctx, "Expected error: %+v", err)
							return nil
						} else {
							logger.Error(ctx, "%s", err.Error())
						}
					}
					return err
				})
				if err != nil {
					logger.Error(ctx, "Send task to pilot failed: %+v", err)
					return err
				}

			case vmbased.ActionRunCommandOnDrone:
				request := new(pbtypes.RunCommandOnDroneRequest)
				err = jsonutil.Decode([]byte(task.Directive), request)
				if err != nil {
					logger.Error(ctx, "Decode task directive [%s] failed: %+v", task.Directive, err)
					return err
				}

				err = retryutil.Retry(3, 0, func() error {
					_, err = pilotClient.RunCommandOnDrone(withTimeoutCtx, request)
					return err
				})
				if err != nil {
					logger.Error(ctx, "Send task to pilot failed: %+v", err)
					return err
				}

			case vmbased.ActionRunCommandOnFrontgateNode:
				request := new(pbtypes.RunCommandOnFrontgateRequest)
				err = jsonutil.Decode([]byte(task.Directive), request)
				if err != nil {
					logger.Error(ctx, "Decode task directive [%s] failed: %+v", task.Directive, err)
					return err
				}
				err = retryutil.Retry(3, 0, func() error {
					_, err = pilotClient.RunCommandOnFrontgateNode(withTimeoutCtx, request)
					return err
				})
				if err != nil {
					logger.Error(ctx, "Send task to pilot failed: %+v", err)
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
					logger.Error(ctx, "Failed to handle task to pilot: %+v", err)
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
					logger.Error(ctx, "Failed to handle task to pilot: %+v", err)
					return err
				}
				err = pilotClient.WaitSubtask(
					ctx, task.TaskId, task.GetTimeout(constants.WaitTaskTimeout), constants.WaitTaskInterval)
				if err != nil {
					logger.Error(ctx, "Failed to wait task: %+v", err)
					return err
				}

			default:
				logger.Error(ctx, "Unknown task action [%s]", task.TaskAction)
			}
		} else {
			providerInterface, err := plugins.GetProviderPlugin(ctx, task.Target)
			if err != nil {
				logger.Error(ctx, "No such runtime [%s]. ", task.Target)
				return err
			}
			err = providerInterface.HandleSubtask(ctx, task)
			if err != nil {
				logger.Error(ctx, "Failed to handle subtask in runtime [%s]: %+v", task.Target, err)
				return err
			}
			err = providerInterface.WaitSubtask(ctx, task)
			if err != nil {
				logger.Error(ctx, "Failed to wait subtask in runtime [%s]: %+v", task.Target, err)
				return err
			}

			logger.Debug(ctx, "After wait subtask directive: %s", task.Directive)
		}

		if err != nil {
			return err
		}

		err = processor.Post(ctx)
		if err != nil {
			logger.Error(ctx, "Executing task post processor failed: %+v", err)
		}
		return err
	}()
	var status = constants.StatusSuccessful
	if err != nil {
		status = constants.StatusFailed

	}
	err = c.updateTaskAttributes(ctx, task.TaskId, map[string]interface{}{
		"status":      status,
		"status_time": time.Now(),
	})
	if err != nil {
		logger.Error(ctx, "Failed to update task: %+v", err)
	}

	return err
}

func (c *Controller) HandleTasks(ctx context.Context) {
	for taskId := range c.runningTasks {
		mutex.Lock()
		c.runningCount++
		mutex.Unlock()

		go c.HandleTask(ctx, taskId, func() {
			mutex.Lock()
			c.runningCount--
			mutex.Unlock()
		})
	}
}

func (c *Controller) Serve() {
	ctx := context.Background()
	err := c.UpdateWorkingTasks(ctx)
	if err != nil {
		logger.Error(ctx, "Failed to update working tasks: %+v", err)
	}
	go c.ExtractTasks(ctx)
	go c.HandleTasks(ctx)
}
