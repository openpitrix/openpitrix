// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package repo

import (
	"context"
	neturl "net/url"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"openpitrix.io/openpitrix/pkg/constants"
	"openpitrix.io/openpitrix/pkg/db"
	"openpitrix.io/openpitrix/pkg/manager"
	"openpitrix.io/openpitrix/pkg/models"
	"openpitrix.io/openpitrix/pkg/pb"
	"openpitrix.io/openpitrix/pkg/utils"
	"openpitrix.io/openpitrix/pkg/utils/sender"
	"openpitrix.io/openpitrix/pkg/utils/stringutil"
)

func (p *Server) DescribeRepos(ctx context.Context, req *pb.DescribeReposRequest) (*pb.DescribeReposResponse, error) {
	// TODO: validate params
	var repos []*models.Repo
	offset := utils.GetOffsetFromRequest(req)
	limit := utils.GetLimitFromRequest(req)

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

	query = db.AddJoinFilterWithMap(query, models.RepoTableName, models.RepoLabelTableName, models.ColumnRepoId,
		models.ColumnLabelKey, models.ColumnLabelValue, labelMap)
	query = db.AddJoinFilterWithMap(query, models.RepoTableName, models.RepoSelectorTableName, models.ColumnRepoId,
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

	s := sender.GetSenderFromContext(ctx)
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
	err = p.createLabels(newRepo.RepoId, req.GetLabels().GetValue())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "CreateRepo: %+v", err)
	}
	err = p.createSelectors(newRepo.RepoId, req.GetSelectors().GetValue())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "CreateRepo: %+v", err)
	}

	repo, err := p.formatRepo(newRepo)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "CreateRepo: %+v", err)
	}

	res := &pb.CreateRepoResponse{
		Repo: repo,
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
	_, err = p.Db.
		Update(models.RepoTableName).
		SetMap(attributes).
		Where(db.Eq(models.ColumnRepoId, repoId)).
		Exec()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "ModifyRepo: %+v", err)
	}

	if len(providers) > 0 {
		providers = stringutil.Unique(providers)
		err = p.modifyProviders(repoId, providers)
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

func (p *Server) DeleteRepo(ctx context.Context, req *pb.DeleteRepoRequest) (*pb.DeleteRepoResponse, error) {
	// TODO: check resource permission
	repoId := req.GetRepoId().GetValue()
	_, err := p.getRepo(repoId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to get repo [%s]", repoId)
	}

	_, err = p.Db.
		Update(models.RepoTableName).
		Set(models.ColumnStatus, constants.StatusDeleted).
		Where(db.Eq(models.ColumnRepoId, repoId)).
		Exec()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "DeleteRepo: %+v", err)
	}

	repo, err := p.getRepo(repoId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to get repo [%s]", repoId)
	}

	pbRepo, err := p.formatRepo(repo)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "DeleteRepo: %+v", repoId)
	}

	return &pb.DeleteRepoResponse{
		Repo: pbRepo,
	}, nil
}

func (p *Server) CreateRepoLabel(ctx context.Context, req *pb.CreateRepoLabelRequest) (*pb.CreateRepoLabelResponse, error) {
	repoId := req.GetRepoId().GetValue()
	_, err := p.getRepo(repoId)
	if err != nil {
		// TODO: err_code should be implementation
		return nil, status.Errorf(codes.Internal, "Failed to get repo [%s]", repoId)
	}

	newRepoLabel := models.NewRepoLabel(
		repoId,
		req.GetLabelKey().GetValue(),
		req.GetLabelValue().GetValue(),
	)

	_, err = p.Db.
		InsertInto(models.RepoLabelTableName).
		Columns(models.RepoLabelColumns...).
		Record(newRepoLabel).
		Exec()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "CreateRepoLabel: %+v", err)
	}

	res := &pb.CreateRepoLabelResponse{
		RepoLabel: models.RepoLabelToPb(newRepoLabel),
	}
	return res, nil
}

