// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package account

import (
	"context"

	pbam "openpitrix.io/iam/pkg/pb"
	"openpitrix.io/openpitrix/pkg/client/iam/am"
	"openpitrix.io/openpitrix/pkg/gerr"
	"openpitrix.io/openpitrix/pkg/logger"
	"openpitrix.io/openpitrix/pkg/models"
	"openpitrix.io/openpitrix/pkg/pb"
	"openpitrix.io/openpitrix/pkg/util/ctxutil"
	"openpitrix.io/openpitrix/pkg/util/stringutil"
)

var (
	_           pb.AccessManagerServer = (*Server)(nil)
	amClient, _                        = am.NewClient()
)

func validateRoleUser(ctx context.Context, roleId, userId string) bool {
	response, err := amClient.GetRoleWithUser(ctx, &pbam.GetRoleRequest{
		RoleId: roleId,
	})
	if err != nil {
		logger.Error(ctx, "Get role with user failed: %+v", err)
		return false
	}
	if stringutil.StringIn(userId, response.Role.UserIdSet) {
		return true
	} else {
		return false
	}
}

func (p *Server) CanDo(ctx context.Context, req *pb.CanDoRequest) (*pb.CanDoResponse, error) {
	v, err := amClient.CanDo(ctx, &pbam.CanDoRequest{
		UserId:    req.UserId,
		Url:       req.Url,
		UrlMethod: req.UrlMethod,
		ApiMethod: req.ApiMethod,
	})
	if err != nil {
		return nil, err
	}

	return &pb.CanDoResponse{
		UserId:     v.UserId,
		AccessPath: v.AccessPath,
		OwnerPath:  v.OwnerPath,
	}, nil
}

func (p *Server) GetRoleModule(ctx context.Context, req *pb.GetRoleModuleRequest) (*pb.GetRoleModuleResponse, error) {
	response, err := amClient.GetRoleModule(ctx, &pbam.GetRoleModuleRequest{
		RoleId: req.RoleId,
	})
	if err != nil {
		return nil, err
	}

	return &pb.GetRoleModuleResponse{
		RoleId: req.RoleId,
		Module: models.ToPbModule(response.Module),
	}, nil
}

func (p *Server) ModifyRoleModule(ctx context.Context, req *pb.ModifyRoleModuleRequest) (*pb.ModifyRoleModuleResponse, error) {
	_, err := amClient.ModifyRoleModule(ctx, &pbam.ModifyRoleModuleRequest{
		RoleId: req.RoleId,
		Module: models.ToAmModule(req.Module),
	})
	if err != nil {
		return nil, err
	}

	return &pb.ModifyRoleModuleResponse{
		RoleId: req.RoleId,
	}, nil
}

func (p *Server) CreateRole(ctx context.Context, req *pb.CreateRoleRequest) (*pb.CreateRoleResponse, error) {
	s := ctxutil.GetSender(ctx)
	response, err := amClient.CreateRole(ctx, &pbam.CreateRoleRequest{
		RoleName:    req.RoleName,
		Description: req.Description,
		Portal:      req.Portal,
		Owner:       s.GetOwnerPath().Owner(),
		OwnerPath:   string(s.GetOwnerPath()),
	})
	if err != nil {
		return nil, err
	}

	return &pb.CreateRoleResponse{
		RoleId: response.RoleId,
	}, nil
}

func (p *Server) DeleteRoles(ctx context.Context, req *pb.DeleteRolesRequest) (*pb.DeleteRolesResponse, error) {
	_, err := amClient.DeleteRoles(ctx, &pbam.DeleteRolesRequest{
		RoleId: req.RoleId,
	})
	if err != nil {
		return nil, err
	}

	return &pb.DeleteRolesResponse{
		RoleId: req.RoleId,
	}, nil
}

func (p *Server) ModifyRole(ctx context.Context, req *pb.ModifyRoleRequest) (*pb.ModifyRoleResponse, error) {
	_, err := amClient.ModifyRole(ctx, &pbam.ModifyRoleRequest{
		RoleId:      req.RoleId,
		RoleName:    req.RoleName,
		Description: req.Description,
	})
	if err != nil {
		return nil, err
	}

	return &pb.ModifyRoleResponse{
		RoleId: req.RoleId,
	}, nil
}

func (p *Server) GetRole(ctx context.Context, req *pb.GetRoleRequest) (*pb.GetRoleResponse, error) {
	res, err := amClient.GetRole(ctx, &pbam.GetRoleRequest{
		RoleId: req.RoleId,
	})
	if err != nil {
		return nil, err
	}

	return &pb.GetRoleResponse{
		Role: models.ToPbRole(res.Role),
	}, nil
}

func (p *Server) DescribeRoles(ctx context.Context, req *pb.DescribeRolesRequest) (*pb.DescribeRolesResponse, error) {
	response, err := amClient.DescribeRoles(ctx, &pbam.DescribeRolesRequest{
		SearchWord:     req.SearchWord,
		SortKey:        req.SortKey,
		Reverse:        req.Reverse,
		Offset:         req.Offset,
		Limit:          req.Limit,
		RoleId:         req.RoleId,
		RoleName:       req.RoleName,
		Portal:         req.Portal,
		Status:         req.Status,
		ActionBundleId: req.ActionBundleId,
	})
	if err != nil {
		return nil, err
	}

	var roles []*pb.Role
	for _, role := range response.RoleSet {
		roles = append(roles, models.ToPbRole(role))
	}

	reply := &pb.DescribeRolesResponse{
		TotalCount: response.Total,
		RoleSet:    roles,
	}

	return reply, nil
}

func (p *Server) BindUserRole(ctx context.Context, req *pb.BindUserRoleRequest) (*pb.BindUserRoleResponse, error) {
	_, err := amClient.BindUserRole(ctx, &pbam.BindUserRoleRequest{
		RoleId: req.RoleId,
		UserId: req.UserId,
	})
	if err != nil {
		return nil, err
	}

	return &pb.BindUserRoleResponse{
		UserId: req.UserId,
		RoleId: req.RoleId,
	}, nil
}

func (p *Server) UnbindUserRole(ctx context.Context, req *pb.UnbindUserRoleRequest) (*pb.UnbindUserRoleResponse, error) {
	return nil, gerr.New(ctx, gerr.PermissionDenied, gerr.ErrorPermissionDenied)
}
