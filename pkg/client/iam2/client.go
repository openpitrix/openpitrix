// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package iam2

import (
	"openpitrix.io/iam/pkg/pb/am"
	"openpitrix.io/iam/pkg/pb/im"
	"openpitrix.io/openpitrix/pkg/constants"
	"openpitrix.io/openpitrix/pkg/manager"
)

func NewClient() (pbim.AccountManagerClient, error) {
	conn, err := manager.NewClient(constants.IAM2ServiceHost, constants.IAM2ServicePort)
	if err != nil {
		return nil, err
	}

	client := pbim.NewAccountManagerClient(conn)
	return client, nil
}

func NewAMClient() (pbam.AccessManagerClient, error) {
	conn, err := manager.NewClient(constants.AMServiceHost, constants.AMServicePort)
	if err != nil {
		return nil, err
	}

	client := pbam.NewAccessManagerClient(conn)
	return client, nil
}
