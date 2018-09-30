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

func NewGroupId() string {
	return idutil.GetUuid("grp-")
}

type Group struct {
	GroupId     string
	Name        string
	Description string

	Status     string
	CreateTime time.Time
	UpdateTime time.Time
	StatusTime time.Time
}

var GroupColumns = db.GetColumnsFromStruct(&Group{})

func NewGroup(name, description string) *Group {
	return &Group{
		GroupId:     NewGroupId(),
		Name:        name,
		Description: description,
		Status:      constants.StatusActive,
		CreateTime:  time.Now(),
		UpdateTime:  time.Now(),
		StatusTime:  time.Now(),
	}
}

func GroupToPb(p *Group) *pb.Group {
	q := new(pb.Group)
	q.GroupId = pbutil.ToProtoString(p.GroupId)
	q.Name = pbutil.ToProtoString(p.Name)
	q.Description = pbutil.ToProtoString(p.Description)
	q.Status = pbutil.ToProtoString(p.Status)
	q.CreateTime = pbutil.ToProtoTimestamp(p.CreateTime)
	q.UpdateTime = pbutil.ToProtoTimestamp(p.UpdateTime)
	q.StatusTime = pbutil.ToProtoTimestamp(p.StatusTime)
	return q
}

func GroupsToPbs(p []*Group) (pbs []*pb.Group) {
	for _, v := range p {
		pbs = append(pbs, GroupToPb(v))
	}
	return
}
