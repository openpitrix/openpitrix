// Copyright 2017 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package runtime

import (
	"context"
	"time"

	"openpitrix.io/openpitrix/pkg/constants"
	"openpitrix.io/openpitrix/pkg/db"
	"openpitrix.io/openpitrix/pkg/gerr"
	"openpitrix.io/openpitrix/pkg/logger"
	"openpitrix.io/openpitrix/pkg/manager"
	"openpitrix.io/openpitrix/pkg/models"
	"openpitrix.io/openpitrix/pkg/pb"
	"openpitrix.io/openpitrix/pkg/pi"
	"openpitrix.io/openpitrix/pkg/plugins"
	"openpitrix.io/openpitrix/pkg/util/labelutil"
	"openpitrix.io/openpitrix/pkg/util/pbutil"
	"openpitrix.io/openpitrix/pkg/util/senderutil"
)

func (p *Server) CreateRuntime(ctx context.Context, req *pb.CreateRuntimeRequest) (*pb.CreateRuntimeResponse, error) {
	s := senderutil.GetSenderFromContext(ctx)
	// validate req
	err := validateCreateRuntimeRequest(ctx, req)
	// TODO: refactor create runtime params
	if err != nil {
		if gerr.IsGRPCError(err) {
			return nil, err
		} else {
			return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorValidateFailed)
		}
	}

	// create runtime credential
	runtimeCredentialId, err := createRuntimeCredential(ctx, req.Provider.GetValue(), req.RuntimeCredential.GetValue())
	if err != nil {
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorCreateResourcesFailed)
	}

	// create runtime
	runtimeId, err := createRuntime(
		ctx,
		req.GetName().GetValue(),
		req.GetDescription().GetValue(),
		req.Provider.GetValue(),
		req.GetRuntimeUrl().GetValue(),
		runtimeCredentialId,
		req.Zone.GetValue(),
		s.UserId)
	if err != nil {
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorCreateResourcesFailed)
	}

	if req.GetLabels() != nil {
		err = labelutil.SyncRuntimeLabels(ctx, runtimeId, req.GetLabels().GetValue())
		if err != nil {
			return nil, err
		}
	}

	res := &pb.CreateRuntimeResponse{
		RuntimeId: pbutil.ToProtoString(runtimeId),
	}
	return res, nil
}

func (p *Server) DescribeRuntimes(ctx context.Context, req *pb.DescribeRuntimesRequest) (*pb.DescribeRuntimesResponse, error) {
	// TODO: refactor validate req
	err := validateDescribeRuntimesRequest(ctx, req)
	if err != nil {
		if gerr.IsGRPCError(err) {
			return nil, err
		} else {
			return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorValidateFailed)
		}
	}
	offset := pbutil.GetOffsetFromRequest(req)
	limit := pbutil.GetLimitFromRequest(req)
	selectorMap, err := SelectorStringToMap(req.Label.GetValue())
	if err != nil {
		return nil, gerr.NewWithDetail(ctx, gerr.InvalidArgument, err, gerr.ErrorParameterParseFailed, "label")
	}

	var runtimes []*models.Runtime
	var count uint32
	query := pi.Global().DB(ctx).
		Select(models.RuntimeColumnsWithTablePrefix...).
		From(constants.TableRuntime).
		Offset(offset).
		Limit(limit).
		Where(manager.BuildFilterConditionsWithPrefix(req, constants.TableRuntime))

	query = manager.AddQueryJoinWithMap(query, constants.TableRuntime, constants.TableRuntimeLabel, constants.ColumnRuntimeId,
		constants.ColumnLabelKey, constants.ColumnLabelValue, selectorMap)
	query = manager.AddQueryOrderDir(query, req, constants.ColumnCreateTime)
	_, err = query.Load(&runtimes)
	if err != nil {
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorDescribeResourcesFailed)
	}

	count, err = query.Count()
	if err != nil {
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorDescribeResourcesFailed)
	}
	pbRuntime, err := formatRuntimeSet(ctx, runtimes)
	if err != nil {
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorDescribeResourcesFailed)
	}
	res := &pb.DescribeRuntimesResponse{
		RuntimeSet: pbRuntime,
		TotalCount: count,
	}
	return res, nil
}

