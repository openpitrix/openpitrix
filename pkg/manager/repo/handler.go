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

func (p *Server) getRepo(repoId string) (*models.Repo, error) {
	repo := &models.Repo{}
	err := p.Db.
		Select(models.RepoColumns...).
		From(models.RepoTableName).
		Where(db.Eq("repo_id", repoId)).
		LoadOne(&repo)
	if err != nil {
		return nil, err
	}
	return repo, nil
}

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

	res := &pb.DescribeReposResponse{
		RepoSet:    models.ReposToPbs(repos),
		TotalCount: count,
	}
	return res, nil
}

func (p *Server) CreateRepo(ctx context.Context, req *pb.CreateRepoRequest) (*pb.CreateRepoResponse, error) {
	repoType := req.GetType().GetValue()
	url := req.GetUrl().GetValue()
	credential := req.GetCredential().GetValue()
	visibility := req.GetVisibility().GetValue()

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

	res := &pb.CreateRepoResponse{
		Repo: models.RepoToPb(newRepo),
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

	res := &pb.ModifyRepoResponse{
		Repo: models.RepoToPb(repo),
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

	return &pb.DeleteRepoResponse{
		Repo: models.RepoToPb(repo),
	}, nil
}

func (p *Server) getRepoLabel(repoLabelId string) (*models.RepoLabel, error) {
	repoLabel := &models.RepoLabel{}
	err := p.Db.
		Select(models.RepoLabelColumns...).
		From(models.RepoLabelTableName).
		Where(db.Eq("repo_label_id", repoLabelId)).
		LoadOne(&repoLabel)
	if err != nil {
		return nil, err
	}
	return repoLabel, nil
}

func (p *Server) DescribeRepoLabels(ctx context.Context, req *pb.DescribeRepoLabelsRequest) (*pb.DescribeRepoLabelsResponse, error) {
	var repoLabels []*models.RepoLabel
	offset := utils.GetOffsetFromRequest(req)
	limit := utils.GetLimitFromRequest(req)
	query := p.Db.
		Select(models.RepoLabelColumns...).
		From(models.RepoLabelTableName).
		Offset(offset).
		Limit(limit).
		Where(manager.BuildFilterConditions(req, models.RepoLabelTableName))
	_, err := query.Load(&repoLabels)
	if err != nil {
		// TODO: err_code should be implementation
		return nil, status.Errorf(codes.Internal, "DescribeRepoLabels: %+v", err)
	}
	count, err := query.Count()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "DescribeRepoLabels: %+v", err)
	}

	res := &pb.DescribeRepoLabelsResponse{
		RepoLabelSet: models.RepoLabelsToPbs(repoLabels),
		TotalCount:   count,
	}
	return res, nil
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

	attributes := manager.BuildUpdateAttributes(req, "label_key", "label_value")
	_, err = p.Db.
		Update(models.RepoLabelTableName).
		SetMap(attributes).
		Where(db.Eq("repo_label_id", repoLabelId)).
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
	_, err := p.getRepoLabel(repoLabelId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to get repo label [%s]", repoLabelId)
	}

	_, err = p.Db.
		Update(models.RepoLabelTableName).
		Set("status", constants.StatusDeleted).
		Where(db.Eq("repo_label_id", repoLabelId)).
		Exec()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "DeleteRepoLabel: %+v", err)
	}

	repoLabel, err := p.getRepoLabel(repoLabelId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to get repo label [%s]", repoLabelId)
	}

	return &pb.DeleteRepoLabelResponse{
		RepoLabel: models.RepoLabelToPb(repoLabel),
	}, nil
}

func (p *Server) getRepoSelector(repoSelectorId string) (*models.RepoSelector, error) {
	repoSelector := &models.RepoSelector{}
	err := p.Db.
		Select(models.RepoSelectorColumns...).
		From(models.RepoSelectorTableName).
		Where(db.Eq("repo_selector_id", repoSelectorId)).
		LoadOne(&repoSelector)
	if err != nil {
		return nil, err
	}
	return repoSelector, nil
}

func (p *Server) DescribeRepoSelectors(ctx context.Context, req *pb.DescribeRepoSelectorsRequest) (*pb.DescribeRepoSelectorsResponse, error) {
	var repoSelectors []*models.RepoSelector
	offset := utils.GetOffsetFromRequest(req)
	limit := utils.GetLimitFromRequest(req)

	query := p.Db.
		Select(models.RepoSelectorColumns...).
		From(models.RepoSelectorTableName).
		Offset(offset).
		Limit(limit).
		Where(manager.BuildFilterConditions(req, models.RepoSelectorTableName))
	_, err := query.Load(&repoSelectors)
	if err != nil {
		// TODO: err_code should be implementation
		return nil, status.Errorf(codes.Internal, "DescribeRepoSelectors: %+v", err)
	}
	count, err := query.Count()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "DescribeRepoSelectors: %+v", err)
	}

	res := &pb.DescribeRepoSelectorsResponse{
		RepoSelectorSet: models.RepoSelectorsToPbs(repoSelectors),
		TotalCount:      count,
	}
	return res, nil
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

	attributes := manager.BuildUpdateAttributes(req, "selector_key", "selector_value")
	_, err = p.Db.
		Update(models.RepoSelectorTableName).
		SetMap(attributes).
		Where(db.Eq("repo_selector_id", repoSelectorId)).
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
	_, err := p.getRepoSelector(repoSelectorId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to get repo selector [%s]", repoSelectorId)
	}

	_, err = p.Db.
		Update(models.RepoSelectorTableName).
		Set("status", constants.StatusDeleted).
		Where(db.Eq("repo_selector_id", repoSelectorId)).
		Exec()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "DeleteRepoSelector: %+v", err)
	}

	repoSelector, err := p.getRepoSelector(repoSelectorId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to get repo selector [%s]", repoSelectorId)
	}

	return &pb.DeleteRepoSelectorResponse{
		RepoSelector: models.RepoSelectorToPb(repoSelector),
	}, nil
}