func (p *Server) ModifyRepoLabel(ctx context.Context, req *pb.ModifyRepoLabelRequest) (*pb.ModifyRepoLabelResponse, error) {
	// TODO: check resource permission
	repoLabelId := req.GetRepoLabelId().GetValue()
	repoLabel, err := p.getRepoLabel(repoLabelId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to get repo label [%s]", repoLabelId)
	}

	attributes := manager.BuildUpdateAttributes(req, models.ColumnLabelKey, models.ColumnLabelValue)
	_, err = p.Db.
		Update(models.RepoLabelTableName).
		SetMap(attributes).
		Where(db.Eq(models.ColumnRepoLabelId, repoLabelId)).
		Exec()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "ModifyRepoLabel: %+v", err)
	}
	repoLabel, err = p.getRepoLabel(repoLabelId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to get repo label [%s]", repoLabelId)
	}

	res := &pb.ModifyRepoLabelResponse{
		RepoLabel: models.RepoLabelToPb(repoLabel),
	}
	return res, nil
}

func (p *Server) DeleteRepoLabel(ctx context.Context, req *pb.DeleteRepoLabelRequest) (*pb.DeleteRepoLabelResponse, error) {
	// TODO: check resource permission
	repoLabelId := req.GetRepoLabelId().GetValue()
	repoLabel, err := p.getRepoLabel(repoLabelId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to get repo label [%s]", repoLabelId)
	}

	_, err = p.Db.
		DeleteFrom(models.RepoLabelTableName).
		Where(db.Eq(models.ColumnRepoLabelId, repoLabelId)).
		Exec()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "DeleteRepoLabel: %+v", err)
	}

	return &pb.DeleteRepoLabelResponse{
		RepoLabel: models.RepoLabelToPb(repoLabel),
	}, nil
}

func (p *Server) CreateRepoSelector(ctx context.Context, req *pb.CreateRepoSelectorRequest) (*pb.CreateRepoSelectorResponse, error) {
	repoId := req.GetRepoId().GetValue()
	_, err := p.getRepo(repoId)
	if err != nil {
		// TODO: err_code should be implementation
		return nil, status.Errorf(codes.Internal, "Failed to get repo [%s]", repoId)
	}

	newRepoSelector := models.NewRepoSelector(
		repoId,
		req.GetSelectorKey().GetValue(),
		req.GetSelectorValue().GetValue(),
	)

	_, err = p.Db.
		InsertInto(models.RepoSelectorTableName).
		Columns(models.RepoSelectorColumns...).
		Record(newRepoSelector).
		Exec()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "CreateRepoSelector: %+v", err)
	}

	res := &pb.CreateRepoSelectorResponse{
		RepoSelector: models.RepoSelectorToPb(newRepoSelector),
	}
	return res, nil
}

func (p *Server) ModifyRepoSelector(ctx context.Context, req *pb.ModifyRepoSelectorRequest) (*pb.ModifyRepoSelectorResponse, error) {
	// TODO: check resource permission
	repoSelectorId := req.GetRepoSelectorId().GetValue()
	repoSelector, err := p.getRepoSelector(repoSelectorId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to get repo selector [%s]", repoSelectorId)
	}

	attributes := manager.BuildUpdateAttributes(req, models.ColumnSelectorKey, models.ColumnSelectorValue)
	_, err = p.Db.
		Update(models.RepoSelectorTableName).
		SetMap(attributes).
		Where(db.Eq(models.ColumnRepoSelectorId, repoSelectorId)).
		Exec()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "ModifyRepoSelector: %+v", err)
	}
	repoSelector, err = p.getRepoSelector(repoSelectorId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to get repo selector [%s]", repoSelectorId)
	}

	res := &pb.ModifyRepoSelectorResponse{
		RepoSelector: models.RepoSelectorToPb(repoSelector),
	}
	return res, nil
}

func (p *Server) DeleteRepoSelector(ctx context.Context, req *pb.DeleteRepoSelectorRequest) (*pb.DeleteRepoSelectorResponse, error) {
	// TODO: check resource permission
	repoSelectorId := req.GetRepoSelectorId().GetValue()
	repoSelector, err := p.getRepoSelector(repoSelectorId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to get repo selector [%s]", repoSelectorId)
	}

	_, err = p.Db.
		DeleteFrom(models.RepoSelectorTableName).
		Where(db.Eq(models.ColumnRepoSelectorId, repoSelectorId)).
		Exec()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "DeleteRepoSelector: %+v", err)
	}

	return &pb.DeleteRepoSelectorResponse{
		RepoSelector: models.RepoSelectorToPb(repoSelector),
	}, nil
}
