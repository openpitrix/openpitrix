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
	indexer *Indexer
}

func Serve(cfg *config.Config) {
	p := pi.NewPi(cfg)
	indexer := NewIndexer(p)
	s := Server{Pi: p, indexer: indexer}
	go indexer.Serve()
	manager.NewGrpcServer("repo-indexer", constants.RepoIndexerPort).Serve(func(server *grpc.Server) {
		pb.RegisterRepoIndexerServer(server, &s)
	})
}
