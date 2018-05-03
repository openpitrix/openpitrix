// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package pilot

import (
	"time"

	"github.com/golang/protobuf/proto"
	"google.golang.org/grpc"

	"openpitrix.io/openpitrix/pkg/manager"
	"openpitrix.io/openpitrix/pkg/pb/pilot"
	"openpitrix.io/openpitrix/pkg/pb/types"
)

type Server struct {
	cfg           *pbtypes.PilotConfig
	fgClientMgr   *FrontgateClientManager
	taskStatusMgr *TaskStatusManager
}

func Serve(cfg *pbtypes.PilotConfig, opts ...Options) {
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

	manager.NewGrpcServer("pilot-service", int(p.cfg.ListenPort)).Serve(func(server *grpc.Server) {
		pbpilot.RegisterPilotServiceServer(server, p)
	})
}
