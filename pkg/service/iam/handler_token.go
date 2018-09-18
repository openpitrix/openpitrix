// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package iam

import (
	"context"

	"openpitrix.io/openpitrix/pkg/pb"
)

var (
	_ pb.TokenManagerServer = (*Server)(nil)
)

func (p *Server) CreateClient(ctx context.Context, req *pb.CreateClientRequest) (*pb.CreateClientResponse, error) {
	return nil, nil
}

func (p *Server) Auth(ctx context.Context, req *pb.AuthRequest) (*pb.AuthResponse, error) {
	return nil, nil
}

func (p *Server) Token(ctx context.Context, req *pb.TokenRequest) (*pb.TokenResponse, error) {
	return nil, nil
}
