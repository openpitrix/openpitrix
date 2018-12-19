// Copyright 2017 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package runtime

import (
	"context"
	"fmt"
	"time"

	clusterclient "openpitrix.io/openpitrix/pkg/client/cluster"
	"openpitrix.io/openpitrix/pkg/constants"
	"openpitrix.io/openpitrix/pkg/db"
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
	runtimeId := models.NewRuntimeId()
	runtimeCredentialId := req.GetRuntimeCredentialId().GetValue()
	zone := req.GetZone().GetValue()

	runtimeCredential, err := CheckRuntimeCredentialPermission(ctx, runtimeCredentialId)
	if err != nil {
		return nil, err
	}

	if runtimeCredential.Status == constants.StatusDeleted {
		logger.Error(ctx, "Runtime credential [%s] has been deleted", runtimeCredentialId)
		return nil, gerr.NewWithDetail(ctx, gerr.PermissionDenied, err, gerr.ErrorResourceAlreadyDeleted, runtimeCredentialId)
	}

	if runtimeCredential.Provider != req.GetProvider().GetValue() {
		logger.Error(ctx, "Runtime credential [%s] provider is [%s] not [%s]", runtimeCredentialId,
			runtimeCredential.Provider, req.GetProvider().GetValue())
		return nil, gerr.NewWithDetail(ctx, gerr.PermissionDenied, err, gerr.ErrorCreateResourcesFailed)
	}

	runtimeCredential.RuntimeCredentialContent, err = encodeRuntimeCredentialContent(
		runtimeCredential.Provider, runtimeCredential.RuntimeCredentialContent)
	if err != nil {
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorInternalError)
	}

	err = ValidateRuntime(ctx, runtimeId, zone, runtimeCredential, true)
	if err != nil {
		if gerr.IsGRPCError(err) {
			return nil, err
		} else {
			return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorValidateFailed)
		}
	}

	newRuntime := models.NewRuntime(
		runtimeId,
		req.GetName().GetValue(),
		req.GetDescription().GetValue(),
		req.GetProvider().GetValue(),
		runtimeCredentialId,
		req.GetZone().GetValue(),
		s.UserId,
	)
	_, err = pi.Global().DB(ctx).
		InsertInto(constants.TableRuntime).
		Record(newRuntime).
		Exec()
	if err != nil {
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorCreateResourcesFailed)
	}

	res := &pb.CreateRuntimeResponse{
		RuntimeId: pbutil.ToProtoString(runtimeId),
	}
	return res, nil
}

func (p *Server) DescribeRuntimes(ctx context.Context, req *pb.DescribeRuntimesRequest) (*pb.DescribeRuntimesResponse, error) {
	offset := pbutil.GetOffsetFromRequest(req)
	limit := pbutil.GetLimitFromRequest(req)

	var runtimes []*models.Runtime
	var count uint32
	query := pi.Global().DB(ctx).
		Select(models.RuntimeColumns...).
		From(constants.TableRuntime).
		Offset(offset).
		Limit(limit).
		Where(manager.BuildFilterConditions(req, constants.TableRuntime))

	query = manager.AddQueryOrderDir(query, req, constants.ColumnCreateTime)
	_, err := query.Load(&runtimes)
	if err != nil {
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorDescribeResourcesFailed)
	}

	count, err = query.Count()
	if err != nil {
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorDescribeResourcesFailed)
	}
	res := &pb.DescribeRuntimesResponse{
		RuntimeSet: models.RuntimeToPbs(runtimes),
		TotalCount: count,
	}
	return res, nil
}

