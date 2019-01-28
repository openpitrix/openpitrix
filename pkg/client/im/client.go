// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package im

import (
	pbim "openpitrix.io/iam/pkg/pb/im"
	"openpitrix.io/openpitrix/pkg/constants"
	"openpitrix.io/openpitrix/pkg/manager"
)

func NewClient() (pbim.AccountManagerClient, error) {
	conn, err := manager.NewClient(constants.IMServiceHost, constants.IMServicePort)
	if err != nil {
		return nil, err
	}

	client := pbim.NewAccountManagerClient(conn)
	return client, nil
}
