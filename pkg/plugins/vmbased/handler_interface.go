// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package vmbased

import (
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

	WaitFrontgateAvailable(task *models.Task) error

	DescribeSubnet(runtimeId, subnetId string) (*models.Subnet, error)
	DescribeVpc(runtimeId, vpcId string) (*models.Vpc, error)
}
