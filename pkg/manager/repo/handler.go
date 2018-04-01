// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package repo

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"openpitrix.io/openpitrix/pkg/constants"
	"openpitrix.io/openpitrix/pkg/db"
	"openpitrix.io/openpitrix/pkg/manager"
	"openpitrix.io/openpitrix/pkg/models"
	"openpitrix.io/openpitrix/pkg/pb"
	"openpitrix.io/openpitrix/pkg/utils"
	"openpitrix.io/openpitrix/pkg/utils/sender"
)

func (p *Server) DescribeRepos(ctx context.Context, req *pb.DescribeReposRequest) (*pb.DescribeReposResponse, error) {
	var repos []*models.Repo
	offset := utils.GetOffsetFromRequest(req)
	limit := utils.GetLimitFromRequest(req)

	labelMap, err := GetLabelMapFromRequest(req)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "DescribeRepos: GetLabelMapFromRequest: %+v", err)
	}

	selectorMap, err := GetSelectorMapFromRequest(req)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "DescribeRepos: GetSelectorMapFromRequest: %+v", err)
	}

	query := p.Db.
		Select(models.RepoColumnsWithTablePrefix...).
		From(models.RepoTableName).
		Offset(offset).
		Limit(limit).
		Where(manager.BuildFilterConditionsWithPrefix(req, models.RepoTableName))

	query = GenerateSelectQuery(query, models.RepoTableName, models.RepoLabelTableName,
		"label_key", "label_value", labelMap)
	query = GenerateSelectQuery(query, models.RepoTableName, models.RepoSelectorTableName,
		"selector_key", "selector_value", selectorMap)

	_, err = query.Load(&repos)
	if err != nil {
		// TODO: err_code should be implementation
		return nil, status.Errorf(codes.Internal, "DescribeRepos: %+v", err)
	}

	query = p.Db.
		Select(models.RepoColumnsWithTablePrefix...).
		From(models.RepoTableName).
		Offset(offset).
		Limit(limit).
		Where(manager.BuildFilterConditionsWithPrefix(req, models.RepoTableName))
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
	repoType := req.GetType().GetValue()
	url := req.GetUrl().GetValue()
	credential := req.GetCredential().GetValue()
	visibility := req.GetVisibility().GetValue()
	provider := req.GetProvider()

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

	err = p.createProvider(newRepo.RepoId, provider)
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
	url := req.GetUrl().GetValue()
	credential := req.GetCredential().GetValue()
	visibility := req.GetVisibility().GetValue()

	err := validate(repoType, url, credential, visibility)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "ModifyRepo: Validate failed, %+v", err)
	}

	// TODO: check resource permission
	repoId := req.GetRepoId().GetValue()
	repo, err := p.getRepo(repoId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to get repo [%s]", repoId)
	}

	attributes := manager.BuildUpdateAttributes(req,
		"name", "description", "type", "url",
		"credential", "visibility")
	_, err = p.Db.
		Update(models.RepoTableName).
		SetMap(attributes).
		Where(db.Eq("repo_id", repoId)).
		Exec()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "ModifyRepo: %+v", err)
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
		Set("status", constants.StatusDeleted).
		Where(db.Eq("repo_id", repoId)).
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
