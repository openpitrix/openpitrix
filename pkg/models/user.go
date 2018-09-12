// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package models

import (
	"openpitrix.io/openpitrix/pkg/pb/iam"
	"openpitrix.io/openpitrix/pkg/util/idutil"
	"openpitrix.io/openpitrix/pkg/util/pbutil"
)

const UserTableName = "user"

func NewUserId() string {
	return idutil.GetUuid("usr-")
}

type User struct {
	UserId      string
	Username    string
	Password    string
	Email       string
	Role        string
	Description string

	Status string
}

var UserColumns = GetColumnsFromStruct(&User{})

func NewUser(username, password, email, role, description string) *User {
	return &User{
		UserId:      NewUserId(),
		Username:    username,
		Password:    password,
		Email:       email,
		Role:        role,
		Description: description,
	}
}

func UserToPb(p *User) *pbiam.User {
	q := new(pbiam.User)
	q.UserId = pbutil.ToProtoString(p.UserId)
	q.Username = pbutil.ToProtoString(p.Username)
	q.Email = pbutil.ToProtoString(p.Email)
	q.Role = pbutil.ToProtoString(p.Role)
	q.Description = pbutil.ToProtoString(p.Description)
	return q
}

func UsersToPbs(p []*User) (pbs []*pbiam.User) {
	for _, v := range p {
		pbs = append(pbs, UserToPb(v))
	}
	return
}
