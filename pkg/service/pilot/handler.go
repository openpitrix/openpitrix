// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package pilot

import (
	"context"
	"fmt"
	"math/rand"

	pb_empty "github.com/golang/protobuf/ptypes/empty"

	"openpitrix.io/openpitrix/pkg/logger"
	pb "openpitrix.io/openpitrix/pkg/pb"
	pb_frontgate "openpitrix.io/openpitrix/pkg/pb/frontgate"
)

var (
	_ pb.PilotServiceServer             = (*Server)(nil)
	_ pb.PilotServiceForFrontgateServer = (*Server)(nil)
)

func (p *Server) HandleSubtask(context.Context, *pb.HandleSubtaskRequest) (*pb.HandleSubtaskResponse, error) {
	panic("todo")
}

func (p *Server) GetSubtaskStatus(context.Context, *pb.GetSubtaskStatusRequest) (*pb.GetSubtaskStatusResponse, error) {
	panic("todo")
}

func (p *Server) CloseChannel(context.Context, *pb_empty.Empty) (*pb_empty.Empty, error) {
	panic("TODO")
}

func (p *Server) Channel(ch pb.PilotServiceForFrontgate_ChannelServer) error {
	logger.Debug("Pilot.Channel begin")
	defer logger.Debug("Pilot.Channel end")

	c := pb_frontgate.NewFrontgateServiceClient(
		NewFrontgateChannelFromServer(ch),
	)

	_ = c

	/*
		info, err := c.GetInfo(&pb_frontgate.Empty{})
		if err != nil {
			c.CloseChannel(&pb_frontgate.Empty{})
			logger.Debug(err)
			return err
		}

		p.putFrontgateClient(info.FrontgateId, c)
	*/

	return nil
}

func (p *Server) GetFrontgateClient(id string) (*pb_frontgate.FrontgateServiceClient, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if cs := p.clientMap[id]; len(cs) > 0 {
		return cs[rand.Intn(len(cs))], nil
	}

	return nil, fmt.Errorf("not found")
}

func (p *Server) CloseFrontgateClient(id string) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	var lastErr error
	for _, c := range p.clientMap[id] {
		_, err := c.CloseChannel(&pb_frontgate.Empty{})
		if err != nil {
			lastErr = err
		}
	}

	delete(p.clientMap, id)
	return lastErr
}

func (p *Server) GetPilotInfo(context.Context, *pb_empty.Empty) (*pb.PilotInfo, error) {
	panic("TODO")
}
func (p *Server) GetFrontgateInfo(context.Context, *pb_empty.Empty) (*pb.FrontgateInfo, error) {
	panic("TODO")
}

func (p *Server) GetDroneInfo(context.Context, *pb_empty.Empty) (*pb.DroneInfo, error) {
	panic("TODO")
}
func (p *Server) GetConfdInfo(context.Context, *pb.GetConfdInfoRequest) (*pb.ConfdInfo, error) {
	panic("TODO")
}

func (p *Server) StartConfd(context.Context, *pb.StartConfdRequest) (*pb_empty.Empty, error) {
	panic("TODO")
}
func (p *Server) StopConfd(context.Context, *pb.StopConfdRequest) (*pb_empty.Empty, error) {
	panic("TODO")
}

func (p *Server) RegisterMetadata(context.Context, *pb.RegisterMetadataRequest) (*pb_empty.Empty, error) {
	panic("TODO")
}
func (p *Server) DeregisterMetadata(context.Context, *pb.DeregisterMetadataRequest) (*pb_empty.Empty, error) {
	panic("TODO")
}

func (p *Server) RegisterCmd(context.Context, *pb.RegisterCmdRequest) (*pb_empty.Empty, error) {
	panic("TODO")
}
func (p *Server) DeregisterCmd(context.Context, *pb.DeregisterCmdRequest) (*pb_empty.Empty, error) {
	panic("TODO")
}

func (p *Server) putFrontgateClient(id string, c *pb_frontgate.FrontgateServiceClient) {
	p.mu.Lock()
	defer p.mu.Unlock()

	cs := p.clientMap[id]
	for _, t := range cs {
		if t == c {
			logger.Error("putFrontgateClient: exists")
			return
		}
	}

	p.clientMap[id] = append(cs, c)
	return
}
