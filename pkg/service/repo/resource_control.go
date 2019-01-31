// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package repo

import (
	"context"

	"openpitrix.io/openpitrix/pkg/constants"
	"openpitrix.io/openpitrix/pkg/db"
	"openpitrix.io/openpitrix/pkg/gerr"
	"openpitrix.io/openpitrix/pkg/models"
	"openpitrix.io/openpitrix/pkg/pb"
	"openpitrix.io/openpitrix/pkg/pi"
	"openpitrix.io/openpitrix/pkg/service/category/categoryutil"
	"openpitrix.io/openpitrix/pkg/util/stringutil"
)

func (p *Server) validateRepoName(ctx context.Context, name string) error {
	if len(name) == 0 {
		return gerr.New(ctx, gerr.InvalidArgument, gerr.ErrorParameterShouldNotBeEmpty, "name")
	}

	var repo []*models.Repo
	_, err := pi.Global().DB(ctx).
		Select(models.RepoColumns...).
		From(constants.TableRepo).
		Where(db.Eq(constants.ColumnName, name)).
		Where(db.Eq(constants.ColumnStatus, constants.StatusActive)).
		Limit(1).
		Load(&repo)
	if err != nil {
		return gerr.New(ctx, gerr.Internal, gerr.ErrorDescribeResourcesFailed)
	}

	if len(repo) > 0 {
		return gerr.New(ctx, gerr.Internal, gerr.ErrorConflictRepoName, name)
	}

	return nil
}

