// Copyright 2017 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package runtime

import (
	"context"

	"openpitrix.io/openpitrix/pkg/constants"
	"openpitrix.io/openpitrix/pkg/gerr"
	"openpitrix.io/openpitrix/pkg/logger"
	"openpitrix.io/openpitrix/pkg/manager"
	"openpitrix.io/openpitrix/pkg/models"
	"openpitrix.io/openpitrix/pkg/pb"
	"openpitrix.io/openpitrix/pkg/pi"
	"openpitrix.io/openpitrix/pkg/plugins"
	"openpitrix.io/openpitrix/pkg/util/pbutil"
	"openpitrix.io/openpitrix/pkg/util/senderutil"
)

func (p *Server) CreateRuntime(ctx context.Context, req *pb.CreateRuntimeRequest) (*pb.CreateRuntimeResponse, error) {
	s := senderutil.GetSenderFromContext(ctx)
	// validate req
	err := validateCreateRuntimeRequest(req)
	// TODO: refactor create runtime params
	if err != nil {
		if gerr.IsGRPCError(err) {
			return nil, err
		} else {
			return nil, gerr.NewWithDetail(gerr.Internal, err, gerr.ErrorValidateFailed)
		}
	}

	// create runtime credential
	runtimeCredentialId, err := p.createRuntimeCredential(req.Provider.GetValue(), req.RuntimeCredential.GetValue())
	if err != nil {
		return nil, gerr.NewWithDetail(gerr.Internal, err, gerr.ErrorCreateResourcesFailed)
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
		return nil, gerr.NewWithDetail(gerr.Internal, err, gerr.ErrorCreateResourcesFailed)
	}

	// create labels
	err = p.createRuntimeLabels(runtimeId, req.Labels.GetValue())
	if err != nil {
		return nil, gerr.NewWithDetail(gerr.Internal, err, gerr.ErrorCreateResourcesFailed)
	}

	res := &pb.CreateRuntimeResponse{
		RuntimeId: pbutil.ToProtoString(runtimeId),
	}
	return res, nil
}

