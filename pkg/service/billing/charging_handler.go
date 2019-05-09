// Copyright 2017 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package mbing

import (
	"context"

	"openpitrix.io/openpitrix/pkg/models"
	"openpitrix.io/openpitrix/pkg/pb"
)

func (s *Server) DescribeCharges(ctx context.Context, req *pb.DescribeChargesRequest) (*pb.DescribeChargesResponse, error) {
	//TODO: impl DescribeCharges
	return &pb.DescribeChargesResponse{}, nil
}

func (s *Server) DescribeRefunds(ctx context.Context, req *pb.DescribeRefundsRequest) (*pb.DescribeRefundsResponse, error) {
	//TODO: impl DescribeRefunds
	return &pb.DescribeRefundsResponse{}, nil
}

func (s *Server) CreateAccount(ctx context.Context, req *pb.CreateAccountRequest) (*pb.CreateAccountResponse, error) {
	//TODO: impl CreateAccount
	return &pb.CreateAccountResponse{}, nil
}

func (s *Server) DescribeAccounts(ctx context.Context, req *pb.DescribeAccountsRequest) (*pb.DescribeAccountsResponse, error) {
	//TODO: impl DescribeAccounts
	return &pb.DescribeAccountsResponse{}, nil
}

func (s *Server) ModifyAccount(ctx context.Context, req *pb.ModifyAccountRequest) (*pb.ModifyAccountResponse, error) {
	//TODO: impl ModifyAccount
	return &pb.ModifyAccountResponse{}, nil
}

func (s *Server) DeleteAccounts(ctx context.Context, req *pb.DeleteAccountsRequest) (*pb.DeleteAccountsResponse, error) {
	//TODO: impl DeleteAccounts
	return &pb.DeleteAccountsResponse{}, nil
}

func (s *Server) DescribeIncomes(ctx context.Context, req *pb.DescribeIncomesRequest) (*pb.DescribeIncomesResponse, error) {
	//TODO: impl DescribeIncomes
	return &pb.DescribeIncomesResponse{}, nil
}

func (s *Server) CreateRecharge(ctx context.Context, req *pb.CreateRechargeRequest) (*pb.CreateRechargeResponse, error) {
	//TODO: impl CreateRecharge
	return &pb.CreateRechargeResponse{}, nil
}

func (s *Server) DescribeRecharges(ctx context.Context, req *pb.DescribeRechargesRequest) (*pb.DescribeRechargesResponse, error) {
	//TODO: impl DescribeRecharges
	return &pb.DescribeRechargesResponse{}, nil
}

func (s *Server) CreateWithdraw(ctx context.Context, req *pb.CreateWithdrawRequest) (*pb.CreateWithdrawResponse, error) {
	//TODO: impl CreateWithdraw
	return &pb.CreateWithdrawResponse{}, nil
}

func (s *Server) DescribeWithdraws(ctx context.Context, req *pb.DescribeWithdrawsRequest) (*pb.DescribeWithdrawsResponse, error) {
	//TODO: impl DescribeWithdraws
	return &pb.DescribeWithdrawsResponse{}, nil
}

func charge(contract *models.LeasingContract) (string, error) {
	//TODO: generate Charge and set status to updating
	//      get balance of Account by contract.UserId
	//      if balance > contract.DueFee, update Balance/Status of user.Account and Income of owner.Account
	//      else return BALANCE_NOT_ENOUGH Error.
	return "Charge.id", nil
}

func refund(contract *models.LeasingContract) (string, error) {
	//TODO: generate Refund
	//Update Account.Balance
	return "refund.id", nil
}
