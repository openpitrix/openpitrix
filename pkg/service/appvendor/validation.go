// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package appvendor

import (
	"context"
	"regexp"

	"github.com/asaskevich/govalidator"

	"openpitrix.io/openpitrix/pkg/gerr"
	"openpitrix.io/openpitrix/pkg/models"
)

//Url
func VerifyUrl(ctx context.Context, urlStr string) (bool, error) {
	if !govalidator.IsURL(urlStr) {
		return false, gerr.New(ctx, gerr.InvalidArgument, gerr.ErrorValidateFailed, urlStr)
	}
	return true, nil
}

//Email
func VerifyEmailFmt(ctx context.Context, emailStr string) (bool, error) {
	pattern := `\w+([-+.]\w+)*@\w+([-.]\w+)*\.\w+([-.]\w+)*`
	reg := regexp.MustCompile(pattern)
	result := reg.MatchString(emailStr)
	if result {
		return true, nil
	} else {
		return false, gerr.New(ctx, gerr.InvalidArgument, gerr.ErrorValidateFailed, emailStr)
	}

}

//mobilephone, no prefix, the length is 11.
func VerifyPhoneFmt(ctx context.Context, phoneNumberStr string) (bool, error) {
	pattern := `^1[0-9]{10}$`
	reg := regexp.MustCompile(pattern)
	result := reg.MatchString(phoneNumberStr)
	if result {
		return true, nil
	} else {
		return false, gerr.New(ctx, gerr.InvalidArgument, gerr.ErrorValidateFailed, phoneNumberStr)
	}
}

//BankAccountNumber
func VerifyBankAccountNumberFmt(ctx context.Context, bankAccountNumberStr string) (bool, error) {
	pattern := `\d{12}|\d{15}|\d{16}|\d{17}|\d{18}|\d{19}`
	reg := regexp.MustCompile(pattern)
	result := reg.MatchString(bankAccountNumberStr)
	if result {
		return true, nil
	} else {
		return false, gerr.New(ctx, gerr.InvalidArgument, gerr.ErrorValidateFailed, bankAccountNumberStr)
	}
}

//Status
func VerifyStatus(ctx context.Context, strStatus ...string) (bool, error) {
	for _, s := range strStatus {
		if (s != models.StatusNew) && (s != models.StatusSubmitted) && (s != models.StatusPassed) && (s != models.StatusRejected) {
			return false, gerr.New(ctx, gerr.InvalidArgument, gerr.ErrorValidateFailed, strStatus)
		}
	}
	return true, nil
}
