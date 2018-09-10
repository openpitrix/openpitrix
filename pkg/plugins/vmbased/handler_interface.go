// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package vmbased

import (
	"context"
	"fmt"

	"openpitrix.io/openpitrix/pkg/logger"
	"openpitrix.io/openpitrix/pkg/models"
)

type ProviderHandlerInterface interface {
	RunInstances(task *models.Task) error
	WaitRunInstances(task *models.Task) error

	StopInstances(task *models.Task) error
	WaitStopInstances(task *models.Task) error

	StartInstances(task *models.Task) error
	WaitStartInstances(task *models.Task) error

	DeleteInstances(task *models.Task) error
	WaitDeleteInstances(task *models.Task) error

	CreateVolumes(task *models.Task) error
	WaitCreateVolumes(task *models.Task) error

	DetachVolumes(task *models.Task) error
	WaitDetachVolumes(task *models.Task) error

	AttachVolumes(task *models.Task) error
	WaitAttachVolumes(task *models.Task) error

	DeleteVolumes(task *models.Task) error
	WaitDeleteVolumes(task *models.Task) error

	ResizeInstances(task *models.Task) error
	WaitResizeInstances(task *models.Task) error

	ResizeVolumes(task *models.Task) error
	WaitResizeVolumes(task *models.Task) error

	WaitFrontgateAvailable(task *models.Task) error
}

func HandleSubtask(ctx context.Context, task *models.Task, handler ProviderHandlerInterface) error {
	switch task.TaskAction {
	case ActionRunInstances:
		return handler.RunInstances(task)
	case ActionStopInstances:
		return handler.StopInstances(task)
	case ActionStartInstances:
		return handler.StartInstances(task)
	case ActionTerminateInstances:
		return handler.DeleteInstances(task)
	case ActionResizeInstances:
		return handler.ResizeInstances(task)
	case ActionCreateVolumes:
		return handler.CreateVolumes(task)
	case ActionDetachVolumes:
		return handler.DetachVolumes(task)
	case ActionAttachVolumes:
		return handler.AttachVolumes(task)
	case ActionDeleteVolumes:
		return handler.DeleteVolumes(task)
	case ActionResizeVolumes:
		return handler.ResizeVolumes(task)
	case ActionWaitFrontgateAvailable:
		return nil
	default:
		logger.Error(ctx, "Unknown task action [%s]", task.TaskAction)
		return fmt.Errorf("unknown task action [%s]", task.TaskAction)
	}
}

func WaitSubtask(ctx context.Context, task *models.Task, handler ProviderHandlerInterface) error {
	switch task.TaskAction {
	case ActionRunInstances:
		return handler.WaitRunInstances(task)
	case ActionStopInstances:
		return handler.WaitStopInstances(task)
	case ActionStartInstances:
		return handler.WaitStartInstances(task)
	case ActionTerminateInstances:
		return handler.WaitDeleteInstances(task)
	case ActionResizeInstances:
		return handler.WaitResizeInstances(task)
	case ActionCreateVolumes:
		return handler.WaitCreateVolumes(task)
	case ActionDetachVolumes:
		return handler.WaitDetachVolumes(task)
	case ActionAttachVolumes:
		return handler.WaitAttachVolumes(task)
	case ActionDeleteVolumes:
		return handler.WaitDeleteVolumes(task)
	case ActionResizeVolumes:
		return handler.WaitResizeVolumes(task)
	case ActionWaitFrontgateAvailable:
		return handler.WaitFrontgateAvailable(task)
	default:
		logger.Error(ctx, "Unknown task action [%s]", task.TaskAction)
		return fmt.Errorf("unknown task action [%s]", task.TaskAction)
	}
}
