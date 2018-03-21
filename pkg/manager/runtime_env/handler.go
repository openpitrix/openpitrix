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
	"openpitrix.io/openpitrix/pkg/utils/sender"
	"openpitrix.io/openpitrix/pkg/logger"
)

func (p *Server) CreateRuntimeEnv(ctx context.Context, req *pb.CreateRuntimeEnvRequest) (*pb.CreateRuntimeEnvResponse, error) {
	s := sender.GetSenderFromContext(ctx)
	logger.Debugf("Got sender: %+v", s)
	logger.Debugf("Got req: %+v", req)

	//create runtime env
	runtimeEnvId, err := p.createRuntimeEnv(
		req.GetName().GetValue(),
		req.GetDescription().GetValue(),
		req.GetRuntimeEnvUrl().GetValue(),
		s.UserId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "CreateRuntimeEnv: %+v", err)
	}

	//create labels
	err = p.createRuntimeEnvLabels(runtimeEnvId, req.Labels.GetValue())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "CreateRuntimeEnv: %+v", err)
	}

	//get response
	pbRuntimeEnv, err := p.getRuntimeEnvPbWithLabel(runtimeEnvId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "CreateRuntimeEnv: %+v", err)
	}
	res := &pb.CreateRuntimeEnvResponse{
		RuntimeEnv: pbRuntimeEnv,
	}
	return res, nil
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
