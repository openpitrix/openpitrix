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
	"openpitrix.io/openpitrix/pkg/utils"
)

type Request interface {
	Reset()
	String() string
	ProtoMessage()
	Descriptor() ([]byte, []int)
}

const TagName = "protobuf"

func BuildFilterConditions(req Request, columns ...string) dbr.Builder {
	var filter []dbr.Builder
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
				filter = append(filter, db.Eq(column, v))
			case *wrappers.StringValue:
				filter = append(filter, db.Eq(column, v.GetValue()))
			case []string:
				var set []string
				for _, i := range v {
					set = append(set, i)
				}
				filter = append(filter, db.Eq(column, set))
			}
		}
	}
	if len(filter) == 0 {
		return nil
	}
	return db.And(filter...)
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
