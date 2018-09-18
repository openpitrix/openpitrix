// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package models

import (
	"time"

	"openpitrix.io/openpitrix/pkg/constants"
	"openpitrix.io/openpitrix/pkg/db"
	"openpitrix.io/openpitrix/pkg/pb/iam"
	"openpitrix.io/openpitrix/pkg/util/idutil"
	"openpitrix.io/openpitrix/pkg/util/pbutil"
)

func NewUserClientId() string {
	return idutil.GetUuid("usr-client-")
}

type UserClient struct {
	ClientId     string
	ClientSecret string
	UserId       string
	Description  string

	Status     string
	CreateTime time.Time
	UpdateTime time.Time
	StatusTime time.Time
}

var UserClientColumns = db.GetColumnsFromStruct(&UserClient{})

func NewUserClient(user_id, client_id, client_password, description string) *UserClient {
	return &UserClient{
		UserId:       user_id,
		ClientId:     client_id,
		ClientSecret: client_password,
		Description:  description,
		Status:       constants.StatusActive,
		CreateTime:   time.Now(),
		UpdateTime:   time.Now(),
		StatusTime:   time.Now(),
	}
}

func UserClientToPb(p *UserClient) *pbiam.UserClient {
	q := new(pbiam.UserClient)
	q.UserId = pbutil.ToProtoString(p.UserId)
	q.ClientId = pbutil.ToProtoString(p.ClientId)
	q.ClientSecret = pbutil.ToProtoString(p.ClientSecret)
	q.Description = pbutil.ToProtoString(p.Description)
	q.Status = pbutil.ToProtoString(p.Status)
	q.CreateTime = pbutil.ToProtoTimestamp(p.CreateTime)
	q.UpdateTime = pbutil.ToProtoTimestamp(p.UpdateTime)
	q.StatusTime = pbutil.ToProtoTimestamp(p.StatusTime)
	return q
}

func UserClientsToPbs(p []*UserClient) (pbs []*pbiam.UserClient) {
	for _, v := range p {
		pbs = append(pbs, UserClientToPb(v))
	}
	return
}
