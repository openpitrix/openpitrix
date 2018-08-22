// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package models

import (
	"openpitrix.io/openpitrix/pkg/logger"
	"openpitrix.io/openpitrix/pkg/util/jsonutil"
)

type Volume struct {
	VolumeId         string
	NodeId           string
	InstanceId       string
	Device           string
	MountPoint       string
	MountOptions     string
	FileSystem       string
	Name             string
	Size             int
	Status           string
	TransitionStatus string
	Zone             string
	RuntimeId        string
	TargetJobId      string // target cloud job id
	Timeout          int    `json:"timeout"`
}

func NewVolume(data string) (*Volume, error) {
	volume := &Volume{}
	err := jsonutil.Decode([]byte(data), volume)
	if err != nil {
		logger.Error(nil, "Decode [%s] into volume failed: %+v", data, err)
	}
	return volume, err
}
