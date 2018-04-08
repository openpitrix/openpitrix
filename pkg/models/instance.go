// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package models

import (
	"encoding/json"

	"openpitrix.io/openpitrix/pkg/logger"
)

type Instance struct {
	InstanceId       string
	NodeId           string
	Name             string
	ImageId          string
	Cpu              int
	Memory           int
	Gpu              int
	PrivateIp        string
	VolumeId         string
	Device           string
	LoginPasswd      string
	Subnet           string
	UserDataValue    string
	UserdataPath     string
	Zone             string
	Status           string
	TransitionStatus string
	RuntimeId        string
	TargetJobId      string // target cloud job id
	Timeout          int
}

func NewInstance(data string) (*Instance, error) {
	instance := &Instance{}
	err := json.Unmarshal([]byte(data), instance)
	if err != nil {
		logger.Errorf("Unmarshal into instance failed: %+v", err)
	}
	return instance, err
}

func (i *Instance) ToString() (string, error) {
	result, err := json.Marshal(i)
	if err != nil {
		logger.Errorf("Marshal instance with instance id [%s] failed: %+v",
			i.InstanceId, err)
	}
	return string(result), err
}
