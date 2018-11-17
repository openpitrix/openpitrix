// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package manager

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/fatih/structs"
	"github.com/gocraft/dbr"
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/golang/protobuf/ptypes/wrappers"

	"openpitrix.io/openpitrix/pkg/constants"
	"openpitrix.io/openpitrix/pkg/db"
	"openpitrix.io/openpitrix/pkg/logger"
	"openpitrix.io/openpitrix/pkg/util/pbutil"
	"openpitrix.io/openpitrix/pkg/util/stringutil"
)

type Request interface {
	Reset()
	String() string
	ProtoMessage()
	Descriptor() ([]byte, []int)
}
type RequestWithSortKey interface {
	Request
	GetSortKey() *wrappers.StringValue
}
type RequestWithReverse interface {
	RequestWithSortKey
	GetReverse() *wrappers.BoolValue
}

const (
	TagName              = "json"
	SearchWordColumnName = "search_word"
)

func getSearchFilter(tableName string, value interface{}, exclude ...string) dbr.Builder {
	if v, ok := value.(string); ok {
		var ops []dbr.Builder
		for _, column := range constants.SearchColumns[tableName] {
			if stringutil.StringIn(column, exclude) {
				continue
			}
			// if column suffix is _id, must exact match
			if strings.HasSuffix(column, "_id") {
				ops = append(ops, db.Eq(column, v))
			} else {
				ops = append(ops, db.Like(column, v))
			}
		}
		if len(ops) == 0 {
			return nil
		}
		return db.Or(ops...)
	} else if value != nil {
		logger.Warn(nil, "search_word [%+v] is not string", value)
	}
	return nil
}

func getReqValue(param interface{}) interface{} {
	switch value := param.(type) {
	case string:
		if value == "" {
			return nil
		}
		return value
	case *wrappers.StringValue:
		if value == nil {
			return nil
		}
		return value.GetValue()
	case *wrappers.Int32Value:
		if value == nil {
			return nil
		}
		return value.GetValue()
	case []string:
		var values []string
		for _, v := range value {
			if v != "" {
				values = append(values, v)
			}
		}
		if len(values) == 0 {
			return nil
		}
		return values
	}
	return nil
}

func BuildFilterConditions(req Request, tableName string, exclude ...string) dbr.Builder {
	return buildFilterConditions(false, req, tableName, exclude...)
}

func BuildFilterConditionsWithPrefix(req Request, tableName string, exclude ...string) dbr.Builder {
	return buildFilterConditions(true, req, tableName, exclude...)
}

func getFieldName(field *structs.Field) string {
	tag := field.Tag(TagName)
	t := strings.Split(tag, ",")
	if len(t) == 0 {
		return "-"
	}
	return t[0]
}

func buildFilterConditions(withPrefix bool, req Request, tableName string, exclude ...string) dbr.Builder {
	var conditions []dbr.Builder
	for _, field := range structs.Fields(req) {
		column := getFieldName(field)
		param := field.Value()
		indexedColumns, ok := constants.IndexedColumns[tableName]
		if ok && stringutil.StringIn(column, indexedColumns) {
			value := getReqValue(param)
			if value != nil {
				key := column
				if withPrefix {
					key = tableName + "." + key
				}
				conditions = append(conditions, db.Eq(key, value))
			}
		}
		// TODO: search column
		if column == SearchWordColumnName && stringutil.StringIn(tableName, constants.SearchWordColumnTable) {
			value := getReqValue(param)
			condition := getSearchFilter(tableName, value, exclude...)
			if condition != nil {
				conditions = append(conditions, condition)
			}
		}
	}
	if len(conditions) == 0 {
		return nil
	}
	return db.And(conditions...)
}

func BuildUpdateAttributes(req Request, columns ...string) map[string]interface{} {
	attributes := make(map[string]interface{})
	for _, field := range structs.Fields(req) {
		column := getFieldName(field)
		f := field.Value()
		v := reflect.ValueOf(f)
		if stringutil.FindString(columns, column) > -1 && !v.IsNil() {
			switch v := f.(type) {
			case *wrappers.StringValue:
				attributes[column] = v.GetValue()
			case *wrappers.BoolValue:
				attributes[column] = v.GetValue()
			case *wrappers.Int32Value:
				attributes[column] = v.GetValue()
			case *wrappers.UInt32Value:
				attributes[column] = v.GetValue()
			case *timestamp.Timestamp:
				attributes[column] = pbutil.FromProtoTimestamp(v)

			default:
				attributes[column] = v
			}
		}
	}
	return attributes
}

func AddQueryOrderDir(query *db.SelectQuery, req Request, defaultColumn string) *db.SelectQuery {
	isAsc := false
	if r, ok := req.(RequestWithReverse); ok {
		reverse := r.GetReverse()
		if reverse != nil {
			isAsc = reverse.GetValue()
		}
	}
	if r, ok := req.(RequestWithSortKey); ok {
		s := r.GetSortKey()
		if s != nil {
			defaultColumn = s.GetValue()
		}
	}
	query = query.OrderDir(defaultColumn, isAsc)
	return query
}

func AddQueryJoinWithMap(query *db.SelectQuery, table, joinTable, primaryKey, keyField, valueField string, filterMap map[string][]string) *db.SelectQuery {
	var whereCondition []dbr.Builder
	for key, values := range filterMap {
		aliasTableName := fmt.Sprintf("table_label_%d", query.JoinCount)
		onCondition := fmt.Sprintf("%s.%s = %s.%s", aliasTableName, primaryKey, table, primaryKey)
		query = query.Join(dbr.I(joinTable).As(aliasTableName), onCondition)
		whereCondition = append(whereCondition, db.And(db.Eq(aliasTableName+"."+keyField, key), db.Eq(aliasTableName+"."+valueField, values)))
		query.JoinCount++
	}
	if len(whereCondition) > 0 {
		query = query.Where(db.And(whereCondition...))
	}
	return query
}
