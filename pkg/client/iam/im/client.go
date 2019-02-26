// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package im

import (
	pbim "kubesphere.io/im/pkg/pb"

	"openpitrix.io/openpitrix/pkg/constants"
	"openpitrix.io/openpitrix/pkg/manager"
)

func NewClient() (pbim.IdentityManagerClient, error) {
	conn, err := manager.NewClient(constants.IMServiceHost, constants.IMServicePort)
	if err != nil {
		return nil, err
	}

	client := pbim.NewIdentityManagerClient(conn)
	return client, nil
}
