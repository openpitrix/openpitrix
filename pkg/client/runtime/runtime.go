// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package runtime

import (
	"context"
	"strings"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	clientutil "openpitrix.io/openpitrix/pkg/client"
	"openpitrix.io/openpitrix/pkg/logger"
	"openpitrix.io/openpitrix/pkg/models"
	"openpitrix.io/openpitrix/pkg/pb"
)

type Runtime struct {
	models.Runtime
	Credential string
}

func NewRuntime(ctx context.Context, runtimeId string) (*Runtime, error) {
	runtime, err := getRuntime(ctx, runtimeId)
	if err != nil {
		return nil, err
	}
	result := &Runtime{
		Runtime:    *models.PbToRuntime(runtime.Runtime),
		Credential: runtime.GetRuntimeCredential().GetValue(),
	}
	return result, nil
}

func getRuntime(ctx context.Context, runtimeId string) (*pb.RuntimeDetail, error) {
	runtimeIds := []string{runtimeId}
	ctx = clientutil.SetSystemUserToContext(ctx)
	client, err := NewRuntimeManagerClient()
	if err != nil {
		return nil, err
	}
	response, err := client.DescribeRuntimeDetails(ctx, &pb.DescribeRuntimesRequest{
		RuntimeId: runtimeIds,
	})
	if err != nil {
		logger.Error(ctx, "Describe runtime [%s] failed: %+v",
			strings.Join(runtimeIds, ","), err)
		return nil, status.Errorf(codes.Internal, "Describe runtime [%s] failed: %+v",
			strings.Join(runtimeIds, ","), err)
	}

	if response.GetTotalCount() == 0 {
		logger.Error(ctx, "Runtime [%s] not found", strings.Join(runtimeIds, ","))
		return nil, status.Errorf(codes.PermissionDenied, "Runtime [%s] not found: %+v",
			strings.Join(runtimeIds, ","), err)
	}

	return response.RuntimeDetailSet[0], nil
}
