// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package iam

import (
	"context"
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
	config.IAMConfig
}

func initIAMClient() {
	clientId := os.Getenv("IAM_INIT_CLIENT_ID")
	clientSecret := os.Getenv("IAM_INIT_CLIENT_SECRET")
	const userId = "system"

	if clientId == "" || clientSecret == "" {
		return
	}
	_, err := pi.Global().DB(context.Background()).InsertBySql(
		`insert into user_client (client_id, user_id, client_secret, status, description)
values (?, ?, ?, 'active', '')
on duplicate key update user_id = ?, client_secret = ?, status = 'active';`,
		clientId, userId, clientSecret,
		userId, clientSecret,
	).Exec()
	if err != nil {
		logger.Error(nil, "Init default IAM client [%s] [%s] failed", clientId, clientSecret)
		panic(err)
	}
	logger.Info(nil, "Init IAM client [%s] done", clientId)
}

func Serve(cfg *config.Config) {
	pi.SetGlobal(cfg)

	go initIAMClient()

	s := Server{cfg.IAM}

	manager.NewGrpcServer("iam-service", constants.IAMServicePort).
		ShowErrorCause(cfg.Grpc.ShowErrorCause).
		WithChecker(s.Checker).
		WithBuilder(s.Builder).
		Serve(func(server *grpc.Server) {
			pb.RegisterAccountManagerServer(server, &s)
			pb.RegisterTokenManagerServer(server, &s)
		})
}
