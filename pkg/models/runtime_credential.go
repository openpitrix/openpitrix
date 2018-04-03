// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package models

import (
	"encoding/json"
	"time"

	"openpitrix.io/openpitrix/pkg/constants"
	"openpitrix.io/openpitrix/pkg/logger"
	"openpitrix.io/openpitrix/pkg/pb"
	"openpitrix.io/openpitrix/pkg/utils"
)

const RuntimeCredentialTableName = "runtime_credential"
const RuntimeCredentialPrifix = "runtimec-"

type RuntimeCredential struct {
	RuntimeCredentialId string
	Name                string
	Description         string
	Owner               string
	Content             string
	Status              string
	CreateTime          time.Time
	StatusTime          time.Time
}

func NewRuntimeCrentialId() string {
	return utils.GetUuid(RuntimeCredentialPrifix)
}

var RuntimeCredentialColumns = GetColumnsFromStruct(&RuntimeCredential{})

func NewRuntimeCredential(name, description, owner string, content map[string]string) *RuntimeCredential {
	return &RuntimeCredential{
		RuntimeCredentialId: NewRuntimeCrentialId(),
		Name:                name,
		Description:         description,
		Owner:               owner,
		Content:             RuntimeCredentialContentMapToString(content),
		Status:              constants.StatusActive,
		CreateTime:          time.Now(),
		StatusTime:          time.Now(),
	}
}

func RuntimeCredentialToPb(runtimeCredential *RuntimeCredential) *pb.RuntimeCredential {
	pbRuntimeCredential := pb.RuntimeCredential{}
	pbRuntimeCredential.RuntimeCredentialId = utils.ToProtoString(runtimeCredential.RuntimeCredentialId)
	pbRuntimeCredential.Name = utils.ToProtoString(runtimeCredential.Name)
	pbRuntimeCredential.Description = utils.ToProtoString(runtimeCredential.Description)
	pbRuntimeCredential.Owner = utils.ToProtoString(runtimeCredential.Owner)
	pbRuntimeCredential.Status = utils.ToProtoString(runtimeCredential.Status)
	pbRuntimeCredential.CreateTime = utils.ToProtoTimestamp(runtimeCredential.CreateTime)
	pbRuntimeCredential.StatusTime = utils.ToProtoTimestamp(runtimeCredential.StatusTime)
	pbRuntimeCredential.Content = RuntimeCredentialContentStringToMap(runtimeCredential.Content)
	return &pbRuntimeCredential
}

func RuntimeCredentialToPbs(runtimeCredentials []*RuntimeCredential) (pbRuntimeCredentials []*pb.RuntimeCredential) {
	for _, runtimeCredential := range runtimeCredentials {
		pbRuntimeCredentials = append(pbRuntimeCredentials, RuntimeCredentialToPb(runtimeCredential))
	}
	return
}

func RuntimeCredentialContentStringToMap(stringContent string) map[string]string {
	var mapContent map[string]string
	err := json.Unmarshal([]byte(stringContent), &mapContent)
	if err != nil {
		logger.Errorf("unexpected error, unmarshal fail: %v ", stringContent)
		panic(err)
	}
	return mapContent
}

func RuntimeCredentialContentMapToString(mapContent map[string]string) string {
	stringContent, err := json.Marshal(mapContent)
	if err != nil {
		logger.Errorf("unexpected error, marshal map[string]string fail: %v ", mapContent)
		panic(err)
	}
	return string(stringContent)
}