func (p *Server) DescribeRuntimeDetails(ctx context.Context, req *pb.DescribeRuntimesRequest) (*pb.DescribeRuntimeDetailsResponse, error) {
	// TODO: refactor validate req
	err := validateDescribeRuntimesRequest(ctx, req)
	if err != nil {
		if gerr.IsGRPCError(err) {
			return nil, err
		} else {
			return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorValidateFailed)
		}
	}
	offset := pbutil.GetOffsetFromRequest(req)
	limit := pbutil.GetLimitFromRequest(req)
	selectorMap, err := SelectorStringToMap(req.Label.GetValue())
	if err != nil {
		return nil, gerr.NewWithDetail(ctx, gerr.InvalidArgument, err, gerr.ErrorParameterParseFailed, "label")
	}

	var runtimes []*models.Runtime
	var count uint32
	query := pi.Global().DB(ctx).
		Select(models.RuntimeColumnsWithTablePrefix...).
		From(constants.TableRuntime).
		Offset(offset).
		Limit(limit).
		Where(manager.BuildFilterConditionsWithPrefix(req, constants.TableRuntime))

	query = manager.AddQueryJoinWithMap(query, constants.TableRuntime, constants.TableRuntimeLabel, constants.ColumnRuntimeId,
		constants.ColumnLabelKey, constants.ColumnLabelValue, selectorMap)
	query = manager.AddQueryOrderDir(query, req, constants.ColumnCreateTime)
	_, err = query.Load(&runtimes)
	if err != nil {
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorDescribeResourcesFailed)
	}

	count, err = query.Count()
	if err != nil {
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorDescribeResourcesFailed)
	}
	pbRuntimeDetails, err := formatRuntimeDetailSet(ctx, runtimes)
	if err != nil {
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorDescribeResourcesFailed)
	}
	res := &pb.DescribeRuntimeDetailsResponse{
		RuntimeDetailSet: pbRuntimeDetails,
		TotalCount:       count,
	}
	return res, nil
}

func (p *Server) ModifyRuntime(ctx context.Context, req *pb.ModifyRuntimeRequest) (*pb.ModifyRuntimeResponse, error) {
	// validate req
	err := validateModifyRuntimeRequest(ctx, req)
	if err != nil {
		if gerr.IsGRPCError(err) {
			return nil, err
		} else {
			return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorValidateFailed)
		}
	}
	// check runtime can be modified
	runtimeId := req.GetRuntimeId().GetValue()
	runtime, err := getRuntime(ctx, runtimeId)
	if err != nil {
		return nil, gerr.NewWithDetail(ctx, gerr.FailedPrecondition, err, gerr.ErrorResourceNotFound, runtimeId)
	}
	if req.RuntimeCredential != nil {
		err = ValidateCredential(
			ctx,
			runtime.Provider,
			runtime.RuntimeUrl,
			req.RuntimeCredential.GetValue(),
			runtime.Zone)
		if err != nil {
			if gerr.IsGRPCError(err) {
				return nil, err
			} else {
				return nil, gerr.NewWithDetail(ctx, gerr.PermissionDenied, err, gerr.ErrorValidateFailed)
			}
		}
	}
	if runtime.Status == constants.StatusDeleted {
		logger.Error(ctx, "runtime has been deleted [%s]", runtimeId)
		return nil, gerr.NewWithDetail(ctx, gerr.FailedPrecondition, err, gerr.ErrorResourceAlreadyDeleted, runtimeId)
	}
	// update runtime
	err = updateRuntime(ctx, req)
	if err != nil {
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorModifyResourcesFailed)
	}

	if req.GetLabels() != nil {
		err = labelutil.SyncRuntimeLabels(ctx, runtimeId, req.GetLabels().GetValue())
		if err != nil {
			return nil, err
		}
	}

	// update runtime credential
	if req.RuntimeCredential != nil {
		err := updateRuntimeCredential(ctx, runtime.RuntimeCredentialId, runtime.Provider, req.RuntimeCredential.GetValue())
		if err != nil {
			return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorModifyResourcesFailed)
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
	err := deleteRuntimes(ctx, runtimeIds)
	if err != nil {
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorDeleteResourcesFailed)
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
	err := ValidateCredential(ctx, provider, url, credential, "")
	if err != nil {
		if gerr.IsGRPCError(err) {
			return nil, err
		} else {
			return nil, gerr.NewWithDetail(ctx, gerr.PermissionDenied, err, gerr.ErrorValidateFailed)
		}
	}

	providerInterface, err := plugins.GetProviderPlugin(ctx, provider)
	if err != nil {
		return nil, gerr.NewWithDetail(ctx, gerr.NotFound, err, gerr.ErrorProviderNotFound, provider)
	}
	zones, err := providerInterface.DescribeRuntimeProviderZones(url, credential)
	if err != nil {
		return nil, gerr.NewWithDetail(ctx, gerr.PermissionDenied, err, gerr.ErrorDescribeResourceFailed)
	}
	return &pb.DescribeRuntimeProviderZonesResponse{
		Provider: req.Provider,
		Zone:     zones,
	}, nil
}

