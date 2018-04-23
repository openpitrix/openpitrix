// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package repo_indexer

import (
	"context"
	"fmt"
	"runtime/debug"
	"time"

	"openpitrix.io/openpitrix/pkg/client"
	repoClient "openpitrix.io/openpitrix/pkg/client/repo"
	"openpitrix.io/openpitrix/pkg/constants"
	"openpitrix.io/openpitrix/pkg/db"
	"openpitrix.io/openpitrix/pkg/etcd"
	"openpitrix.io/openpitrix/pkg/logger"
	"openpitrix.io/openpitrix/pkg/manager/repo_indexer/indexer"
	"openpitrix.io/openpitrix/pkg/models"
	"openpitrix.io/openpitrix/pkg/pb"
	"openpitrix.io/openpitrix/pkg/pi"
	"openpitrix.io/openpitrix/pkg/utils"
)

type eventChannel chan *models.RepoEvent

type EventController struct {
	*pi.Pi
	queue        *etcd.Queue
	channel      eventChannel
	runningCount utils.Counter
}

func NewEventController(pi *pi.Pi) *EventController {
	return &EventController{
		Pi:           pi,
		queue:        pi.Etcd.NewQueue("repo-indexer-event"),
		channel:      make(eventChannel),
		runningCount: utils.Counter(0),
	}
}

func (i *EventController) NewRepoEvent(repoId, owner string) (*models.RepoEvent, error) {
	var repoEventId string
	err := i.Etcd.Dlock(context.Background(), constants.RepoIndexPrefix+repoId, func() error {
		count, err := i.Db.Select(models.RepoEventColumns...).
			From(models.RepoEventTableName).
			Where(db.Eq("repo_id", repoId)).
			Where(db.Eq("status", []string{constants.StatusWorking, constants.StatusPending})).
			Count()
		if err != nil {
			return err
		}
		if count > 0 {
			return fmt.Errorf("repo [%s] had running index event", repoId)
		}
		repoEvent := models.NewRepoEvent(repoId, owner)
		_, err = i.Db.InsertInto(models.RepoEventTableName).
			Columns(models.RepoEventColumns...).
			Record(repoEvent).
			Exec()
		if err != nil {
			return err
		}
		repoEventId = repoEvent.RepoEventId
		err = i.queue.Enqueue(repoEventId)
		return err
	})
	if err != nil {
		return nil, err
	}
	var repoEvent models.RepoEvent
	err = i.Db.Select(models.RepoEventColumns...).
		From(models.RepoEventTableName).
		Where(db.Eq("repo_event_id", repoEventId)).
		LoadOne(&repoEvent)
	if err != nil {
		return nil, err
	}
	return &repoEvent, nil
}

func (i *EventController) updateRepoEventStatus(repoEventId, status, result string) error {
	_, err := i.Db.
		Update(models.RepoEventTableName).
		Set("status", status).
		Set("result", result).
		Where(db.Eq("repo_event_id", repoEventId)).
		Exec()
	if err != nil {
		logger.Panicf("Failed to set repo event [&s] status to [%s] result to [%s], %+v", repoEventId, status, result, err)
	}
	return err
}

func (i *EventController) ExecuteEvent(repoEvent *models.RepoEvent, cb func()) {
	defer cb()
	defer func() {
		if err := recover(); err != nil {
			logger.Panic(err)
			i.updateRepoEventStatus(repoEvent.RepoEventId, constants.StatusFailed, fmt.Sprintf("%+v", err))
		}
	}()
	logger.Infof("Got repo event: %+v", repoEvent)
	err := func() (err error) {
		ctx := client.GetSystemUserContext()
		repoManagerClient, err := repoClient.NewRepoManagerClient(ctx)
		if err != nil {
			return
		}
		repoId := repoEvent.RepoId
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
		repo := res.RepoSet[0]
		err = indexer.GetIndexer(repo).IndexRepo()
		if err != nil {
			logger.Errorf("Failed to index repo [%s]", repoId)
		}
		return
	}()
	if err != nil {
		// FIXME: remove panic log
		logger.Panicf("Failed to execute repo event: %+v", err)
		logger.Panic(string(debug.Stack()))
		i.updateRepoEventStatus(repoEvent.RepoEventId, constants.StatusFailed, fmt.Sprintf("%+v", err))
	} else {
		i.updateRepoEventStatus(repoEvent.RepoEventId, constants.StatusSuccessful, "")
	}
}

func (i *EventController) getRepoEvent(repoEventId string) (repoEvent models.RepoEvent, err error) {
	err = i.Db.
		Select(models.RepoEventColumns...).
		From(models.RepoEventTableName).
		Where(db.Eq("repo_event_id", repoEventId)).
		LoadOne(&repoEvent)
	return
}

func (i *EventController) getRepoEventFromQueue() (repoEvent models.RepoEvent, err error) {
	repoEventId, err := i.queue.Dequeue()
	if err != nil {
		return
	}
	repoEvent, err = i.getRepoEvent(repoEventId)
	return
}

func (i *EventController) GetEventLength() int32 {
	return constants.RepoEventLength
}

func (i *EventController) Dequeue() {
	for {
		if i.runningCount.Get() > i.GetEventLength() {
			logger.Errorf("Sleep 10s, running event count exceed [%d/%d]", i.runningCount.Get(), i.GetEventLength())
			time.Sleep(10 * time.Second)
			continue
		}
		repoEvent, err := i.getRepoEventFromQueue()
		if err != nil {
			logger.Errorf("Failed to get repo event from etcd: %+v", err)
			time.Sleep(10 * time.Second)
			continue
		}
		i.channel <- &repoEvent
	}
}

func (i *EventController) Serve() {
	go i.Dequeue()
	for event := range i.channel {
		i.runningCount.Add(1)
		go i.ExecuteEvent(event, func() {
			i.runningCount.Add(-1)
		})
	}
}
