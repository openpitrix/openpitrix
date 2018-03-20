// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package models

const RuntimeEnvAttachedCredentialTableName = "runtime_env_attached_credential"

type RuntimeEnvAttachedCredential struct {
	RuntimeEnvId           string
	RuntimeEnvCredentialId string
}

var RuntimeEnvAttachedCredentialColumns = GetColumnsFromStruct(&RuntimeEnvAttachedCredential{})

func NewRuntimeEnvAttachedCredential(runtimeEnvId, runtimeEnvCredentialId string) *RuntimeEnvAttachedCredential {
	return &RuntimeEnvAttachedCredential{
		RuntimeEnvId:           runtimeEnvId,
		RuntimeEnvCredentialId: runtimeEnvCredentialId,
	}
}
