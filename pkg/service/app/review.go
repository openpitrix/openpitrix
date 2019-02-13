// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package app

import (
	"context"
	"time"

	"openpitrix.io/openpitrix/pkg/constants"
	"openpitrix.io/openpitrix/pkg/db"
	"openpitrix.io/openpitrix/pkg/gerr"
	"openpitrix.io/openpitrix/pkg/models"
	"openpitrix.io/openpitrix/pkg/pi"
	"openpitrix.io/openpitrix/pkg/util/ctxutil"
)

var reviewSupportRoles = []string{
	constants.RoleIsv,
	constants.RoleDevelopAdmin,
	constants.RoleBusinessAdmin,
}

func submitAppVersionReview(ctx context.Context, version *models.AppVersion) error {
	var (
		s        = ctxutil.GetSender(ctx)
		operator = s.UserId
		role     = constants.RoleDeveloper
		status   = constants.StatusSubmitted
		action   = Submit
		message  = ""
	)

	err := checkAppVersionHandlePermission(ctx, action, version)
	if err != nil {
		return err
	}
	versionReview := models.NewAppVersionReview(version.VersionId, version.AppId, status, s.GetOwnerPath())
	versionReview.UpdatePhase(role, status, operator, message)

	_, err = pi.Global().DB(ctx).
		InsertInto(constants.TableAppVersionReview).
		Record(versionReview).
		Exec()
	if err != nil {
		return gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorCreateResourceFailed, versionReview.ReviewId)
	}
	err = updateVersionStatus(ctx, version, status, map[string]interface{}{
		constants.ColumnReviewId: versionReview.ReviewId,
	})
	if err != nil {
		return err
	}
	version.ReviewId = versionReview.ReviewId
	err = addAppVersionAudit(ctx, version, status, role, message)
	if err != nil {
		return err
	}
	return nil
}

func getAppVersionReview(ctx context.Context, version *models.AppVersion) (*models.AppVersionReview, error) {
	var err error
	var versionReview = &models.AppVersionReview{}
	var reviewId = version.ReviewId
	if reviewId == "" {
		versionReview = models.NewAppVersionReview(version.VersionId, version.AppId, constants.StatusSubmitted, version.OwnerPath)

		_, err = pi.Global().DB(ctx).
			InsertInto(constants.TableAppVersionReview).
			Record(versionReview).
			Exec()
		if err != nil {
			return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorCreateResourceFailed, versionReview.ReviewId)
		}
	} else {
		err = pi.Global().DB(ctx).
			Select(models.AppVersionReviewColumns...).
			From(constants.TableAppVersionReview).
			Where(db.Eq(constants.ColumnReviewId, reviewId)).LoadOne(versionReview)
		if err != nil {
			return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorDescribeResourceFailed, reviewId)
		}
	}
	return versionReview, nil
}

var reviewActionStatusMap = map[Action]string{
	Review: constants.StatusInReview,
	Pass:   constants.StatusPassed,
	Reject: constants.StatusRejected,
	Cancel: constants.StatusDraft,
}

func execAppVersionReview(ctx context.Context, version *models.AppVersion, action Action, role, message string) error {
	var (
		s        = ctxutil.GetSender(ctx)
		operator = s.UserId
		status   = reviewActionStatusMap[action]
	)
	err := checkAppVersionHandlePermission(ctx, action, version)
	if err != nil {
		return err
	}

	versionReview, err := getAppVersionReview(ctx, version)
	if err != nil {
		return err
	}
	p, ok := versionReview.Phase[role]
	switch action {
	case Review:
		if ok {
			return gerr.New(ctx,
				gerr.FailedPrecondition, gerr.ErrorAppVersionIncorrectStatus, version.VersionId, version.Status)
		}
		switch role {
		case constants.RoleBusinessAdmin, constants.RoleDevelopAdmin:
			if _, ok = versionReview.Phase[constants.RoleIsv]; !ok {
				return gerr.New(ctx,
					gerr.FailedPrecondition, gerr.ErrorAppVersionIncorrectStatus, version.VersionId, version.Status)
			}
		}
		if role == constants.RoleDevelopAdmin {
			if _, ok = versionReview.Phase[constants.RoleBusinessAdmin]; !ok {
				return gerr.New(ctx,
					gerr.FailedPrecondition, gerr.ErrorAppVersionIncorrectStatus, version.VersionId, version.Status)
			}
		}
	case Pass, Reject:
		if p.Status != constants.StatusInReview {
			return gerr.New(ctx,
				gerr.FailedPrecondition, gerr.ErrorAppVersionIncorrectStatus, version.VersionId, version.Status)
		}
	case Cancel:

	}

	versionReview.UpdatePhase(role, status, operator, message)

	var reviewStatus = ""
	switch role {
	case constants.RoleIsv:
		reviewStatus = "isv-"
	case constants.RoleBusinessAdmin:
		reviewStatus = "business-"
	case constants.RoleDevelopAdmin:
		reviewStatus = "dev-"
	}

	var reviewer = ""
	if action == Review {
		reviewer = operator
	}

	_, err = pi.Global().DB(ctx).
		Update(constants.TableAppVersionReview).
		Set(constants.ColumnStatus, reviewStatus+status).
		Set(constants.ColumnStatusTime, time.Now()).
		Set(constants.ColumnPhase, versionReview.Phase).
		Set(constants.ColumnReviewer, reviewer).
		Where(db.Eq(constants.ColumnReviewId, versionReview.ReviewId)).
		Exec()
	if err != nil {
		return gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorUpdateResourceFailed, versionReview.ReviewId)
	}
	return nil
}

func cancelAppVersionReview(ctx context.Context, version *models.AppVersion, role string) error {
	err := execAppVersionReview(ctx, version, Cancel, role, "")
	if err != nil {
		return err
	}
	err = updateVersionStatus(ctx, version, constants.StatusDraft, map[string]interface{}{
		constants.ColumnReviewId: "",
	})
	if err != nil {
		return err
	}
	err = addAppVersionAudit(ctx, version, constants.StatusDraft, role, "")
	if err != nil {
		return err
	}
	return nil
}

func startAppVersionReview(ctx context.Context, version *models.AppVersion, role string) error {
	err := execAppVersionReview(ctx, version, Review, role, "")
	if err != nil {
		return err
	}
	if role == constants.RoleIsv {
		err = updateVersionStatus(ctx, version, constants.StatusInReview)
		if err != nil {
			return err
		}
		err = addAppVersionAudit(ctx, version, constants.StatusInReview, role, "")
		if err != nil {
			return err
		}
	}
	return nil
}

func passAppVersionReview(ctx context.Context, version *models.AppVersion, role string) error {
	err := execAppVersionReview(ctx, version, Pass, role, "")
	if err != nil {
		return err
	}
	if role == constants.RoleDevelopAdmin {
		err = updateVersionStatus(ctx, version, constants.StatusPassed)
		if err != nil {
			return err
		}
		err = addAppVersionAudit(ctx, version, constants.StatusPassed, role, "")
		if err != nil {
			return err
		}
	}
	return nil
}

func rejectAppVersionReview(ctx context.Context, version *models.AppVersion, role string, message string) error {
	err := execAppVersionReview(ctx, version, Reject, role, message)
	if err != nil {
		return err
	}
	err = updateVersionStatus(ctx, version, constants.StatusRejected)
	if err != nil {
		return err
	}
	return addAppVersionAudit(ctx, version, constants.StatusRejected, role, message)
}
