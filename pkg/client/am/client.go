// Copyright 2019 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package am

import (
	pbam "openpitrix.io/iam/pkg/pb/am"
	"openpitrix.io/openpitrix/pkg/constants"
	"openpitrix.io/openpitrix/pkg/manager"
)

func NewClient() (pbam.AccessManagerClient, error) {
	conn, err := manager.NewClient(constants.AMServiceHost, constants.AMServicePort)
	if err != nil {
		return nil, err
	}

	client := pbam.NewAccessManagerClient(conn)
	return client, nil
}
