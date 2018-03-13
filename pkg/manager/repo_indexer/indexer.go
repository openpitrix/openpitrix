// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package repo_indexer

import (
	"context"
	"fmt"
	"time"

	"openpitrix.io/openpitrix/pkg/constants"
	"openpitrix.io/openpitrix/pkg/db"
	"openpitrix.io/openpitrix/pkg/etcd"
	"openpitrix.io/openpitrix/pkg/logger"
	"openpitrix.io/openpitrix/pkg/models"
	"openpitrix.io/openpitrix/pkg/pi"
	"openpitrix.io/openpitrix/pkg/utils"
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

func (i *Indexer) IndexRepo(repoTask *models.RepoTask, cb func()) {
	defer cb()
	logger.Infof("Got repo task: %+v", repoTask)
	// TODO: get index.yaml from repo
	_, err := i.Db.
		Update(models.RepoTaskTableName).
		Set("status", constants.StatusSuccessful).
		Where(db.Eq("repo_task_id", repoTask.RepoTaskId)).
		Exec()
	if err != nil {
		logger.Panicf("Cannot set repo task [&s] status to [success]", repoTask.RepoTaskId)
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
