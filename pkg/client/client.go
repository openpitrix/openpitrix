// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package client

import (
	"context"

	"openpitrix.io/openpitrix/pkg/util/senderutil"
)

func GetSystemUserContext() context.Context {
	return senderutil.NewContext(context.Background(), senderutil.GetSystemUser())
}
