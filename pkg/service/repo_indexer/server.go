// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package repo_indexer

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
	controller *EventController
}

func Serve(cfg *config.Config) {
	pi.SetGlobalPi(cfg)
	p := pi.Global()
	controller := NewEventController(p)
	s := Server{Pi: p, controller: controller}
	go controller.Serve()
	go s.Cron()
	manager.NewGrpcServer("repo-indexer", constants.RepoIndexerPort).Serve(func(server *grpc.Server) {
		pb.RegisterRepoIndexerServer(server, &s)
	})
}
