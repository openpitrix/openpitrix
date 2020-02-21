// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package repo

import (
	"context"
	neturl "net/url"

	clientutil "openpitrix.io/openpitrix/pkg/client"
	indexerclient "openpitrix.io/openpitrix/pkg/client/repo_indexer"
	"openpitrix.io/openpitrix/pkg/constants"
	"openpitrix.io/openpitrix/pkg/db"
	"openpitrix.io/openpitrix/pkg/gerr"
	"openpitrix.io/openpitrix/pkg/logger"
	"openpitrix.io/openpitrix/pkg/manager"
	"openpitrix.io/openpitrix/pkg/models"
	"openpitrix.io/openpitrix/pkg/pb"
	"openpitrix.io/openpitrix/pkg/pi"
	"openpitrix.io/openpitrix/pkg/service/category/categoryutil"
	"openpitrix.io/openpitrix/pkg/util/ctxutil"
	"openpitrix.io/openpitrix/pkg/util/labelutil"
	"openpitrix.io/openpitrix/pkg/util/pbutil"
	"openpitrix.io/openpitrix/pkg/util/stringutil"
)

func (p *Server) DescribeRepos(ctx context.Context, req *pb.DescribeReposRequest) (*pb.DescribeReposResponse, error) {
	// TODO: validate params
	var repos []*models.Repo
	offset := pbutil.GetOffsetFromRequest(req)
	limit := pbutil.GetLimitFromRequest(req)
	categoryIds := req.GetCategoryId()

	labelMap, err := neturl.ParseQuery(req.GetLabel().GetValue())
	if err != nil {
		return nil, gerr.NewWithDetail(ctx, gerr.InvalidArgument, err, gerr.ErrorParameterParseFailed, "label")
	}

	selectorMap, err := neturl.ParseQuery(req.GetSelector().GetValue())
	if err != nil {
		return nil, gerr.NewWithDetail(ctx, gerr.InvalidArgument, err, gerr.ErrorParameterParseFailed, "selector")
	}

	query := pi.Global().DB(ctx).
		Select(models.RepoColumnsWithTablePrefix...).
		From(constants.TableRepo).
		Offset(offset).
		Limit(limit).
		Where(manager.BuildPermissionFilter(ctx)).
		Where(manager.BuildFilterConditionsWithPrefix(req, constants.TableRepo))

	if len(req.UserId) > 0 {
		query = query.Where(db.Or(
			db.Eq(constants.ColumnVisibility, constants.VisibilityPublic),
			db.And(
				db.Eq(constants.ColumnVisibility, constants.VisibilityPrivate),
				db.Eq(constants.ColumnOwner, req.UserId),
			),
		))
	}

	if len(categoryIds) > 0 {
		subqueryStmt := pi.Global().DB(ctx).
			Select(constants.ColumnResouceId).
			From(constants.TableCategoryResource).
			Where(db.Eq(constants.ColumnStatus, constants.StatusEnabled)).
			Where(db.Eq(constants.ColumnCategoryId, categoryIds))
		query = query.Where(db.Eq(db.WithPrefix(
			constants.TableRepo,
			constants.ColumnRepoId,
		), []*db.SelectQuery{subqueryStmt}))
	}

	query = manager.AddQueryJoinWithMap(query, constants.TableRepo, constants.TableRepoLabel, constants.ColumnRepoId,
		constants.ColumnLabelKey, constants.ColumnLabelValue, labelMap)
	query = manager.AddQueryJoinWithMap(query, constants.TableRepo, constants.TableRepoSelector, constants.ColumnRepoId,
		constants.ColumnSelectorKey, constants.ColumnSelectorValue, selectorMap)
	query = manager.AddQueryOrderDir(query, req, constants.ColumnCreateTime)
	query = query.Distinct()

	_, err = query.Load(&repos)
	if err != nil {
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorDescribeResourcesFailed)
	}

	count, err := query.Count()
	if err != nil {
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorDescribeResourcesFailed)
	}

	repoSet, err := p.formatRepoSet(ctx, repos)
	if err != nil {
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorDescribeResourcesFailed)
	}

	res := &pb.DescribeReposResponse{
		RepoSet:    repoSet,
		TotalCount: count,
	}
	return res, nil
}

