// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package senderutil

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const (
	key      = "VQaKWe0MkrKtzuvOs1Hzjgf1UIzlu3sfr2zYe4vwG7TkVtsAEJddRmxswMxR56B0"
	testUser = "test_user"
	testRole = "test_role"
)

func TestValidate(t *testing.T) {
	var i = 0
	for {
		if i > 10 {
			break
		}
		i++

		jwt, err := Generate(key, 2*time.Hour, testUser, testRole)
		assert.NoError(t, err)

		sender, err := Validate(key, jwt)
		assert.NoError(t, err)

		assert.Equal(t, testUser, sender.UserId)
		assert.Equal(t, testRole, sender.Role)

		t.Log(jwt)
		time.Sleep(1 * time.Second)
	}
}
