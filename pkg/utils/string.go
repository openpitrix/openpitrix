// Copyright 2017 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package utils

import "github.com/asaskevich/govalidator"

func CamelCaseToUnderscore(str string) string {
	return govalidator.CamelCaseToUnderscore(str)
}

func FindString(array []string, str string) int {
	for index, s := range array {
		if str == s {
			return index
		}
	}
	return -1
}
