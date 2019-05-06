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
		Name: nil,
	}
	err = NewChecker(context.Background(), req).Required("name").Exec()

	assert.Error(t, err)

	req = &pb.CreateAppRequest{
		Name: pbutil.ToProtoString(""),
	}
	err = NewChecker(context.Background(), req).Required("name").Exec()

	assert.Error(t, err)

	req = &pb.CreateAppRequest{
		Name: pbutil.ToProtoString("111"),
	}
	err = NewChecker(context.Background(), req).Required("name").Exec()

	assert.NoError(t, err)

	req = &pb.CreateRepoRequest{
		Providers: []string{"qingcloud", "aws"},
	}
	err = NewChecker(context.Background(), req).
		Required("providers").
		StringChosen("providers", []string{"qingcloud", "aws", "k8s"}).
		Exec()

	assert.NoError(t, err)

	req = &pb.CreateRepoRequest{
		Providers: []string{"qingcloud", "xxxx"},
	}
	err = NewChecker(context.Background(), req).
		Required("providers").
		StringChosen("providers", []string{"qingcloud", "aws", "k8s"}).
		Exec()

	assert.Error(t, err)

	//req = &pb.CreateRepoRequest{}
	//ctx := ctxutil.ContextWithSender(context.Background(), sender.GetSystemSender())
	//err = NewChecker(ctx, req).OperatorType([]string{"developer"}).Exec()
	//assert.Error(t, err)
}
