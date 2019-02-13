// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package job

import (
	"os"

	"google.golang.org/grpc"

	"openpitrix.io/openpitrix/pkg/config"
	"openpitrix.io/openpitrix/pkg/constants"
	"openpitrix.io/openpitrix/pkg/logger"
	"openpitrix.io/openpitrix/pkg/manager"
	"openpitrix.io/openpitrix/pkg/pb"
	"openpitrix.io/openpitrix/pkg/pi"
)

type Server struct {
	controller *Controller
}

func Serve(cfg *config.Config) {
	hostname, err := os.Hostname()
	if err != nil {
		logger.Critical(nil, "Failed to get os hostname: %+v", err)
		return
	}
	pi.SetGlobal(cfg)
	jobController := NewController(hostname)
	s := Server{controller: jobController}
	go jobController.Serve()

	manager.NewGrpcServer("job-controller", constants.JobManagerPort).
		ShowErrorCause(cfg.Grpc.ShowErrorCause).
		Serve(func(server *grpc.Server) {
			pb.RegisterJobManagerServer(server, &s)
		})
}
