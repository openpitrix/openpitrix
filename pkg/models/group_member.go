// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package models

import "time"

type GroupMember struct {
	GroupId    string
	UserId     string
	CreateTime time.Time
}

func NewGroupMember(gid, uid string) *GroupMember {
	return &GroupMember{
		GroupId:    gid,
		UserId:     uid,
		CreateTime: time.Now(),
	}
}
