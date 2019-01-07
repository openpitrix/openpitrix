// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package senderutil

import (
	"context"
	"testing"

	"google.golang.org/grpc/metadata"
)

func TestGetSenderFromContext(t *testing.T) {
	s := GetSenderFromContext(context.Background())
	if s != nil {
		t.Fatalf("GetSenderFromContext(context.Background()) should be nil")
	}
	t.Logf("GetSenderFromContext(context.Background()) passed")

	user1 := &Sender{UserId: "user1"}
	md := metadata.MD{}
	md["sender"] = []string{user1.ToJson()}
	user1ctx := metadata.NewIncomingContext(context.Background(), md)
	s = GetSenderFromContext(user1ctx)
	if s == nil {
		t.Fatalf("GetSenderFromContext(user1ctx) should not be nil")
	}
	if s.UserId != user1.UserId {
		t.Fatalf("GetSenderFromContext(user1ctx) should be user1")
	}
	t.Logf("GetSenderFromContext(user1ctx) passed")
}
