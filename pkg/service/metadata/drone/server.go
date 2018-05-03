// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package drone

import (
	"context"
	"sync"

	"github.com/gogo/protobuf/proto"
	"google.golang.org/grpc"

	"openpitrix.io/openpitrix/pkg/manager"
	"openpitrix.io/openpitrix/pkg/pb/drone"
	"openpitrix.io/openpitrix/pkg/pb/types"
)

type Server struct {
	mu  sync.Mutex
	cfg *pbtypes.DroneConfig

	confd *ConfdServer
	fg    *FrontgateController
}

func NewServer(cfg *pbtypes.DroneConfig, cfgConfd *pbtypes.ConfdConfig, opts ...Options) *Server {
	if cfg != nil {
		cfg = proto.Clone(cfg).(*pbtypes.DroneConfig)
	} else {
		cfg = NewDefaultConfig()
	}

	for _, fn := range opts {
		fn(cfg)
	}

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

func Serve(cfg *pbtypes.DroneConfig, cfgConfd *pbtypes.ConfdConfig, opts ...Options) {
	s := NewServer(cfg, cfgConfd, opts...)

	manager.NewGrpcServer("drone-service", int(s.cfg.ListenPort)).Serve(func(server *grpc.Server) {
		pbdrone.RegisterDroneServiceServer(server, s)
	})
}
