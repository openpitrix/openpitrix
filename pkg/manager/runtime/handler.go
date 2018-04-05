// Copyright 2017 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package runtime

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"openpitrix.io/openpitrix/pkg/db"
	"openpitrix.io/openpitrix/pkg/logger"
	"openpitrix.io/openpitrix/pkg/manager"
	"openpitrix.io/openpitrix/pkg/models"
	"openpitrix.io/openpitrix/pkg/pb"
	"openpitrix.io/openpitrix/pkg/utils"
	"openpitrix.io/openpitrix/pkg/utils/sender"
)

func (p *Server) CreateRuntime(ctx context.Context, req *pb.CreateRuntimeRequest) (*pb.CreateRuntimeResponse, error) {
	s := sender.GetSenderFromContext(ctx)
	// validate req
	err := validateCreateRuntimeRequest(req)
	if err != nil {
		logger.Errorf("CreateRuntime: %+v", err)
		return nil, status.Errorf(codes.InvalidArgument, "CreateRuntime: %+v", err)
	}

	// create runtime credential
	runtimeCredentialId, err := p.createRuntimeCredential(req.Provider.GetValue(), req.RuntimeCredential.GetValue())
	if err != nil {
		logger.Errorf("CreateRuntime: %+v", err)
		return nil, status.Errorf(codes.Internal, "CreateRuntime: %+v", err)
	}

	// create runtime
	runtimeId, err := p.createRuntime(
		req.GetName().GetValue(),
		req.GetDescription().GetValue(),
		req.Provider.GetValue(),
		req.GetRuntimeUrl().GetValue(),
		runtimeCredentialId,
		req.Zone.GetValue(),
		s.UserId)
	if err != nil {
		logger.Errorf("CreateRuntime: %+v", err)
		return nil, status.Errorf(codes.Internal, "CreateRuntime: %+v", err)
	}

	// create labels
	err = p.createRuntimeLabels(runtimeId, req.Labels.GetValue())
	if err != nil {
		logger.Errorf("CreateRuntime: %+v", err)
		return nil, status.Errorf(codes.Internal, "CreateRuntime: %+v", err)
	}

	// get response
	runtime, err := p.getRuntime(runtimeId)
	if err != nil {
		return nil, err
	}
	pbRuntime, err := p.formatRuntime(runtime)
	if err != nil {
		logger.Errorf("CreateRuntime: %+v", err)
		return nil, status.Errorf(codes.Internal, "CreateRuntime: %+v", err)
	}
	res := &pb.CreateRuntimeResponse{
		Runtime: pbRuntime,
	}
	return res, nil
}

func (p *Server) DescribeRuntimes(ctx context.Context, req *pb.DescribeRuntimesRequest) (*pb.DescribeRuntimesResponse, error) {
	// validate req
	err := validateDescribeRuntimesRequest(req)
	offset := utils.GetOffsetFromRequest(req)
	limit := utils.GetLimitFromRequest(req)
	if err != nil {
		logger.Errorf("DescribeRuntimes: %+v", err)
		return nil, status.Errorf(codes.InvalidArgument, "DescribeRuntimes: %+v", err)
	}
	selectorMap, err := SelectorStringToMap(req.Label.GetValue())
	if err != nil {
		logger.Errorf("DescribeRuntimes: %+v", err)
		return nil, status.Errorf(codes.InvalidArgument, "DescribeRuntimes: %+v", err)
	}

	var runtimes []*models.Runtime
	var count uint32
	query := p.Db.
		Select(models.RuntimeColumnsWithTablePrefix...).
		From(models.RuntimeTableName).
		Offset(offset).
		Limit(limit).
		Where(manager.BuildFilterConditionsWithPrefix(req, models.RuntimeTableName))

	query = db.AddJoinFilterWithMap(query, models.RuntimeTableName, models.RuntimeLabelTableName, RuntimeIdColumn,
		models.ColumnLabelKey, models.ColumnLabelValue, selectorMap)

	_, err = query.Load(&runtimes)
	if err != nil {
		logger.Errorf("DescribeRuntimes: %+v", err)
		return nil, status.Errorf(codes.Internal, "DescribeRuntimes: %+v", err)
	}

	count, err = query.Count()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "DescribeRuntimes: %+v", err)
	}
	pbRuntime, err := p.formatRuntimeSet(runtimes)
	if err != nil {
		logger.Errorf("DescribeRuntimes: %+v", err)
		return nil, status.Errorf(codes.Internal, "DescribeRuntimes %+v", err)
	}
	res := &pb.DescribeRuntimesResponse{
		RuntimeSet: pbRuntime,
		TotalCount: count,
	}
	return res, nil
}

