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

func (p *Server) GetInfo(in *pb_frontgate.Empty, out *pb_frontgate.Info) error {
	panic("todo")
}

func (p *Server) CloseChannel(in *pb_frontgate.Empty, out *pb_frontgate.Empty) error {
	return p.conn.Close()
}

func (p *Server) StartConfd(in *pb_frontgate.Task, out *pb_frontgate.Empty) error {
	ctx := context.Background()

	conn, err := manager.NewClient(ctx,
		in.GetDirective().GetDroneIp(),
		constants.DroneServicePort,
	)
	if err != nil {
		return err
	}

	client := pb_drone.NewDroneServiceClient(conn)
	_, err = client.StartConfd(ctx, &pb_drone.Task{
		Id:     in.GetId(),
		Action: in.GetAction(),
		Target: in.GetTarget(),
		Directive: &pb_drone.TaskDirective{
			DroneIp:               in.GetDirective().GetDroneIp(),
			FrontgateId:           in.GetDirective().GetFrontgateId(),
			Command:               in.GetDirective().GetCommand(),
			CommandRetryTimes:     in.GetDirective().GetCommandRetryTimes(),
			CommandTimeoutSeconds: in.GetDirective().GetCommandTimeoutSeconds(),
		},
	})
	if err != nil {
		return err
	}

	return nil
}

func (p *Server) RegisterMetadata(in *pb_frontgate.Task, out *pb_frontgate.Empty) error {
	panic("todo")
}
func (p *Server) DeregisterMetadata(in *pb_frontgate.Task, out *pb_frontgate.Empty) error {
	panic("todo")
}

func (p *Server) RegisterCmd(in *pb_frontgate.Task, out *pb_frontgate.Empty) error {
	panic("todo")
}
func (p *Server) DeregisterCmd(in *pb_frontgate.Task, out *pb_frontgate.Empty) error {
	panic("todo")
}
