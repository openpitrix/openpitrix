// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package cluster

import (
	"strings"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"openpitrix.io/openpitrix/pkg/client/runtime_env"
	"openpitrix.io/openpitrix/pkg/constants"
	"openpitrix.io/openpitrix/pkg/logger"
	"openpitrix.io/openpitrix/pkg/pb"
	"openpitrix.io/openpitrix/pkg/plugins"
)

type Runtime struct {
	RuntimeEnvId     string
	Runtime          string
	RuntimeInterface plugins.RuntimeInterface
}

func NewRuntime(runtimeEnvId string) (*Runtime, error) {
	runtimeEnv, err := getRuntimeEnv(runtimeEnvId)
	if err != nil {
		return nil, err
	}
	runtime := getRuntime(runtimeEnv)
	runtimeInterface, err := plugins.GetRuntimePlugin(runtime)
	if err != nil {
		logger.Errorf("No such runtime [%s]. ", runtime)
		return nil, err
	}

	result := &Runtime{
		RuntimeEnvId:     runtimeEnvId,
		Runtime:          runtime,
		RuntimeInterface: runtimeInterface,
	}
	return result, nil
}

func getRuntimeEnv(runtimeEnvId string) (*pb.RuntimeEnv, error) {
	runtimeEnvIds := []string{runtimeEnvId}
	response, err := runtime_env.DescribeRuntimeEnvs(&pb.DescribeRuntimeEnvsRequest{
		RuntimeEnvId: runtimeEnvIds,
	})
	if err != nil {
		logger.Errorf("Describe runtime env [%s] failed: %+v",
			strings.Join(runtimeEnvIds, ","), err)
		return nil, status.Errorf(codes.Internal, "Describe runtime env [%s] failed: %+v",
			strings.Join(runtimeEnvIds, ","), err)
	}

	if response.GetTotalCount() == 0 {
		logger.Errorf("Runtime env [%s] not found", strings.Join(runtimeEnvIds, ","))
		return nil, status.Errorf(codes.PermissionDenied, "Runtime env [%s] not found",
			strings.Join(runtimeEnvIds, ","), err)
	}

	return response.RuntimeEnvSet[0], nil
}

func getRuntime(runtimeEnv *pb.RuntimeEnv) string {
	// TODO: need to parse runtime
	return constants.RuntimeQingCloud
}
