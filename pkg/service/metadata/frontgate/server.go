// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package frontgate

import (
	"context"
	"crypto/tls"
	"fmt"
	"os"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/status"

	"openpitrix.io/openpitrix/pkg/logger"
	"openpitrix.io/openpitrix/pkg/pb/metadata/frontgate"
	"openpitrix.io/openpitrix/pkg/pb/metadata/types"
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

	go ServeReverseRpcServerForPilot(cfg.Get(), tlsPilotConfig, p)
	go func() {
		err := pbfrontgate.ListenAndServeFrontgateService("tcp",
			fmt.Sprintf(":%d", cfg.Get().ListenPort),
			p,
		)
		if err != nil {
			logger.Critical(nil, "%+v", err)
			os.Exit(1)
		}
	}()

	<-make(chan bool)
}

func ServeReverseRpcServerForPilot(
	cfg *pbtypes.FrontgateConfig, tlsConfig *tls.Config,
	service pbfrontgate.FrontgateService,
) {
	logger.Info(nil, "ReverseRpcServerForPilot begin")
	defer logger.Info(nil, "ReverseRpcServerForPilot end")

	var lastErrCode = codes.OK

	for {
		ch, conn, err := pilotutil.DialFrontgateChannelTLS(
			context.Background(), fmt.Sprintf("%s:%d", cfg.PilotHost, cfg.PilotPort),
			tlsConfig,
			grpc.WithBlock(),
			grpc.WithKeepaliveParams(keepalive.ClientParameters{
				Time:                30 * time.Second,
				Timeout:             10 * time.Second,
				PermitWithoutStream: true,
			}),
			grpc.WithTimeout(20*time.Second),
		)
		if err != nil {
			gerr, ok := status.FromError(err)
			if !ok {
				logger.Error(nil, "err shoule be grpc error type")
				time.Sleep(time.Second)
				continue
			}

			if gerr.Code() != lastErrCode {
				logger.Error(nil, "did not connect: %v", gerr.Err())
			}

			lastErrCode = gerr.Code()
			continue
		} else {
			if lastErrCode == codes.Unavailable {
				logger.Info(nil, "pilot connect ok")
			}

			lastErrCode = codes.OK
		}

		pbfrontgate.ServeFrontgateService(ch, service)
		conn.Close()
	}
}