func (p *Server) ModifyRuntime(ctx context.Context, req *pb.ModifyRuntimeRequest) (*pb.ModifyRuntimeResponse, error) {
	// validate req
	err := validateModifyRuntimeRequest(req)
	if err != nil {
		logger.Errorf("ModifyRuntime: %+v", err)
		return nil, status.Errorf(codes.InvalidArgument, "ModifyRuntime: %+v", err)
	}
	// check runtime can be modified
	runtimeId := req.GetRuntimeId().GetValue()
	deleted, err := p.checkRuntimeDeleted(runtimeId)
	if err != nil {
		logger.Errorf("ModifyRuntime: %+v", err)
		return nil, status.Errorf(codes.Internal, "ModifyRuntime: %+v", err)
	}
	if deleted {
		logger.Errorf("ModifyRuntime: runtime has been deleted [%+v]", runtimeId)
		return nil, status.Errorf(codes.Internal,
			"ModifyRuntime: runtime has been deleted [%+v]", runtimeId)
	}
	// update runtime
	err = p.updateRuntime(req)
	if err != nil {
		logger.Errorf("ModifyRuntime: %+v", err)
		return nil, status.Errorf(codes.Internal, "ModifyRuntime: %+v", err)
	}

	// update runtime label
	if req.Labels != nil {
		err := p.updateRuntimeLabels(runtimeId, req.Labels.GetValue())
		if err != nil {
			logger.Errorf("ModifyRuntime: %+v", err)
			return nil, status.Errorf(codes.Internal, "ModifyRuntime: %+v", err)
		}
	}

	// get response
	runtime, err := p.getRuntime(runtimeId)
	if err != nil {
		return nil, err
	}
	pbRuntime, err := p.formatRuntime(runtime)
	if err != nil {
		logger.Errorf("ModifyRuntime: %+v", err)
		return nil, status.Errorf(codes.Internal, "ModifyRuntime: %+v", err)
	}
	res := &pb.ModifyRuntimeResponse{
		Runtime: pbRuntime,
	}

	return res, nil
}

func (p *Server) DeleteRuntime(ctx context.Context, req *pb.DeleteRuntimeRequest) (*pb.DeleteRuntimeResponse, error) {
	// validate req
	err := validateDeleteRuntimeRequest(req)
	if err != nil {
		logger.Errorf("DeleteRuntime: %+v", err)
		return nil, status.Errorf(codes.InvalidArgument, "DeleteRuntime: %+v", err)
	}

	// check runtime can be deleted
	runtimeId := req.GetRuntimeId().GetValue()
	deleted, err := p.checkRuntimeDeleted(runtimeId)
	if err != nil {
		logger.Errorf("DeleteRuntime: %+v", err)
		return nil, status.Errorf(codes.Internal,
			"DeleteRuntime: %+v", err)
	}
	if deleted {
		logger.Errorf("DeleteRuntime: runtime has been deleted [%+v]", runtimeId)
		return nil, status.Errorf(codes.Internal,
			"DeleteRuntime: runtime has been deleted [%+v]", runtimeId)
	}
	// deleted runtime
	err = p.deleteRuntime(runtimeId)
	if err != nil {
		logger.Errorf("DeleteRuntime: %+v", err)
		return nil, status.Errorf(codes.Internal, "DeleteRuntime: %+v", err)
	}

	// get runtime
	runtime, err := p.getRuntime(runtimeId)
	if err != nil {
		return nil, err
	}
	pbRuntime, err := p.formatRuntime(runtime)
	if err != nil {
		logger.Errorf("DeleteRuntime: %+v", err)
		return nil, status.Errorf(codes.Internal, "DeleteRuntime: %+v", err)
	}
	res := &pb.DeleteRuntimeResponse{
		Runtime: pbRuntime,
	}
	return res, nil
}

func (p *Server) DescribeRuntimeProviderZones(ctx context.Context, req *pb.DescribeRuntimeProviderZonesRequest) (*pb.DescribeRuntimeProviderZonesResponse, error) {
	err := ValidateCredential(req.Provider.GetValue(), req.RuntimeUrl.GetValue(), req.RuntimeCredential.GetValue())
	if err != nil {
		logger.Errorf("DescribeRuntimeProviderZones: %+v", err)
		return nil, status.Errorf(codes.Internal, "DescribeRuntimeProviderZones: %+v", err)
	}
	// TODO : DescribeRuntimeProviderZones by provider
	return nil, nil
}
