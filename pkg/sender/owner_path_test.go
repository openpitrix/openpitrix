// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package sender

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestOwnerPath_CheckPermission(t *testing.T) {
	require.Equal(t, true, OwnerPath("grp-1:usr-1").CheckPermission(GetSystemSender()))
}