func (p *Server) createProviders(ctx context.Context, repoId string, providers []string) error {
	if len(providers) == 0 {
		return nil
	}
	insert := pi.Global().DB(ctx).InsertInto(constants.TableRepoProvider)
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

func (p *Server) deleteProviders(ctx context.Context, repoId string, providers []string) error {
	if len(providers) == 0 {
		return nil
	}
	_, err := pi.Global().DB(ctx).
		DeleteFrom(constants.TableRepoProvider).
		Where(db.Eq(constants.ColumnRepoId, repoId)).
		Where(db.Eq(constants.ColumnProvider, providers)).
		Exec()
	return err
}

func (p *Server) modifyProviders(ctx context.Context, repoId string, providers []string) error {
	providersMap, err := p.getProvidersMap(ctx, []string{repoId})
	if err != nil {
		return err
	}

	var currentProviders []string
	for _, repoProvider := range providersMap[repoId] {
		currentProviders = append(currentProviders, repoProvider.Provider)
	}
	deleteProviders := stringutil.Diff(currentProviders, providers)
	addProviders := stringutil.Diff(providers, currentProviders)
	err = p.createProviders(ctx, repoId, addProviders)
	if err != nil {
		return err
	}
	err = p.deleteProviders(ctx, repoId, deleteProviders)
	return err
}

func (p *Server) modifySelectors(ctx context.Context, repoId string, selectors []*pb.RepoSelector) error {
	selectorsMap, err := p.getSelectorsMap(ctx, []string{repoId})
	if err != nil {
		return err
	}
	currentSelectors := selectorsMap[repoId]
	currentSelectorsLength := len(currentSelectors)
	selectorsLength := len(selectors)
	firstSelector := selectors[0]
	if selectorsLength == 1 &&
		firstSelector.GetSelectorValue().GetValue() == "" &&
		firstSelector.GetSelectorKey().GetValue() == "" {
		_, err = pi.Global().DB(ctx).
			DeleteFrom(constants.TableRepoSelector).
			Where(db.Eq(constants.ColumnRepoId, repoId)).
			Exec()
		return err
	}
	// create new selectors
	if selectorsLength > currentSelectorsLength {
		err = p.createSelectors(ctx, repoId, selectors[currentSelectorsLength:])
		if err != nil {
			return err
		}
	}
	// update current selectors
	for i, currentSelector := range currentSelectors {
		var err error
		whereCondition := db.Eq(constants.ColumnRepoSelectorId, currentSelector.RepoSelectorId)
		if i+1 <= selectorsLength {
			// if current selectors exist, update it to new key/value
			targetSelector := selectors[i]
			_, err = pi.Global().DB(ctx).
				Update(constants.TableRepoSelector).
				Set(constants.ColumnSelectorKey, targetSelector.GetSelectorKey().GetValue()).
				Set(constants.ColumnSelectorValue, targetSelector.GetSelectorValue().GetValue()).
				Where(whereCondition).
				Exec()
		} else {
			// if current selectors more than arguments, delete it
			_, err = pi.Global().DB(ctx).
				DeleteFrom(constants.TableRepoSelector).
				Where(whereCondition).
				Exec()
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func (p *Server) createSelectors(ctx context.Context, repoId string, selectors []*pb.RepoSelector) error {
	if len(selectors) == 0 {
		return nil
	}
	insert := pi.Global().DB(ctx).InsertInto(constants.TableRepoSelector)
	for _, selector := range selectors {
		repoSelector := models.NewRepoSelector(repoId, selector.GetSelectorKey().GetValue(), selector.GetSelectorValue().GetValue())
		insert = insert.Record(repoSelector)
	}
	_, err := insert.Exec()
	return err
}

func (p *Server) getProvidersMap(ctx context.Context, repoIds []string) (providersMap map[string][]*models.RepoProvider, err error) {
	if len(repoIds) == 0 {
		return
	}
	var repoProviders []*models.RepoProvider
	_, err = pi.Global().DB(ctx).
		Select(models.RepoProviderColumns...).
		From(constants.TableRepoProvider).
		Where(db.Eq(constants.ColumnRepoId, repoIds)).
		Load(&repoProviders)
	if err != nil {
		return
	}
	providersMap = models.RepoProvidersMap(repoProviders)
	return
}

func (p *Server) getSelectorsMap(ctx context.Context, repoIds []string) (selectorsMap map[string][]*models.RepoSelector, err error) {
	if len(repoIds) == 0 {
		return
	}
	var repoSelectors []*models.RepoSelector
	_, err = pi.Global().DB(ctx).
		Select(models.RepoSelectorColumns...).
		From(constants.TableRepoSelector).
		Where(db.Eq(constants.ColumnRepoId, repoIds)).
		OrderDir(constants.ColumnCreateTime, true).
		Load(&repoSelectors)
	if err != nil {
		return
	}
	selectorsMap = models.RepoSelectorsMap(repoSelectors)
	return
}

func getReposLabelsMap(ctx context.Context, repoIds []string) (reposLabelsMap map[string][]*models.RepoLabel, err error) {
	if len(repoIds) == 0 {
		return
	}
	var repoLabels []*models.RepoLabel
	_, err = pi.Global().DB(ctx).
		Select(models.RepoLabelColumns...).
		From(constants.TableRepoLabel).
		Where(db.Eq(constants.ColumnRepoId, repoIds)).
		OrderDir(constants.ColumnCreateTime, true).
		Load(&repoLabels)
	if err != nil {
		return
	}
	reposLabelsMap = models.RepoLabelsMap(repoLabels)
	return
}

func (p *Server) formatRepo(ctx context.Context, repo *models.Repo) (*pb.Repo, error) {
	pbRepos, err := p.formatRepoSet(ctx, []*models.Repo{repo})
	if err != nil {
		return nil, err
	}
	return pbRepos[0], nil
}

func (p *Server) formatRepoSet(ctx context.Context, repos []*models.Repo) (pbRepos []*pb.Repo, err error) {
	pbRepos = models.ReposToPbs(repos)
	var repoIds []string
	for _, repo := range repos {
		repoIds = append(repoIds, repo.RepoId)
	}
	var providersMap map[string][]*models.RepoProvider
	providersMap, err = p.getProvidersMap(ctx, repoIds)
	if err != nil {
		return
	}

	var labelsMap map[string][]*models.RepoLabel
	labelsMap, err = getReposLabelsMap(ctx, repoIds)
	if err != nil {
		return
	}

	var selectorsMap map[string][]*models.RepoSelector
	selectorsMap, err = p.getSelectorsMap(ctx, repoIds)
	if err != nil {
		return
	}

	var rcmap map[string][]*pb.ResourceCategory
	rcmap, err = categoryutil.GetResourcesCategories(ctx, pi.Global().DB(ctx), repoIds)
	if err != nil {
		return
	}

	//sender := ctxutil.GetSender(ctx)
	for _, pbRepo := range pbRepos {
		repoId := pbRepo.GetRepoId().GetValue()
		pbRepo.Labels = models.RepoLabelsToPbs(labelsMap[repoId])
		pbRepo.Selectors = models.RepoSelectorsToPbs(selectorsMap[repoId])
		for _, p := range providersMap[repoId] {
			pbRepo.Providers = append(pbRepo.Providers, p.Provider)
		}
		if categorySet, ok := rcmap[pbRepo.GetRepoId().GetValue()]; ok {
			pbRepo.CategorySet = categorySet
		}
		//if !sender.IsGlobalAdmin() {
		//	pbRepo.Credential = pbutil.ToProtoString("{}")
		//}
	}
	return
}
