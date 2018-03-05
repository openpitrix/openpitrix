// Copyright 2017 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package models

import (
	"time"

	"openpitrix.io/openpitrix/pkg/utils"
)

func NewAppVersionId() string {
	return utils.GetUuid("appv-")
}

type AppVersionModel struct {
	VersionId   string
	AppId       string
	Owner       string
	Name        string
	Description string
	PackageName string
	Status      string
	CreateTime  time.Time
	StatusTime  time.Time
}
