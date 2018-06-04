// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package gerr

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/status"
)

func TestNew(t *testing.T) {
	var e error
	e = New(Internal, ErrorCreateResourceFailed)
	assert.Equal(t, e.Error(), "rpc error: code = Internal desc = create resource failed")
	e = New(InvalidArgument, ErrorMissingParameter, "name")
	assert.Equal(t, e.Error(), "rpc error: code = InvalidArgument desc = missing parameter [name]")
}

func TestNewWithDetail(t *testing.T) {
	var e error
	e = NewWithDetail(InvalidArgument, fmt.Errorf("test with error detail"), ErrorCreateResourceFailed)
	assert.Equal(t, e.Error(), "rpc error: code = InvalidArgument desc = create resource failed")
	ge := status.Convert(e)
	assert.Equal(t, ge.Code().String(), "InvalidArgument")
	assert.Equal(t, ge.Err().Error(), "rpc error: code = InvalidArgument desc = create resource failed")
	assert.Equal(t, fmt.Sprint(ge.Details()), "[error_name:\"create_resource_failed\" cause:\"test with error detail\" ]")
	//t.Log(ge.Code(), ge.Err(), ge.Details())
}

func TestClearErrorCause(t *testing.T) {
	var e = NewWithDetail(InvalidArgument, fmt.Errorf("test with error detail"), ErrorCreateResourceFailed)

	ge := status.Convert(e)
	assert.Equal(t, fmt.Sprint(ge.Details()), "[error_name:\"create_resource_failed\" cause:\"test with error detail\" ]")

	e = ClearErrorCause(e)
	ge = status.Convert(e)
	assert.Equal(t, fmt.Sprint(ge.Details()), "[error_name:\"create_resource_failed\" ]")
}
