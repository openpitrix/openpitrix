// Copyright 2017 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package pbutil

import (
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/golang/protobuf/ptypes/wrappers"

	"openpitrix.io/openpitrix/pkg/db"
	"openpitrix.io/openpitrix/pkg/logger"
)

type RequestHadOffset interface {
	GetOffset() uint32
}

type RequestHadLimit interface {
	GetLimit() uint32
}

const (
	DefaultOffset = uint64(0)
	DefaultLimit  = uint64(20)
)

func GetOffsetFromRequest(req RequestHadOffset) uint64 {
	n := req.GetOffset()
	if n == 0 {
		return DefaultOffset
	}
	return db.GetOffset(uint64(n))
}

func GetLimitFromRequest(req RequestHadLimit) uint64 {
	n := req.GetLimit()
	if n == 0 {
		return DefaultLimit
	}
	return db.GetLimit(uint64(n))
}

func FromProtoTimestamp(t *timestamp.Timestamp) (tt time.Time) {
	tt, err := ptypes.Timestamp(t)
	if err != nil {
		logger.Critical(nil, "Cannot convert timestamp [T] to time.Time [%+v]: %+v", t, err)
		panic(err)
	}
	return
}

func ToProtoTimestamp(t time.Time) (tt *timestamp.Timestamp) {
	if t.IsZero() {
		return nil
	}
	tt, err := ptypes.TimestampProto(t)
	if err != nil {
		logger.Critical(nil, "Cannot convert time.Time [%+v] to ToProtoTimestamp[T]: %+v", t, err)
		panic(err)
	}
	return
}

func ToProtoString(str string) *wrappers.StringValue {
	return &wrappers.StringValue{Value: str}
}

func ToProtoUInt32(uint32 uint32) *wrappers.UInt32Value {
	return &wrappers.UInt32Value{Value: uint32}
}

func ToProtoInt32(i int32) *wrappers.Int32Value {
	return &wrappers.Int32Value{Value: i}
}

func ToProtoBool(bool bool) *wrappers.BoolValue {
	return &wrappers.BoolValue{Value: bool}
}

func ToProtoBytes(bytes []byte) *wrappers.BytesValue {
	return &wrappers.BytesValue{Value: bytes}
}
