// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package models

import (
	"encoding/json"

	"openpitrix.io/openpitrix/pkg/logger"
)

type Volume struct {
	VolumeId         string
	NodeId           string
	Name             string
	Size             int
	Status           string
	TransitionStatus string
	Zone             string
	RuntimeId        string
	TargetJobId      string // target cloud job id
	Timeout          int
}

func NewVolume(data string) (*Volume, error) {
	volume := &Volume{}
	err := json.Unmarshal([]byte(data), volume)
	if err != nil {
		logger.Errorf("Unmarshal into volume failed: %+v", err)
	}
	return volume, err
}

func (v *Volume) ToString() (string, error) {
	result, err := json.Marshal(v)
	if err != nil {
		logger.Errorf("Marshal volume with volume id [%s] failed: %+v",
			v.VolumeId, err)
	}
	return string(result), err
}
