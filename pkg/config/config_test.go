// Copyright 2019 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package config

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetConf(t *testing.T) {
	conf1 := GetConf()
	conf1.DisableGops = true
	conf2 := GetConf()
	require.False(t, conf2.DisableGops)
}
