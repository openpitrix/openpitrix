// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package access

import (
	"openpitrix.io/openpitrix/pkg/constants"
	"openpitrix.io/openpitrix/pkg/manager"
	"openpitrix.io/openpitrix/pkg/pb"
)

type Client struct {
	pb.AccessManagerClient
}

func NewClient() (*Client, error) {
	conn, err := manager.NewClient(constants.AccountServiceHost, constants.AccountServicePort)
	if err != nil {
		return nil, err
	}
	return &Client{
		AccessManagerClient: pb.NewAccessManagerClient(conn),
	}, nil
}
