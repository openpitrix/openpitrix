// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package category

import (
	"openpitrix.io/openpitrix/pkg/constants"
	"openpitrix.io/openpitrix/pkg/manager"
	"openpitrix.io/openpitrix/pkg/pb"
)

func NewCategoryManagerClient() (pb.CategoryManagerClient, error) {
	conn, err := manager.NewClient(constants.CategoryManagerHost, constants.CategoryManagerPort)
	if err != nil {
		return nil, err
	}
	return pb.NewCategoryManagerClient(conn), err
}
