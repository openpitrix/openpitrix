// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package sender

import (
	"encoding/json"
	"fmt"
)

type Sender struct {
	UserId     string    `json:"user_id,omitempty"`
	OwnerPath  OwnerPath `json:"owner_path,omitempty"`
	AccessPath OwnerPath `json:"access_path,omitempty"`
}

func GetSystemSender() *Sender {
	return &Sender{
		UserId:     "system",
		OwnerPath:  ":system",
		AccessPath: "",
	}
}

func New(userId string, ownerPath, accessPath OwnerPath) *Sender {
	return &Sender{
		UserId:     userId,
		OwnerPath:  ownerPath,
		AccessPath: accessPath,
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
	// system can access all data
	if s.UserId == "system" {
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
