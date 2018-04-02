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

func RuntimeToPb(runtime *Runtime) *pb.RuntimeEnv {
	pbRuntimeEnv := pb.RuntimeEnv{}
	pbRuntimeEnv.RuntimeEnvId = utils.ToProtoString(runtime.RuntimeId)
	pbRuntimeEnv.Name = utils.ToProtoString(runtime.Name)
	pbRuntimeEnv.Description = utils.ToProtoString(runtime.Description)
	pbRuntimeEnv.RuntimeEnvUrl = utils.ToProtoString(runtime.RuntimeUrl)
	pbRuntimeEnv.Owner = utils.ToProtoString(runtime.Owner)
	pbRuntimeEnv.Status = utils.ToProtoString(runtime.Status)
	pbRuntimeEnv.CreateTime = utils.ToProtoTimestamp(runtime.CreateTime)
	pbRuntimeEnv.StatusTime = utils.ToProtoTimestamp(runtime.StatusTime)
	return &pbRuntimeEnv
}

func RuntimeEnvToPbs(runtimeEnvs []*Runtime) (pbRuntimeEnvs []*pb.RuntimeEnv) {
	for _, runtimeEnv := range runtimeEnvs {
		pbRuntimeEnvs = append(pbRuntimeEnvs, RuntimeToPb(runtimeEnv))
	}
	return
}
