// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package pilot

import (
	"sync"

	"google.golang.org/grpc"

	"openpitrix.io/openpitrix/pkg/config"
	"openpitrix.io/openpitrix/pkg/constants"
	"openpitrix.io/openpitrix/pkg/manager"
	"openpitrix.io/openpitrix/pkg/pb"
	pb_frontgate "openpitrix.io/openpitrix/pkg/pb/frontgate"
	"openpitrix.io/openpitrix/pkg/pi"
)

type Server struct {
	*pi.Pi

	clientMap map[string][]*pb_frontgate.FrontgateServiceClient
	mu        sync.Mutex
}

func Serve(cfg *config.Config) {
	s := Server{
		Pi:        pi.NewPi(cfg),
		clientMap: make(map[string][]*pb_frontgate.FrontgateServiceClient),
	}
	manager.NewGrpcServer("pilot-manager", constants.PilotManagerPort).Serve(func(server *grpc.Server) {
		pb.RegisterPilotManagerServer(server, &s)
	})
}
