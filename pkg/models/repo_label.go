// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package models

import (
	"time"

	"openpitrix.io/openpitrix/pkg/db"
	"openpitrix.io/openpitrix/pkg/pb"
	"openpitrix.io/openpitrix/pkg/util/idutil"
	"openpitrix.io/openpitrix/pkg/util/pbutil"
)

func NewRepoLabelId() string {
	return idutil.GetUuid("repol-")
}

type RepoLabel struct {
	RepoLabelId string
	RepoId      string
	LabelKey    string
	LabelValue  string

	CreateTime time.Time
}

var RepoLabelColumns = db.GetColumnsFromStruct(&RepoLabel{})

func NewRepoLabel(repoId, labelKey, labelValue string) *RepoLabel {
	return &RepoLabel{
		RepoLabelId: NewRepoLabelId(),
		RepoId:      repoId,
		LabelKey:    labelKey,
		LabelValue:  labelValue,

		CreateTime: time.Now(),
	}
}

func RepoLabelToPb(repoLabel *RepoLabel) *pb.RepoLabel {
	return &pb.RepoLabel{
		LabelKey:   pbutil.ToProtoString(repoLabel.LabelKey),
		LabelValue: pbutil.ToProtoString(repoLabel.LabelValue),
		CreateTime: pbutil.ToProtoTimestamp(repoLabel.CreateTime),
	}
}

func RepoLabelsToPbs(repoLabels []*RepoLabel) (pbRepoLabels []*pb.RepoLabel) {
	for _, repoLabel := range repoLabels {
		pbRepoLabels = append(pbRepoLabels, RepoLabelToPb(repoLabel))
	}
	return
}

func RepoLabelsMap(repoLabels []*RepoLabel) map[string][]*RepoLabel {
	labelsMap := make(map[string][]*RepoLabel)
	for _, l := range repoLabels {
		repoId := l.RepoId
		labelsMap[repoId] = append(labelsMap[repoId], l)
	}
	return labelsMap
}
