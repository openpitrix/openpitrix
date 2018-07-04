// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package repo

import (
	"context"
	neturl "net/url"
	"strings"

	clientutil "openpitrix.io/openpitrix/pkg/client"
	indexerclient "openpitrix.io/openpitrix/pkg/client/repo_indexer"
	"openpitrix.io/openpitrix/pkg/constants"
	"openpitrix.io/openpitrix/pkg/db"
	"openpitrix.io/openpitrix/pkg/gerr"
	"openpitrix.io/openpitrix/pkg/logger"
	"openpitrix.io/openpitrix/pkg/manager"
	"openpitrix.io/openpitrix/pkg/models"
	"openpitrix.io/openpitrix/pkg/pb"
	"openpitrix.io/openpitrix/pkg/service/category/categoryutil"
	"openpitrix.io/openpitrix/pkg/util/pbutil"
	"openpitrix.io/openpitrix/pkg/util/senderutil"
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
		return nil, gerr.NewWithDetail(gerr.InvalidArgument, err, gerr.ErrorParameterParseFailed, "label")
	}

	selectorMap, err := neturl.ParseQuery(req.GetSelector().GetValue())
	if err != nil {
		return nil, gerr.NewWithDetail(gerr.InvalidArgument, err, gerr.ErrorParameterParseFailed, "selector")
	}

	query := p.Db.
		Select(models.RepoColumnsWithTablePrefix...).
		From(models.RepoTableName).
		Offset(offset).
		Limit(limit).
		Where(manager.BuildFilterConditionsWithPrefix(req, models.RepoTableName))

	if len(categoryIds) > 0 {
		subqueryStmt := p.Db.
			Select(models.ColumnResouceId).
			From(models.CategoryResourceTableName).
			Where(db.Eq(models.ColumnCategoryId, categoryIds))
		query = query.Where(db.Eq(models.ColumnAppId, []*db.SelectQuery{subqueryStmt}))
	}

	query = manager.AddQueryJoinWithMap(query, models.RepoTableName, models.RepoLabelTableName, models.ColumnRepoId,
		models.ColumnLabelKey, models.ColumnLabelValue, labelMap)
	query = manager.AddQueryJoinWithMap(query, models.RepoTableName, models.RepoSelectorTableName, models.ColumnRepoId,
		models.ColumnSelectorKey, models.ColumnSelectorValue, selectorMap)
	query = manager.AddQueryOrderDir(query, req, models.ColumnCreateTime)
	query = query.Distinct()

	_, err = query.Load(&repos)
	if err != nil {
		// TODO: err_code should be implementation
		return nil, gerr.NewWithDetail(gerr.Internal, err, gerr.ErrorDescribeResourcesFailed)
	}

	count, err := query.Count()
	if err != nil {
		return nil, gerr.NewWithDetail(gerr.Internal, err, gerr.ErrorDescribeResourcesFailed)
	}

	repoSet, err := p.formatRepoSet(repos)
	if err != nil {
		return nil, gerr.NewWithDetail(gerr.Internal, err, gerr.ErrorDescribeResourcesFailed)
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

	err := validate(repoType, url, credential, providers)
	if err != nil {
		return nil, gerr.NewWithDetail(gerr.InvalidArgument, err, gerr.ErrorValidateFailed)
	}

	name := req.GetName().GetValue()
	err = p.validateRepoName(name)
	if err != nil {
		return nil, err
	}

	s := senderutil.GetSenderFromContext(ctx)
	newRepo := models.NewRepo(
		name,
		req.GetDescription().GetValue(),
		repoType,
		url,
		credential,
		visibility,
		s.UserId)

	_, err = p.Db.
		InsertInto(models.RepoTableName).
		Columns(models.RepoColumns...).
		Record(newRepo).
		Exec()
	if err != nil {
		return nil, gerr.NewWithDetail(gerr.Internal, err, gerr.ErrorCreateResourcesFailed)
	}

	err = p.createProviders(newRepo.RepoId, providers)
	if err != nil {
		return nil, gerr.NewWithDetail(gerr.Internal, err, gerr.ErrorCreateResourcesFailed)
	}
	if len(req.GetLabels()) > 0 {
		err = p.createLabels(newRepo.RepoId, req.GetLabels())
		if err != nil {
			return nil, gerr.NewWithDetail(gerr.Internal, err, gerr.ErrorCreateResourcesFailed)
		}
	}
	if len(req.GetSelectors()) > 0 {
		err = p.createSelectors(newRepo.RepoId, req.GetSelectors())
		if err != nil {
			return nil, gerr.NewWithDetail(gerr.Internal, err, gerr.ErrorCreateResourcesFailed)
		}
	}

	err = categoryutil.SyncResourceCategories(p.Db, newRepo.RepoId, strings.Split(req.GetCategoryId().GetValue(), ","))
	if err != nil {
		return nil, gerr.NewWithDetail(gerr.Internal, err, gerr.ErrorCreateResourcesFailed)
	}

	res := &pb.CreateRepoResponse{
		RepoId: pbutil.ToProtoString(newRepo.RepoId),
	}

	ctx = clientutil.GetSystemUserContext()
	repoIndexerClient, err := indexerclient.NewRepoIndexerClient()
	if err != nil {
		logger.Warn("Could not get repo indexer client, %+v", err)
		return res, nil
	}

	indexRequest := pb.IndexRepoRequest{
		RepoId: pbutil.ToProtoString(newRepo.RepoId),
	}
	_, err = repoIndexerClient.IndexRepo(ctx, &indexRequest)
	if err != nil {
		logger.Warn("Call index repo service failed, %+v", err)
	}

	return res, nil
}

