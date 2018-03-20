// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package models

import (
	"openpitrix.io/openpitrix/pkg/pb"
	"openpitrix.io/openpitrix/pkg/utils"
)

const RuntimeEnvLabelTableName = "runtime_env_label"

func NewRuntimeEnvLabelId() string {
	return utils.GetUuid("rel-")
}

type RuntimeEnvLabel struct {
	RuntimeEnvLabelId string
	RuntimeEnvId      string
	LabelKey          string
	LabelValue        string
}

var RuntimeEnvLabelColumns = GetColumnsFromStruct(&RuntimeEnvLabel{})

func NewRuntimeEnvLabel(runtimeEnvId, labelKey, labelValue string) *RuntimeEnvLabel {
	return &RuntimeEnvLabel{
		RuntimeEnvLabelId: NewRuntimeEnvLabelId(),
		RuntimeEnvId:      runtimeEnvId,
		LabelKey:          labelKey,
		LabelValue:        labelValue,
	}
}

func RuntimeEnvLabelToPb(runtimeEnvLabel *RuntimeEnvLabel) *pb.RuntimeEnvLabel {
	return &pb.RuntimeEnvLabel{
		RuntimeEnvLabelId: utils.ToProtoString(runtimeEnvLabel.RuntimeEnvLabelId),
		RuntimeEnvId:      utils.ToProtoString(runtimeEnvLabel.RuntimeEnvId),
		LabelKey:          utils.ToProtoString(runtimeEnvLabel.LabelKey),
		LabelValue:        utils.ToProtoString(runtimeEnvLabel.LabelValue),
	}
}

func RuntimeEnvLabelsToPbs(runtimeEnvLabels []*RuntimeEnvLabel) (pbRuntimeEnvLabels []*pb.RuntimeEnvLabel) {
	for _, runtimeEnvLabel := range runtimeEnvLabels {
		pbRuntimeEnvLabels = append(pbRuntimeEnvLabels, RuntimeEnvLabelToPb(runtimeEnvLabel))
	}
	return pbRuntimeEnvLabels
}
