// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package testutil

import (
	"fmt"
	"os/exec"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func ExecCmd(t *testing.T, cmd string) string {
	t.Logf("run command [%s]", cmd)
	var err error
	var output []byte
	for i := 0; i < 10; i++ {
		c := exec.Command("/bin/sh", "-c", cmd)
		output, err = c.CombinedOutput()
		if err == nil {
			return string(output)
		}
		time.Sleep(2 * time.Second)
		t.Log(string(output))
		t.Log(err)
		t.Log("sleep 2 second...")
	}
	require.NoError(t, err)
	return string(output)
}

func NoError(t *testing.T, err error, services []string, msgAndArgs ...interface{}) {
	if err != nil {
		for _, service := range services {
			fmt.Print(ExecCmd(t, "docker-compose logs "+service))
		}
	}
	require.NoError(t, err, msgAndArgs)
}
