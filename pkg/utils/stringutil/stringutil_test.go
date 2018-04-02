// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package stringutil

import "testing"

func TestDiff(t *testing.T) {
	testCase := [][]string{
		{"foo", "bar", "hello"},
		{"foo", "bar", "world"},
	}
	result := Diff(testCase[0], testCase[1])
	if len(result) != 1 || result[0] != "hello" {
		t.Fatalf("Diff failed")
	}
}
