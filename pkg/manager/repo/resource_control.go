// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package repo

import (
	"net/url"

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

func (p *Server) createProviders(repoId string, providers []string) error {
	if providers == nil {
		return nil
	}
	insert := p.Db.InsertInto(models.RepoProviderTableName).Columns(models.RepoProviderColumns...)
	for _, provider := range providers {
		record := models.RepoProvider{
			RepoId:   repoId,
			Provider: provider,
		}
		insert = insert.Record(record)
	}
	_, err := insert.Exec()
	return err
}

func (p *Server) deleteProviders(repoId string, providers []string) error {
	if providers == nil {
		return nil
	}
	_, err := p.Db.
		DeleteFrom(models.RepoProviderTableName).
		Where(db.Eq("repo_id", repoId)).
		Where(db.Eq("provider", providers)).
		Exec()
	return err
}

func (p *Server) createLabels(repoId string, labels string) error {
	labelsValue, err := url.ParseQuery(labels)
	if err != nil {
		return err
	}
	insert := p.Db.InsertInto(models.RepoLabelTableName).Columns(models.RepoLabelColumns...)
	for key, values := range labelsValue {
		for _, value := range values {
			repoLabel := models.NewRepoLabel(repoId, key, value)
			insert = insert.Record(repoLabel)
		}
	}
	_, err = insert.Exec()
	return err
}

func (p *Server) createSelectors(repoId string, selectors string) error {
	selectorsValue, err := url.ParseQuery(selectors)
	if err != nil {
		return err
	}
	insert := p.Db.InsertInto(models.RepoSelectorTableName).Columns(models.RepoSelectorColumns...)
	for key, values := range selectorsValue {
		for _, value := range values {
			repoSelector := models.NewRepoSelector(repoId, key, value)
			insert = insert.Record(repoSelector)
		}
	}
	_, err = insert.Exec()
	return err
}

func (p *Server) getProvidersMap(repoIds []string) (providersMap map[string][]*models.RepoProvider, err error) {
	var repoProviders []*models.RepoProvider
	_, err = p.Db.
		Select(models.RepoProviderColumns...).
		From(models.RepoProviderTableName).
		Where(db.Eq("repo_id", repoIds)).
		Load(&repoProviders)
	if err != nil {
		return
	}
	providersMap = models.RepoProvidersMap(repoProviders)
	return
}

func (p *Server) getSelectorsMap(repoIds []string) (selectorsMap map[string][]*models.RepoSelector, err error) {
	var repoSelectors []*models.RepoSelector
	_, err = p.Db.
		Select(models.RepoSelectorColumns...).
		From(models.RepoSelectorTableName).
		Where(db.Eq("repo_id", repoIds)).
		Load(&repoSelectors)
	if err != nil {
		return
	}
	selectorsMap = models.RepoSelectorsMap(repoSelectors)
	return
}

func (p *Server) getLabelsMap(repoIds []string) (labelsMap map[string][]*models.RepoLabel, err error) {
	var repoLabels []*models.RepoLabel
	_, err = p.Db.
		Select(models.RepoLabelColumns...).
		From(models.RepoLabelTableName).
		Where(db.Eq("repo_id", repoIds)).
		Load(&repoLabels)
	if err != nil {
		return
	}
	labelsMap = models.RepoLabelsMap(repoLabels)
	return
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

	var providersMap map[string][]*models.RepoProvider
	providersMap, err = p.getProvidersMap(repoIds)
	if err != nil {
		return
	}

	var labelsMap map[string][]*models.RepoLabel
	labelsMap, err = p.getLabelsMap(repoIds)
	if err != nil {
		return
	}

	var selectorsMap map[string][]*models.RepoSelector
	selectorsMap, err = p.getSelectorsMap(repoIds)
	if err != nil {
		return
	}

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
