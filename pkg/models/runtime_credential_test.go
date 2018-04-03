// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRuntimeCredentialContentConvert(t *testing.T) {
	testMaps := []map[string]string{
		{
			"11": "22",
			"33": "44",
		},
		{
			"test":       "aa",
			"openpitrix": "bb",
		},
		{
			"test": "11",
			"11":   "test",
		},
	}
	for n, testMap := range testMaps {
		stringContent := RuntimeCredentialContentMapToString(testMap)
		mapContent := RuntimeCredentialContentStringToMap(stringContent)
		assert.Equal(t, testMaps[n], mapContent)
	}

}
