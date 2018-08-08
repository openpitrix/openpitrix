// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package pilot

import (
	"crypto/tls"
	"time"

	"github.com/golang/protobuf/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	"openpitrix.io/openpitrix/pkg/manager"
	"openpitrix.io/openpitrix/pkg/pb/metadata/pilot"
	"openpitrix.io/openpitrix/pkg/pb/metadata/types"
)

type Server struct {
	cfg           *pbtypes.PilotConfig
	fgClientMgr   *FrontgateClientManager
	taskStatusMgr *TaskStatusManager
}

func Serve(cfg *pbtypes.PilotConfig, tlsCfg *tls.Config, opts ...Options) {
	if cfg != nil {
		cfg = proto.Clone(cfg).(*pbtypes.PilotConfig)
	} else {
		cfg = NewDefaultConfig()
	}

	for _, fn := range opts {
		fn(cfg)
	}

	p := &Server{
		cfg:           cfg,
		fgClientMgr:   NewFrontgateClientManager(),
		taskStatusMgr: NewTaskStatusManager(),
	}

	go func() {
		for {
			p.fgClientMgr.CheckAllClient()
			time.Sleep(time.Second * 10)
		}
	}()

	// internal service
	go manager.NewGrpcServer("pilot-service", int(p.cfg.ListenPort)).Serve(
		func(server *grpc.Server) {
			pbpilot.RegisterPilotServiceServer(server, p)
		},
	)

	// tls for public service
	if tlsCfg != nil {
		manager.NewGrpcServer("pilot-service-for-frontgate", int(p.cfg.ForFrontgateListenPort)).Serve(
			func(server *grpc.Server) {
				pbpilot.RegisterPilotServiceForFrontgateServer(server, p)
			},
			grpc.Creds(credentials.NewTLS(tlsCfg)),
		)
	} else {
		manager.NewGrpcServer("pilot-service-for-frontgate", int(p.cfg.ForFrontgateListenPort)).Serve(
			func(server *grpc.Server) {
				pbpilot.RegisterPilotServiceForFrontgateServer(server, p)
			},
		)
	}
}
