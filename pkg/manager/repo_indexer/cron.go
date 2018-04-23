// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package repo_indexer

import (
	"github.com/robfig/cron"

	"openpitrix.io/openpitrix/pkg/client"
	repoClient "openpitrix.io/openpitrix/pkg/client/repo"
	"openpitrix.io/openpitrix/pkg/config"
	"openpitrix.io/openpitrix/pkg/constants"
	"openpitrix.io/openpitrix/pkg/logger"
	"openpitrix.io/openpitrix/pkg/pb"
)

type repoInfos map[string]string // repoId & owner

func getRepos() (repoInfos, error) {
	ctx := client.GetSystemUserContext()
	repoManagerClient, err := repoClient.NewRepoManagerClient(ctx)
	if err != nil {
		return nil, err
	}
	limit := uint32(50)
	offset := uint32(0)
	rs := make(repoInfos)
	for {
		req := pb.DescribeReposRequest{
			Limit:  limit,
			Offset: offset,
			Status: []string{constants.StatusActive}}
		res, err := repoManagerClient.DescribeRepos(ctx, &req)
		if err != nil {
			return nil, err
		}
		for _, r := range res.GetRepoSet() {
			rs[r.GetRepoId().GetValue()] = r.GetOwner().GetValue()
		}
		// In most cases, len(res.GetRepoSet()) <= limit
		if len(res.GetRepoSet()) >= int(limit) {
			offset += uint32(len(res.GetRepoSet()))
		} else {
			return rs, nil
		}
	}
}

func (p *Server) autoIndex() error {
	repos, err := getRepos()
	if err != nil {
		return err
	}
	logger.Infof("Got repos [%+v]", repos)
	for repoId, owner := range repos {
		repoEvent, err := p.controller.NewRepoEvent(repoId, owner)
		if err != nil {
			return err
		}
		logger.Infof("Repo [%s] submit repo event [%+v] success", repoId, repoEvent)
	}
	return nil
}

func (p *Server) startCron(repoCron string) *cron.Cron {
	c := cron.New()
	if repoCron != "" {
		c.AddFunc(repoCron, func() {
			logger.Debugf("Start auto index, current cron is [%s]", repoCron)
			err := p.autoIndex()
			if err != nil {
				logger.Panicf("failed to auto index repos, [%+v]", err)
			}
		})
	}
	c.Start()
	logger.Debugf("Repo cron had started")
	return c
}

func (p *Server) Cron() {
	repoCron := p.GlobalConfig().Repo.Cron
	c := p.startCron(repoCron)
	p.ThreadWatchGlobalConfig(func(globalConfig *config.GlobalConfig) {
		currentRepoCron := globalConfig.Repo.Cron
		if currentRepoCron != repoCron {
			logger.Debugf("Repo cron had update to [%s], stop old cron job [%s]", currentRepoCron, repoCron)
			c.Stop()
			repoCron = currentRepoCron
			c = p.startCron(repoCron)
		}
	})
}
