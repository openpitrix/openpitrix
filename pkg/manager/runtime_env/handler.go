// Copyright 2017 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package runtime_env

import (
	"context"
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"openpitrix.io/openpitrix/pkg/pb"
)

func (p *Server) CreateRuntimeEnv(ctx context.Context, req *pb.CreateRuntimeEnvRequest) (*pb.CreateRuntimeEnvResponse, error) {
	return nil, status.Errorf(codes.Internal, "CreateRuntimeEnv: %+v", fmt.Errorf("hello world"))
}

func (p *Server) DescribeRuntimeEnvs(ctx context.Context, req *pb.DescribeRuntimeEnvsRequest) (*pb.DescribeRuntimeEnvsResponse, error) {
	return nil, status.Errorf(codes.Internal, "DescribeRuntimeEnvs: %+v", fmt.Errorf("hello world"))
}
func (p *Server) ModifyRuntimeEnv(ctx context.Context, req *pb.ModifyRuntimeEnvRequest) (*pb.ModifyRuntimeEnvResponse, error) {
	return nil, status.Errorf(codes.Internal, "ModifyRuntimeEnv: %+v", fmt.Errorf("hello world"))
}
func (p *Server) DeleteRuntimeEnv(ctx context.Context, req *pb.DeleteRuntimeEnvRequest) (*pb.DeleteRuntimeEnvResponse, error) {
	return nil, status.Errorf(codes.Internal, "DeleteRuntimeEnv: %+v", fmt.Errorf("hello world"))
}
func (p *Server) CreateRuntimeEnvCredential(ctx context.Context, req *pb.CreateRuntimeEnvCredentialRequset) (*pb.CreateRuntimeEnvCredentialResponse, error) {
	return nil, status.Errorf(codes.Internal, "CreateRuntimeEnvCredential: %+v", fmt.Errorf("hello world"))
}
func (p *Server) DescribeRuntimeEnvCredentials(ctx context.Context, req *pb.DescribeRuntimeEnvCredentialsRequset) (*pb.DescribeRuntimeEnvCredentialsResponse, error) {
	return nil, status.Errorf(codes.Internal, "DescribeRuntimeEnvCredentials: %+v", fmt.Errorf("hello world"))
}
func (p *Server) ModifyRuntimeEnvCredential(ctx context.Context, req *pb.ModifyRuntimeEnvCredentialRequest) (*pb.ModifyRuntimeEnvCredentialResponse, error) {
	return nil, status.Errorf(codes.Internal, "ModifyRuntimeEnvCredential: %+v", fmt.Errorf("hello world"))
}
func (p *Server) DeleteRuntimeEnvCredential(ctx context.Context, req *pb.DeleteRuntimeEnvCredentialRequset) (*pb.DeleteRuntimeEnvCredentialResponse, error) {
	return nil, status.Errorf(codes.Internal, "DeleteRuntimeEnvCredential: %+v", fmt.Errorf("hello world"))
}
func (p *Server) AttachCredentialToRuntimeEnv(ctx context.Context, req *pb.AttachCredentialToRuntimeEnvRequset) (*pb.AttachCredentialToRuntimeEnvResponse, error) {
	return nil, status.Errorf(codes.Internal, "AttachCredentialToRuntimeEnv: %+v", fmt.Errorf("hello world"))
}
func (p *Server) DetachCredentialFromRuntimeEnv(ctx context.Context, req *pb.DetachCredentialFromRuntimeEnvRequset) (*pb.DetachCredentialFromRuntimeEnvResponse, error) {
	return nil, status.Errorf(codes.Internal, "DetachCredentialFromRuntimeEnv: %+v", fmt.Errorf("hello world"))
}
