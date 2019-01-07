// Copyright 2017 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package senderutil

import (
	"context"
	"encoding/json"

	"google.golang.org/grpc/metadata"

	"openpitrix.io/openpitrix/pkg/constants"
	"openpitrix.io/openpitrix/pkg/pb"
	"openpitrix.io/openpitrix/pkg/util/ctxutil"
	"openpitrix.io/openpitrix/pkg/util/jsonutil"
)

const (
	SenderKey = "sender"
	TokenType = "Bearer"
)

type Sender struct {
	UserId string `json:"user_id,omitempty"`
	Role   string `json:"role,omitempty"`
}

func GetSystemSender() *Sender {
	return &Sender{UserId: "system", Role: constants.RoleGlobalAdmin}
}

func GetSender(user *pb.User) *Sender {
	return &Sender{UserId: user.GetUserId().GetValue(), Role: user.GetRole().GetValue()}
}

func (s *Sender) ToJson() string {
	return jsonutil.ToString(s)
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

func GetSenderFromContext(ctx context.Context) *Sender {
	values := ctxutil.GetValueFromContext(ctx, SenderKey)
	if len(values) == 0 || len(values[0]) == 0 {
		return nil
	}
	sender := Sender{}
	err := json.Unmarshal([]byte(values[0]), &sender)
	if err != nil {
		panic(err)
	}
	return &sender
}

func ContextWithSender(ctx context.Context, user *Sender) context.Context {
	if user == nil {
		return ctx
	}
	ctx = context.WithValue(ctx, SenderKey, []string{user.ToJson()})
	md, ok := metadata.FromOutgoingContext(ctx)
	if !ok {
		md = metadata.MD{}
	}
	md[SenderKey] = []string{user.ToJson()}
	return metadata.NewOutgoingContext(ctx, md)
}
