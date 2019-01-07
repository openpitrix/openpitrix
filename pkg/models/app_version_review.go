// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package models

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	"openpitrix.io/openpitrix/pkg/db"
	"openpitrix.io/openpitrix/pkg/pb"
	"openpitrix.io/openpitrix/pkg/sender"
	"openpitrix.io/openpitrix/pkg/util/idutil"
	"openpitrix.io/openpitrix/pkg/util/pbutil"
)

func NewAppVersionReviewId() string {
	return idutil.GetUuid("appvr-")
}

type AppVersionReviewPhase struct {
	Status     string
	Operator   string
	Role       string
	Message    string
	ReviewTime time.Time
	StatusTime time.Time
}
type AppVersionReviewPhases map[string]AppVersionReviewPhase

func (p *AppVersionReviewPhases) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	b, ok := value.([]byte)
	if ok {
		return json.Unmarshal(b, p)
	}
	s, _ := value.(string)
	return json.Unmarshal([]byte(s), p)
}

// Value implements the driver Valuer interface.
func (p AppVersionReviewPhases) Value() (driver.Value, error) {
	b, err := json.Marshal(p)
	if err != nil {
		return nil, err
	}
	return string(b), nil
}

//func (p AppVersionReviewPhases) String() string {
//	return jsonutil.ToString(p)
//}

type AppVersionReview struct {
	ReviewId   string
	VersionId  string
	AppId      string
	OwnerPath  sender.OwnerPath
	Owner      string
	Status     string
	Phase      AppVersionReviewPhases
	StatusTime time.Time
	Reviewer   string
}

type AppVersionReviews []*AppVersionReview

func (avr *AppVersionReview) UpdatePhase(role, status, operator, message string) {
	p, ok := avr.Phase[role]
	if !ok {
		p = AppVersionReviewPhase{
			ReviewTime: time.Now(),
		}
	}

	p.Status = status
	p.Operator = operator
	p.Role = role
	p.Message = message
	p.StatusTime = time.Now()

	avr.Phase[role] = p
}

var AppVersionReviewColumns = db.GetColumnsFromStruct(&AppVersionReview{})

func NewAppVersionReview(versionId, appId, status string, ownerPath sender.OwnerPath) *AppVersionReview {
	return &AppVersionReview{
		ReviewId:   NewAppVersionReviewId(),
		VersionId:  versionId,
		AppId:      appId,
		Status:     status,
		Owner:      ownerPath.Owner(),
		OwnerPath:  ownerPath,
		Phase:      make(AppVersionReviewPhases),
		StatusTime: time.Now(),
	}
}

func AppVersionReviewToPb(versionReview *AppVersionReview) *pb.AppVersionReview {
	if versionReview == nil {
		return nil
	}
	pbAppVersionReview := pb.AppVersionReview{}
	pbAppVersionReview.ReviewId = pbutil.ToProtoString(versionReview.ReviewId)
	pbAppVersionReview.VersionId = pbutil.ToProtoString(versionReview.VersionId)
	pbAppVersionReview.AppId = pbutil.ToProtoString(versionReview.AppId)
	pbAppVersionReview.Status = pbutil.ToProtoString(versionReview.Status)
	pbAppVersionReview.StatusTime = pbutil.ToProtoTimestamp(versionReview.StatusTime)
	pbAppVersionReview.Reviewer = pbutil.ToProtoString(versionReview.Reviewer)

	pbAppVersionReview.Phase = make(map[string]*pb.AppVersionReviewPhase)
	for role, p := range versionReview.Phase {
		pbPhase := pb.AppVersionReviewPhase{}
		pbPhase.Role = pbutil.ToProtoString(p.Role)
		pbPhase.Status = pbutil.ToProtoString(p.Status)
		pbPhase.Operator = pbutil.ToProtoString(p.Operator)
		pbPhase.Message = pbutil.ToProtoString(p.Message)
		pbPhase.StatusTime = pbutil.ToProtoTimestamp(p.StatusTime)
		pbPhase.ReviewTime = pbutil.ToProtoTimestamp(p.ReviewTime)

		pbAppVersionReview.Phase[role] = &pbPhase
	}
	return &pbAppVersionReview
}

func AppVersionReviewsToPbs(versionReviews AppVersionReviews) (pbAppVersionReviews []*pb.AppVersionReview) {
	for _, appVersionAudit := range versionReviews {
		pbAppVersionReviews = append(pbAppVersionReviews, AppVersionReviewToPb(appVersionAudit))
	}
	return
}
