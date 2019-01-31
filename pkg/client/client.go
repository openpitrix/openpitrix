// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package client

import (
	"context"

	"openpitrix.io/openpitrix/pkg/sender"
	"openpitrix.io/openpitrix/pkg/util/ctxutil"
)

func SetSystemUserToContext(ctx context.Context) context.Context {
	return ctxutil.ContextWithSender(ctx, sender.GetSystemSender())
}
