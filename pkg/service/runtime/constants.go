// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package runtime

const (
	RuntimeCredentialIdColumn      = "runtime_credential_id"
	RuntimeCredentialContentColumn = "content"
)

const (
	NameColumn        = "name"
	DescriptionColumn = "description"
	StatusColumn      = "status"
	StatusTimeColumn  = "status_time"
)

const (
	NameMinLength       = "1"
	NameMaxLength       = "255"
	ZoneMinLength       = "1"
	ZoneMaxLength       = "255"
	CredentialMinLength = 1
	LabelKeyMinLength   = "1"
	LabelKeyMaxLength   = "50"
	LabelValueMinLength = "1"
	LabelValueMaxLength = "255"
	LabelKeyFmt         = "^[a-zA-Z]([-_a-zA-Z0-9]*[a-zA-Z0-9])?$"
)
