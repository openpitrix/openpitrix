// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package db

import (
	"fmt"

	"github.com/gocraft/dbr"
)

func AddJoinFilterWithMap(q *SelectQuery, table, joinTable, primaryKey, keyField, valueField string, filterMap map[string][]string) *SelectQuery {
	var whereCondition []dbr.Builder
	for key, values := range filterMap {
		aliasTableName := fmt.Sprintf("table_label_%d", q.joinCount)
		onCondition := fmt.Sprintf("%s.%s = %s.%s", aliasTableName, primaryKey, table, primaryKey)
		q = q.Join(dbr.I(joinTable).As(aliasTableName), onCondition)
		whereCondition = append(whereCondition, And(Eq(aliasTableName+"."+keyField, key), Eq(aliasTableName+"."+valueField, values)))
		q.joinCount++
	}
	if len(whereCondition) > 0 {
		q = q.Where(And(whereCondition...))
	}
	return q
}
