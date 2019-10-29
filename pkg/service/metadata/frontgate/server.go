// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package frontgate

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"os"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/status"

	"openpitrix.io/openpitrix/pkg/constants"
	"openpitrix.io/openpitrix/pkg/logger"
	pbfrontgate "openpitrix.io/openpitrix/pkg/pb/metadata/frontgate"
	pbtypes "openpitrix.io/openpitrix/pkg/pb/metadata/types"
	"openpitrix.io/openpitrix/pkg/service/metadata/pilot/pilotutil"
)

type Server struct {
	cfg            *ConfigManager
	tlsPilotConfig *tls.Config
	etcd           *EtcdClientManager

	ch   *pilotutil.FrameChannel
	conn *grpc.ClientConn
	err  error
}

func Serve(cfg *ConfigManager, tlsPilotConfig *tls.Config) {
	p := &Server{
		cfg:            cfg,
		tlsPilotConfig: tlsPilotConfig,
		etcd:           NewEtcdClientManager(),
	}

	go func() {
		logger.Info(nil, "Starting file server on :%d", constants.FrontgateFileServerPort)
		err := http.ListenAndServe(fmt.Sprintf(":%d", constants.FrontgateFileServerPort), http.FileServer(http.Dir(HttpServePath)))
		if err != nil {
			logger.Critical(nil, "Start file server failed: %+v", err)
			os.Exit(1)
		}
	}()

	go ServeReverseRpcServerForPilot(cfg.Get(), tlsPilotConfig, p)
	go func() {
		err := pbfrontgate.ListenAndServeFrontgateService("tcp",
			fmt.Sprintf(":%d", cfg.Get().ListenPort),
			p,
		)
		if err != nil {
			logger.Critical(nil, "Listen and serve frontgate service failed: %+v", err)
			os.Exit(1)
		}
	}()

	<-make(chan bool)
}

func ServeReverseRpcServerForPilot(
	cfg *pbtypes.FrontgateConfig, tlsConfig *tls.Config,
	service pbfrontgate.FrontgateService,
) {
	target := fmt.Sprintf("%s:%d", cfg.PilotHost, cfg.PilotPort)
	logger.Info(nil, "Serve reverse rpc server for pilot [%s] begin.", target)
	defer logger.Info(nil, "Serve reverse rpc server for pilot [%s] finished.", target)

	var lastErrCode = codes.OK
	ctx := context.Background()
	for {
		logger.Info(nil, "Dial pilot [%s] channel begin.", target)
		ch, conn, err := pilotutil.DialPilotChannelTLS(
			ctx,
			target,
			tlsConfig,
			grpc.WithBlock(),
			grpc.WithKeepaliveParams(keepalive.ClientParameters{
				Time:                30 * time.Second,
				Timeout:             10 * time.Second,
				PermitWithoutStream: true,
			}),
		)
		if err != nil {
			logger.Error(nil, "Dial pilot [%s] channel failed: %+v", target, err)
			gerr, ok := status.FromError(err)
			if !ok {
				logger.Error(nil, "The error is not grpc error type.")
				time.Sleep(time.Second)
				continue
			}

			if gerr.Code() != lastErrCode {
				logger.Error(nil, "Failed to connect: %v", gerr.Err())
			}

			lastErrCode = gerr.Code()
			continue
		} else {
			if lastErrCode == codes.Unavailable {
				logger.Info(nil, "Pilot connect ok")
			}

			lastErrCode = codes.OK
		}
		logger.Info(nil, "Dial pilot [%s] channel finished.", target)

		updater := NewUpdater(conn, cfg)
		go updater.Serve()

		logger.Info(nil, "Serving frontgate service...")
		// will long run util err
		pbfrontgate.ServeFrontgateService(ch, service)

		// close all
		conn.Close()
		ch.Close()
		updater.Close()
		logger.Info(nil, "Serve frontgate service closed.")
	}
}
