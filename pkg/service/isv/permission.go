// Copyright 2019 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package isv

import (
	"context"

	"openpitrix.io/openpitrix/pkg/gerr"
	"openpitrix.io/openpitrix/pkg/models"
	"openpitrix.io/openpitrix/pkg/util/ctxutil"
)

func CheckAppVendorPermission(ctx context.Context, appVendorUserID string) (*models.VendorVerifyInfo, error) {
	if len(appVendorUserID) == 0 {
		return nil, nil
	}
	var sender = ctxutil.GetSender(ctx)
	appVendor, err := GetVendorVerifyInfo(ctx, appVendorUserID)
	if err != nil {
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorInternalError)
	}
	if sender != nil {
		if !appVendor.OwnerPath.CheckPermission(sender) {
			return nil, gerr.New(ctx, gerr.PermissionDenied, gerr.ErrorResourceAccessDenied, appVendor.UserId)
		}
	}
	return appVendor, nil
}
