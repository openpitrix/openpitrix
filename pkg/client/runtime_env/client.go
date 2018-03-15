package runtime_env

import (
	"context"

	"openpitrix.io/openpitrix/pkg/constants"
	"openpitrix.io/openpitrix/pkg/manager"
	"openpitrix.io/openpitrix/pkg/pb"
)

func NewRuntimeEnvManagerClient(ctx context.Context) (pb.RuntimeEnvManagerClient, error) {
	conn, err := manager.NewClient(ctx, constants.RuntimeEnvManagerHost, constants.RuntimeEnvManagerPort)
	if err != nil {
		return nil, err
	}
	return pb.NewRuntimeEnvManagerClient(conn), err
}

func DescribeRuntimeEnvs(request *pb.DescribeRuntimeEnvsRequest) (*pb.DescribeRuntimeEnvsResponse, error) {
	ctx := context.Background()
	client, err := NewRuntimeEnvManagerClient(ctx)
	if err != nil {
		return nil, err
	}
	response, err := client.DescribeRuntimeEnvs(ctx, request)
	if err != nil {
		return nil, err
	}
	return response, err
}
