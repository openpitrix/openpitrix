// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package models

import (
	"openpitrix.io/openpitrix/pkg/logger"
	"openpitrix.io/openpitrix/pkg/util/jsonutil"
)

type Meta struct {
	FrontgateId string `json:"frontgate_id"`
	Timeout     int    `json:"timeout"`
	Cnodes      string `json:"cnodes"`
	DroneIp     string `json:"drone_ip"`
	NodeId      string `json:"node_id"`
	ClusterId   string `json:"cluster_id"`
}

func NewMeta(data string) (*Meta, error) {
	meta := &Meta{}
	err := jsonutil.Decode([]byte(data), meta)
	if err != nil {
		logger.Error(nil, "Decode [%s] into meta failed: %+v", data, err)
	}
	return meta, err
}
