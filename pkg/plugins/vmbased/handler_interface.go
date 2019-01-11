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
	RunInstances(ctx context.Context, task *models.Task) (*models.Task, error)
	WaitRunInstances(ctx context.Context, task *models.Task) (*models.Task, error)

	StopInstances(ctx context.Context, task *models.Task) (*models.Task, error)
	WaitStopInstances(ctx context.Context, task *models.Task) (*models.Task, error)

	StartInstances(ctx context.Context, task *models.Task) (*models.Task, error)
	WaitStartInstances(ctx context.Context, task *models.Task) (*models.Task, error)

	DeleteInstances(ctx context.Context, task *models.Task) (*models.Task, error)
	WaitDeleteInstances(ctx context.Context, task *models.Task) (*models.Task, error)

	CreateVolumes(ctx context.Context, task *models.Task) (*models.Task, error)
	WaitCreateVolumes(ctx context.Context, task *models.Task) (*models.Task, error)

	DetachVolumes(ctx context.Context, task *models.Task) (*models.Task, error)
	WaitDetachVolumes(ctx context.Context, task *models.Task) (*models.Task, error)

	AttachVolumes(ctx context.Context, task *models.Task) (*models.Task, error)
	WaitAttachVolumes(ctx context.Context, task *models.Task) (*models.Task, error)

	DeleteVolumes(ctx context.Context, task *models.Task) (*models.Task, error)
	WaitDeleteVolumes(ctx context.Context, task *models.Task) (*models.Task, error)

	ResizeInstances(ctx context.Context, task *models.Task) (*models.Task, error)
	WaitResizeInstances(ctx context.Context, task *models.Task) (*models.Task, error)

	ResizeVolumes(ctx context.Context, task *models.Task) (*models.Task, error)
	WaitResizeVolumes(ctx context.Context, task *models.Task) (*models.Task, error)

	WaitFrontgateAvailable(ctx context.Context, task *models.Task) (*models.Task, error)
}

func HandleSubtask(ctx context.Context, task *models.Task, handler ProviderHandlerInterface) (*models.Task, error) {
	switch task.TaskAction {
	case ActionRunInstances:
		return handler.RunInstances(ctx, task)
	case ActionStopInstances:
		return handler.StopInstances(ctx, task)
	case ActionStartInstances:
		return handler.StartInstances(ctx, task)
	case ActionTerminateInstances:
		return handler.DeleteInstances(ctx, task)
	case ActionResizeInstances:
		return handler.ResizeInstances(ctx, task)
	case ActionCreateVolumes:
		return handler.CreateVolumes(ctx, task)
	case ActionDetachVolumes:
		return handler.DetachVolumes(ctx, task)
	case ActionAttachVolumes:
		return handler.AttachVolumes(ctx, task)
	case ActionDeleteVolumes:
		return handler.DeleteVolumes(ctx, task)
	case ActionResizeVolumes:
		return handler.ResizeVolumes(ctx, task)
	case ActionWaitFrontgateAvailable:
		return task, nil
	default:
		logger.Error(ctx, "Unknown task action [%s]", task.TaskAction)
		return nil, fmt.Errorf("unknown task action [%s]", task.TaskAction)
	}
}

func WaitSubtask(ctx context.Context, task *models.Task, handler ProviderHandlerInterface) (*models.Task, error) {
	switch task.TaskAction {
	case ActionRunInstances:
		return handler.WaitRunInstances(ctx, task)
	case ActionStopInstances:
		return handler.WaitStopInstances(ctx, task)
	case ActionStartInstances:
		return handler.WaitStartInstances(ctx, task)
	case ActionTerminateInstances:
		return handler.WaitDeleteInstances(ctx, task)
	case ActionResizeInstances:
		return handler.WaitResizeInstances(ctx, task)
	case ActionCreateVolumes:
		return handler.WaitCreateVolumes(ctx, task)
	case ActionDetachVolumes:
		return handler.WaitDetachVolumes(ctx, task)
	case ActionAttachVolumes:
		return handler.WaitAttachVolumes(ctx, task)
	case ActionDeleteVolumes:
		return handler.WaitDeleteVolumes(ctx, task)
	case ActionResizeVolumes:
		return handler.WaitResizeVolumes(ctx, task)
	case ActionWaitFrontgateAvailable:
		return handler.WaitFrontgateAvailable(ctx, task)
	default:
		logger.Error(ctx, "Unknown task action [%s]", task.TaskAction)
		return nil, fmt.Errorf("unknown task action [%s]", task.TaskAction)
	}
}
