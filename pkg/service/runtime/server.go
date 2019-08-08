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
}

func Serve(cfg *config.Config) {
	pi.SetGlobal(cfg)
	s := Server{}
	manager.NewGrpcServer("runtime-manager", constants.RuntimeManagerPort).
		ShowErrorCause(cfg.Grpc.ShowErrorCause).
		WithChecker(s.Checker).
		WithMysqlConfig(cfg.Mysql).
		Serve(func(server *grpc.Server) {
			pb.RegisterRuntimeManagerServer(server, &s)
		})
}
