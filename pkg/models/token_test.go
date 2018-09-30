// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package models

import (
	"testing"
)

func TestNewToken(t *testing.T) {
	var i = 0
	for {
		if i == 100 {
			break
		}
		i++
		t.Log(NewToken("", "", ""))
	}

}