type runtimeStatistic struct {
	Date  string `db:"DATE_FORMAT(create_time, '%Y-%m-%d')"`
	Count uint32 `db:"COUNT(runtime_id)"`
}
type providerStatistic struct {
	Provider string `db:"provider"`
	Count    uint32 `db:"COUNT(runtime_id)"`
}

func (p *Server) GetRuntimeStatistics(ctx context.Context, req *pb.GetRuntimeStatisticsRequest) (*pb.GetRuntimeStatisticsResponse, error) {
	res := &pb.GetRuntimeStatisticsResponse{
		LastTwoWeekCreated: make(map[string]uint32),
		TopTenProviders:    make(map[string]uint32),
	}
	runtimeCount, err := pi.Global().DB(ctx).
		Select(constants.ColumnRuntimeId).
		From(constants.TableRuntime).
		Where(db.Neq(constants.ColumnStatus, constants.StatusDeleted)).
		Count()
	if err != nil {
		logger.Error(ctx, "Failed to get runtime count, error: %+v", err)
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorDescribeResourcesFailed)
	}
	res.RuntimeCount = runtimeCount

	err = pi.Global().DB(ctx).
		Select("COUNT(DISTINCT provider)").
		From(constants.TableRuntime).
		Where(db.Neq(constants.ColumnStatus, constants.StatusDeleted)).
		LoadOne(&res.ProviderCount)
	if err != nil {
		logger.Error(ctx, "Failed to get provider count, error: %+v", err)
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorDescribeResourcesFailed)
	}

	time2week := time.Now().Add(-14 * 24 * time.Hour)
	var rs []*runtimeStatistic
	_, err = pi.Global().DB(ctx).
		Select("DATE_FORMAT(create_time, '%Y-%m-%d')", "COUNT(runtime_id)").
		From(constants.TableRuntime).
		GroupBy("DATE_FORMAT(create_time, '%Y-%m-%d')").
		Where(db.Gte(constants.ColumnCreateTime, time2week)).
		Limit(14).Load(&rs)

	if err != nil {
		logger.Error(ctx, "Failed to get runtime statistics, error: %+v", err)
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorDescribeResourcesFailed)
	}
	for _, a := range rs {
		res.LastTwoWeekCreated[a.Date] = a.Count
	}

	var ps []*providerStatistic
	_, err = pi.Global().DB(ctx).
		Select("provider", "COUNT(runtime_id)").
		From(constants.TableRuntime).
		Where(db.Neq(constants.ColumnStatus, constants.StatusDeleted)).
		GroupBy(constants.ColumnProvider).
		OrderDir("COUNT(runtime_id)", false).
		Limit(10).Load(&ps)

	if err != nil {
		logger.Error(ctx, "Failed to get provider statistics, error: %+v", err)
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorDescribeResourcesFailed)
	}
	for _, a := range ps {
		res.TopTenProviders[a.Provider] = a.Count
	}

	return res, nil
}
