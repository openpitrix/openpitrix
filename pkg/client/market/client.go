// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file

package market

import (
	"openpitrix.io/openpitrix/pkg/constants"
	"openpitrix.io/openpitrix/pkg/manager"
	"openpitrix.io/openpitrix/pkg/pb"
)

func NewMarketManagerClient() (pb.MarketManagerClient, error) {
	conn, err := manager.NewClient(constants.MarketManagerHost, constants.MarketManagerPort)
	if err != nil {
		return nil, err
	}
	return pb.NewMarketManagerClient(conn), err
}
