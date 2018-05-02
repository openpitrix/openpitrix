// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package repo

import (
	"openpitrix.io/openpitrix/pkg/db"
	"openpitrix.io/openpitrix/pkg/models"
	"openpitrix.io/openpitrix/pkg/pb"
	"openpitrix.io/openpitrix/pkg/util/stringutil"
)

func (p *Server) getRepo(repoId string) (*models.Repo, error) {
	repo := &models.Repo{}
	err := p.Db.
		Select(models.RepoColumns...).
		From(models.RepoTableName).
		Where(db.Eq(models.ColumnRepoId, repoId)).
		LoadOne(&repo)
	if err != nil {
		return nil, err
	}
	return repo, nil
}

func (p *Server) createProviders(repoId string, providers []string) error {
	if len(providers) == 0 {
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
	if len(providers) == 0 {
		return nil
	}
	_, err := p.Db.
		DeleteFrom(models.RepoProviderTableName).
		Where(db.Eq(models.ColumnRepoId, repoId)).
		Where(db.Eq(models.ColumnProvider, providers)).
		Exec()
	return err
}

func (p *Server) modifyProviders(repoId string, providers []string) error {
	providersMap, err := p.getProvidersMap([]string{repoId})

	var currentProviders []string
	for _, repoProvider := range providersMap[repoId] {
		currentProviders = append(currentProviders, repoProvider.Provider)
	}
	deleteProviders := stringutil.Diff(currentProviders, providers)
	addProviders := stringutil.Diff(providers, currentProviders)
	err = p.createProviders(repoId, addProviders)
	if err != nil {
		return err
	}
	err = p.deleteProviders(repoId, deleteProviders)
	return err
}

func (p *Server) modifyLabels(repoId string, labels []*pb.RepoLabel) error {
	labelsMap, err := p.getLabelsMap([]string{repoId})
	if err != nil {
		return err
	}
	currentLabels := labelsMap[repoId]
	currentLabelsLength := len(currentLabels)
	labelsLength := len(labels)
	// create new labels
	if labelsLength > currentLabelsLength {
		err = p.createLabels(repoId, labels[currentLabelsLength:])
		if err != nil {
			return err
		}
	}
	// update current labels
	for i, currentLabel := range currentLabels {
		var err error
		whereCondition := db.Eq(models.ColumnRepoLabelId, currentLabel.RepoLabelId)
		if i+1 <= labelsLength {
			// if current label exist, update it to new key/value
			targetLabel := labels[i]
			_, err = p.Db.
				Update(models.RepoLabelTableName).
				Set(models.ColumnLabelKey, targetLabel.GetLabelKey().GetValue()).
				Set(models.ColumnLabelValue, targetLabel.GetLabelValue().GetValue()).
				Where(whereCondition).
				Exec()
		} else {
			// if current label more than arguments, delete it
			_, err = p.Db.
				DeleteFrom(models.RepoLabelTableName).
				Where(whereCondition).
				Exec()
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func (p *Server) createLabels(repoId string, labels []*pb.RepoLabel) error {
	if len(labels) == 0 {
		return nil
	}
	insert := p.Db.InsertInto(models.RepoLabelTableName).Columns(models.RepoLabelColumns...)
	for _, label := range labels {
		repoLabel := models.NewRepoLabel(repoId, label.GetLabelKey().GetValue(), label.GetLabelValue().GetValue())
		insert = insert.Record(repoLabel)
	}
	_, err := insert.Exec()
	return err
}

func (p *Server) modifySelectors(repoId string, selectors []*pb.RepoSelector) error {
	selectorsMap, err := p.getSelectorsMap([]string{repoId})
	if err != nil {
		return err
	}
	currentSelectors := selectorsMap[repoId]
	currentSelectorsLength := len(currentSelectors)
	selectorsLength := len(selectors)
	// create new selectors
	if selectorsLength > currentSelectorsLength {
		err = p.createSelectors(repoId, selectors[currentSelectorsLength:])
		if err != nil {
			return err
		}
	}
	// update current selectors
	for i, currentSelector := range currentSelectors {
		var err error
		whereCondition := db.Eq(models.ColumnRepoSelectorId, currentSelector.RepoSelectorId)
		if i+1 <= selectorsLength {
			// if current selectors exist, update it to new key/value
			targetSelector := selectors[i]
			_, err = p.Db.
				Update(models.RepoSelectorTableName).
				Set(models.ColumnSelectorKey, targetSelector.GetSelectorKey().GetValue()).
				Set(models.ColumnSelectorValue, targetSelector.GetSelectorValue().GetValue()).
				Where(whereCondition).
				Exec()
		} else {
			// if current selectors more than arguments, delete it
			_, err = p.Db.
				DeleteFrom(models.RepoSelectorTableName).
				Where(whereCondition).
				Exec()
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func (p *Server) createSelectors(repoId string, selectors []*pb.RepoSelector) error {
	if len(selectors) == 0 {
		return nil
	}
	insert := p.Db.InsertInto(models.RepoSelectorTableName).Columns(models.RepoSelectorColumns...)
	for _, selector := range selectors {
		repoSelector := models.NewRepoSelector(repoId, selector.GetSelectorKey().GetValue(), selector.GetSelectorValue().GetValue())
		insert = insert.Record(repoSelector)
	}
	_, err := insert.Exec()
	return err
}

func (p *Server) getProvidersMap(repoIds []string) (providersMap map[string][]*models.RepoProvider, err error) {
	var repoProviders []*models.RepoProvider
	_, err = p.Db.
		Select(models.RepoProviderColumns...).
		From(models.RepoProviderTableName).
		Where(db.Eq(models.ColumnRepoId, repoIds)).
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
		Where(db.Eq(models.ColumnRepoId, repoIds)).
		OrderDir(models.ColumnCreateTime, true).
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
		Where(db.Eq(models.ColumnRepoId, repoIds)).
		OrderDir(models.ColumnCreateTime, true).
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
