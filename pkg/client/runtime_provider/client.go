// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package runtime_provider

import (
	"openpitrix.io/openpitrix/pkg/constants"
	"openpitrix.io/openpitrix/pkg/manager"
	"openpitrix.io/openpitrix/pkg/pb"
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
