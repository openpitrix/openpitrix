// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package models

import (
	"fmt"

	"openpitrix.io/openpitrix/pkg/logger"
	"openpitrix.io/openpitrix/pkg/util/jsonutil"
)

type Cmd struct {
	Id      string `json:"id"`
	Cmd     string `json:"cmd"`
	Timeout int    `json:"timeout"`
}

type CmdCnodes map[string]map[string]*Cmd

func (c CmdCnodes) Format(ip, taskId string) error {
	cnodes, exist := c[""]
	if exist {
		c[ip] = cnodes
		delete(c, "")
	}
	cnodes, exist = c[ip]
	if exist {
		cmd := cnodes["cmd"]
		cmd.Id = taskId
	} else {
		logger.Error("Ip [%s] not in Cnodes [%+v]", ip, c)
		return fmt.Errorf("ip [%s] not in Cnodes [%+v]", ip, c)
	}
	return nil
}

func NewCmdCnodes(data string) (*CmdCnodes, error) {
	cmdCnodes := new(CmdCnodes)
	err := jsonutil.Decode([]byte(data), cmdCnodes)
	if err != nil {
		logger.Error("Decode [%s] into cmd failed: %+v", data, err)
	}
	return cmdCnodes, err
}
