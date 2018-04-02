// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package runtime

import (
	"google.golang.org/grpc"

	"openpitrix.io/openpitrix/pkg/config"
	"openpitrix.io/openpitrix/pkg/constants"
	"openpitrix.io/openpitrix/pkg/manager"
	"openpitrix.io/openpitrix/pkg/pb"
	"openpitrix.io/openpitrix/pkg/pi"
)

type Server struct {
	*pi.Pi
}

func Serve(cfg *config.Config) {
	s := Server{pi.NewPi(cfg)}
	manager.NewGrpcServer("runtime-manager", constants.RuntimeEnvManagerPort).Serve(func(server *grpc.Server) {
		pb.RegisterRuntimeEnvManagerServer(server, &s)
	})
}
