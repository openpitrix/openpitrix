// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package models

import (
	"time"

	"openpitrix.io/openpitrix/pkg/constants"
	"openpitrix.io/openpitrix/pkg/pb"
	"openpitrix.io/openpitrix/pkg/utils"
)

const RuntimeTableName = "runtime"
const RuntimeIdPrifix = "runtime-"

func NewRuntimeId() string {
	return utils.GetUuid(RuntimeIdPrifix)
}

type Runtime struct {
	RuntimeId   string
	Name        string
	Description string
	RuntimeUrl  string
	Owner       string
	Status      string
	CreateTime  time.Time
	StatusTime  time.Time
}

var RuntimeColumns = GetColumnsFromStruct(&Runtime{})

func NewRuntime(name, description, runtimeUrl, owner string) *Runtime {
	return &Runtime{
		RuntimeId:   NewRuntimeId(),
		Name:        name,
		Description: description,
		RuntimeUrl:  runtimeUrl,
		Owner:       owner,
		Status:      constants.StatusActive,
		CreateTime:  time.Now(),
		StatusTime:  time.Now(),
	}
}

func RuntimeToPb(runtime *Runtime) *pb.Runtime {
	pbRuntime := pb.Runtime{}
	pbRuntime.RuntimeId = utils.ToProtoString(runtime.RuntimeId)
	pbRuntime.Name = utils.ToProtoString(runtime.Name)
	pbRuntime.Description = utils.ToProtoString(runtime.Description)
	pbRuntime.RuntimeUrl = utils.ToProtoString(runtime.RuntimeUrl)
	pbRuntime.Owner = utils.ToProtoString(runtime.Owner)
	pbRuntime.Status = utils.ToProtoString(runtime.Status)
	pbRuntime.CreateTime = utils.ToProtoTimestamp(runtime.CreateTime)
	pbRuntime.StatusTime = utils.ToProtoTimestamp(runtime.StatusTime)
	return &pbRuntime
}

func RuntimeToPbs(runtimes []*Runtime) (pbRuntimes []*pb.Runtime) {
	for _, runtime := range runtimes {
		pbRuntimes = append(pbRuntimes, RuntimeToPb(runtime))
	}
	return
}
