// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package models

import (
	"openpitrix.io/openpitrix/pkg/logger"
	"openpitrix.io/openpitrix/pkg/util/jsonutil"
)

type Cmd struct {
	Id      string `json:"id"`
	Cmd     string `json:"cmd"`
	Timeout int    `json:"timeout"`
}

func NewCmd(data string) (*Cmd, error) {
	cmd := new(Cmd)
	err := jsonutil.Decode([]byte(data), cmd)
	if err != nil {
		logger.Error("Decode [%s] into cmd failed: %+v", data, err)
	}
	return cmd, err
}
