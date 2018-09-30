// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package models

import (
	"time"

	"openpitrix.io/openpitrix/pkg/constants"
	"openpitrix.io/openpitrix/pkg/db"
	"openpitrix.io/openpitrix/pkg/pb"
	"openpitrix.io/openpitrix/pkg/util/idutil"
	"openpitrix.io/openpitrix/pkg/util/pbutil"
)

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

	Status     string
	CreateTime time.Time
	UpdateTime time.Time
	StatusTime time.Time
}

var UserColumns = db.GetColumnsFromStruct(&User{})

func NewUser(username, password, email, role, description string) *User {
	return &User{
		UserId:      NewUserId(),
		Username:    username,
		Password:    password,
		Email:       email,
		Role:        role,
		Description: description,
		Status:      constants.StatusActive,
		CreateTime:  time.Now(),
		UpdateTime:  time.Now(),
		StatusTime:  time.Now(),
	}
}

func UserToPb(p *User) *pb.User {
	q := new(pb.User)
	q.UserId = pbutil.ToProtoString(p.UserId)
	q.Username = pbutil.ToProtoString(p.Username)
	q.Email = pbutil.ToProtoString(p.Email)
	q.Role = pbutil.ToProtoString(p.Role)
	q.Description = pbutil.ToProtoString(p.Description)
	q.Status = pbutil.ToProtoString(p.Status)
	q.CreateTime = pbutil.ToProtoTimestamp(p.CreateTime)
	q.UpdateTime = pbutil.ToProtoTimestamp(p.UpdateTime)
	q.StatusTime = pbutil.ToProtoTimestamp(p.StatusTime)
	return q
}

func UsersToPbs(p []*User) (pbs []*pb.User) {
	for _, v := range p {
		pbs = append(pbs, UserToPb(v))
	}
	return
}
