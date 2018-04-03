// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package runtime

import (
	"context"
	"strings"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"openpitrix.io/openpitrix/pkg/constants"
	"openpitrix.io/openpitrix/pkg/logger"
	"openpitrix.io/openpitrix/pkg/pb"
	"openpitrix.io/openpitrix/pkg/plugins"
)

type Runtime struct {
	RuntimeId         string
	Provider          string
	Zone              string
	ProviderInterface plugins.ProviderInterface
	Credential        map[string]string
	Url               string
}

func NewRuntime(runtimeId string) (*Runtime, error) {
	runtime, err := getRuntime(runtimeId)
	if err != nil {
		return nil, err
	}
	provider := getProvider(runtime)
	zone := getZone(runtime)
	providerInterface, err := plugins.GetProviderPlugin(provider)
	if err != nil {
		logger.Errorf("No such provider [%s]. ", provider)
		return nil, err
	}

	result := &Runtime{
		RuntimeId:         runtimeId,
		Provider:          provider,
		Zone:              zone,
		ProviderInterface: providerInterface,
	}
	return result, nil
}

func getRuntime(runtimeId string) (*pb.Runtime, error) {
	runtimeIds := []string{runtimeId}
	ctx := context.Background()
	client, err := NewRuntimeManagerClient(ctx)
	if err != nil {
		return nil, err
	}
	response, err := client.DescribeRuntimes(ctx, &pb.DescribeRuntimesRequest{
		RuntimeId: runtimeIds,
	})
	if err != nil {
		logger.Errorf("Describe runtime [%s] failed: %+v",
			strings.Join(runtimeIds, ","), err)
		return nil, status.Errorf(codes.Internal, "Describe runtime [%s] failed: %+v",
			strings.Join(runtimeIds, ","), err)
	}

	if response.GetTotalCount() == 0 {
		logger.Errorf("Runtime [%s] not found", strings.Join(runtimeIds, ","))
		return nil, status.Errorf(codes.PermissionDenied, "Runtime [%s] not found",
			strings.Join(runtimeIds, ","), err)
	}

	return response.RuntimeSet[0], nil
}

func getProvider(runtime *pb.Runtime) string {
	// TODO: need to parse runtime
	return constants.ProviderQingCloud
}

func getZone(runtime *pb.Runtime) string {
	// TODO: need to parse runtime
	return "testing"
}
