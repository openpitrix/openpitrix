// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package pilotutil

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	grpc "google.golang.org/grpc"

	"openpitrix.io/openpitrix/pkg/logger"
	"openpitrix.io/openpitrix/pkg/pb/metadata/pilot"
	"openpitrix.io/openpitrix/pkg/pb/metadata/types"
)

func MustLoadPilotConfig(path string) *pbtypes.PilotConfig {
	p, err := LoadPilotConfig(path)
	if err != nil {
		logger.Critical("%+v", err)
		os.Exit(1)
	}
	return p
}

func LoadPilotConfig(path string) (*pbtypes.PilotConfig, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	p := new(pbtypes.PilotConfig)
	if err := json.Unmarshal(data, p); err != nil {
		return nil, err
	}

	return p, nil
}

func DialPilotService(ctx context.Context, host string, port int) (
	client pbpilot.PilotServiceClient,
	conn *grpc.ClientConn,
	err error,
) {
	conn, err = grpc.Dial(fmt.Sprintf("%s:%d", host, port), grpc.WithInsecure())
	if err != nil {
		return
	}

	client = pbpilot.NewPilotServiceClient(conn)
	return
}
