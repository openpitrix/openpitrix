// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package drone

import (
	"context"

	pb_drone "openpitrix.io/openpitrix/pkg/pb/drone"
)

var (
	_ pb_drone.DroneServiceServer             = (*Server)(nil)
	_ pb_drone.DroneServiceForFrontgateServer = (*Server)(nil)
)

func (p *Server) GetInfo(context.Context, *pb_drone.Empty) (*pb_drone.Info, error) {
	panic("todo")
}
func (p *Server) GetConfdConfig(context.Context, *pb_drone.Empty) (*pb_drone.ConfdConfig, error) {
	panic("todo")
}
func (p *Server) GetBackendConfig(context.Context, *pb_drone.Empty) (*pb_drone.ConfdBackendConfig, error) {
	panic("todo")
}
func (p *Server) StartConfd(context.Context, *pb_drone.StartConfdRequest) (*pb_drone.Empty, error) {
	panic("todo")
}
func (p *Server) StopConfd(context.Context, *pb_drone.Empty) (*pb_drone.Empty, error) {
	panic("todo")
}
func (p *Server) GetConfdStatus(context.Context, *pb_drone.Empty) (*pb_drone.ConfdStatus, error) {
	panic("todo")
}

func (p *Server) GetTemplateFiles(context.Context, *pb_drone.GetTemplateFilesRequest) (*pb_drone.GetTemplateFilesResponse, error) {
	panic("todo")
}
func (p *Server) GetValues(context.Context, *pb_drone.GetValuesRequest) (*pb_drone.GetValuesResponse, error) {
	panic("todo")
}

func (p *Server) SubscribeCmdStatus(*pb_drone.SubscribeCmdStatusRequest, pb_drone.DroneServiceForFrontgate_SubscribeCmdStatusServer) error {
	panic("todo")
}

func (p *Server) GetRegisterCmdStatus(context.Context, *pb_drone.GetRegisterCmdStatusRequest) (*pb_drone.GetRegisterCmdStatusResponse, error) {
	panic("todo")
}
