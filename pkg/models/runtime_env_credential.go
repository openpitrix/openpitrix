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

const RuntimeEnvCredentialTableName = "runtime_env_credential"

type RuntimeEnvCredential struct {
	RuntimeEnvCredentialId string
	Name                   string
	Description            string
	Owner                  string
	Content                string
	Status                 string
	CreateTime             time.Time
	StatusTime             time.Time
}

func NewRuntimeEnvCrentialId() string {
	return utils.GetUuid("rec-")
}

var RuntimeEnvCredentialColumns = GetColumnsFromStruct(&RuntimeEnvCredential{})

func NewRuntimeEnvCredential(name, description, owner string, content map[string]string) *RuntimeEnvCredential {
	return &RuntimeEnvCredential{
		RuntimeEnvCredentialId: NewRuntimeEnvCrentialId(),
		Name:        name,
		Description: description,
		Owner:       owner,
		Content:     RuntimeEnvCredentialContentMapToString(content),
		Status:      constants.StatusActive,
		CreateTime:  time.Now(),
		StatusTime:  time.Now(),
	}
}

func RuntimeEnvCredentialToPb(runtimeEnvCredential *RuntimeEnvCredential) *pb.RuntimeEnvCredential {
	pbRuntimeEnvCredential := pb.RuntimeEnvCredential{}
	pbRuntimeEnvCredential.RuntimeEnvCredentialId = utils.ToProtoString(runtimeEnvCredential.RuntimeEnvCredentialId)
	pbRuntimeEnvCredential.Name = utils.ToProtoString(runtimeEnvCredential.Name)
	pbRuntimeEnvCredential.Description = utils.ToProtoString(runtimeEnvCredential.Description)
	pbRuntimeEnvCredential.Owner = utils.ToProtoString(runtimeEnvCredential.Owner)
	pbRuntimeEnvCredential.Status = utils.ToProtoString(runtimeEnvCredential.Status)
	pbRuntimeEnvCredential.CreateTime = utils.ToProtoTimestamp(runtimeEnvCredential.CreateTime)
	pbRuntimeEnvCredential.StatusTime = utils.ToProtoTimestamp(runtimeEnvCredential.StatusTime)
	pbRuntimeEnvCredential.Content = RuntimeEnvCredentialContentStringToMap(runtimeEnvCredential.Content)
	return &pbRuntimeEnvCredential
}

func RuntimeEnvCredentialToPbs(runtimeEnvCredentials []*RuntimeEnvCredential) (pbRuntimeEnvCredentials []*pb.RuntimeEnvCredential) {
	for _, runtimeEnvCredential := range runtimeEnvCredentials {
		pbRuntimeEnvCredentials = append(pbRuntimeEnvCredentials, RuntimeEnvCredentialToPb(runtimeEnvCredential))
	}
	return
}

func RuntimeEnvCredentialContentStringToMap(stringContent string) map[string]string {
	var mapContent map[string]string
	err := json.Unmarshal([]byte(stringContent), &mapContent)
	if err != nil {
		logger.Errorf("unexpected error, unmarshal fail: %v ", stringContent)
		panic(err)
	}
	return mapContent
}

func RuntimeEnvCredentialContentMapToString(mapContent map[string]string) string {
	stringContent, err := json.Marshal(mapContent)
	if err != nil {
		logger.Errorf("unexpected error, marshal map[string]string fail: %v ", mapContent)
		panic(err)
	}
	return string(stringContent)
}
