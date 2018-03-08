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
	"openpitrix.io/openpitrix/pkg/logger"
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
	s := sender.GetSenderFromContext(ctx)
	logger.Infof("Got sender: %+v", s)
	logger.Debugf("Got req: %+v", req)
	var repos []*models.Repo
	offset := utils.GetOffsetFromRequest(req)
	limit := utils.GetLimitFromRequest(req)

	logger.Debugf("Filter condition is: %+v", manager.BuildFilterConditions(req, models.RepoTableName))

	query := p.Db.
		Select(models.RepoColumns...).
		From(models.RepoTableName).
		Offset(offset).
		Limit(limit).
		Where(manager.BuildFilterConditions(req, models.RepoTableName))
	_, err := query.Load(&repos)
	if err != nil {
		// TODO: err_code should be implementation
		return nil, status.Errorf(codes.Internal, "DescribeRepos: %+v", err)
	}
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
	// TODO: validate CreateRepoRequest
	s := sender.GetSenderFromContext(ctx)
	newRepo := models.NewRepo(
		req.GetName().GetValue(),
		req.GetDescription().GetValue(),
		req.GetUrl().GetValue(),
		req.GetCredential().GetValue(),
		req.GetVisibility().GetValue(),
		s.UserId)

	_, err := p.Db.
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
	// TODO: check resource permission
	repoId := req.GetRepoId().GetValue()
	repo, err := p.getRepo(repoId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to get repo [%s]", repoId)
	}

	attributes := manager.BuildUpdateAttributes(req,
		"repo_id", "name", "description", "url",
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
