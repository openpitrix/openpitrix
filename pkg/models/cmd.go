// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package models

import (
	"encoding/json"

	"openpitrix.io/openpitrix/pkg/logger"
)

type Cmd struct {
	Id      string `json:"id"`
	Cmd     string `json:"cmd"`
	Timeout int    `json:"timeout"`
}

func NewCmd(data string) (*Cmd, error) {
	cmd := new(Cmd)
	err := json.Unmarshal([]byte(data), cmd)
	if err != nil {
		logger.Error("Unmarshal into cmd failed: %+v", err)
	}
	return cmd, err
}

func (c *Cmd) ToString() (string, error) {
	result, err := json.Marshal(c)
	if err != nil {
		logger.Error("Marshal cmd with cmd id [%s] failed: %+v", c.Id, err)
	}
	return string(result), err
}
