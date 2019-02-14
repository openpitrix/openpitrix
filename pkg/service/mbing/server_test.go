// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package mbing

import (
	"context"
	"testing"
	"time"

)

func common(t *testing.T) (*Server, context.Context, context.CancelFunc) {
	if !*tTestingEnvEnabled {
		t.Skip("testing env disabled")
	}
	InitGlobelSetting()
	s, _ := NewServer()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	return s, ctx, cancel
}

func TestNewServer(t *testing.T) {
	s, ctx, cancle := common(t)
	t.Logf("TestNewServer Passed, server: %v, context: %v, cancleFunc: %v", s, ctx, cancle)
}

