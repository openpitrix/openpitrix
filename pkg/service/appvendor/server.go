// Copyright 2017 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package appvendor

import (
	"context"

	"google.golang.org/grpc"

	"openpitrix.io/openpitrix/pkg/config"
	"openpitrix.io/openpitrix/pkg/constants"
	"openpitrix.io/openpitrix/pkg/manager"
	"openpitrix.io/openpitrix/pkg/pb"
	"openpitrix.io/openpitrix/pkg/pi"
)

type Server struct {
	vendorhandler Handler
}

func NewServer() (*Server, error) {
	server := &Server{}
	vendorhandler := Handler{}
	server.vendorhandler = vendorhandler
	return server, nil
}

func Serve(cfg *config.Config) {
	pi.SetGlobal(cfg)
	s, _ := NewServer()
	manager.NewGrpcServer("appvendor-manager", constants.VendorManagerPort).
		ShowErrorCause(cfg.Grpc.ShowErrorCause).
		WithChecker(s.Checker).
		WithBuilder(nil).
		Serve(func(server *grpc.Server) {
			pb.RegisterAppVendorManagerServer(server, s)
		})
}

//=======================================================================================

func (s *Server) DescribeVendorVerifyInfos(ctx context.Context, req *pb.DescribeVendorVerifyInfosRequest) (*pb.DescribeVendorVerifyInfosResponse, error) {
	res, err := s.vendorhandler.DescribeVendorVerifyInfos(ctx, req)
	return res, err
}

func (s *Server) GetVendorVerifyInfo(ctx context.Context, req *pb.GetVendorVerifyInfoRequest) (*pb.VendorVerifyInfo, error) {
	res, err := s.vendorhandler.GetVendorVerifyInfo(ctx, req)
	return res, err
}

func (s *Server) SubmitVendorVerifyInfo(ctx context.Context, req *pb.SubmitVendorVerifyInfoRequest) (*pb.SubmitVendorVerifyInfoResponse, error) {
	res, err := s.vendorhandler.SubmitVendorVerifyInfo(ctx, req)
	return res, err
}

func (s *Server) PassVendorVerifyInfo(ctx context.Context, req *pb.PassVendorVerifyInfoRequest) (*pb.PassVendorVerifyInfoResponse, error) {
	res, err := s.vendorhandler.PassVendorVerifyInfo(ctx, req)
	return res, err
}

func (s *Server) RejectVendorVerifyInfo(ctx context.Context, req *pb.RejectVendorVerifyInfoRequest) (*pb.RejectVendorVerifyInfoResponse, error) {
	res, err := s.vendorhandler.RejectVendorVerifyInfo(ctx, req)
	return res, err
}

func (s *Server) Checker(ctx context.Context, req interface{}) error {
	switch r := req.(type) {
	case *pb.SubmitVendorVerifyInfoRequest:
		return manager.NewChecker(ctx, r).
			Required("user_id", "company_name", "company_website", "company_profile", "authorizer_name", "authorizer_email", "authorizer_phone", "bank_name", "bank_account_name", "bank_account_number").
			Exec()
	case *pb.DescribeVendorVerifyInfosRequest:
		return manager.NewChecker(ctx, r).
			Required().
			Exec()
	case *pb.GetVendorVerifyInfoRequest:
		return manager.NewChecker(ctx, r).
			Required("user_id").
			Exec()
	case *pb.PassVendorVerifyInfoRequest:
		return manager.NewChecker(ctx, r).
			Required("user_id").
			Exec()
	case *pb.RejectVendorVerifyInfoRequest:
		return manager.NewChecker(ctx, r).
			Required("user_id").
			Exec()
	}
	return nil
}
