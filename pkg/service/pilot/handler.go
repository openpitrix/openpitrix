// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package pilot

import (
	"context"
	"fmt"
	"math/rand"

	"openpitrix.io/openpitrix/pkg/logger"
	pb "openpitrix.io/openpitrix/pkg/pb"
	pb_frontgate "openpitrix.io/openpitrix/pkg/pb/frontgate"
)

func (p *Server) HandleSubtask(context.Context, *pb.HandleSubtaskRequest) (*pb.HandleSubtaskResponse, error) {
	panic("todo")
}

func (p *Server) GetSubtaskStatus(context.Context, *pb.GetSubtaskStatusRequest) (*pb.GetSubtaskStatusResponse, error) {
	panic("todo")
}

func (p *Server) FrontgateChannel(ch pb.PilotManager_FrontgateChannelServer) error {
	logger.Debug("Pilot.FrontgateChannel begin")
	defer logger.Debug("Pilot.FrontgateChannel end")

	c := pb_frontgate.NewFrontgateServiceClient(
		NewFrontgateChannelFromServer(ch),
	)
	info, err := c.GetInfo(&pb_frontgate.Empty{})
	if err != nil {
		c.CloseChannel(&pb_frontgate.Empty{})
		logger.Debug(err)
		return err
	}

	p.putFrontgateClient(info.FrontgateId, c)
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
