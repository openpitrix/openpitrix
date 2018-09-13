// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package models

import (
	"time"

	"openpitrix.io/openpitrix/pkg/constants"
	"openpitrix.io/openpitrix/pkg/util/idutil"
)

const UserPasswordResetTableName = "user_password_reset"

func NewUserPasswordResetId() string {
	return idutil.GetUuid("reset-id-")
}

type UserPasswordReset struct {
	ResetId string
	UserId  string

	Status     string
	CreateTime time.Time
}

var UserPasswordResetColumns = GetColumnsFromStruct(&UserPasswordReset{})

func NewUserPasswordReset(user_id string) *UserPasswordReset {
	return &UserPasswordReset{
		ResetId:    NewUserPasswordResetId(),
		UserId:     user_id,
		Status:     constants.StatusActive,
		CreateTime: time.Now(),
	}
}
