// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package drone

import (
	"context"
	"fmt"

	"google.golang.org/grpc"

	"openpitrix.io/openpitrix/pkg/manager"
	pbdrone "openpitrix.io/openpitrix/pkg/pb/drone"
)

type DroneServiceClient interface {
	pbdrone.DroneServiceClient
	pbdrone.DroneServiceForFrontgateClient
}

type Server struct {
	opt   *Options
	confd *ConfdServer
}

func NewServer(opt *Options, opts ...func(opt *Options)) *Server {
	if opt == nil {
		opt = NewDefaultOptions()
	} else {
		opt = opt.Clone()
	}

	for _, fn := range opts {
		fn(opt)
	}

	p := &Server{
		opt:   opt,
		confd: NewConfdServer(),
	}

	return p
}

func Serve(opt *Options, opts ...func(opt *Options)) {
	s := NewServer(opt, opts...)

	manager.NewGrpcServer("drone-service", s.opt.Port).Serve(func(server *grpc.Server) {
		pbdrone.RegisterDroneServiceServer(server, s)
		pbdrone.RegisterDroneServiceForFrontgateServer(server, s)
	})
}

func DialDroneService(ctx context.Context, host string, port int) (
	client DroneServiceClient,
	conn *grpc.ClientConn,
	err error,
) {
	conn, err = grpc.Dial(fmt.Sprintf("%s:%d", host, port), grpc.WithInsecure())
	if err != nil {
		return
	}

	type _DroneServiceClient struct {
		pbdrone.DroneServiceClient
		pbdrone.DroneServiceForFrontgateClient
	}

	client = &_DroneServiceClient{
		DroneServiceClient:             pbdrone.NewDroneServiceClient(conn),
		DroneServiceForFrontgateClient: pbdrone.NewDroneServiceForFrontgateClient(conn),
	}
	return
}
