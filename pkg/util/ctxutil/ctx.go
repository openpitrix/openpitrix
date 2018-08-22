// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package ctxutil

import (
	"context"

	"google.golang.org/grpc/metadata"
)

var messageIdKey struct{}

func GetMessageId(ctx context.Context) []string {
	if ctx == nil {
		return []string{}
	}
	m, ok := ctx.Value(messageIdKey).([]string)
	if !ok {
		return []string{}
	}
	return m
}

func SetMessageId(ctx context.Context, messageId []string) context.Context {
	if ctx == nil {
		return ctx
	}
	return context.WithValue(ctx, messageIdKey, messageId)
}

func AddMessageId(ctx context.Context, messageId ...string) context.Context {
	if ctx == nil {
		return ctx
	}
	m := GetMessageId(ctx)
	for _, mi := range messageId {
		m = append(m, mi)
	}
	return SetMessageId(ctx, m)
}

func ClearMessageId(ctx context.Context) context.Context {
	if ctx == nil {
		return ctx
	}
	return SetMessageId(ctx, []string{})
}

func Copy(src, dst context.Context) context.Context {
	if src != nil {
		dst = SetMessageId(dst, GetMessageId(src))
	}
	return dst
}

func GetRequestId(ctx context.Context) string {
	if ctx == nil {
		return ""
	}
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return ""
	}
	rid, ok := md["x-request-id"]
	if !ok || len(rid) == 0 {
		return ""
	}
	return rid[0]
}
