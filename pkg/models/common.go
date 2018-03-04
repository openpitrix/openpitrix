// Copyright 2017 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package models

import (
	"github.com/fatih/structs"

	"openpitrix.io/openpitrix/pkg/utils"
)

func GetColumnsFromStruct(s interface{}) []string {
	names := structs.Names(s)
	for i, name := range names {
		names[i] = utils.CamelCaseToUnderscore(name)
	}
	return names
}
