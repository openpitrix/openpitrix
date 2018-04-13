// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package models

import (
	"time"

	"openpitrix.io/openpitrix/pkg/constants"
	"openpitrix.io/openpitrix/pkg/pb"
	"openpitrix.io/openpitrix/pkg/utils"
	"openpitrix.io/openpitrix/pkg/utils/idtool"
)

const RuntimeTableName = "runtime"
const RuntimeIdPrefix = "runtime-"

func NewRuntimeId() string {
	return idtool.GetUuid(RuntimeIdPrefix)
}

type Runtime struct {
	RuntimeId           string
	Name                string
	Description         string
	Provider            string
	RuntimeUrl          string
	Zone                string
	RuntimeCredentialId string
	Owner               string
	Status              string
	CreateTime          time.Time
	StatusTime          time.Time
}

var RuntimeColumnsWithTablePrefix = GetColumnsFromStructWithPrefix(RuntimeTableName, &Runtime{})
var RuntimeColumns = GetColumnsFromStruct(&Runtime{})

func NewRuntime(name, description, provider, runtimeUrl, runtimeCredentialId, zone, owner string) *Runtime {
	return &Runtime{
		RuntimeId:           NewRuntimeId(),
		Name:                name,
		Description:         description,
		Provider:            provider,
		RuntimeUrl:          runtimeUrl,
		Zone:                zone,
		RuntimeCredentialId: runtimeCredentialId,
		Owner:               owner,
		Status:              constants.StatusActive,
		CreateTime:          time.Now(),
		StatusTime:          time.Now(),
	}
}

func RuntimeToPb(runtime *Runtime) *pb.Runtime {
	pbRuntime := pb.Runtime{}
	pbRuntime.RuntimeId = utils.ToProtoString(runtime.RuntimeId)
	pbRuntime.Name = utils.ToProtoString(runtime.Name)
	pbRuntime.Description = utils.ToProtoString(runtime.Description)
	pbRuntime.Provider = utils.ToProtoString(runtime.Provider)
	pbRuntime.RuntimeUrl = utils.ToProtoString(runtime.RuntimeUrl)
	pbRuntime.Zone = utils.ToProtoString(runtime.Zone)
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
