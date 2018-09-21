// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package topic

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/mvcc/mvccpb"

	"openpitrix.io/openpitrix/pkg/etcd"
	"openpitrix.io/openpitrix/pkg/logger"
	"openpitrix.io/openpitrix/pkg/util/idutil"
	"openpitrix.io/openpitrix/pkg/util/jsonutil"
)

const expireTime = 60 // second

func formatTopic(uid string, eventId uint64) string {
	return fmt.Sprintf("%s/%s/%d", topicPrefix, uid, eventId)
}

func parseTopic(topic string) (uid string, eventId uint64) {
	t := strings.Split(topic, "/")
	uid = t[1]
	eid, _ := strconv.Atoi(t[2])
	eventId = uint64(eid)
	return
}

func pushEvent(ctx context.Context, e *etcd.Etcd, uid string, msg Message) error {
	var err error
	var eventId = idutil.GetIntId()
	var key = formatTopic(uid, eventId)
	value, err := jsonutil.Encode(msg)
	if err != nil {
		logger.Error(ctx, "Encode message [%+v] to json failed", msg)
		return err
	}

	resp, err := e.Grant(ctx, expireTime)
	if err != nil {
		logger.Error(ctx, "Grant ttl from etcd failed: %+v", err)
		return err
	}

	_, err = e.Put(ctx, key, string(value), clientv3.WithLease(resp.ID))
	if err != nil {
		logger.Error(ctx, "Push user [%s] event [%d] [%s] to etcd failed: %+v", uid, eventId, string(value), err)
		return err
	}
	return nil
}

func watchEvents(e *etcd.Etcd) chan userMessage {
	var c = make(chan userMessage, 255)
	go func() {
		logger.Debug(nil, "Start watch events")
		watchRes := e.Watch(context.Background(), topicPrefix+"/", clientv3.WithPrefix())
		for res := range watchRes {
			for _, ev := range res.Events {
				if ev.Type == mvccpb.PUT {
					var message Message
					key := string(ev.Kv.Key)
					userId, eventId := parseTopic(key)
					err := jsonutil.Decode(ev.Kv.Value, &message)
					if err != nil {
						logger.Error(nil, "Decode event [%s] [%d] [%s] failed: %+v",
							userId, eventId, string(ev.Kv.Value), err)
					} else {
						logger.Debug(nil, "Got event [%s] [%d] [%s]", userId, eventId, string(ev.Kv.Value))
						c <- userMessage{
							UserId:  userId,
							Message: message,
						}
					}
				}
			}
		}
	}()
	return c
}

func PushEvent(ctx context.Context, e *etcd.Etcd, uid string, t messageType, model model) error {
	msg := Message{
		Type:     t,
		Resource: model.GetTopicResource(),
	}
	return pushEvent(ctx, e, uid, msg)
}
