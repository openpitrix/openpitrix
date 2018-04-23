// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package client

import (
	"context"

	"openpitrix.io/openpitrix/pkg/utils/sender"
)

func GetSystemUserContext() context.Context {
	return sender.NewContext(context.Background(), sender.GetSystemUser())
}
