// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package models

import (
	"time"

	"openpitrix.io/openpitrix/pkg/pb"
	"openpitrix.io/openpitrix/pkg/utils"
	"openpitrix.io/openpitrix/pkg/utils/idtool"
)

const RuntimeLabelTableName = "runtime_label"

func NewRuntimeLabelId() string {
	return idtool.GetUuid("runtimel-")
}

type RuntimeLabel struct {
	RuntimeLabelId string
	RuntimeId      string
	LabelKey       string
	LabelValue     string
	CreateTime     time.Time
}

var RuntimeLabelColumns = GetColumnsFromStruct(&RuntimeLabel{})

func NewRuntimeLabel(runtimeId, labelKey, labelValue string) *RuntimeLabel {
	return &RuntimeLabel{
		RuntimeLabelId: NewRuntimeLabelId(),
		RuntimeId:      runtimeId,
		LabelKey:       labelKey,
		LabelValue:     labelValue,
		CreateTime:     time.Now(),
	}
}

func RuntimeLabelToPb(runtimeLabel *RuntimeLabel) *pb.RuntimeLabel {
	return &pb.RuntimeLabel{
		RuntimeLabelId: utils.ToProtoString(runtimeLabel.RuntimeLabelId),
		RuntimeId:      utils.ToProtoString(runtimeLabel.RuntimeId),
		LabelKey:       utils.ToProtoString(runtimeLabel.LabelKey),
		LabelValue:     utils.ToProtoString(runtimeLabel.LabelValue),
		CreateTime:     utils.ToProtoTimestamp(runtimeLabel.CreateTime),
	}
}

func RuntimeLabelsToPbs(runtimeLabels []*RuntimeLabel) (pbRuntimeLabels []*pb.RuntimeLabel) {
	for _, runtimeLabel := range runtimeLabels {
		pbRuntimeLabels = append(pbRuntimeLabels, RuntimeLabelToPb(runtimeLabel))
	}
	return pbRuntimeLabels
}

func RuntimeLabelsMap(runtimeLabels []*RuntimeLabel) map[string][]*RuntimeLabel {
	labelsMap := make(map[string][]*RuntimeLabel)
	for _, l := range runtimeLabels {
		runtimeId := l.RuntimeId
		labelsMap[runtimeId] = append(labelsMap[runtimeId], l)
	}
	return labelsMap
}
