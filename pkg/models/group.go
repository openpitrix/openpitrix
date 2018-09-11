// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package models

import (
	"openpitrix.io/openpitrix/pkg/pb/iam"
	"openpitrix.io/openpitrix/pkg/util/idutil"
	"openpitrix.io/openpitrix/pkg/util/pbutil"
)

const GroupTableName = "group"

func NewGroupId() string {
	return idutil.GetUuid("grp-")
}

type Group struct {
	GroupId     string
	Name        string
	Description string

	Status string
}

var GroupColumns = GetColumnsFromStruct(&Group{})

func NewGroup(name, description string) *Group {
	return &Group{
		GroupId:     NewGroupId(),
		Name:        name,
		Description: description,
	}
}

func GroupToPb(p *Group) *pbiam.Group {
	q := new(pbiam.Group)
	q.GroupId = pbutil.ToProtoString(p.GroupId)
	q.Name = pbutil.ToProtoString(p.Name)
	q.Description = pbutil.ToProtoString(p.Description)
	return q
}

func GroupsToPbs(p []*Group) (pbs []*pbiam.Group) {
	for _, v := range p {
		pbs = append(pbs, GroupToPb(v))
	}
	return
}
