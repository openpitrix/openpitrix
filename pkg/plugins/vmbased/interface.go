// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package vmbased

import "openpitrix.io/openpitrix/pkg/models"

type VmBasedInterface interface {
	RunInstance(instance *models.Instance) error
}
