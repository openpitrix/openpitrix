// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package appvendor

import (
	"google.golang.org/grpc"

	"openpitrix.io/openpitrix/pkg/config"
	"openpitrix.io/openpitrix/pkg/constants"
	"openpitrix.io/openpitrix/pkg/manager"
	"openpitrix.io/openpitrix/pkg/pb"
	"openpitrix.io/openpitrix/pkg/pi"
)

type Server struct {
}

func NewServer() (*Server, error) {
	server := &Server{}
	return server, nil
}

func Serve(cfg *config.Config) {
	pi.SetGlobal(cfg)
	s, _ := NewServer()
	manager.NewGrpcServer("appvendor-manager", constants.VendorManagerPort).
		ShowErrorCause(cfg.Grpc.ShowErrorCause).
		WithChecker(s.Checker).
		WithBuilder(nil).
		Serve(func(server *grpc.Server) {
			pb.RegisterAppVendorManagerServer(server, s)
		})
}