func (p *Server) DescribeRuntimeDetails(ctx context.Context, req *pb.DescribeRuntimesRequest) (*pb.DescribeRuntimeDetailsResponse, error) {
	offset := pbutil.GetOffsetFromRequest(req)
	limit := pbutil.GetLimitFromRequest(req)

	var runtimes []*models.Runtime
	var count uint32
	query := pi.Global().DB(ctx).
		Select(models.RuntimeColumns...).
		From(constants.TableRuntime).
		Offset(offset).
		Limit(limit).
		Where(manager.BuildFilterConditions(req, constants.TableRuntime))
	query = manager.AddQueryOrderDir(query, req, constants.ColumnCreateTime)
	_, err := query.Load(&runtimes)
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
	runtimeId := req.GetRuntimeId().GetValue()
	runtime, err := CheckRuntimePermission(ctx, runtimeId)
	if err != nil {
		return nil, err
	}

	if runtime.Status == constants.StatusDeleted {
		logger.Error(ctx, "Runtime [%s] has been deleted", runtimeId)
		return nil, gerr.NewWithDetail(ctx, gerr.PermissionDenied, err, gerr.ErrorResourceAlreadyDeleted, runtimeId)
	}

	attributes := manager.BuildUpdateAttributes(req, constants.ColumnName, constants.ColumnDescription)
	attributes[constants.ColumnStatusTime] = time.Now()

	runtimeCredentialId := req.GetRuntimeCredentialId().GetValue()
	if len(runtimeCredentialId) > 0 {
		runtimeCredential, err := CheckRuntimeCredentialPermission(ctx, runtimeCredentialId)
		if err != nil {
			return nil, err
		}
		if runtimeCredential.Status == constants.StatusDeleted {
			logger.Error(ctx, "Runtime credential [%s] has been deleted", runtimeCredentialId)
			return nil, gerr.NewWithDetail(ctx, gerr.PermissionDenied, err, gerr.ErrorResourceAlreadyDeleted, runtimeCredentialId)
		}

		err = ValidateRuntime(ctx, runtime.RuntimeId, runtime.Zone, runtimeCredential, false)
		if err != nil {
			if gerr.IsGRPCError(err) {
				return nil, err
			} else {
				return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorValidateFailed)
			}
		}

		attributes[constants.ColumnRuntimeCredentialId] = runtimeCredentialId
	}

	_, err = pi.Global().DB(ctx).
		Update(constants.TableRuntime).
		SetMap(attributes).
		Where(db.Eq(constants.ColumnRuntimeId, runtimeId)).
		Exec()
	if err != nil {
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorModifyResourceFailed, runtimeId)
	}

	res := &pb.ModifyRuntimeResponse{
		RuntimeId: req.GetRuntimeId(),
	}

	return res, nil
}

func (p *Server) DeleteRuntimes(ctx context.Context, req *pb.DeleteRuntimesRequest) (*pb.DeleteRuntimesResponse, error) {
	runtimeIds := req.GetRuntimeId()
	runtimes, err := CheckRuntimesPermission(ctx, runtimeIds)
	if err != nil {
		return nil, err
	}

	for _, runtime := range runtimes {
		if runtime.Status == constants.StatusDeleted {
			logger.Error(ctx, "Runtime [%s] has been deleted", runtime.RuntimeId)
			return nil, gerr.NewWithDetail(ctx, gerr.PermissionDenied, err, gerr.ErrorResourceAlreadyDeleted, runtime.RuntimeId)
		}
	}

	// There can be no cluster in the runtimes
	clusterClient, err := clusterclient.NewClient()
	if err != nil {
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorDeleteResourcesFailed)
	}
	clusters, err := clusterClient.DescribeClusters(ctx, &pb.DescribeClustersRequest{
		RuntimeId: runtimeIds,
		Status: []string{
			constants.StatusActive,
			constants.StatusStopped,
			constants.StatusSuspended,
			constants.StatusPending,
		},
	})
	if err != nil {
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorDeleteResourcesFailed)
	}

	if clusters.TotalCount > 0 {
		err = fmt.Errorf("there are still [%d] clusters in the runtime", clusters.TotalCount)
		return nil, gerr.NewWithDetail(ctx, gerr.PermissionDenied, err, gerr.ErrorDeleteResourcesFailed)
	}

	// deleted runtimes
	_, err = pi.Global().DB(ctx).
		Update(constants.TableRuntime).
		Set(constants.ColumnStatus, constants.StatusDeleted).
		Set(constants.ColumnStatusTime, time.Now()).
		Where(db.Eq(constants.ColumnRuntimeId, runtimeIds)).
		Exec()
	if err != nil {
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorDeleteResourcesFailed)
	}

	res := &pb.DeleteRuntimesResponse{
		RuntimeId: runtimeIds,
	}
	return res, nil
}

func (p *Server) CreateRuntimeCredential(ctx context.Context, req *pb.CreateRuntimeCredentialRequest) (*pb.CreateRuntimeCredentialResponse, error) {
	s := senderutil.GetSenderFromContext(ctx)

	runtimeCredentialId := models.NewRuntimeCredentialId()
	provider := req.Provider.GetValue()
	runtimeUrl := req.RuntimeUrl.GetValue()
	runtimeCredentialContent := req.RuntimeCredentialContent.GetValue()

	runtimeCredential := &models.RuntimeCredential{
		RuntimeUrl:               runtimeUrl,
		Provider:                 provider,
		RuntimeCredentialContent: runtimeCredentialContent,
	}
	err := ValidateRuntime(ctx, "", "", runtimeCredential, false)
	if err != nil {
		if gerr.IsGRPCError(err) {
			return nil, err
		} else {
			return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorValidateFailed)
		}
	}

	content, err := decodeRuntimeCredentialContent(provider, runtimeCredentialContent)
	if err != nil {
		return nil, gerr.NewWithDetail(ctx, gerr.InvalidArgument, err, gerr.ErrorCreateResourcesFailed)
	}

	newRuntimeCredential := models.NewRuntimeCredential(
		runtimeCredentialId,
		req.GetName().GetValue(),
		req.GetDescription().GetValue(),
		req.GetRuntimeUrl().GetValue(),
		content,
		req.GetProvider().GetValue(),
		s.UserId,
	)

	_, err = pi.Global().DB(ctx).
		InsertInto(constants.TableRuntimeCredential).
		Record(newRuntimeCredential).
		Exec()
	if err != nil {
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorCreateResourcesFailed)
	}

	res := &pb.CreateRuntimeCredentialResponse{
		RuntimeCredentialId: pbutil.ToProtoString(newRuntimeCredential.RuntimeCredentialId),
	}
	return res, nil
}

