// Copyright 2017 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package app

import (
	"google.golang.org/grpc"

	"openpitrix.io/openpitrix/pkg/config"
	"openpitrix.io/openpitrix/pkg/constants"
	"openpitrix.io/openpitrix/pkg/manager"
	"openpitrix.io/openpitrix/pkg/pb"
	"openpitrix.io/openpitrix/pkg/pi"
)

type Server struct {
}

func Serve(cfg *config.Config) {
	pi.GlobalContext = pi.NewPi(cfg)
	s := new(Server)
	manager.NewGrpcServer("app-manager", constants.AppManagerPort).Serve(func(server *grpc.Server) {
		pb.RegisterAppManagerServer(server, s)
	})
}