func (p *Server) CreateRepo(ctx context.Context, req *pb.CreateRepoRequest) (*pb.CreateRepoResponse, error) {
	repoType := req.GetType().GetValue()
	url := req.GetUrl().GetValue()
	credential := req.GetCredential().GetValue()
	visibility := req.GetVisibility().GetValue()
	providers := req.GetProviders()

	err := validate(ctx, repoType, url, credential)
	if err != nil {
		return nil, gerr.NewWithDetail(ctx, gerr.InvalidArgument, err, gerr.ErrorValidateFailed)
	}

	name := req.GetName().GetValue()
	err = p.validateRepoName(ctx, name)
	if err != nil {
		return nil, err
	}

	s := ctxutil.GetSender(ctx)
	newRepo := models.NewRepo(
		name,
		req.GetDescription().GetValue(),
		repoType,
		url,
		credential,
		visibility,
		s.GetOwnerPath())

	newRepo.AppDefaultStatus = req.GetAppDefaultStatus().GetValue()

	_, err = pi.Global().DB(ctx).
		InsertInto(constants.TableRepo).
		Record(newRepo).
		Exec()
	if err != nil {
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorCreateResourcesFailed)
	}

	err = p.createProviders(ctx, newRepo.RepoId, providers)
	if err != nil {
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorCreateResourcesFailed)
	}
	if req.GetLabels() != nil {
		err = labelutil.SyncRepoLabels(ctx, newRepo.RepoId, req.GetLabels().GetValue())
		if err != nil {
			return nil, err
		}
	}
	if req.GetSelectors() != nil {
		err = labelutil.SyncRepoSelectors(ctx, newRepo.RepoId, req.GetSelectors().GetValue())
		if err != nil {
			return nil, err
		}
	}

	err = categoryutil.SyncResourceCategories(
		ctx,
		pi.Global().DB(ctx),
		newRepo.RepoId,
		categoryutil.DecodeCategoryIds(req.GetCategoryId().GetValue()),
	)
	if err != nil {
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorCreateResourcesFailed)
	}

	res := &pb.CreateRepoResponse{
		RepoId: pbutil.ToProtoString(newRepo.RepoId),
	}

	ctx = clientutil.SetSystemUserToContext(ctx)
	repoIndexerClient, err := indexerclient.NewRepoIndexerClient()
	if err != nil {
		logger.Warn(ctx, "Could not get repo indexer client, %+v", err)
		return res, nil
	}

	indexRequest := pb.IndexRepoRequest{
		RepoId: pbutil.ToProtoString(newRepo.RepoId),
	}
	_, err = repoIndexerClient.IndexRepo(ctx, &indexRequest)
	if err != nil {
		logger.Warn(ctx, "Call index repo service failed, %+v", err)
	}

	return res, nil
}

