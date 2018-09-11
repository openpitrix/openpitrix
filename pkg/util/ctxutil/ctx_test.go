// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package ctxutil

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAddMessageId(t *testing.T) {
	ctx := context.TODO()
	ctx = SetMessageId(ctx, []string{"1", "2", "3"})

	messageId := GetMessageId(ctx)
	require.Equal(t, messageId, []string{"1", "2", "3"})

	ctx = AddMessageId(ctx, "4")

	messageId = GetMessageId(ctx)
	require.Equal(t, messageId, []string{"1", "2", "3", "4"})

	ctx = ClearMessageId(ctx)

	messageId = GetMessageId(ctx)
	require.Equal(t, messageId, []string{})
}

func TestGetRequestId(t *testing.T) {
	ctx := context.TODO()
	requestId := "abcdef"
	ctx = SetRequestId(ctx, requestId)

	require.Equal(t, requestId, GetRequestId(ctx))

	ctx = context.TODO()
	requestId = "12345"
	ctx = SetRequestId(ctx, requestId)

	require.Equal(t, requestId, GetRequestId(ctx))

	ctx = context.TODO()
	requestId = "qwert"
	ctx = SetRequestId(ctx, requestId)

	require.Equal(t, requestId, GetRequestId(ctx))
}
