// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package models

import (
	"time"

	"openpitrix.io/openpitrix/pkg/constants"
	"openpitrix.io/openpitrix/pkg/pb"
	"openpitrix.io/openpitrix/pkg/util/idutil"
	"openpitrix.io/openpitrix/pkg/util/pbutil"
)

const RuntimeIdPrefix = "runtime-"

func NewRuntimeId() string {
	return idutil.GetUuid(RuntimeIdPrefix)
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

var RuntimeColumnsWithTablePrefix = GetColumnsFromStructWithPrefix(constants.TableRuntime, &Runtime{})
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
	pbRuntime.RuntimeId = pbutil.ToProtoString(runtime.RuntimeId)
	pbRuntime.Name = pbutil.ToProtoString(runtime.Name)
	pbRuntime.Description = pbutil.ToProtoString(runtime.Description)
	pbRuntime.Provider = pbutil.ToProtoString(runtime.Provider)
	pbRuntime.RuntimeUrl = pbutil.ToProtoString(runtime.RuntimeUrl)
	pbRuntime.Zone = pbutil.ToProtoString(runtime.Zone)
	pbRuntime.Owner = pbutil.ToProtoString(runtime.Owner)
	pbRuntime.Status = pbutil.ToProtoString(runtime.Status)
	pbRuntime.CreateTime = pbutil.ToProtoTimestamp(runtime.CreateTime)
	pbRuntime.StatusTime = pbutil.ToProtoTimestamp(runtime.StatusTime)
	return &pbRuntime
}

func PbToRuntime(pbRuntime *pb.Runtime) *Runtime {
	runtime := Runtime{}
	runtime.RuntimeId = pbRuntime.GetRuntimeId().GetValue()
	runtime.Name = pbRuntime.GetName().GetValue()
	runtime.Description = pbRuntime.GetDescription().GetValue()
	runtime.Provider = pbRuntime.GetProvider().GetValue()
	runtime.RuntimeUrl = pbRuntime.GetRuntimeUrl().GetValue()
	runtime.Zone = pbRuntime.GetZone().GetValue()
	runtime.Owner = pbRuntime.GetOwner().GetValue()
	runtime.Status = pbRuntime.GetStatus().GetValue()
	runtime.CreateTime = pbutil.FromProtoTimestamp(pbRuntime.GetCreateTime())
	runtime.StatusTime = pbutil.FromProtoTimestamp(pbRuntime.GetStatusTime())
	return &runtime
}

func RuntimeToPbs(runtimes []*Runtime) (pbRuntimes []*pb.Runtime) {
	for _, runtime := range runtimes {
		pbRuntimes = append(pbRuntimes, RuntimeToPb(runtime))
	}
	return
}
