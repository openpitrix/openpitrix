// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package manager

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"openpitrix.io/openpitrix/pkg/pb"
	"openpitrix.io/openpitrix/pkg/util/pbutil"
)

func TestChecker(t *testing.T) {
	var req Request
	var err error
	req = &pb.CreateAppRequest{
		Owner:  pbutil.ToProtoString(""),
		RepoId: nil,
	}
	err = NewChecker(context.Background(), req).Required("repo_id").Exec()

	assert.Error(t, err)

	req = &pb.CreateAppRequest{
		RepoId: pbutil.ToProtoString(""),
	}
	err = NewChecker(context.Background(), req).Required("repo_id").Exec()

	assert.Error(t, err)

	req = &pb.CreateAppRequest{
		RepoId: pbutil.ToProtoString("111"),
	}
	err = NewChecker(context.Background(), req).Required("repo_id").Exec()

	assert.NoError(t, err)
}
