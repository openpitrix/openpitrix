// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package drone

import (
	"google.golang.org/grpc"

	"openpitrix.io/openpitrix/pkg/manager"
	pbdrone "openpitrix.io/openpitrix/pkg/pb/metadata/drone"
)

type Server struct {
	cfg   *ConfigManager
	confd *ConfdServer
	fg    *FrontgateController
}

func NewServer(cfg *ConfigManager, confd *ConfdServer) *Server {
	p := &Server{
		cfg:   cfg,
		confd: confd,
		fg:    NewFrontgateController(),
	}

	return p
}

func Serve(cfg *ConfigManager, confd *ConfdServer) {
	s := NewServer(cfg, confd)

	manager.NewGrpcServer("drone-service", int(s.cfg.Get().ListenPort)).Serve(func(server *grpc.Server) {
		pbdrone.RegisterDroneServiceServer(server, s)
	})
}
