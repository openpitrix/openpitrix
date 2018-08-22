// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package models

import (
	"openpitrix.io/openpitrix/pkg/logger"
	"openpitrix.io/openpitrix/pkg/util/jsonutil"
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
	Eip              string
	VolumeId         string
	Device           string
	LoginPasswd      string
	Subnet           string
	NeedUserData     int
	UserDataValue    string
	UserdataFile     string
	Zone             string
	Status           string
	TransitionStatus string
	RuntimeId        string
	TargetJobId      string // target cloud job id
	Timeout          int    `json:"timeout"`
	Hostname         string
}

func NewInstance(data string) (*Instance, error) {
	instance := &Instance{}
	err := jsonutil.Decode([]byte(data), instance)
	if err != nil {
		logger.Error(nil, "Decode [%s] into instance failed: %+v", data, err)
	}
	return instance, err
}
