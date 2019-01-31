// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package account

import (
	"context"

	pbam "openpitrix.io/iam/pkg/pb/am"
	"openpitrix.io/openpitrix/pkg/client/am"
	"openpitrix.io/openpitrix/pkg/gerr"
	"openpitrix.io/openpitrix/pkg/pb"
	"openpitrix.io/openpitrix/pkg/util/ctxutil"
)

var (
	_           pb.AccessManagerServer = (*Server)(nil)
	amClient, _                        = am.NewClient()
)

func pbRole(p *pbam.Role) *pb.Role {
	return &pb.Role{
		RoleId:      p.RoleId,
		RoleName:    p.RoleName,
		Description: p.Description,
		Portal:      p.Portal,
		Owner:       p.Owner,
		OwnerPath:   p.OwnerPath,
		Status:      p.Status,

		CreateTime: p.CreateTime,
		UpdateTime: p.UpdateTime,
		StatusTime: p.StatusTime,

		UserId: p.UserId,
	}
}
func pbamRole(p *pb.Role) *pbam.Role {
	return &pbam.Role{
		RoleId:      p.RoleId,
		RoleName:    p.RoleName,
		Description: p.Description,
		Portal:      p.Portal,
		Owner:       p.Owner,
		OwnerPath:   p.OwnerPath,
		Status:      p.Status,

		CreateTime: p.CreateTime,
		UpdateTime: p.UpdateTime,
		StatusTime: p.StatusTime,

		UserId: p.UserId,
	}
}

func pbRoleModule(p *pbam.RoleModule) *pb.RoleModule {
	return &pb.RoleModule{
		RoleId: p.RoleId,
		Module: pbRoleModuleElemList(p.Module),
	}
}
func pbamRoleModule(p *pb.RoleModule) *pbam.RoleModule {
	return &pbam.RoleModule{
		RoleId: p.RoleId,
		Module: pbamRoleModuleElemList(p.Module),
	}
}

func pbRoleModuleElemList(ps []*pbam.RoleModuleElem) []*pb.RoleModuleElem {
	var results []*pb.RoleModuleElem
	for _, p := range ps {
		results = append(results, &pb.RoleModuleElem{
			ModuleId:   p.ModuleId,
			ModuleName: p.ModuleName,
			Feature:    pbModuleFeatureList(p.Feature),
			Owner:      p.Owner,
			DataLevel:  p.DataLevel,
			IsCheckAll: p.IsCheckAll,
		})
	}
	return results
}
func pbamRoleModuleElemList(ps []*pb.RoleModuleElem) []*pbam.RoleModuleElem {
	var results []*pbam.RoleModuleElem
	for _, p := range ps {
		results = append(results, &pbam.RoleModuleElem{
			ModuleId:   p.ModuleId,
			ModuleName: p.ModuleName,
			Feature:    pbamModuleFeatureList(p.Feature),
			Owner:      p.Owner,
			DataLevel:  p.DataLevel,
			IsCheckAll: p.IsCheckAll,
		})
	}
	return results
}

func pbModuleFeatureList(ps []*pbam.ModuleFeature) []*pb.ModuleFeature {
	var results []*pb.ModuleFeature
	for _, p := range ps {
		results = append(results, &pb.ModuleFeature{
			FeatureId:       p.FeatureId,
			FeatureName:     p.FeatureName,
			Action:          pbModuleFeatureActionList(p.Action),
			CheckedActionId: p.CheckedActionId,
		})
	}
	return results
}

func pbamModuleFeatureList(ps []*pb.ModuleFeature) []*pbam.ModuleFeature {
	var results []*pbam.ModuleFeature
	for _, p := range ps {
		results = append(results, &pbam.ModuleFeature{
			FeatureId:       p.FeatureId,
			FeatureName:     p.FeatureName,
			Action:          pbamModuleFeatureActionList(p.Action),
			CheckedActionId: p.CheckedActionId,
		})
	}
	return results
}

func pbModuleFeatureActionList(ps []*pbam.ModuleFeatureActionBundle) []*pb.ModuleFeatureActionBundle {
	var results []*pb.ModuleFeatureActionBundle
	for _, p := range ps {
		results = append(results, &pb.ModuleFeatureActionBundle{
			RoleId:         p.RoleId,
			RoleName:       p.RoleName,
			Portal:         p.Portal,
			ModuleId:       p.ModuleId,
			ModuleName:     p.ModuleName,
			DataLevel:      p.DataLevel,
			Owner:          p.Owner,
			FeatureId:      p.FeatureId,
			FeatureName:    p.FeatureName,
			ActionId:       p.ActionId,
			ActionName:     p.ActionName,
			ActionEnabled:  p.ActionEnabled,
			ApiId:          p.ApiId,
			ApiMethod:      p.ApiMethod,
			ApiDescription: p.ApiDescription,
			Url:            p.Url,
			UrlMethod:      p.UrlMethod,
		})
	}
	return results
}

