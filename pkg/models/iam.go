// Copyright 2019 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package models

import (
	pbim "kubesphere.io/im/pkg/pb"

	pbam "openpitrix.io/iam/pkg/pb"
	"openpitrix.io/openpitrix/pkg/pb"
	"openpitrix.io/openpitrix/pkg/util/pbutil"
)

func ToPbRole(role *pbam.Role) *pb.Role {
	if role == nil {
		return new(pb.Role)
	}
	return &pb.Role{
		RoleId:      role.RoleId,
		RoleName:    role.RoleName,
		Description: role.Description,
		Portal:      role.Portal,
		Owner:       role.Owner,
		OwnerPath:   role.OwnerPath,
		Status:      role.Status,
		Controller:  role.Controller,
		CreateTime:  role.CreateTime,
		UpdateTime:  role.UpdateTime,
		StatusTime:  role.StatusTime,
	}
}

func ToPbRoles(roles []*pbam.Role) []*pb.Role {
	var pbRoles []*pb.Role
	for _, role := range roles {
		pbRoles = append(pbRoles, ToPbRole(role))
	}
	return pbRoles
}

func ToPbGroup(group *pbim.Group) *pb.Group {
	if group == nil {
		return new(pb.Group)
	}
	return &pb.Group{
		ParentGroupId: pbutil.ToProtoString(group.ParentGroupId),
		GroupId:       pbutil.ToProtoString(group.GroupId),
		GroupPath:     pbutil.ToProtoString(group.GroupPath),
		Name:          pbutil.ToProtoString(group.GroupName),
		Description:   pbutil.ToProtoString(group.Description),
		Status:        pbutil.ToProtoString(group.Status),
		CreateTime:    group.CreateTime,
		UpdateTime:    group.UpdateTime,
		StatusTime:    group.StatusTime,
	}
}

func ToPbGroups(groups []*pbim.Group) []*pb.Group {
	var pbGroups []*pb.Group
	for _, group := range groups {
		pbGroups = append(pbGroups, ToPbGroup(group))
	}
	return pbGroups
}

func ToPbUsers(users []*pbim.User) []*pb.User {
	var pbUsers []*pb.User
	for _, user := range users {
		pbUsers = append(pbUsers, ToPbUser(user))
	}
	return pbUsers
}

func ToPbUser(user *pbim.User) *pb.User {
	if user == nil {
		return new(pb.User)
	}
	return &pb.User{
		UserId:      pbutil.ToProtoString(user.UserId),
		Username:    pbutil.ToProtoString(user.Username),
		Email:       pbutil.ToProtoString(user.Email),
		PhoneNumber: pbutil.ToProtoString(user.PhoneNumber),
		Description: pbutil.ToProtoString(user.Description),
		Status:      pbutil.ToProtoString(user.Status),
		CreateTime:  user.CreateTime,
		UpdateTime:  user.UpdateTime,
		StatusTime:  user.StatusTime,
	}
}

func ToPbModule(p *pbam.Module) *pb.Module {
	if p == nil {
		return new(pb.Module)
	}
	return &pb.Module{
		ModuleElemSet: ToPbModuleElems(p.ModuleElemSet),
	}
}

func ToAmModule(p *pb.Module) *pbam.Module {
	if p == nil {
		return new(pbam.Module)
	}
	return &pbam.Module{
		ModuleElemSet: ToAmModuleElems(p.ModuleElemSet),
	}
}

func ToPbModuleElems(ps []*pbam.ModuleElem) []*pb.ModuleElem {
	var results []*pb.ModuleElem
	for _, p := range ps {
		results = append(results, &pb.ModuleElem{
			ModuleId:   p.ModuleId,
			ModuleName: p.ModuleName,
			FeatureSet: ToPbFeatures(p.FeatureSet),
			DataLevel:  p.DataLevel,
			IsCheckAll: p.IsCheckAll,
		})
	}
	return results
}

func ToAmModuleElems(ps []*pb.ModuleElem) []*pbam.ModuleElem {
	var results []*pbam.ModuleElem
	for _, p := range ps {
		results = append(results, &pbam.ModuleElem{
			ModuleId:   p.ModuleId,
			ModuleName: p.ModuleName,
			FeatureSet: ToAmFeatures(p.FeatureSet),
			DataLevel:  p.DataLevel,
			IsCheckAll: p.IsCheckAll,
		})
	}
	return results
}

func ToPbFeatures(ps []*pbam.Feature) []*pb.Feature {
	var results []*pb.Feature
	for _, p := range ps {
		results = append(results, &pb.Feature{
			FeatureId:                p.FeatureId,
			FeatureName:              p.FeatureName,
			ActionBundleSet:          ToPbActionBundles(p.ActionBundleSet),
			CheckedActionBundleIdSet: p.CheckedActionBundleIdSet,
		})
	}
	return results
}

func ToAmFeatures(ps []*pb.Feature) []*pbam.Feature {
	var results []*pbam.Feature
	for _, p := range ps {
		results = append(results, &pbam.Feature{
			FeatureId:                p.FeatureId,
			FeatureName:              p.FeatureName,
			ActionBundleSet:          ToAmActionBundles(p.ActionBundleSet),
			CheckedActionBundleIdSet: p.CheckedActionBundleIdSet,
		})
	}
	return results
}

func ToPbActionBundles(ps []*pbam.ActionBundle) []*pb.ActionBundle {
	var results []*pb.ActionBundle
	for _, p := range ps {
		results = append(results, &pb.ActionBundle{
			ActionBundleId:   p.ActionBundleId,
			ActionBundleName: p.ActionBundleName,
			ApiSet:           ToPbApis(p.ApiSet),
		})
	}
	return results
}

func ToAmActionBundles(ps []*pb.ActionBundle) []*pbam.ActionBundle {
	var results []*pbam.ActionBundle
	for _, p := range ps {
		results = append(results, &pbam.ActionBundle{
			ActionBundleId:   p.ActionBundleId,
			ActionBundleName: p.ActionBundleName,
			ApiSet:           ToAmApis(p.ApiSet),
		})
	}
	return results
}

func ToPbApis(ps []*pbam.Api) []*pb.Api {
	var results []*pb.Api
	for _, p := range ps {
		results = append(results, &pb.Api{
			ApiId:     p.ApiId,
			ApiMethod: p.ApiMethod,
			UrlMethod: p.UrlMethod,
			Url:       p.Url,
		})
	}
	return results
}

func ToAmApis(ps []*pb.Api) []*pbam.Api {
	var results []*pbam.Api
	for _, p := range ps {
		results = append(results, &pbam.Api{
			ApiId:     p.ApiId,
			ApiMethod: p.ApiMethod,
			UrlMethod: p.UrlMethod,
			Url:       p.Url,
		})
	}
	return results
}
