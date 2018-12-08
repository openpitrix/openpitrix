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
	cmd := &Cmd{}
	err := jsonutil.Decode([]byte(data), cmd)
	if err != nil {
		logger.Error(nil, "Decode [%s] into cmd failed: %+v", data, err)
	}
	return cmd, err
}

type CmdCnodes struct {
	RootPath   string
	ClusterId  string
	CmdKey     string
	InstanceId string
	Cmd        *Cmd
}

/*
{
  	clusters: {
		<cluster_id>: {
			cmd: {
				<instance_id>: <cmd>
			}
		}
	}
}
*/
func (c *CmdCnodes) Format() map[string]interface{} {
	if c.Cmd == nil {
		return map[string]interface{}{
			c.RootPath: map[string]interface{}{
				c.ClusterId: map[string]interface{}{
					c.CmdKey: map[string]interface{}{
						c.InstanceId: "",
					},
				},
			},
		}
	} else {
		return map[string]interface{}{
			c.RootPath: map[string]interface{}{
				c.ClusterId: map[string]interface{}{
					c.CmdKey: map[string]interface{}{
						c.InstanceId: c.Cmd,
					},
				},
			},
		}
	}
}

func NewCmdCnodes(data string) (*CmdCnodes, error) {
	cmdCnodes := new(CmdCnodes)
	err := jsonutil.Decode([]byte(data), cmdCnodes)
	if err != nil {
		logger.Error(nil, "Decode [%s] into cmd cnodes failed: %+v", data, err)
	}
	return cmdCnodes, err
}
