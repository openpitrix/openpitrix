// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package repo

import (
	"openpitrix.io/openpitrix/pkg/db"
	"openpitrix.io/openpitrix/pkg/models"
	"openpitrix.io/openpitrix/pkg/pb"
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

func (p *Server) createProvider(repoId string, provider []string) error {
	if provider == nil {
		return nil
	}
	insert := p.Db.InsertInto(models.RepoProviderTableName).Columns(models.RepoProviderColumns...)
	for _, p := range provider {
		record := models.RepoProvider{
			RepoId:   repoId,
			Provider: p,
		}
		insert = insert.Record(record)
	}
	_, err := insert.Exec()
	return err
}

func (p *Server) formatRepo(repo *models.Repo) (*pb.Repo, error) {
	pbRepos, err := p.formatRepoSet([]*models.Repo{repo})
	if err != nil {
		return nil, err
	}
	return pbRepos[0], nil
}

func (p *Server) formatRepoSet(repos []*models.Repo) (pbRepos []*pb.Repo, err error) {
	pbRepos = models.ReposToPbs(repos)
	var repoIds []string
	for _, repo := range repos {
		repoIds = append(repoIds, repo.RepoId)
	}

	var repoProviders []*models.RepoProvider
	_, err = p.Db.
		Select(models.RepoProviderColumns...).
		From(models.RepoProviderTableName).
		Where(db.Eq("repo_id", repoIds)).
		Load(&repoProviders)
	if err != nil {
		return
	}
	providersMap := models.RepoProvidersMap(repoProviders)

	var repoLabels []*models.RepoLabel
	_, err = p.Db.
		Select(models.RepoLabelColumns...).
		From(models.RepoLabelTableName).
		Where(db.Eq("repo_id", repoIds)).
		Load(&repoLabels)
	if err != nil {
		return
	}
	labelsMap := models.RepoLabelsMap(repoLabels)

	var repoSelectors []*models.RepoSelector
	_, err = p.Db.
		Select(models.RepoSelectorColumns...).
		From(models.RepoSelectorTableName).
		Where(db.Eq("repo_id", repoIds)).
		Load(&repoSelectors)
	if err != nil {
		return
	}
	selectorsMap := models.RepoSelectorsMap(repoSelectors)

	for _, pbRepo := range pbRepos {
		repoId := pbRepo.GetRepoId().GetValue()
		pbRepo.Labels = models.RepoLabelsToPbs(labelsMap[repoId])
		pbRepo.Selectors = models.RepoSelectorsToPbs(selectorsMap[repoId])
		for _, p := range providersMap[repoId] {
			pbRepo.Providers = append(pbRepo.Providers, p.Provider)
		}
	}
	return
}
