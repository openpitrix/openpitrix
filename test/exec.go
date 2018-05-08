// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package test

import (
	"os/exec"
	"testing"

	"github.com/stretchr/testify/require"
)

func execCmd(t *testing.T, cmd string) string {
	t.Logf("run command [%s]", cmd)
	c := exec.Command("/bin/sh", "-c", cmd)
	output, err := c.CombinedOutput()
	require.NoError(t, err)
	return string(output)
}