func (p *Server) DescribeRuntimeCredentials(ctx context.Context, req *pb.DescribeRuntimeCredentialsRequest) (*pb.DescribeRuntimeCredentialsResponse, error) {
	offset := pbutil.GetOffsetFromRequest(req)
	limit := pbutil.GetLimitFromRequest(req)

	var runtimeCredentials []*models.RuntimeCredential
	var count uint32
	query := pi.Global().DB(ctx).
		Select(models.RuntimeCredentialColumns...).
		From(constants.TableRuntimeCredential).
		Offset(offset).
		Limit(limit).
		Where(manager.BuildFilterConditions(req, constants.TableRuntimeCredential))

	query = manager.AddQueryOrderDir(query, req, constants.ColumnCreateTime)
	_, err := query.Load(&runtimeCredentials)
	if err != nil {
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorDescribeResourcesFailed)
	}

	count, err = query.Count()
	if err != nil {
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorDescribeResourcesFailed)
	}

	// encrypt runtime credential contents
	for _, runtimeCredential := range runtimeCredentials {
		runtimeCredential.RuntimeCredentialContent = "******"
	}

	res := &pb.DescribeRuntimeCredentialsResponse{
		RuntimeCredentialSet: models.RuntimeCredentialToPbs(runtimeCredentials),
		TotalCount:           count,
	}
	return res, nil
}

func (p *Server) DeleteRuntimeCredentials(ctx context.Context, req *pb.DeleteRuntimeCredentialsRequest) (*pb.DeleteRuntimeCredentialsResponse, error) {
	runtimeCredentialIds := req.GetRuntimeCredentialId()
	runtimeCredentials, err := CheckRuntimeCredentialsPermission(ctx, runtimeCredentialIds)
	if err != nil {
		return nil, err
	}

	for _, runtimeCredential := range runtimeCredentials {
		if runtimeCredential.Status == constants.StatusDeleted {
			logger.Info(ctx, "Runtime credential [%s] has been deleted", runtimeCredential.RuntimeCredentialId)
			continue
		}

		count, err := pi.Global().DB(ctx).
			Select(models.RuntimeCredentialColumns...).
			From(constants.TableRuntime).
			Where(db.Eq(constants.ColumnRuntimeCredentialId, runtimeCredential.RuntimeCredentialId)).
			Where(db.Eq(constants.ColumnStatus, constants.StatusActive)).
			Count()
		if err != nil {
			return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorDeleteResourcesFailed)
		}
		if count > 0 {
			err = fmt.Errorf("there are still [%d] runtimes use credential [%s]", count, runtimeCredential.RuntimeCredentialId)
			return nil, gerr.NewWithDetail(ctx, gerr.PermissionDenied, err, gerr.ErrorDeleteResourcesFailed)
		}
	}

	_, err = pi.Global().DB(ctx).
		Update(constants.TableRuntimeCredential).
		Set(constants.ColumnStatus, constants.StatusDeleted).
		Set(constants.ColumnStatusTime, time.Now()).
		Where(db.Eq(constants.ColumnRuntimeCredentialId, runtimeCredentialIds)).
		Exec()
	if err != nil {
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorDeleteResourcesFailed)
	}

	res := &pb.DeleteRuntimeCredentialsResponse{
		RuntimeCredentialId: runtimeCredentialIds,
	}
	return res, nil
}

