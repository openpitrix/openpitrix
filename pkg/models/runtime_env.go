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

const RuntimeEnvTableName = "runtime_env"
const RuntimeEnvIdPrifix = "re-"

func NewRuntimeEnvId() string {
	return utils.GetUuid(RuntimeEnvIdPrifix)
}

type RuntimeEnv struct {
	RuntimeEnvId  string
	Name          string
	Description   string
	RuntimeEnvUrl string
	Owner         string
	Status        string
	CreateTime    time.Time
	StatusTime    time.Time
}

var RuntimeEnvColumns = GetColumnsFromStruct(&RuntimeEnv{})

func NewRuntimeEnv(name, description, runtimeEnvUrl, owner string) *RuntimeEnv {
	return &RuntimeEnv{
		RuntimeEnvId:  NewRuntimeEnvId(),
		Name:          name,
		Description:   description,
		RuntimeEnvUrl: runtimeEnvUrl,
		Owner:         owner,
		Status:        constants.StatusActive,
		CreateTime:    time.Now(),
		StatusTime:    time.Now(),
	}
}

func RuntimeEnvToPb(runtimeEnv *RuntimeEnv) *pb.RuntimeEnv {
	pbRuntimeEnv := pb.RuntimeEnv{}
	pbRuntimeEnv.RuntimeEnvId = utils.ToProtoString(runtimeEnv.RuntimeEnvId)
	pbRuntimeEnv.Name = utils.ToProtoString(runtimeEnv.Name)
	pbRuntimeEnv.Description = utils.ToProtoString(runtimeEnv.Description)
	pbRuntimeEnv.RuntimeEnvUrl = utils.ToProtoString(runtimeEnv.RuntimeEnvUrl)
	pbRuntimeEnv.Owner = utils.ToProtoString(runtimeEnv.Owner)
	pbRuntimeEnv.Status = utils.ToProtoString(runtimeEnv.Status)
	pbRuntimeEnv.CreateTime = utils.ToProtoTimestamp(runtimeEnv.CreateTime)
	pbRuntimeEnv.StatusTime = utils.ToProtoTimestamp(runtimeEnv.StatusTime)
	return &pbRuntimeEnv
}

func RuntimeEnvToPbs(runtimeEnvs []*RuntimeEnv) (pbRuntimeEnvs []*pb.RuntimeEnv) {
	for _, runtimeEnv := range runtimeEnvs {
		pbRuntimeEnvs = append(pbRuntimeEnvs, RuntimeEnvToPb(runtimeEnv))
	}
	return
}
