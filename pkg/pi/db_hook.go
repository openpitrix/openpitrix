// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package pi

import (
	"context"

	"github.com/gocraft/dbr"

	"openpitrix.io/openpitrix/pkg/constants"
	"openpitrix.io/openpitrix/pkg/db"
	"openpitrix.io/openpitrix/pkg/topic"
	"openpitrix.io/openpitrix/pkg/util/stringutil"
)

func getResourceIds(key string, whereCond []dbr.Builder) []string {
	var rids []string
	for _, where := range whereCond {
		w, ok := where.(*db.EqCondition)
		if !ok {
			continue
		}
		if w.Column == key {
			switch rid := w.Value.(type) {
			case string:
				rids = append(rids, rid)
			case []string:
				rids = rid
			}
		}
	}
	return rids
}

func (p *Pi) GetUpdateHook(ctx context.Context) db.UpdateHook {
	return func(query *db.UpdateQuery) {
		table := query.Table
		whereCond := query.WhereCond
		columns, ok := constants.PushEventTables[table]
		if !ok {
			return
		}
		key, columns := columns[0], columns[1:]
		rids := getResourceIds(key, whereCond)
		if len(rids) == 0 {
			return
		}
		for _, rid := range rids {
			var owner string
			err := p.DB(ctx).Select(constants.ColumnOwner).From(table).Where(db.Eq(key, rids)).LoadOne(&owner)
			if err != nil {
				continue
			}
			if owner == "" {
				continue
			}
			event := topic.NewResource(table, rid)
			for _, column := range columns {
				if value, ok := query.Value[column]; ok {
					event.WithValue(column, value)
				}
			}
			topic.PushEvent(ctx, p.Etcd(ctx), owner, topic.Update, event)
		}
	}
}

func (p *Pi) GetDeleteHook(ctx context.Context) db.DeleteHook {
	return func(query *db.DeleteQuery) {
		table := query.Table
		whereCond := query.WhereCond
		columns, ok := constants.PushEventTables[table]
		if !ok {
			return
		}
		key, columns := columns[0], columns[1:]
		rids := getResourceIds(key, whereCond)
		if len(rids) == 0 {
			return
		}
		for _, rid := range rids {
			var owner string
			err := p.DB(ctx).Select(constants.ColumnOwner).From(table).Where(db.Eq(key, rid)).LoadOne(&owner)
			if err != nil {
				continue
			}
			if owner == "" {
				continue
			}
			topic.PushEvent(ctx, p.Etcd(ctx), owner, topic.Delete, topic.NewResource(table, rid))
		}
	}
}

func (p *Pi) GetInsertHook(ctx context.Context) db.InsertHook {
	return func(query *db.InsertQuery) {
		table := query.Table
		columns, ok := constants.PushEventTables[table]
		if !ok {
			return
		}
		key, columns := columns[0], columns[1:]
		var keyIdx = -1
		var columnsMap = make(map[string]int)
		for idx, c := range query.Column {
			if c == key {
				keyIdx = idx
			}
			if stringutil.StringIn(c, columns) {
				columnsMap[c] = idx
			}
		}
		var resources = make(map[string]map[string]interface{})
		for _, v := range query.Value {
			rid := v[keyIdx].(string)
			resources[rid] = make(map[string]interface{})
			for column, idx := range columnsMap {
				resources[rid][column] = v[idx]
			}
		}
		for rid, resource := range resources {
			var owner string
			err := p.DB(ctx).Select(constants.ColumnOwner).From(table).Where(db.Eq(key, rid)).LoadOne(&owner)
			if err != nil {
				return
			}
			if owner == "" {
				return
			}
			event := topic.NewResource(table, rid)
			for key, value := range resource {
				event.WithValue(key, value)
			}
			topic.PushEvent(ctx, p.Etcd(ctx), owner, topic.Create, event)
		}
	}
}