func (p *Server) ModifyRuntimeCredential(ctx context.Context, req *pb.ModifyRuntimeCredentialRequest) (*pb.ModifyRuntimeCredentialResponse, error) {
	runtimeCredentialId := req.GetRuntimeCredentialId().GetValue()
	runtimeCredentialContent := req.GetRuntimeCredentialContent().GetValue()

	runtimeCredential, err := CheckRuntimeCredentialPermission(ctx, runtimeCredentialId)
	if err != nil {
		return nil, err
	}

	if runtimeCredential.Status == constants.StatusDeleted {
		logger.Error(ctx, "Runtime credential [%s] has been deleted", runtimeCredentialId)
		return nil, gerr.NewWithDetail(ctx, gerr.PermissionDenied, err, gerr.ErrorResourceAlreadyDeleted, runtimeCredentialId)
	}

	attributes := manager.BuildUpdateAttributes(req, constants.ColumnName, constants.ColumnDescription)
	attributes[constants.ColumnStatusTime] = time.Now()

	if len(runtimeCredentialContent) > 0 {
		var runtimes []*models.Runtime
		query := pi.Global().DB(ctx).
			Select(models.RuntimeColumns...).
			From(constants.TableRuntime).
			Where(db.Eq(constants.ColumnRuntimeCredentialId, runtimeCredentialId)).
			Where(db.Eq(constants.ColumnStatus, constants.StatusActive))
		_, err = query.Load(&runtimes)
		if err != nil {
			return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorModifyResourceFailed, runtimeCredentialId)
		}

		for _, runtime := range runtimes {
			newRuntimeCredential := &models.RuntimeCredential{
				RuntimeUrl:               runtimeCredential.RuntimeUrl,
				RuntimeCredentialContent: runtimeCredentialContent,
				Provider:                 runtimeCredential.Provider,
			}
			err = ValidateRuntime(ctx, runtime.RuntimeId, runtime.Zone, newRuntimeCredential, false)
			if err != nil {
				if gerr.IsGRPCError(err) {
					return nil, err
				} else {
					return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorValidateFailed)
				}
			}
		}

		newContent, err := decodeRuntimeCredentialContent(runtimeCredential.Provider, runtimeCredentialContent)
		if err != nil {
			return nil, gerr.NewWithDetail(ctx, gerr.InvalidArgument, err, gerr.ErrorModifyResourceFailed, runtimeCredentialId)
		}
		attributes[constants.ColumnRuntimeCredentialContent] = newContent
	}

	_, err = pi.Global().DB(ctx).
		Update(constants.TableRuntimeCredential).
		SetMap(attributes).
		Where(db.Eq(constants.ColumnRuntimeCredentialId, runtimeCredentialId)).
		Exec()
	if err != nil {
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorModifyResourceFailed, runtimeCredentialId)
	}

	res := &pb.ModifyRuntimeCredentialResponse{
		RuntimeCredentialId: req.GetRuntimeCredentialId(),
	}

	return res, nil
}

func (p *Server) ValidateRuntimeCredential(ctx context.Context, req *pb.ValidateRuntimeCredentialRequest) (*pb.ValidateRuntimeCredentialResponse, error) {
	runtimeCredential := &models.RuntimeCredential{
		RuntimeUrl:               req.GetRuntimeUrl().GetValue(),
		RuntimeCredentialContent: req.GetRuntimeCredentialContent().GetValue(),
		Provider:                 req.GetProvider().GetValue(),
	}
	err := ValidateRuntime(ctx, "", "", runtimeCredential, false)
	if err != nil {
		if gerr.IsGRPCError(err) {
			return &pb.ValidateRuntimeCredentialResponse{
				Ok: pbutil.ToProtoBool(false),
			}, err
		} else {
			return &pb.ValidateRuntimeCredentialResponse{
				Ok: pbutil.ToProtoBool(false),
			}, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorValidateFailed)
		}
	}
	return &pb.ValidateRuntimeCredentialResponse{
		Ok: pbutil.ToProtoBool(true),
	}, nil
}

func (p *Server) DescribeRuntimeProviderZones(ctx context.Context, req *pb.DescribeRuntimeProviderZonesRequest) (*pb.DescribeRuntimeProviderZonesResponse, error) {
	runtimeCredentialId := req.GetRuntimeCredentialId().GetValue()
	runtimeCredential, err := CheckRuntimeCredentialPermission(ctx, runtimeCredentialId)
	if err != nil {
		return nil, err
	}

	providerInterface, err := plugins.GetProviderPlugin(ctx, runtimeCredential.Provider)
	if err != nil {
		return nil, gerr.NewWithDetail(ctx, gerr.NotFound, err, gerr.ErrorProviderNotFound, runtimeCredential.Provider)
	}
	zones, err := providerInterface.DescribeRuntimeProviderZones(ctx, runtimeCredential)
	if err != nil {
		return nil, gerr.NewWithDetail(ctx, gerr.PermissionDenied, err, gerr.ErrorDescribeResourceFailed)
	}
	return &pb.DescribeRuntimeProviderZonesResponse{
		RuntimeCredentialId: req.GetRuntimeCredentialId(),
		Zone:                zones,
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
