// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package models

import (
	"encoding/json"

	"openpitrix.io/openpitrix/pkg/logger"
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
	err := json.Unmarshal([]byte(data), meta)
	if err != nil {
		logger.Error("Unmarshal into meta failed: %+v", err)
	}
	return meta, err
}

func (m *Meta) ToString() (string, error) {
	result, err := json.Marshal(m)
	if err != nil {
		logger.Error("Marshal meta with frontgate id [%s] failed: %+v",
			m.FrontgateId, err)
	}
	return string(result), err
}
