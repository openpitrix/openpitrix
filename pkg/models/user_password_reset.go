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

func NewUserPasswordResetId() string {
	return idutil.GetUuid("reset-id-")
}

type UserPasswordReset struct {
	ResetId string
	UserId  string

	Status     string
	CreateTime time.Time
}

var UserPasswordResetColumns = db.GetColumnsFromStruct(&UserPasswordReset{})

func NewUserPasswordReset(userId string) *UserPasswordReset {
	return &UserPasswordReset{
		ResetId:    NewUserPasswordResetId(),
		UserId:     userId,
		Status:     constants.StatusActive,
		CreateTime: time.Now(),
	}
}
