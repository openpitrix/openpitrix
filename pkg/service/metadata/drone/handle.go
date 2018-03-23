// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package drone

import (
	"context"

	pb_drone "openpitrix.io/openpitrix/pkg/pb/drone"
)

var _ pb_drone.DroneServiceServer = (*Server)(nil)

func (p *Server) StartConfd(context.Context, *pb_drone.Task) (*pb_drone.Empty, error) {
	panic("todo")
}
