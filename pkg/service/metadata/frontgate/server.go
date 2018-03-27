// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package frontgate

import (
	"context"
	"fmt"

	"google.golang.org/grpc"

	"openpitrix.io/openpitrix/pkg/config"
	"openpitrix.io/openpitrix/pkg/constants"
	"openpitrix.io/openpitrix/pkg/logger"
	pb_frontgate "openpitrix.io/openpitrix/pkg/pb/frontgate"
	"openpitrix.io/openpitrix/pkg/pi"
	"openpitrix.io/openpitrix/pkg/service/pilot"
)

type Server struct {
	*pi.Pi

	ch   *pilot.FrameChannel
	conn *grpc.ClientConn
	err  error
}

func Serve(cfg *config.Config) {
	s := &Server{
		Pi: pi.NewPi(cfg),
	}

	s.ch, s.conn, s.err = pilot.DialFrontgateChannel(
		context.Background(), fmt.Sprintf("%s:%d", "pilot-manager", constants.PilotManagerPort),
		grpc.WithInsecure(),
	)
	if s.err != nil {
		logger.Errorf("did not connect: %v", s.err)
		return
	}

	pb_frontgate.ServeFrontgateService(s.ch, s)
}
