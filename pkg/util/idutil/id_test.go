// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package idutil

import (
	"fmt"
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetUuid(t *testing.T) {
	fmt.Println(GetUuid(""))
}

func TestGetUuid36(t *testing.T) {
	fmt.Println(GetUuid36(""))
}

func TestGetManyUuid(t *testing.T) {
	var strSlice []string
	for i := 0; i < 10000; i++ {
		testId := GetUuid("")
		strSlice = append(strSlice, testId)
	}
	sort.Strings(strSlice)
}

func TestRandString(t *testing.T) {
	str := randString(Alphabet62, 50)
	assert.Equal(t, 50, len(str))
	t.Log(str)

	str = randString(Alphabet62, 255)
	assert.Equal(t, 255, len(str))
	t.Log(str)
}
