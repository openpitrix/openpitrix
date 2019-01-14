// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package runtime_provider

import (
	"context"
	"fmt"

	"openpitrix.io/openpitrix/pkg/constants"
	"openpitrix.io/openpitrix/pkg/manager"
	"openpitrix.io/openpitrix/pkg/pb"
	"openpitrix.io/openpitrix/pkg/util/pbutil"
	"openpitrix.io/openpitrix/pkg/util/retryutil"
)

func NewRuntimeProviderManagerClient() (pb.RuntimeProviderManagerClient, error) {
	conn, err := manager.NewClient(constants.RuntimeProviderManagerHost, constants.RuntimeProviderManagerPort)
	if err != nil {
		return nil, err
	}
	return pb.NewRuntimeProviderManagerClient(conn), err
}

func NewRuntimeProviderClient(host string, port int) (pb.RuntimeProviderManagerClient, error) {
	conn, err := manager.NewClient(host, port)
	if err != nil {
		return nil, err
	}
	return pb.NewRuntimeProviderManagerClient(conn), err
}

func RegisterRuntimeProvider(provider, config string) error {
	// wait 5 min at most
	err := retryutil.Retry(60, 5, func() error {
		providerClient, err := NewRuntimeProviderManagerClient()
		if err != nil {
			return err
		}
		response, err := providerClient.RegisterRuntimeProvider(
			context.Background(),
			&pb.RegisterRuntimeProviderRequest{
				Provider: pbutil.ToProtoString(provider),
				Config:   pbutil.ToProtoString(config),
			})
		if err != nil {
			return err
		} else if !response.Ok.GetValue() {
			return fmt.Errorf("response is not ok")
		}
		return nil
	})
	return err
}
