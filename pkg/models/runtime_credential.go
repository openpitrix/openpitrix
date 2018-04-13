// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package models

import (
	"time"

	"openpitrix.io/openpitrix/pkg/utils/idtool"
)

const RuntimeCredentialTableName = "runtime_credential"
const RuntimeCredentialPrifix = "runtimec-"

type RuntimeCredential struct {
	RuntimeCredentialId string
	Content             string
	CreateTime          time.Time
}

func NewRuntimeCrentialId() string {
	return idtool.GetUuid(RuntimeCredentialPrifix)
}

var RuntimeCredentialColumns = GetColumnsFromStruct(&RuntimeCredential{})

func NewRuntimeCredential(content string) *RuntimeCredential {
	return &RuntimeCredential{
		RuntimeCredentialId: NewRuntimeCrentialId(),
		Content:             content,
		CreateTime:          time.Now(),
	}
}

func RuntimeCredentialMap(runtimeCredentials []*RuntimeCredential) map[string]*RuntimeCredential {
	credentialMap := make(map[string]*RuntimeCredential)
	for _, credential := range runtimeCredentials {
		credentialMap[credential.RuntimeCredentialId] = credential
	}
	return credentialMap
}
