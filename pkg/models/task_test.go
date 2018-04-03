// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package models

import (
	"testing"
	"time"
)

func TestGetTimeout(t *testing.T) {
	timeout := 60 * time.Second
	instance := &Instance{
		Timeout: int(timeout / time.Second),
	}
	directive, err := instance.ToString()
	if err != nil {
		t.Errorf("Error: %+v", err)
	}

	task := &Task{
		Directive: directive,
	}

	taskTimeout := task.GetTimeout(20 * time.Second)
	if timeout != taskTimeout {
		t.Errorf("Expect timeout %d, get timeout %d", timeout, taskTimeout)
	}
}
