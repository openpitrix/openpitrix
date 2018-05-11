// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package repo

import (
	"context"
	neturl "net/url"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	clientutil "openpitrix.io/openpitrix/pkg/client"
	indexerclient "openpitrix.io/openpitrix/pkg/client/repo_indexer"
	"openpitrix.io/openpitrix/pkg/constants"
	"openpitrix.io/openpitrix/pkg/db"
	"openpitrix.io/openpitrix/pkg/logger"
	"openpitrix.io/openpitrix/pkg/manager"
	"openpitrix.io/openpitrix/pkg/models"
	"openpitrix.io/openpitrix/pkg/pb"
	"openpitrix.io/openpitrix/pkg/util/pbutil"
	"openpitrix.io/openpitrix/pkg/util/senderutil"
	"openpitrix.io/openpitrix/pkg/util/stringutil"
)

func (p *Server) DescribeRepos(ctx context.Context, req *pb.DescribeReposRequest) (*pb.DescribeReposResponse, error) {
	// TODO: validate params
	var repos []*models.Repo
	offset := pbutil.GetOffsetFromRequest(req)
	limit := pbutil.GetLimitFromRequest(req)

	labelMap, err := neturl.ParseQuery(req.GetLabel().GetValue())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "DescribeRepos: GetLabelMapFromRequest %+v", err)
	}

	selectorMap, err := neturl.ParseQuery(req.GetSelector().GetValue())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "DescribeRepos: GetSelectorMapFromRequest %+v", err)
	}

	query := p.Db.
		Select(models.RepoColumnsWithTablePrefix...).
		From(models.RepoTableName).
		Offset(offset).
		Limit(limit).
		Where(manager.BuildFilterConditionsWithPrefix(req, models.RepoTableName))

	query = manager.AddQueryJoinWithMap(query, models.RepoTableName, models.RepoLabelTableName, models.ColumnRepoId,
		models.ColumnLabelKey, models.ColumnLabelValue, labelMap)
	query = manager.AddQueryJoinWithMap(query, models.RepoTableName, models.RepoSelectorTableName, models.ColumnRepoId,
		models.ColumnSelectorKey, models.ColumnSelectorValue, selectorMap)
	query = query.Distinct()

	_, err = query.Load(&repos)
	if err != nil {
		// TODO: err_code should be implementation
		return nil, status.Errorf(codes.Internal, "DescribeRepos: %+v", err)
	}

	count, err := query.Count()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "DescribeRepos: %+v", err)
	}

	repoSet, err := p.formatRepoSet(repos)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "DescribeRepos: %+v", err)
	}

	res := &pb.DescribeReposResponse{
		RepoSet:    repoSet,
		TotalCount: count,
	}
	return res, nil
}

func (p *Server) CreateRepo(ctx context.Context, req *pb.CreateRepoRequest) (*pb.CreateRepoResponse, error) {
	// TODO: common validate
	repoType := req.GetType().GetValue()
	url := req.GetUrl().GetValue()
	credential := req.GetCredential().GetValue()
	visibility := req.GetVisibility().GetValue()
	providers := req.GetProviders()

	err := validate(repoType, url, credential, visibility)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "CreateRepo: Validate failed, %+v", err)
	}

	s := senderutil.GetSenderFromContext(ctx)
	newRepo := models.NewRepo(
		req.GetName().GetValue(),
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
		return nil, status.Errorf(codes.Internal, "CreateRepo: %+v", err)
	}

	err = p.createProviders(newRepo.RepoId, providers)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "CreateRepo: %+v", err)
	}
	if len(req.GetLabels()) > 0 {
		err = p.createLabels(newRepo.RepoId, req.GetLabels())
		if err != nil {
			return nil, status.Errorf(codes.Internal, "CreateRepo: %+v", err)
		}
	}
	if len(req.GetSelectors()) > 0 {
		err = p.createSelectors(newRepo.RepoId, req.GetSelectors())
		if err != nil {
			return nil, status.Errorf(codes.Internal, "CreateRepo: %+v", err)
		}
	}

	repo, err := p.formatRepo(newRepo)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "CreateRepo: %+v", err)
	}

	res := &pb.CreateRepoResponse{
		Repo: repo,
	}

	ctx = clientutil.GetSystemUserContext()
	repoIndexerClient, err := indexerclient.NewRepoIndexerClient(ctx)
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
		return nil, status.Errorf(codes.Internal, "Failed to get repo [%s]", repoId)
	}
	url := repo.Url
	credential := repo.Credential
	visibility := repo.Visibility
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
		visibility = req.GetVisibility().GetValue()
		needValidate = true
	}
	if needValidate {
		err = validate(repoType, url, credential, visibility)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "ModifyRepo: Validate failed, %+v", err)
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
			return nil, status.Errorf(codes.Internal, "ModifyRepo: %+v", err)
		}
	}

	if len(providers) > 0 {
		providers = stringutil.Unique(providers)
		err = p.modifyProviders(repoId, providers)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "ModifyRepo: %+v", err)
		}
	}
	if len(req.GetLabels()) > 0 {
		err = p.modifyLabels(repoId, req.GetLabels())
		if err != nil {
			return nil, status.Errorf(codes.Internal, "ModifyRepo: %+v", err)
		}
	}
	if len(req.GetSelectors()) > 0 {
		err = p.modifySelectors(repoId, req.GetSelectors())
		if err != nil {
			return nil, status.Errorf(codes.Internal, "ModifyRepo: %+v", err)
		}
	}

	repo, err = p.getRepo(repoId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to get repo [%s]", repoId)
	}
	pbRepo, err := p.formatRepo(repo)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "ModifyRepo: %+v", repoId)
	}

	res := &pb.ModifyRepoResponse{
		Repo: pbRepo,
	}
	return res, nil
}

func (p *Server) DeleteRepos(ctx context.Context, req *pb.DeleteReposRequest) (*pb.DeleteReposResponse, error) {
	// TODO: check resource permission
	err := manager.CheckParamsRequired(req, "repo_id")
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	repoIds := req.GetRepoId()

	_, err = p.Db.
		Update(models.RepoTableName).
		Set(models.ColumnStatus, constants.StatusDeleted).
		Where(db.Eq(models.ColumnRepoId, repoIds)).
		Exec()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "DeleteRepos: %+v", err)
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
	visibility := "public"

	err := validate(repoType, url, credential, visibility)
	if err != nil {
		return &pb.ValidateRepoResponse{
			Ok: pbutil.ToProtoBool(false),
		}, nil
	}

	return &pb.ValidateRepoResponse{
		Ok: pbutil.ToProtoBool(true),
	}, nil
}
