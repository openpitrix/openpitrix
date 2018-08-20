// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package iam

import (
	"google.golang.org/grpc"

	"openpitrix.io/openpitrix/pkg/config"
	"openpitrix.io/openpitrix/pkg/constants"
	"openpitrix.io/openpitrix/pkg/manager"
	"openpitrix.io/openpitrix/pkg/pb/iam"
	"openpitrix.io/openpitrix/pkg/pi"
)

type Server struct {
	config.TokenConfig
}

func Serve(cfg *config.Config) {
	pi.SetGlobal(cfg)
	s := Server{cfg.Token}

	manager.NewGrpcServer("iam-manager", constants.IamListenPort).Serve(func(server *grpc.Server) {
		pbiam.RegisterAccountManagerServer(server, &s)
		pbiam.RegisterTokenManagerServer(server, &s)
	})
}
