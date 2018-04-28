// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package pilot

import (
	"bytes"
	"encoding/json"

	"openpitrix.io/openpitrix/pkg/constants"
	"openpitrix.io/openpitrix/pkg/logger"
	"openpitrix.io/openpitrix/pkg/pb/types"
)

type Options func(opt *pbtypes.PilotConfig)

func NewDefaultConfigString() string {
	p := &pbtypes.PilotConfig{
		Id:         "openpitrix-pilot-001",
		Host:       "localhost",
		ListenPort: constants.PilotManagerPort,
		LogLevel:   logger.DebugLevel.String(),
	}

	data, err := json.MarshalIndent(p, "", "\t")
	if err != nil {
		panic(err) // unreachable
	}

	data = bytes.Replace(data, []byte("\n"), []byte("\r\n"), -1)
	return string(data)
}

func NewDefaultConfig() *pbtypes.PilotConfig {
	s := NewDefaultConfigString()

	p := new(pbtypes.PilotConfig)
	if err := json.Unmarshal([]byte(s), p); err != nil {
		panic(err) // unreachable
	}
	return p
}

func WithPilotId(id string) func(opt *pbtypes.PilotConfig) {
	return func(opt *pbtypes.PilotConfig) {
		opt.Id = id
	}
}

func WithListenPort(port int) func(opt *pbtypes.PilotConfig) {
	return func(opt *pbtypes.PilotConfig) {
		opt.ListenPort = int32(port)
	}
}
