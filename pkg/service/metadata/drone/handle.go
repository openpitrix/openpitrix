// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package drone

import (
	"context"

	pb_drone "openpitrix.io/openpitrix/pkg/pb/drone"
)

var _ pb_drone.DroneServiceServer = (*Server)(nil)

func (p *Server) GetInfo(context.Context, *pb_drone.Empty) (*pb_drone.Info, error) {
	panic("todo")
}
func (p *Server) GetConfdConfig(context.Context, *pb_drone.Empty) (*pb_drone.ConfdConfig, error) {
	panic("todo")
}
func (p *Server) GetBackendConfig(context.Context, *pb_drone.Empty) (*pb_drone.BackendConfig, error) {
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
