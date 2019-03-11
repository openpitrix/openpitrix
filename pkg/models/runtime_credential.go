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

const RuntimeCredentialPrifix = "runtimec-"

type RuntimeCredential struct {
	RuntimeCredentialId      string
	Name                     string
	Description              string
	RuntimeUrl               string
	RuntimeCredentialContent string
	Owner                    string
	OwnerPath                sender.OwnerPath
	Provider                 string
	Status                   string
	Debug                    bool
	CreateTime               time.Time
	StatusTime               time.Time
}

func NewRuntimeCredentialId() string {
	return idutil.GetUuid(RuntimeCredentialPrifix)
}

var RuntimeCredentialColumns = db.GetColumnsFromStruct(&RuntimeCredential{})

func NewRuntimeCredential(runtimeCredentialId, name, description, runtimeUrl, runtimeCredentialContent, provider string, ownerPath sender.OwnerPath, debug bool) *RuntimeCredential {
	if len(runtimeCredentialId) == 0 {
		runtimeCredentialId = NewRuntimeCredentialId()
	}
	return &RuntimeCredential{
		RuntimeCredentialId:      runtimeCredentialId,
		Name:                     name,
		Description:              description,
		RuntimeUrl:               runtimeUrl,
		RuntimeCredentialContent: runtimeCredentialContent,
		Provider:                 provider,
		Owner:                    ownerPath.Owner(),
		OwnerPath:                ownerPath,
		Status:                   constants.StatusActive,
		CreateTime:               time.Now(),
		StatusTime:               time.Now(),
		Debug:                    debug,
	}
}

func RuntimeCredentialToPb(runtimeCredential *RuntimeCredential) *pb.RuntimeCredential {
	pbRuntimeCredential := pb.RuntimeCredential{}
	pbRuntimeCredential.RuntimeCredentialId = pbutil.ToProtoString(runtimeCredential.RuntimeCredentialId)
	pbRuntimeCredential.Name = pbutil.ToProtoString(runtimeCredential.Name)
	pbRuntimeCredential.Description = pbutil.ToProtoString(runtimeCredential.Description)
	pbRuntimeCredential.RuntimeUrl = pbutil.ToProtoString(runtimeCredential.RuntimeUrl)
	pbRuntimeCredential.RuntimeCredentialContent = pbutil.ToProtoString(runtimeCredential.RuntimeCredentialContent)
	pbRuntimeCredential.Provider = pbutil.ToProtoString(runtimeCredential.Provider)
	pbRuntimeCredential.OwnerPath = runtimeCredential.OwnerPath.ToProtoString()
	pbRuntimeCredential.Owner = pbutil.ToProtoString(runtimeCredential.Owner)
	pbRuntimeCredential.Status = pbutil.ToProtoString(runtimeCredential.Status)
	pbRuntimeCredential.CreateTime = pbutil.ToProtoTimestamp(runtimeCredential.CreateTime)
	pbRuntimeCredential.StatusTime = pbutil.ToProtoTimestamp(runtimeCredential.StatusTime)
	pbRuntimeCredential.Debug = pbutil.ToProtoBool(runtimeCredential.Debug)
	return &pbRuntimeCredential
}

func PbToRuntimeCredential(pbRuntimeCredential *pb.RuntimeCredential) *RuntimeCredential {
	ownerPath := sender.OwnerPath(pbRuntimeCredential.GetOwnerPath().GetValue())
	runtimeCredential := RuntimeCredential{}
	runtimeCredential.RuntimeCredentialId = pbRuntimeCredential.GetRuntimeCredentialId().GetValue()
	runtimeCredential.Name = pbRuntimeCredential.GetName().GetValue()
	runtimeCredential.Description = pbRuntimeCredential.GetDescription().GetValue()
	runtimeCredential.RuntimeUrl = pbRuntimeCredential.GetRuntimeUrl().GetValue()
	runtimeCredential.RuntimeCredentialContent = pbRuntimeCredential.GetRuntimeCredentialContent().GetValue()
	runtimeCredential.Provider = pbRuntimeCredential.GetProvider().GetValue()
	runtimeCredential.OwnerPath = ownerPath
	runtimeCredential.Owner = ownerPath.Owner()
	runtimeCredential.Status = pbRuntimeCredential.GetStatus().GetValue()
	runtimeCredential.Debug = pbRuntimeCredential.GetDebug().GetValue()
	runtimeCredential.CreateTime = pbutil.GetTime(pbRuntimeCredential.CreateTime)
	runtimeCredential.StatusTime = pbutil.GetTime(pbRuntimeCredential.StatusTime)
	return &runtimeCredential
}

func RuntimeCredentialToPbs(runtimeCredentials []*RuntimeCredential) (pbRuntimeCredentials []*pb.RuntimeCredential) {
	for _, runtimeCredential := range runtimeCredentials {
		pbRuntimeCredentials = append(pbRuntimeCredentials, RuntimeCredentialToPb(runtimeCredential))
	}
	return
}

func RuntimeCredentialMap(runtimeCredentials []*RuntimeCredential) map[string]*RuntimeCredential {
	credentialMap := make(map[string]*RuntimeCredential)
	for _, credential := range runtimeCredentials {
		credentialMap[credential.RuntimeCredentialId] = credential
	}
	return credentialMap
}
