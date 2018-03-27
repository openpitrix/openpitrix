// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package drone

import (
	"google.golang.org/grpc"

	"openpitrix.io/openpitrix/pkg/config"
	"openpitrix.io/openpitrix/pkg/constants"
	"openpitrix.io/openpitrix/pkg/manager"
	pb_drone "openpitrix.io/openpitrix/pkg/pb/drone"
	"openpitrix.io/openpitrix/pkg/pi"
)

type Server struct {
	*pi.Pi
}

func Serve(cfg *config.Config) {
	s := Server{
		Pi: pi.NewPi(cfg),
	}
	manager.NewGrpcServer("drone-service", constants.DroneServicePort).Serve(func(server *grpc.Server) {
		pb_drone.RegisterDroneServiceServer(server, &s)
	})
}
