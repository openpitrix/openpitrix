// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package frontgate

import (
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
	panic("todo")
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
