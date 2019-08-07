// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package task

import (
	"context"
	"strings"
	"sync"
	"time"

	pilotclient "openpitrix.io/openpitrix/pkg/client/pilot"
	providerclient "openpitrix.io/openpitrix/pkg/client/runtime_provider"
	"openpitrix.io/openpitrix/pkg/constants"
	"openpitrix.io/openpitrix/pkg/db"
	"openpitrix.io/openpitrix/pkg/etcd"
	"openpitrix.io/openpitrix/pkg/gerr"
	"openpitrix.io/openpitrix/pkg/logger"
	"openpitrix.io/openpitrix/pkg/models"
	"openpitrix.io/openpitrix/pkg/pb"
	pbtypes "openpitrix.io/openpitrix/pkg/pb/metadata/types"
	"openpitrix.io/openpitrix/pkg/pi"
	"openpitrix.io/openpitrix/pkg/plugins/vmbased"
	"openpitrix.io/openpitrix/pkg/sender"
	"openpitrix.io/openpitrix/pkg/util/ctxutil"
	"openpitrix.io/openpitrix/pkg/util/funcutil"
	"openpitrix.io/openpitrix/pkg/util/jsonutil"
	"openpitrix.io/openpitrix/pkg/util/pbutil"
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
		Where(db.Eq(constants.ColumnTaskId, taskId)).
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
		SetMap(map[string]interface{}{constants.ColumnStatus: constants.StatusFailed}).
		Where(db.Eq(constants.ColumnStatus, constants.StatusWorking)).
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
		Where(db.Eq(constants.ColumnTaskId, taskId))

	err := query.LoadOne(&task)
	if err != nil {
		logger.Error(ctx, "Failed to get task [%s]: %+v", task.TaskId, err)
		return err
	}
	ctx = ctxutil.AddMessageId(ctx, task.JobId)

	ctx = ctxutil.ContextWithSender(ctx, sender.New(task.Owner, task.OwnerPath, ""))

	err = c.updateTaskAttributes(ctx, task.TaskId, map[string]interface{}{
		constants.ColumnStatus:   constants.StatusWorking,
		constants.ColumnExecutor: c.hostname,
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
			switch task.TaskAction {
			case vmbased.ActionSetDroneConfig:
				config := new(pbtypes.SetDroneConfigRequest)
				err = jsonutil.Decode([]byte(task.Directive), config)
				if err != nil {
					logger.Error(ctx, "Decode task directive [%s] failed: %+v", task.Directive, err)
					return err
				}
				err = retryutil.RetryWithContext(ctx, constants.PilotTasksRetry, constants.PilotTasksSleep, func() error {
					_, err = pilotClient.SetDroneConfigWithTimeout(ctx, config)
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
				err = retryutil.RetryWithContext(ctx, constants.PilotTasksRetry, constants.PilotTasksSleep, func() error {
					_, err = pilotClient.SetFrontgateConfigWithTimeout(ctx, config)
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
					_, err := pilotClient.PingDroneWithTimeout(ctx, droneEndpoint)
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
					_, err := pilotClient.PingFrontgateWithTimeout(ctx, request)
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
					_, err := pilotClient.PingMetadataBackendWithTimeout(ctx, request)
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
				err = retryutil.RetryWithContext(ctx, constants.PilotTasksRetry, constants.PilotTasksSleep, func() error {
					_, err := pilotClient.HandleSubtaskWithTimeout(ctx,
						&pbtypes.SubTaskMessage{
							TaskId:    pbTask.TaskId.GetValue(),
							Action:    pbTask.TaskAction.GetValue(),
							Directive: pbTask.Directive.GetValue(),
						})
					if err != nil && strings.Contains(err.Error(), "drone: confd is running") {
						logger.Debug(ctx, "Expected error: %+v", err)
						return nil
					}
					return err
				})
				if err != nil {
					logger.Error(ctx, "Failed to handle task to pilot: %+v", err)
					return err
				}

				time.Sleep(1 * time.Second)

			case vmbased.ActionStopConfd:
				pbTask := models.TaskToPb(task)
				err = retryutil.RetryWithContext(ctx, constants.PilotTasksRetry, constants.PilotTasksSleep, func() error {
					_, err := pilotClient.HandleSubtaskWithTimeout(ctx,
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

				err = retryutil.RetryWithContext(ctx, constants.PilotTasksRetry, constants.PilotTasksSleep, func() error {
					_, err = pilotClient.RunCommandOnDroneWithTimeout(ctx, request)
					if err != nil && strings.Contains(err.Error(), "transport is closing") {
						logger.Debug(ctx, "Expected error: %+v", err)
						return nil
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

				err = retryutil.RetryWithContext(ctx, constants.PilotTasksRetry, constants.PilotTasksSleep, func() error {
					_, err = pilotClient.RunCommandOnFrontgateNodeWithTimeout(ctx, request)
					if err != nil && strings.Contains(err.Error(), "context canceled") {
						logger.Debug(ctx, "Expected error: %+v", err)
						return nil
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

				err = retryutil.RetryWithContext(ctx, constants.PilotTasksRetry, constants.PilotTasksSleep, func() error {
					_, err = pilotClient.RunCommandOnDroneWithTimeout(ctx, request)
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
				err = retryutil.RetryWithContext(ctx, constants.PilotTasksRetry, constants.PilotTasksSleep, func() error {
					_, err = pilotClient.RunCommandOnFrontgateNodeWithTimeout(ctx, request)
					return err
				})
				if err != nil {
					logger.Error(ctx, "Send task to pilot failed: %+v", err)
					return err
				}

			case vmbased.ActionRegisterMetadata,
				vmbased.ActionDeregisterCmd,
				vmbased.ActionDeregisterMetadata,
				vmbased.ActionRegisterMetadataMapping,
				vmbased.ActionDeregisterMetadataMapping:
				pbTask := models.TaskToPb(task)
				err = retryutil.RetryWithContext(ctx, constants.PilotTasksRetry, constants.PilotTasksSleep, func() error {
					_, err := pilotClient.HandleSubtaskWithTimeout(ctx,
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
				err = retryutil.RetryWithContext(ctx, constants.PilotTasksRetry, constants.PilotTasksSleep, func() error {
					_, err := pilotClient.HandleSubtaskWithTimeout(ctx,
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
			providerClient, err := providerclient.NewRuntimeProviderManagerClient()
			if err != nil {
				return gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorInternalError)
			}
			handleResponse, err := providerClient.HandleSubtask(ctx, &pb.HandleSubtaskRequest{
				RuntimeId: pbutil.ToProtoString(task.Target),
				Task:      models.TaskToPb(task),
			})
			if err != nil {
				logger.Error(ctx, "Failed to handle subtask in runtime [%s]: %+v", task.Target, err)
				return err
			}
			withTimeoutCtx, cancel := context.WithTimeout(ctx, constants.MaxTaskTimeout)
			defer cancel()
			waitResponse, err := providerClient.WaitSubtask(withTimeoutCtx, &pb.WaitSubtaskRequest{
				RuntimeId: handleResponse.Task.Target,
				Task:      handleResponse.Task,
			})
			if err != nil {
				logger.Error(ctx, "Failed to wait subtask in runtime [%s]: %+v", task.Target, err)
				return err
			}
			processor.Task = models.PbToTask(waitResponse.Task)

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
		constants.ColumnStatus:     status,
		constants.ColumnStatusTime: time.Now(),
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

func (c *Controller) Serve(ctx context.Context) {
	err := c.UpdateWorkingTasks(ctx)
	if err != nil {
		logger.Error(ctx, "Failed to update working tasks: %+v", err)
	}
	go c.ExtractTasks(ctx)
	go c.HandleTasks(ctx)
}
