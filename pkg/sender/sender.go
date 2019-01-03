// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package sender

import (
	"encoding/json"
	"fmt"

	"openpitrix.io/openpitrix/pkg/constants"
	"openpitrix.io/openpitrix/pkg/pb"
)

type Sender struct {
	UserId     string    `json:"user_id,omitempty"`
	Role       string    `json:"role,omitempty"`
	OwnerPath  OwnerPath `json:"owner_path,omitempty"`
	AccessPath OwnerPath `json:"access_path,omitempty"`
}

func GetSystemSender() *Sender {
	return &Sender{
		UserId: "system",
		Role:   constants.RoleGlobalAdmin,
	}
}

func GetSender(user *pb.User) *Sender {
	uid := user.GetUserId().GetValue()
	return &Sender{
		UserId: uid,
		Role:   user.GetRole().GetValue(),
	}
}

func (s Sender) GetOwnerPath() OwnerPath {
	if len(s.OwnerPath) > 0 {
		return s.OwnerPath
	}
	// group1.group2.group3:user1
	return OwnerPath(fmt.Sprintf(":%s", s.UserId))
}

func (s Sender) GetAccessPath() OwnerPath {
	if len(s.AccessPath) > 0 {
		return s.AccessPath
	}
	// global admin can access all data
	if s.IsGlobalAdmin() {
		return OwnerPath("")
	}
	// developer and normal user only can access data created by self
	return OwnerPath(fmt.Sprintf(":%s", s.UserId))
}

func (s *Sender) ToJson() string {
	b, err := json.Marshal(s)
	if err != nil {
		panic(err)
	}
	return string(b)
}

func (s *Sender) IsGlobalAdmin() bool {
	if s == nil {
		return false
	}
	return s.Role == constants.RoleGlobalAdmin
}

func (s *Sender) IsDeveloper() bool {
	if s == nil {
		return false
	}
	return s.Role == constants.RoleGlobalAdmin || s.Role == constants.RoleDeveloper
}

func (s *Sender) IsUser() bool {
	if s == nil {
		return false
	}
	return s.Role == constants.RoleGlobalAdmin || s.Role == constants.RoleDeveloper || s.Role == constants.RoleUser
}

func (s *Sender) IsGuest() bool {
	return true
}