func (p *Server) ModifyRepo(ctx context.Context, req *pb.ModifyRepoRequest) (*pb.ModifyRepoResponse, error) {
	repoId := req.GetRepoId().GetValue()
	repo, err := CheckRepoPermission(ctx, repoId)
	if err != nil {
		return nil, err
	}

	repoType := req.GetType().GetValue()
	providers := req.GetProviders()
	url := repo.Url
	credential := repo.Credential
	needValidate := false
	if req.GetUrl() != nil {
		url = req.GetUrl().GetValue()
		needValidate = true
	}
	if req.GetCredential() != nil {
		credential = req.GetCredential().GetValue()
		needValidate = true
	}
	if req.GetVisibility() != nil {
		needValidate = true
	}
	if needValidate {
		err = validate(ctx, repoType, url, credential)
		if err != nil {
			return nil, gerr.NewWithDetail(ctx, gerr.InvalidArgument, err, gerr.ErrorValidateFailed)
		}
	}
	if req.GetName() != nil && req.GetName().GetValue() != repo.Name {
		err = p.validateRepoName(ctx, req.GetName().GetValue())
		if err != nil {
			return nil, err
		}
	}

	attributes := manager.BuildUpdateAttributes(req,
		constants.ColumnName, constants.ColumnDescription, constants.ColumnType, constants.ColumnUrl,
		constants.ColumnCredential, constants.ColumnVisibility, constants.ColumnAppDefaultStatus)
	if len(attributes) > 0 {
		_, err = pi.Global().DB(ctx).
			Update(constants.TableRepo).
			SetMap(attributes).
			Where(db.Eq(constants.ColumnRepoId, repoId)).
			Exec()
		if err != nil {
			return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorModifyResourcesFailed)
		}
	}

	if len(providers) > 0 {
		providers = stringutil.Unique(providers)
		err = p.modifyProviders(ctx, repoId, providers)
		if err != nil {
			return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorModifyResourcesFailed)
		}
	}
	if req.GetLabels() != nil {
		err = labelutil.SyncRepoLabels(ctx, repoId, req.GetLabels().GetValue())
		if err != nil {
			return nil, err
		}
	}
	if req.GetSelectors() != nil {
		err = labelutil.SyncRepoSelectors(ctx, repoId, req.GetSelectors().GetValue())
		if err != nil {
			return nil, err
		}
	}
	if req.GetCategoryId() != nil {
		err = categoryutil.SyncResourceCategories(
			ctx,
			pi.Global().DB(ctx),
			repoId,
			categoryutil.DecodeCategoryIds(req.GetCategoryId().GetValue()),
		)
		if err != nil {
			return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorCreateResourcesFailed)
		}
	}

	res := &pb.ModifyRepoResponse{
		RepoId: req.GetRepoId(),
	}
	if needValidate {
		ctx = clientutil.SetSystemUserToContext(ctx)
		repoIndexerClient, err := indexerclient.NewRepoIndexerClient()
		if err != nil {
			logger.Warn(ctx, "Could not get repo indexer client, %+v", err)
			return res, nil
		}

		indexRequest := pb.IndexRepoRequest{
			RepoId: pbutil.ToProtoString(repoId),
		}
		_, err = repoIndexerClient.IndexRepo(ctx, &indexRequest)
		if err != nil {
			logger.Warn(ctx, "Call index repo service failed, %+v", err)
		}
	}
	return res, nil
}

func (p *Server) DeleteRepos(ctx context.Context, req *pb.DeleteReposRequest) (*pb.DeleteReposResponse, error) {
	repoIds := req.GetRepoId()
	_, err := CheckReposPermission(ctx, repoIds)
	if err != nil {
		return nil, err
	}

	for _, repoId := range repoIds {
		if stringutil.StringIn(repoId, constants.InternalRepos) {
			return nil, gerr.New(ctx, gerr.InvalidArgument, gerr.ErrorCannotDeleteInternalRepo, repoId)
		}
	}

	_, err = pi.Global().DB(ctx).
		Update(constants.TableRepo).
		Set(constants.ColumnStatus, constants.StatusDeleted).
		Where(db.Eq(constants.ColumnRepoId, repoIds)).
		Exec()
	if err != nil {
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorDeleteResourcesFailed)
	}

	res := &pb.DeleteReposResponse{
		RepoId: repoIds,
	}

	ctx = clientutil.SetSystemUserToContext(ctx)
	repoIndexerClient, err := indexerclient.NewRepoIndexerClient()
	if err != nil {
		logger.Warn(ctx, "Could not get repo indexer client, %+v", err)
		return res, nil
	}

	for _, repoId := range repoIds {
		indexRequest := pb.IndexRepoRequest{
			RepoId: pbutil.ToProtoString(repoId),
		}
		_, err = repoIndexerClient.IndexRepo(ctx, &indexRequest)
		if err != nil {
			logger.Warn(ctx, "Call index repo service failed, %+v", err)
		}
	}

	return res, nil
}

func (p *Server) ValidateRepo(ctx context.Context, req *pb.ValidateRepoRequest) (*pb.ValidateRepoResponse, error) {
	// TODO: check resource permission
	repoType := req.GetType().GetValue()
	url := req.GetUrl().GetValue()
	credential := req.GetCredential().GetValue()

	err := validate(ctx, repoType, url, credential)
	if err != nil {
		e, ok := err.(*ErrorWithCode)
		if !ok {
			return &pb.ValidateRepoResponse{
				Ok:        pbutil.ToProtoBool(false),
				ErrorCode: ErrNotExpect,
			}, nil
		}

		return &pb.ValidateRepoResponse{
			Ok:        pbutil.ToProtoBool(false),
			ErrorCode: e.Code(),
		}, nil
	}

	return &pb.ValidateRepoResponse{
		Ok:        pbutil.ToProtoBool(true),
		ErrorCode: 0,
	}, nil
}
