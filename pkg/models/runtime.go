// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package models

import (
	"time"

	"openpitrix.io/openpitrix/pkg/constants"
	"openpitrix.io/openpitrix/pkg/db"
	"openpitrix.io/openpitrix/pkg/pb"
	"openpitrix.io/openpitrix/pkg/sender"
	"openpitrix.io/openpitrix/pkg/util/idutil"
	"openpitrix.io/openpitrix/pkg/util/pbutil"
)

const RuntimeIdPrefix = "runtime-"

func NewRuntimeId() string {
	return idutil.GetUuid(RuntimeIdPrefix)
}

type RuntimeDetails struct {
	Runtime
	RuntimeCredential
}

type Runtime struct {
	RuntimeId           string
	Name                string
	Description         string
	Provider            string
	Zone                string
	RuntimeCredentialId string
	Owner               string
	OwnerPath           sender.OwnerPath
	Status              string
	Debug               bool
	CreateTime          time.Time
	StatusTime          time.Time
}

var RuntimeColumns = db.GetColumnsFromStruct(&Runtime{})

func NewRuntime(runtimeId, name, description, provider, runtimeCredentialId, zone string, ownerPath sender.OwnerPath, debug bool) *Runtime {
	if len(runtimeId) == 0 {
		runtimeId = NewRuntimeId()
	}
	return &Runtime{
		RuntimeId:           runtimeId,
		Name:                name,
		Description:         description,
		Provider:            provider,
		Zone:                zone,
		RuntimeCredentialId: runtimeCredentialId,
		Owner:               ownerPath.Owner(),
		OwnerPath:           ownerPath,
		Status:              constants.StatusActive,
		CreateTime:          time.Now(),
		StatusTime:          time.Now(),
		Debug:               debug,
	}
}

func RuntimeToPb(runtime *Runtime) *pb.Runtime {
	pbRuntime := pb.Runtime{}
	pbRuntime.RuntimeId = pbutil.ToProtoString(runtime.RuntimeId)
	pbRuntime.Name = pbutil.ToProtoString(runtime.Name)
	pbRuntime.Description = pbutil.ToProtoString(runtime.Description)
	pbRuntime.Provider = pbutil.ToProtoString(runtime.Provider)
	pbRuntime.Zone = pbutil.ToProtoString(runtime.Zone)
	pbRuntime.RuntimeCredentialId = pbutil.ToProtoString(runtime.RuntimeCredentialId)
	pbRuntime.OwnerPath = runtime.OwnerPath.ToProtoString()
	pbRuntime.Status = pbutil.ToProtoString(runtime.Status)
	pbRuntime.Debug = pbutil.ToProtoBool(runtime.Debug)
	pbRuntime.CreateTime = pbutil.ToProtoTimestamp(runtime.CreateTime)
	pbRuntime.StatusTime = pbutil.ToProtoTimestamp(runtime.StatusTime)
	return &pbRuntime
}

func PbToRuntime(pbRuntime *pb.Runtime) *Runtime {
	ownerPath := sender.OwnerPath(pbRuntime.GetOwnerPath().GetValue())
	runtime := Runtime{}
	runtime.RuntimeId = pbRuntime.GetRuntimeId().GetValue()
	runtime.Name = pbRuntime.GetName().GetValue()
	runtime.Description = pbRuntime.GetDescription().GetValue()
	runtime.Provider = pbRuntime.GetProvider().GetValue()
	runtime.Zone = pbRuntime.GetZone().GetValue()
	runtime.RuntimeCredentialId = pbRuntime.GetRuntimeCredentialId().GetValue()
	runtime.OwnerPath = ownerPath
	runtime.Owner = ownerPath.Owner()
	runtime.Status = pbRuntime.GetStatus().GetValue()
	runtime.Debug = pbRuntime.GetDebug().GetValue()
	runtime.CreateTime = pbutil.GetTime(pbRuntime.GetCreateTime())
	runtime.StatusTime = pbutil.GetTime(pbRuntime.GetStatusTime())
	return &runtime
}

func RuntimeToPbs(runtimes []*Runtime) (pbRuntimes []*pb.Runtime) {
	for _, runtime := range runtimes {
		pbRuntimes = append(pbRuntimes, RuntimeToPb(runtime))
	}
	return
}
