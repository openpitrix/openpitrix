// Copyright 2017 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package runtime

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"openpitrix.io/openpitrix/pkg/logger"
	"openpitrix.io/openpitrix/pkg/manager"
	"openpitrix.io/openpitrix/pkg/models"
	"openpitrix.io/openpitrix/pkg/pb"
	"openpitrix.io/openpitrix/pkg/plugins"
	"openpitrix.io/openpitrix/pkg/util/pbutil"
	"openpitrix.io/openpitrix/pkg/util/senderutil"
)

func (p *Server) CreateRuntime(ctx context.Context, req *pb.CreateRuntimeRequest) (*pb.CreateRuntimeResponse, error) {
	s := senderutil.GetSenderFromContext(ctx)
	// validate req
	err := validateCreateRuntimeRequest(req)
	if err != nil {
		logger.Error("CreateRuntime: %+v", err)
		return nil, status.Errorf(codes.InvalidArgument, "CreateRuntime: %+v", err)
	}

	// create runtime credential
	runtimeCredentialId, err := p.createRuntimeCredential(req.Provider.GetValue(), req.RuntimeCredential.GetValue())
	if err != nil {
		logger.Error("CreateRuntime: %+v", err)
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
		logger.Error("CreateRuntime: %+v", err)
		return nil, status.Errorf(codes.Internal, "CreateRuntime: %+v", err)
	}

	// create labels
	err = p.createRuntimeLabels(runtimeId, req.Labels.GetValue())
	if err != nil {
		logger.Error("CreateRuntime: %+v", err)
		return nil, status.Errorf(codes.Internal, "CreateRuntime: %+v", err)
	}

	// get response
	runtime, err := p.getRuntime(runtimeId)
	if err != nil {
		return nil, err
	}
	pbRuntime, err := p.formatRuntime(runtime)
	if err != nil {
		logger.Error("CreateRuntime: %+v", err)
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
	offset := pbutil.GetOffsetFromRequest(req)
	limit := pbutil.GetLimitFromRequest(req)
	if err != nil {
		logger.Error("DescribeRuntimes: %+v", err)
		return nil, status.Errorf(codes.InvalidArgument, "DescribeRuntimes: %+v", err)
	}
	selectorMap, err := SelectorStringToMap(req.Label.GetValue())
	if err != nil {
		logger.Error("DescribeRuntimes: %+v", err)
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

	query = manager.AddQueryJoinWithMap(query, models.RuntimeTableName, models.RuntimeLabelTableName, RuntimeIdColumn,
		models.ColumnLabelKey, models.ColumnLabelValue, selectorMap)

	_, err = query.Load(&runtimes)
	if err != nil {
		logger.Error("DescribeRuntimes: %+v", err)
		return nil, status.Errorf(codes.Internal, "DescribeRuntimes: %+v", err)
	}

	count, err = query.Count()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "DescribeRuntimes: %+v", err)
	}
	pbRuntime, err := p.formatRuntimeSet(runtimes)
	if err != nil {
		logger.Error("DescribeRuntimes: %+v", err)
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
		logger.Error("ModifyRuntime: %+v", err)
		return nil, status.Errorf(codes.InvalidArgument, "ModifyRuntime: %+v", err)
	}
	// check runtime can be modified
	runtimeId := req.GetRuntimeId().GetValue()
	deleted, err := p.checkRuntimeDeleted(runtimeId)
	if err != nil {
		logger.Error("ModifyRuntime: %+v", err)
		return nil, status.Errorf(codes.Internal, "ModifyRuntime: %+v", err)
	}
	if deleted {
		logger.Error("ModifyRuntime: runtime has been deleted [%+v]", runtimeId)
		return nil, status.Errorf(codes.Internal,
			"ModifyRuntime: runtime has been deleted [%+v]", runtimeId)
	}
	// update runtime
	err = p.updateRuntime(req)
	if err != nil {
		logger.Error("ModifyRuntime: %+v", err)
		return nil, status.Errorf(codes.Internal, "ModifyRuntime: %+v", err)
	}

	// update runtime label
	if req.Labels != nil {
		err := p.updateRuntimeLabels(runtimeId, req.Labels.GetValue())
		if err != nil {
			logger.Error("ModifyRuntime: %+v", err)
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
		logger.Error("ModifyRuntime: %+v", err)
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
		logger.Error("DeleteRuntime: %+v", err)
		return nil, status.Errorf(codes.InvalidArgument, "DeleteRuntime: %+v", err)
	}

	// check runtime can be deleted
	runtimeId := req.GetRuntimeId().GetValue()
	deleted, err := p.checkRuntimeDeleted(runtimeId)
	if err != nil {
		logger.Error("DeleteRuntime: %+v", err)
		return nil, status.Errorf(codes.Internal,
			"DeleteRuntime: %+v", err)
	}
	if deleted {
		logger.Error("DeleteRuntime: runtime has been deleted [%+v]", runtimeId)
		return nil, status.Errorf(codes.Internal,
			"DeleteRuntime: runtime has been deleted [%+v]", runtimeId)
	}
	// deleted runtime
	err = p.deleteRuntime(runtimeId)
	if err != nil {
		logger.Error("DeleteRuntime: %+v", err)
		return nil, status.Errorf(codes.Internal, "DeleteRuntime: %+v", err)
	}

	// get runtime
	runtime, err := p.getRuntime(runtimeId)
	if err != nil {
		return nil, err
	}
	pbRuntime, err := p.formatRuntime(runtime)
	if err != nil {
		logger.Error("DeleteRuntime: %+v", err)
		return nil, status.Errorf(codes.Internal, "DeleteRuntime: %+v", err)
	}
	res := &pb.DeleteRuntimeResponse{
		Runtime: pbRuntime,
	}
	return res, nil
}

func (p *Server) DescribeRuntimeProviderZones(ctx context.Context, req *pb.DescribeRuntimeProviderZonesRequest) (*pb.DescribeRuntimeProviderZonesResponse, error) {
	provider := req.Provider.GetValue()
	url := req.RuntimeUrl.GetValue()
	credential := req.RuntimeCredential.GetValue()
	err := ValidateCredential(provider, url, credential)
	if err != nil {
		logger.Error("DescribeRuntimeProviderZones: %+v", err)
		return nil, status.Errorf(codes.Internal, "DescribeRuntimeProviderZones: %+v", err)
	}

	providerInterface, err := plugins.GetProviderPlugin(provider)
	if err != nil {
		logger.Error("No such provider [%s]. ", provider)
		return nil, err
	}
	zones := providerInterface.DescribeRuntimeProviderZones(url, credential)
	return &pb.DescribeRuntimeProviderZonesResponse{
		Provider: req.Provider,
		Zone:     zones,
	}, nil
}
