// Copyright 2017 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package utils

import (
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func GetUuid(l int) string {
	initByte := []byte{}
	result := []byte{}

	for i := '0'; i <= 'z'; i++ {
		switch {
		case i < '9':
			initByte = append(initByte, byte(i))
		case i >= 'A' && i <= 'Z':
			initByte = append(initByte, byte(i))
		case i >= 'a' && i <= 'z':
			initByte = append(initByte, byte(i))
		}

	}

	rand.Seed(time.Now().UnixNano())
	for i := 0; i < l; i++ {
		result = append(result, initByte[rand.Intn(len(initByte))])
	}
	return string(result)
}

func GetLowerAndNumUuid(l int) string {
	var (
		initByte []byte
		result []byte
	)

	for i := '0'; i <= 'z'; i++ {
		switch {
		case i <= '9':
			initByte = append(initByte, byte(i))
		case i >= 'a' && i <= 'z':
			initByte = append(initByte, byte(i))
		}
	}

	for i := 0; i < l; i++ {
		result = append(result, initByte[rand.Intn(len(initByte))])
	}
	return string(result)
}
