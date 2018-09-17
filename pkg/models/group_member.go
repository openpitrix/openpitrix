// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package models

import (
	"openpitrix.io/openpitrix/pkg/pb/iam"
	"openpitrix.io/openpitrix/pkg/util/idutil"
	"openpitrix.io/openpitrix/pkg/util/pbutil"
)

const GroupMemberTableName = "group_member"

func NewGroupMemberId() string {
	return idutil.GetUuid("grp-mem-")
}

type GroupMember struct {
	GroupId string
	UserId  string
}

var GroupMemberColumns = GetColumnsFromStruct(&GroupMember{})

func NewGroupMember(gid, uid string) *GroupMember {
	return &GroupMember{
		GroupId: gid,
		UserId:  uid,
	}
}

func GroupMemberToPb(p *GroupMember) *pbiam.GroupMember {
	q := new(pbiam.GroupMember)
	q.GroupId = pbutil.ToProtoString(p.GroupId)
	q.UserId = pbutil.ToProtoString(p.UserId)
	return q
}

func GroupMembersToPbs(p []*GroupMember) (pbs []*pbiam.GroupMember) {
	for _, v := range p {
		pbs = append(pbs, GroupMemberToPb(v))
	}
	return
}
