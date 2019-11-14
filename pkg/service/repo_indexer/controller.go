// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package repo_indexer

import (
	"context"
	"fmt"
	"time"

	"openpitrix.io/openpitrix/pkg/client"
	appClient "openpitrix.io/openpitrix/pkg/client/app"
	"openpitrix.io/openpitrix/pkg/constants"
	"openpitrix.io/openpitrix/pkg/db"
	"openpitrix.io/openpitrix/pkg/etcd"
	"openpitrix.io/openpitrix/pkg/logger"
	"openpitrix.io/openpitrix/pkg/models"
	"openpitrix.io/openpitrix/pkg/pb"
	"openpitrix.io/openpitrix/pkg/pi"
	"openpitrix.io/openpitrix/pkg/sender"
	"openpitrix.io/openpitrix/pkg/util/atomicutil"
	"openpitrix.io/openpitrix/pkg/util/ctxutil"
)

type eventChannel chan *models.RepoEvent

type EventController struct {
	ctx          context.Context
	queue        *etcd.Queue
	channel      eventChannel
	runningCount atomicutil.Counter
}

func NewEventController(ctx context.Context) *EventController {
	return &EventController{
		ctx:          ctx,
		queue:        pi.Global().Etcd(ctx).NewQueue("repo-indexer-event"),
		channel:      make(eventChannel),
		runningCount: atomicutil.Counter(0),
	}
}

func (i *EventController) NewRepoEvent(repoId string, ownerPath sender.OwnerPath) (*models.RepoEvent, error) {
	var repoEventId string
	err := pi.Global().Etcd(i.ctx).Dlock(context.Background(), constants.RepoIndexPrefix+repoId, func() error {
		count, err := pi.Global().DB(i.ctx).Select(models.RepoEventColumns...).
			From(constants.TableRepoEvent).
			Where(db.Eq(constants.ColumnRepoId, repoId)).
			Where(db.Eq(constants.ColumnStatus, []string{constants.StatusWorking, constants.StatusPending})).
			Count()
		if err != nil {
			return err
		}
		if count > 0 {
			return fmt.Errorf("repo index event [%s] is running", repoId)
		}
		repoEvent := models.NewRepoEvent(repoId, ownerPath)
		_, err = pi.Global().DB(i.ctx).
			InsertInto(constants.TableRepoEvent).
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
	err = pi.Global().DB(i.ctx).Select(models.RepoEventColumns...).
		From(constants.TableRepoEvent).
		Where(db.Eq(constants.ColumnRepoEventId, repoEventId)).
		LoadOne(&repoEvent)
	if err != nil {
		return nil, err
	}
	return &repoEvent, nil
}

func (i *EventController) updateRepoEventStatus(ctx context.Context, repoEvent *models.RepoEvent, status, result string) error {
	_, err := pi.Global().DB(ctx).
		Update(constants.TableRepoEvent).
		Set(constants.ColumnStatus, status).
		Set(constants.ColumnResult, result).
		Where(db.Eq(constants.ColumnRepoEventId, repoEvent.RepoEventId)).
		Exec()
	if err != nil {
		logger.Critical(
			ctx,
			"Failed to set repo event [%s] status to [%s] result to [%s], %+v",
			repoEvent.RepoEventId, status, result, err)
		return err
	}

	return nil
}

func (i *EventController) ExecuteEvent(ctx context.Context, repoEvent *models.RepoEvent, cb func()) {
	ctx = client.SetSystemUserToContext(ctx)
	ctx = ctxutil.AddMessageId(ctx, repoEvent.RepoEventId)
	ctx = ctxutil.AddMessageId(ctx, repoEvent.RepoId)
	defer cb()
	defer func() {
		if err := recover(); err != nil {
			logger.Critical(ctx, "ExecuteEvent [%s] recover with error: %+v", repoEvent.RepoEventId, err)
			i.updateRepoEventStatus(ctx, repoEvent, constants.StatusFailed, fmt.Sprintf("%+v", err))
		}
	}()
	logger.Info(ctx, "Got repo event: %+v", repoEvent)
	err := func() (err error) {
		repoId := repoEvent.RepoId
		appManagerClient, err := appClient.NewAppManagerClient()
		if err != nil {
			return
		}
		req := pb.SyncRepoRequest{
			RepoId: repoId,
		}
		res, err := appManagerClient.SyncRepo(ctx, &req)
		if err != nil {
			return
		}
		if res.Failed {
			err = fmt.Errorf(res.Result)
			logger.Error(ctx, "Failed to index repo [%s], %+v", repoId, err)
		}
		return
	}()
	if err != nil {
		logger.Critical(ctx, "Failed to execute repo event: %+v", err)
		i.updateRepoEventStatus(ctx, repoEvent, constants.StatusFailed, err.Error())
	} else {
		i.updateRepoEventStatus(ctx, repoEvent, constants.StatusSuccessful, "")
	}
}

func (i *EventController) getRepoEvent(repoEventId string) (repoEvent models.RepoEvent, err error) {
	err = pi.Global().DB(i.ctx).
		Select(models.RepoEventColumns...).
		From(constants.TableRepoEvent).
		Where(db.Eq(constants.ColumnRepoEventId, repoEventId)).
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
	if pi.Global().GlobalConfig().Repo.MaxRepoEvents > 0 {
		return pi.Global().GlobalConfig().Repo.MaxRepoEvents
	} else {
		return constants.DefaultMaxRepoEvents
	}
}

func (i *EventController) Dequeue() {
	for {
		if i.runningCount.Get() > i.GetEventLength() {
			logger.Error(i.ctx, "Sleep 10s, running event count exceed [%d/%d]", i.runningCount.Get(), i.GetEventLength())
			time.Sleep(10 * time.Second)
			continue
		}
		repoEvent, err := i.getRepoEventFromQueue()
		if err != nil {
			logger.Error(i.ctx, "Failed to get repo event from etcd: %+v", err)
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
		go i.ExecuteEvent(i.ctx, event, func() {
			i.runningCount.Add(-1)
		})
	}
}
