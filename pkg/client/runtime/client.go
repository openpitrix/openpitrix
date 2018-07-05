// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package runtime

import (
	"openpitrix.io/openpitrix/pkg/constants"
	"openpitrix.io/openpitrix/pkg/manager"
	"openpitrix.io/openpitrix/pkg/pb"
)

func NewRuntimeManagerClient() (pb.RuntimeManagerClient, error) {
	conn, err := manager.NewClient(constants.RuntimeManagerHost, constants.RuntimeManagerPort)
	if err != nil {
		return nil, err
	}
	return pb.NewRuntimeManagerClient(conn), err
}
