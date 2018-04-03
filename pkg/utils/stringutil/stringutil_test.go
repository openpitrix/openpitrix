// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package stringutil

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDiff(t *testing.T) {
	for _, test := range []struct {
		Base   []string
		Target []string
		Output []string
	}{
		{
			Base:   []string{"foo", "bar", "hello"},
			Target: []string{"foo", "bar", "world"},
			Output: []string{"hello"},
		},
		{
			Base:   []string{"a", "b", "c"},
			Target: []string{""},
			Output: []string{"a", "b", "c"},
		},
		{
			Base:   []string{""},
			Target: []string{"foo", "bar", "world"},
			Output: []string{""},
		},
	} {
		assert.Equal(t, test.Output, Diff(test.Base, test.Target))
	}
}

func TestUnique(t *testing.T) {
	for _, test := range []struct {
		Input  []string
		Output []string
	}{
		{
			Input:  []string{""},
			Output: []string{""},
		},
		{
			Input:  []string{"a", "b", "c", "d"},
			Output: []string{"a", "b", "c", "d"},
		},
		{
			Input:  []string{"a", "b", "c", "a"},
			Output: []string{"a", "b", "c"},
		},
		{
			Input:  []string{"a", "a", "a", "a"},
			Output: []string{"a"},
		},
	} {
		assert.Equal(t, test.Output, Unique(test.Input))
	}
}
