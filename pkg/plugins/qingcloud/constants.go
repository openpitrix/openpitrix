// Copyright 2017 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package qingcloud

import "time"

const (
	TIMEOUT_CREATE_CLUSTER  = 600 * time.Second
	TIMEOUT_START_CLUSTER   = 600 * time.Second
	TIMEOUT_STOP_CLUSTER    = 600 * time.Second
	TIMEOUT_DELETE_CLUSTER  = 600 * time.Second
	TIMEOUT_RECOVER_CLUSTER = 600 * time.Second
	TIMEOUT_CEASE_CLUSTER   = 600 * time.Second

	WAIT_INTERVAL = 20 * time.Second

	STATUS_ACTIVE = "active"
)
