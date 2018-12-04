// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package app

import (
	"openpitrix.io/openpitrix/pkg/constants"
	"openpitrix.io/openpitrix/pkg/manager"
	"openpitrix.io/openpitrix/pkg/pb"
)

func NewIsvManagerClient() (pb.VendorManagerClient, error) {
	conn, err := manager.NewClient(constants.VendorManagerHost, constants.VendorManagerPort)
	if err != nil {
		return nil, err
	}
	return pb.NewVendorManagerClient(conn), err
}
