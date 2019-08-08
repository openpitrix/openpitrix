// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package repo_indexer

import (
	"context"

	"google.golang.org/grpc"

	"openpitrix.io/openpitrix/pkg/config"
	"openpitrix.io/openpitrix/pkg/constants"
	"openpitrix.io/openpitrix/pkg/db"
	"openpitrix.io/openpitrix/pkg/manager"
	"openpitrix.io/openpitrix/pkg/pb"
	"openpitrix.io/openpitrix/pkg/pi"
)

type Server struct {
	controller *EventController
}

func Serve(cfg *config.Config) {
	pi.SetGlobal(cfg)
	controller := NewEventController(db.NewContext(context.Background(), cfg.Mysql))
	s := Server{controller: controller}
	go controller.Serve()
	go s.Cron()
	manager.NewGrpcServer("repo-indexer", constants.RepoIndexerPort).
		ShowErrorCause(cfg.Grpc.ShowErrorCause).
		WithChecker(s.Checker).
		WithMysqlConfig(cfg.Mysql).
		Serve(func(server *grpc.Server) {
			pb.RegisterRepoIndexerServer(server, &s)
		})
}