func pbamModuleFeatureActionList(ps []*pb.ModuleFeatureActionBundle) []*pbam.ModuleFeatureActionBundle {
	var results []*pbam.ModuleFeatureActionBundle
	for _, p := range ps {
		results = append(results, &pbam.ModuleFeatureActionBundle{
			RoleId:         p.RoleId,
			RoleName:       p.RoleName,
			Portal:         p.Portal,
			ModuleId:       p.ModuleId,
			ModuleName:     p.ModuleName,
			DataLevel:      p.DataLevel,
			Owner:          p.Owner,
			FeatureId:      p.FeatureId,
			FeatureName:    p.FeatureName,
			ActionId:       p.ActionId,
			ActionName:     p.ActionName,
			ActionEnabled:  p.ActionEnabled,
			ApiId:          p.ApiId,
			ApiMethod:      p.ApiMethod,
			ApiDescription: p.ApiDescription,
			Url:            p.Url,
			UrlMethod:      p.UrlMethod,
		})
	}
	return results
}

func (p *Server) CanDo(ctx context.Context, req *pb.CanDoRequest) (*pb.CanDoResponse, error) {
	userId := ctxutil.GetSender(ctx).UserId
	v, err := amClient.CanDo(ctx, &pbam.CanDoRequest{
		UserId:    userId,
		Url:       req.Url,
		UrlMethod: req.UrlMethod,
	})
	if err != nil {
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorInternalError)
	}

	reply := &pb.CanDoResponse{
		UserId:     v.UserId,
		AccessPath: v.AccessPath,
		OwnerPath:  v.OwnerPath,
	}
	return reply, nil
}

func (p *Server) GetRoleModule(ctx context.Context, req *pb.GetRoleModuleRequest) (*pb.GetRoleModuleResponse, error) {
	// TODO: check permission

	v, err := amClient.GetRoleModule(ctx, &pbam.RoleId{
		RoleId: req.RoleId,
	})
	if err != nil {
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorInternalError)
	}

	reply := &pb.GetRoleModuleResponse{
		RoleModule: pbRoleModule(v),
	}
	return reply, nil
}
func (p *Server) ModifyRoleModule(ctx context.Context, req *pb.ModifyRoleModuleRequest) (*pb.ModifyRoleModuleResponse, error) {
	// TODO: check permission

	v, err := amClient.ModifyRoleModule(ctx, pbamRoleModule(req.RoleModule))
	if err != nil {
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorInternalError)
	}

	reply := &pb.ModifyRoleModuleResponse{
		RoleModule: pbRoleModule(v),
	}

	return reply, nil
}
func (p *Server) CreateRole(ctx context.Context, req *pb.CreateRoleRequest) (*pb.CreateRoleResponse, error) {
	// todo: cando?

	v, err := amClient.CreateRole(ctx, pbamRole(req.Role))
	if err != nil {
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorInternalError)
	}

	reply := &pb.CreateRoleResponse{
		RoleId: v.RoleId,
	}

	return reply, nil
}
func (p *Server) DeleteRoles(ctx context.Context, req *pb.DeleteRolesRequest) (*pb.DeleteRolesResponse, error) {
	// TODO: check permission

	_, err := amClient.DeleteRoles(ctx, &pbam.RoleIdList{
		RoleId: []string{req.RoleId},
	})
	if err != nil {
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorInternalError)
	}

	reply := &pb.DeleteRolesResponse{
		RoleId: req.RoleId,
	}

	return reply, nil
}
func (p *Server) ModifyRole(ctx context.Context, req *pb.ModifyRoleRequest) (*pb.ModifyRoleResponse, error) {
	// todo: cando?

	v, err := amClient.ModifyRole(ctx, pbamRole(req.Role))
	if err != nil {
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorInternalError)
	}

	reply := &pb.ModifyRoleResponse{
		Role: pbRole(v),
	}

	return reply, nil
}
func (p *Server) GetRole(ctx context.Context, req *pb.GetRoleRequest) (*pb.GetRoleResponse, error) {
	// todo: cando?

	v, err := amClient.GetRole(ctx, &pbam.RoleId{RoleId: req.RoleId})
	if err != nil {
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorInternalError)
	}

	reply := &pb.GetRoleResponse{
		Role: pbRole(v),
	}

	return reply, nil
}
func (p *Server) DescribeRoles(ctx context.Context, req *pb.DescribeRolesRequest) (*pb.DescribeRolesResponse, error) {
	// todo: cando?

	v, err := amClient.DescribeRoles(ctx, &pbam.DescribeRolesRequest{
		RoleId:   req.RoleId,
		RoleName: req.RoleName,
		Portal:   req.Portal,
		UserId:   req.UserId,
	})
	if err != nil {
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorInternalError)
	}

	var roles []*pb.Role
	for _, v := range v.Value {
		roles = append(roles, pbRole(v))
	}

	reply := &pb.DescribeRolesResponse{
		Role: roles,
	}

	return reply, nil
}

func (p *Server) BindUserRole(ctx context.Context, req *pb.BindUserRoleRequest) (*pb.BindUserRoleResponse, error) {
	// todo: cando?

	_, err := amClient.BindUserRole(ctx, &pbam.BindUserRoleRequest{
		RoleId: req.RoleId,
		UserId: req.UserId,
	})
	if err != nil {
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorInternalError)
	}

	reply := &pb.BindUserRoleResponse{}
	return reply, nil
}
func (p *Server) UnbindUserRole(ctx context.Context, req *pb.UnbindUserRoleRequest) (*pb.UnbindUserRoleResponse, error) {
	// todo: cando?

	_, err := amClient.UnbindUserRole(ctx, &pbam.UnbindUserRoleRequest{
		RoleId: req.RoleId,
		UserId: req.UserId,
	})
	if err != nil {
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorInternalError)
	}

	reply := &pb.UnbindUserRoleResponse{}
	return reply, nil
}
