// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package drone

import (
	"context"

	"google.golang.org/grpc"

	"openpitrix.io/openpitrix/pkg/manager"
	"openpitrix.io/openpitrix/pkg/pb/drone"
	"openpitrix.io/openpitrix/pkg/pb/types"
)

type Server struct {
	cfg   *ConfigManager
	confd *ConfdServer
	fg    *FrontgateController
}

func NewServer(cfg *ConfigManager, cfgConfd *pbtypes.ConfdConfig) *Server {
	p := &Server{
		cfg:   cfg,
		confd: NewConfdServer(),
		fg:    NewFrontgateController(),
	}

	if cfgConfd != nil {
		p.SetConfdConfig(context.Background(), cfgConfd)
	}

	return p
}

func Serve(cfg *ConfigManager, cfgConfd *pbtypes.ConfdConfig) {
	s := NewServer(cfg, cfgConfd)

	manager.NewGrpcServer("drone-service", int(s.cfg.Get().ListenPort)).Serve(func(server *grpc.Server) {
		pbdrone.RegisterDroneServiceServer(server, s)
	})
}
