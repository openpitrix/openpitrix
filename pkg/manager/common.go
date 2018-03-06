// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package manager

import (
	"reflect"

	"github.com/fatih/structs"
	"github.com/gocraft/dbr"
	"github.com/gogo/protobuf/proto"
	"github.com/golang/protobuf/ptypes/wrappers"

	"openpitrix.io/openpitrix/pkg/db"
	"openpitrix.io/openpitrix/pkg/logger"
	"openpitrix.io/openpitrix/pkg/models"
	"openpitrix.io/openpitrix/pkg/utils"
)

type Request interface {
	Reset()
	String() string
	ProtoMessage()
	Descriptor() ([]byte, []int)
}

const (
	TagName              = "protobuf"
	SearchWordColumnName = "search_word"
)

func getSearchFilter(tableName string, value interface{}, exclude ...string) dbr.Builder {
	if v, ok := value.(string); ok {
		var ops []dbr.Builder
		for _, column := range models.SearchColumns[tableName] {
			if !utils.StringIn(column, exclude) {
				ops = append(ops, db.Like(column, v))
			}
		}
		if len(ops) == 0 {
			return nil
		}
		return db.Or(ops...)
	} else {
		logger.Warnf("search_word [%+v] is not string", value)
	}
	return nil
}

func getStringValue(param interface{}) interface{} {
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
	var conditions []dbr.Builder
	for _, field := range structs.Fields(req) {
		tag := field.Tag(TagName)
		prop := proto.Properties{}
		prop.Parse(tag)
		column := prop.OrigName
		param := field.Value()
		if utils.StringIn(column, models.IndexedColumns[tableName]) {
			value := getStringValue(param)
			if value != nil {
				conditions = append(conditions, db.Eq(column, value))
			}
		}
		if column == SearchWordColumnName && utils.StringIn(tableName, models.SearchWordColumnTable) {
			value := getStringValue(param)
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
		tag := field.Tag(TagName)
		prop := proto.Properties{}
		prop.Parse(tag)
		column := prop.OrigName
		f := field.Value()
		v := reflect.ValueOf(f)
		if utils.FindString(columns, column) > -1 && !v.IsNil() {
			switch v := f.(type) {
			case string:
				attributes[column] = v
			case *wrappers.StringValue:
				attributes[column] = v.GetValue()
			}
		}
	}
	return attributes
}
