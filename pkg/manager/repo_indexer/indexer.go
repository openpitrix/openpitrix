// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package repo_indexer

import (
	"context"
	"fmt"
	"runtime/debug"
	"time"

	"openpitrix.io/openpitrix/pkg/constants"
	"openpitrix.io/openpitrix/pkg/db"
	"openpitrix.io/openpitrix/pkg/etcd"
	"openpitrix.io/openpitrix/pkg/logger"
	"openpitrix.io/openpitrix/pkg/manager/repo"
	"openpitrix.io/openpitrix/pkg/models"
	"openpitrix.io/openpitrix/pkg/pb"
	"openpitrix.io/openpitrix/pkg/pi"
	"openpitrix.io/openpitrix/pkg/utils"
	"openpitrix.io/openpitrix/pkg/utils/sender"
)

type taskChannel chan *models.RepoTask

type Indexer struct {
	*pi.Pi
	queue        *etcd.Queue
	channel      taskChannel
	runningCount utils.Counter
}

func NewIndexer(pi *pi.Pi) *Indexer {
	return &Indexer{
		Pi:           pi,
		queue:        pi.Etcd.NewQueue("repo-indexer-task"),
		channel:      make(taskChannel),
		runningCount: utils.Counter(0),
	}
}

func (i *Indexer) NewRepoTask(repoId, owner string) (*models.RepoTask, error) {
	var repoTaskId string
	err := i.Etcd.Dlock(context.Background(), constants.RepoIndexPrefix+repoId, func() error {
		count, err := i.Db.Select(models.RepoTaskColumns...).
			From(models.RepoTaskTableName).
			Where(db.Eq("repo_id", repoId)).
			Where(db.Eq("status", []string{constants.StatusWorking, constants.StatusPending})).
			Count()
		if err != nil {
			return err
		}
		if count > 0 {
			return fmt.Errorf("repo [%s] had running index task", repoId)
		}
		repoTask := models.NewRepoTask(repoId, owner)
		_, err = i.Db.InsertInto(models.RepoTaskTableName).
			Columns(models.RepoTaskColumns...).
			Record(repoTask).
			Exec()
		if err != nil {
			return err
		}
		repoTaskId = repoTask.RepoTaskId
		err = i.queue.Enqueue(repoTaskId)
		return err
	})
	if err != nil {
		return nil, err
	}
	var repoTask models.RepoTask
	err = i.Db.Select(models.RepoTaskColumns...).
		From(models.RepoTaskTableName).
		Where(db.Eq("repo_task_id", repoTaskId)).
		LoadOne(&repoTask)
	if err != nil {
		return nil, err
	}
	return &repoTask, nil
}

func (i *Indexer) updateRepoTaskStatus(repoTaskId, status, result string) error {
	_, err := i.Db.
		Update(models.RepoTaskTableName).
		Set("status", status).
		Set("result", result).
		Where(db.Eq("repo_task_id", repoTaskId)).
		Exec()
	if err != nil {
		logger.Panicf("Failed to set repo task [&s] status to [%s] result to [%s], %+v", repoTaskId, status, result, err)
	}
	return err
}

func (i *Indexer) IndexRepo(repoTask *models.RepoTask, cb func()) {
	defer cb()
	defer func() {
		if err := recover(); err != nil {
			logger.Panic(err)
			i.updateRepoTaskStatus(repoTask.RepoTaskId, constants.StatusFailed, fmt.Sprintf("%+v", err))
		}
	}()
	logger.Infof("Got repo task: %+v", repoTask)
	err := func() (err error) {
		ctx := sender.NewContext(context.Background(), sender.GetSystemUser())
		repoManagerClient, err := repo.NewRepoManagerClient(ctx)
		if err != nil {
			return
		}
		repoId := repoTask.RepoId
		owner := repoTask.Owner
		req := pb.DescribeReposRequest{
			RepoId: []string{repoId},
		}
		res, err := repoManagerClient.DescribeRepos(ctx, &req)
		if err != nil {
			return
		}
		if res.TotalCount == 0 {
			err = fmt.Errorf("failed to get repo [%s]", repoId)
			return
		}
		repoUrl := res.RepoSet[0].Url.GetValue()
		indexFile, err := GetIndexFile(repoUrl)
		if err != nil {
			return
		}
		for chartName, chartVersions := range indexFile.Entries {
			var appId string
			logger.Debugf("Start index chart [%s]", chartName)
			appId, err = SyncAppInfo(repoId, owner, chartName, &chartVersions)
			if err != nil {
				logger.Errorf("Failed to sync chart [%s] to app info", chartName)
				return
			}
			logger.Infof("Sync chart [%s] to app [%s] success", chartName, appId)
			for _, chartVersion := range chartVersions {
				var versionId string
				versionId, err = SyncAppVersionInfo(appId, owner, chartVersion)
				if err != nil {
					logger.Errorf("Failed to sync chart version [%s] to app version", chartVersion.GetAppVersion())
					return
				}
				logger.Debugf("Chart version [%s] sync to app version [%s]", chartVersion.GetVersion(), versionId)
			}
		}
		return
	}()
	if err != nil {
		// FIXME: remove panic log
		logger.Panicf("Failed to execute repo task: %+v", err)
		logger.Panic(string(debug.Stack()))
		i.updateRepoTaskStatus(repoTask.RepoTaskId, constants.StatusFailed, fmt.Sprintf("%+v", err))
	} else {
		i.updateRepoTaskStatus(repoTask.RepoTaskId, constants.StatusSuccessful, "")
	}
}

func (i *Indexer) getRepoTask(repoTaskId string) (repoTask models.RepoTask, err error) {
	err = i.Db.
		Select(models.RepoTaskColumns...).
		From(models.RepoTaskTableName).
		Where(db.Eq("repo_task_id", repoTaskId)).
		LoadOne(&repoTask)
	return
}

func (i *Indexer) getRepoTaskFromQueue() (repoTask models.RepoTask, err error) {
	repoTaskId, err := i.queue.Dequeue()
	if err != nil {
		return
	}
	repoTask, err = i.getRepoTask(repoTaskId)
	return
}

func (i *Indexer) GetTaskLength() int32 {
	return constants.RepoTaskLength
}

func (i *Indexer) Dequeue() {
	for {
		if i.runningCount.Get() > i.GetTaskLength() {
			logger.Errorf("Sleep 10s, running task count exceed [%d/%d]", i.runningCount.Get(), i.GetTaskLength())
			time.Sleep(10 * time.Second)
			continue
		}
		repoTask, err := i.getRepoTaskFromQueue()
		if err != nil {
			logger.Errorf("Failed to get repo task from etcd: %+v", err)
			time.Sleep(10 * time.Second)
			continue
		}
		i.channel <- &repoTask
	}
}

func (i *Indexer) Serve() {
	go i.Dequeue()
	for task := range i.channel {
		i.runningCount.Add(1)
		go i.IndexRepo(task, func() {
			i.runningCount.Add(-1)
		})
	}
}
