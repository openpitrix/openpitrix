// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package cluster

import (
	"google.golang.org/grpc"

	"openpitrix.io/openpitrix/pkg/pi"

	"openpitrix.io/openpitrix/pkg/config"
	"openpitrix.io/openpitrix/pkg/constants"
	"openpitrix.io/openpitrix/pkg/manager"
	"openpitrix.io/openpitrix/pkg/pb"
)

type Server struct {
}

func Serve(cfg *config.Config) {
	pi.SetGlobal(cfg)
	s := Server{}
	manager.NewGrpcServer("cluster-manager", constants.ClusterManagerPort).
		ShowErrorCause(cfg.Grpc.ShowErrorCause).
		WithChecker(s.Checker).
		Serve(func(server *grpc.Server) {
			pb.RegisterClusterManagerServer(server, &s)
		})
}
