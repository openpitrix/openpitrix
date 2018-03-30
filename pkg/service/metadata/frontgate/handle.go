// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package frontgate

import (
	"context"

	"openpitrix.io/openpitrix/pkg/manager"

	"openpitrix.io/openpitrix/pkg/constants"
	pb_drone "openpitrix.io/openpitrix/pkg/pb/drone"
	pb_frontgate "openpitrix.io/openpitrix/pkg/pb/frontgate"
)

var _ pb_frontgate.FrontgateService = (*Server)(nil)

func (p *Server) CloseChannel(in *pb_frontgate.Empty, out *pb_frontgate.Empty) error {
	return p.conn.Close()
}

func (p *Server) GetConfdInfo(in *pb_frontgate.GetConfdInfoRequest, out *pb_frontgate.ConfdInfo) error {
	panic("todo")
}

func (p *Server) StartConfd(in *pb_frontgate.StartConfdRequest, out *pb_frontgate.Empty) error {
	ctx := context.Background()

	conn, err := manager.NewClient(ctx,
		in.GetDroneIp(), constants.DroneServicePort,
	)
	if err != nil {
		return err
	}

	client := pb_drone.NewDroneServiceClient(conn)
	_, err = client.StartConfd(ctx, &pb_drone.StartConfdRequest{
		ConfdConfig:   &pb_drone.ConfdConfig{},   // todo
		BackendConfig: &pb_drone.BackendConfig{}, // todo
	})
	if err != nil {
		return err
	}

	return nil
}

func (p *Server) StopConfd(in *pb_frontgate.StopConfdRequest, out *pb_frontgate.Empty) error {
	panic("todo")
}

func (p *Server) RegisterMetadata(in *pb_frontgate.RegisterMetadataRequest, out *pb_frontgate.Empty) error {
	panic("todo")
}
func (p *Server) DeregisterMetadata(in *pb_frontgate.DeregisterMetadataRequest, out *pb_frontgate.Empty) error {
	panic("todo")
}

func (p *Server) RegisterCmd(in *pb_frontgate.RegisterCmdRequest, out *pb_frontgate.Empty) error {
	panic("todo")
}
func (p *Server) DeregisterCmd(in *pb_frontgate.DeregisterCmdRequest, out *pb_frontgate.Empty) error {
	panic("todo")
}

func (p *Server) ReportSubTaskResult(in *pb_frontgate.SubTaskResult, out *pb_frontgate.Empty) error {
	panic("todo")
}
