// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package models

const GroupMemberTableName = "group_member"

type GroupMember struct {
	GroupId string
	UserId  string
}

func NewGroupMember(gid, uid string) *GroupMember {
	return &GroupMember{
		GroupId: gid,
		UserId:  uid,
	}
}