func (p *Server) ModifyRepo(ctx context.Context, req *pb.ModifyRepoRequest) (*pb.ModifyRepoResponse, error) {
	repoType := req.GetType().GetValue()
	providers := req.GetProviders()
	// TODO: check resource permission
	repoId := req.GetRepoId().GetValue()
	repo, err := p.getRepo(repoId)
	if err != nil {
		return nil, gerr.NewWithDetail(gerr.InvalidArgument, err, gerr.ErrorResourceNotFound, repoId)
	}
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
		err = validate(repoType, url, credential, providers)
		if err != nil {
			return nil, gerr.NewWithDetail(gerr.InvalidArgument, err, gerr.ErrorValidateFailed)
		}
	}
	if req.GetName() != nil && req.GetName().GetValue() != repo.Name {
		err = p.validateRepoName(req.GetName().GetValue())
		if err != nil {
			return nil, err
		}
	}

	attributes := manager.BuildUpdateAttributes(req,
		models.ColumnName, models.ColumnDescription, models.ColumnType, models.ColumnUrl,
		models.ColumnCredential, models.ColumnVisibility)
	if len(attributes) > 0 {
		_, err = p.Db.
			Update(models.RepoTableName).
			SetMap(attributes).
			Where(db.Eq(models.ColumnRepoId, repoId)).
			Exec()
		if err != nil {
			return nil, gerr.NewWithDetail(gerr.Internal, err, gerr.ErrorModifyResourcesFailed)
		}
	}

	if len(providers) > 0 {
		providers = stringutil.Unique(providers)
		err = p.modifyProviders(repoId, providers)
		if err != nil {
			return nil, gerr.NewWithDetail(gerr.Internal, err, gerr.ErrorModifyResourcesFailed)
		}
	}
	if len(req.GetLabels()) > 0 {
		err = p.modifyLabels(repoId, req.GetLabels())
		if err != nil {
			return nil, gerr.NewWithDetail(gerr.Internal, err, gerr.ErrorModifyResourcesFailed)
		}
	}
	if len(req.GetSelectors()) > 0 {
		err = p.modifySelectors(repoId, req.GetSelectors())
		if err != nil {
			return nil, gerr.NewWithDetail(gerr.Internal, err, gerr.ErrorModifyResourcesFailed)
		}
	}
	err = categoryutil.SyncResourceCategories(p.Db, repoId, strings.Split(req.GetCategoryId().GetValue(), ","))
	if err != nil {
		return nil, gerr.NewWithDetail(gerr.Internal, err, gerr.ErrorCreateResourcesFailed)
	}

	res := &pb.ModifyRepoResponse{
		RepoId: req.GetRepoId(),
	}
	return res, nil
}

func (p *Server) DeleteRepos(ctx context.Context, req *pb.DeleteReposRequest) (*pb.DeleteReposResponse, error) {
	// TODO: check resource permission
	repoIds := req.GetRepoId()

	_, err := p.Db.
		Update(models.RepoTableName).
		Set(models.ColumnStatus, constants.StatusDeleted).
		Where(db.Eq(models.ColumnRepoId, repoIds)).
		Exec()
	if err != nil {
		return nil, gerr.NewWithDetail(gerr.Internal, err, gerr.ErrorDeleteResourcesFailed)
	}

	return &pb.DeleteReposResponse{
		RepoId: repoIds,
	}, nil
}

func (p *Server) ValidateRepo(ctx context.Context, req *pb.ValidateRepoRequest) (*pb.ValidateRepoResponse, error) {
	// TODO: check resource permission
	repoType := req.GetType().GetValue()
	url := req.GetUrl().GetValue()
	credential := req.GetCredential().GetValue()
	providers := []string{"qingcloud"}

	err := validate(repoType, url, credential, providers)
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
