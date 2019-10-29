// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package drone

import (
	"bytes"
	"encoding/json"

	"openpitrix.io/openpitrix/pkg/constants"
	"openpitrix.io/openpitrix/pkg/logger"
	pbtypes "openpitrix.io/openpitrix/pkg/pb/metadata/types"
)

type Options func(opt *pbtypes.DroneConfig)

func NewDefaultConfigString() string {
	p := &pbtypes.DroneConfig{
		Id:             "drone-001",
		Host:           "localhost",
		ListenPort:     constants.DroneServicePort,
		CmdInfoLogPath: "/opt/openpitrix/log/cmd.log",
		ConfdSelfHost:  "127.0.0.1",
		LogLevel:       logger.DebugLevel.String(),
	}

	data, err := json.MarshalIndent(p, "", "\t")
	if err != nil {
		panic(err) // unreachable
	}

	data = bytes.Replace(data, []byte("\n"), []byte("\r\n"), -1)
	return string(data)
}

func NewDefaultConfig() *pbtypes.DroneConfig {
	s := NewDefaultConfigString()

	p := new(pbtypes.DroneConfig)
	if err := json.Unmarshal([]byte(s), p); err != nil {
		panic(err) // unreachable
	}
	return p
}

func WithDrondId(id string) func(opt *pbtypes.DroneConfig) {
	return func(opt *pbtypes.DroneConfig) {
		opt.Id = id
	}
}

func WithListenPort(port int) func(opt *pbtypes.DroneConfig) {
	return func(opt *pbtypes.DroneConfig) {
		opt.ListenPort = int32(port)
	}
}

func WithCmdInfoLogPath(path string) func(opt *pbtypes.DroneConfig) {
	return func(opt *pbtypes.DroneConfig) {
		opt.CmdInfoLogPath = path
	}
}