func (p *Server) DescribeRuntimes(ctx context.Context, req *pb.DescribeRuntimesRequest) (*pb.DescribeRuntimesResponse, error) {
	// TODO: refactor validate req
	err := validateDescribeRuntimesRequest(req)
	if err != nil {
		if gerr.IsGRPCError(err) {
			return nil, err
		} else {
			return nil, gerr.NewWithDetail(gerr.Internal, err, gerr.ErrorValidateFailed)
		}
	}
	offset := pbutil.GetOffsetFromRequest(req)
	limit := pbutil.GetLimitFromRequest(req)
	selectorMap, err := SelectorStringToMap(req.Label.GetValue())
	if err != nil {
		return nil, gerr.NewWithDetail(gerr.InvalidArgument, err, gerr.ErrorParameterParseFailed, "label")
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
	query = manager.AddQueryOrderDir(query, req, models.ColumnCreateTime)
	_, err = query.Load(&runtimes)
	if err != nil {
		return nil, gerr.NewWithDetail(gerr.Internal, err, gerr.ErrorDescribeResourcesFailed)
	}

	count, err = query.Count()
	if err != nil {
		return nil, gerr.NewWithDetail(gerr.Internal, err, gerr.ErrorDescribeResourcesFailed)
	}
	pbRuntime, err := p.formatRuntimeSet(runtimes)
	if err != nil {
		return nil, gerr.NewWithDetail(gerr.Internal, err, gerr.ErrorDescribeResourcesFailed)
	}

	providerCount := 0
	providerQuery := pi.Global().Db.SelectBySql("select count(distinct provider) from runtime;")
	err = providerQuery.LoadOne(&providerCount)
	if err != nil {
		return nil, gerr.NewWithDetail(gerr.Internal, err, gerr.ErrorDescribeResourcesFailed)
	}

	lastTwoWeekCount := 0
	lastTwoWeekQuery := pi.Global().Db.SelectBySql("select count(runtime_id) from runtime where DATE_SUB(CURDATE(), INTERVAL 14 DAY) <= date(create_time);")
	err = lastTwoWeekQuery.LoadOne(&lastTwoWeekCount)
	if err != nil {
		return nil, gerr.NewWithDetail(gerr.Internal, err, gerr.ErrorDescribeResourcesFailed)
	}

	res := &pb.DescribeRuntimesResponse{
		RuntimeSet:       pbRuntime,
		TotalCount:       count,
		ProviderCount:    uint32(providerCount),
		LastTwoWeekCount: uint32(lastTwoWeekCount),
	}
	return res, nil
}

func (p *Server) DescribeRuntimeDetails(ctx context.Context, req *pb.DescribeRuntimesRequest) (*pb.DescribeRuntimeDetailsResponse, error) {
	// TODO: refactor validate req
	err := validateDescribeRuntimesRequest(req)
	if err != nil {
		if gerr.IsGRPCError(err) {
			return nil, err
		} else {
			return nil, gerr.NewWithDetail(gerr.Internal, err, gerr.ErrorValidateFailed)
		}
	}
	offset := pbutil.GetOffsetFromRequest(req)
	limit := pbutil.GetLimitFromRequest(req)
	selectorMap, err := SelectorStringToMap(req.Label.GetValue())
	if err != nil {
		return nil, gerr.NewWithDetail(gerr.InvalidArgument, err, gerr.ErrorParameterParseFailed, "label")
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
	query = manager.AddQueryOrderDir(query, req, models.ColumnCreateTime)
	_, err = query.Load(&runtimes)
	if err != nil {
		return nil, gerr.NewWithDetail(gerr.Internal, err, gerr.ErrorDescribeResourcesFailed)
	}

	count, err = query.Count()
	if err != nil {
		return nil, gerr.NewWithDetail(gerr.Internal, err, gerr.ErrorDescribeResourcesFailed)
	}
	pbRuntimeDetails, err := p.formatRuntimeDetailSet(runtimes)
	if err != nil {
		return nil, gerr.NewWithDetail(gerr.Internal, err, gerr.ErrorDescribeResourcesFailed)
	}
	res := &pb.DescribeRuntimeDetailsResponse{
		RuntimeDetailSet: pbRuntimeDetails,
		TotalCount:       count,
	}
	return res, nil
}

func (p *Server) ModifyRuntime(ctx context.Context, req *pb.ModifyRuntimeRequest) (*pb.ModifyRuntimeResponse, error) {
	// validate req
	err := validateModifyRuntimeRequest(req)
	if err != nil {
		if gerr.IsGRPCError(err) {
			return nil, err
		} else {
			return nil, gerr.NewWithDetail(gerr.Internal, err, gerr.ErrorValidateFailed)
		}
	}
	// check runtime can be modified
	runtimeId := req.GetRuntimeId().GetValue()
	runtime, err := p.getRuntime(runtimeId)
	if err != nil {
		return nil, gerr.NewWithDetail(gerr.FailedPrecondition, err, gerr.ErrorResourceNotFound, runtimeId)
	}
	if runtime.Status == constants.StatusDeleted {
		logger.Error("runtime has been deleted [%s]", runtimeId)
		return nil, gerr.NewWithDetail(gerr.FailedPrecondition, err, gerr.ErrorResourceAlreadyDeleted, runtimeId)
	}
	// update runtime
	err = p.updateRuntime(req)
	if err != nil {
		return nil, gerr.NewWithDetail(gerr.Internal, err, gerr.ErrorModifyResourcesFailed)
	}

	// update runtime label
	if req.Labels != nil {
		err := p.updateRuntimeLabels(runtimeId, req.Labels.GetValue())
		if err != nil {
			return nil, gerr.NewWithDetail(gerr.Internal, err, gerr.ErrorModifyResourcesFailed)
		}
	}

	res := &pb.ModifyRuntimeResponse{
		RuntimeId: req.GetRuntimeId(),
	}

	return res, nil
}

func (p *Server) DeleteRuntimes(ctx context.Context, req *pb.DeleteRuntimesRequest) (*pb.DeleteRuntimesResponse, error) {
	// TODO: check runtime can be deleted
	runtimeIds := req.GetRuntimeId()

	// deleted runtime
	err := p.deleteRuntimes(runtimeIds)
	if err != nil {
		return nil, gerr.NewWithDetail(gerr.Internal, err, gerr.ErrorDeleteResourcesFailed)
	}

	res := &pb.DeleteRuntimesResponse{
		RuntimeId: runtimeIds,
	}
	return res, nil
}

func (p *Server) DescribeRuntimeProviderZones(ctx context.Context, req *pb.DescribeRuntimeProviderZonesRequest) (*pb.DescribeRuntimeProviderZonesResponse, error) {
	provider := req.Provider.GetValue()
	url := req.RuntimeUrl.GetValue()
	credential := req.RuntimeCredential.GetValue()
	err := ValidateCredential(provider, url, credential, "")
	if err != nil {
		if gerr.IsGRPCError(err) {
			return nil, err
		} else {
			return nil, gerr.NewWithDetail(gerr.PermissionDenied, err, gerr.ErrorValidateFailed)
		}
	}

	providerInterface, err := plugins.GetProviderPlugin(provider, nil)
	if err != nil {
		return nil, gerr.NewWithDetail(gerr.NotFound, err, gerr.ErrorProviderNotFound, provider)
	}
	zones, err := providerInterface.DescribeRuntimeProviderZones(url, credential)
	if err != nil {
		return nil, gerr.NewWithDetail(gerr.PermissionDenied, err, gerr.ErrorDescribeResourceFailed)
	}
	return &pb.DescribeRuntimeProviderZonesResponse{
		Provider: req.Provider,
		Zone:     zones,
	}, nil
}
