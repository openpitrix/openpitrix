// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package models

import (
	"encoding/json"

	"openpitrix.io/openpitrix/pkg/logger"
)

type Meta struct {
	FrontgateId string
	Timeout     int
	Cnodes      string
	DroneIp     string
	NodeId      string
	ClusterId   string
}

func NewMeta(data string) (*Meta, error) {
	meta := &Meta{}
	err := json.Unmarshal([]byte(data), meta)
	if err != nil {
		logger.Errorf("Unmarshal into meta failed: %+v", err)
	}
	return meta, err
}

func (m *Meta) ToString() (string, error) {
	result, err := json.Marshal(m)
	if err != nil {
		logger.Errorf("Marshal meta with frontgate id [%s] failed: %+v",
			m.FrontgateId, err)
	}
	return string(result), err
}
