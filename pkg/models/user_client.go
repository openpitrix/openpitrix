// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package models

import (
	"time"

	"openpitrix.io/openpitrix/pkg/constants"
	"openpitrix.io/openpitrix/pkg/db"
	"openpitrix.io/openpitrix/pkg/util/idutil"
)

func NewUserClientId() string {
	return idutil.GetUuid("usrc-")
}

type UserClient struct {
	ClientId     string
	UserId       string
	ClientSecret string
	Status       string
	Description  string
	CreateTime   time.Time
}

var UserClientColumns = db.GetColumnsFromStruct(&UserClient{})

func NewUserClient(userId string) *UserClient {
	return &UserClient{
		ClientId:     NewUserClientId(),
		ClientSecret: idutil.GetSecret(),
		UserId:       userId,
		Description:  "",
		Status:       constants.StatusActive,
		CreateTime:   time.Now(),
	}
}
