// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package isv

import (
	"context"
	"regexp"

	"github.com/asaskevich/govalidator"

	"openpitrix.io/openpitrix/pkg/gerr"
	"openpitrix.io/openpitrix/pkg/logger"
)

func VerifyUrl(ctx context.Context, url string) error {
	if !govalidator.IsURL(url) {
		logger.Error(ctx, "Failed to verify url [%s]", url)
		return gerr.New(ctx, gerr.InvalidArgument, gerr.ErrorIllegalUrlFormat, url)
	}
	return nil
}

func VerifyEmail(ctx context.Context, email string) error {
	pattern := `\w+([-+.]\w+)*@\w+([-.]\w+)*\.\w+([-.]\w+)*`
	reg := regexp.MustCompile(pattern)
	isMatch := reg.MatchString(email)
	if isMatch {
		return nil
	} else {
		logger.Error(ctx, "Failed to verify email [%s]", email)
		return gerr.New(ctx, gerr.InvalidArgument, gerr.ErrorValidateFailed, email)
	}

}

// no prefix, the length is 11.
func VerifyPhoneNumber(ctx context.Context, phoneNumber string) error {
	pattern := `^1[0-9]{10}$`
	reg := regexp.MustCompile(pattern)
	isMatch := reg.MatchString(phoneNumber)
	if isMatch {
		return nil
	} else {
		logger.Error(ctx, "Failed to verify phone number [%s]", phoneNumber)
		return gerr.New(ctx, gerr.InvalidArgument, gerr.ErrorIllegalPhoneNumFormat, phoneNumber)
	}
}

func VerifyBankAccountNumber(ctx context.Context, bankAccountNumber string) error {
	pattern := `\d{12}|\d{15}|\d{16}|\d{17}|\d{18}|\d{19}`
	reg := regexp.MustCompile(pattern)
	isMatch := reg.MatchString(bankAccountNumber)
	if isMatch {
		return nil
	} else {
		logger.Error(ctx, "Failed to verify bank account number [%s]", bankAccountNumber)
		return gerr.New(ctx, gerr.InvalidArgument, gerr.ErrorIllegalBankAccountNumberFormat, bankAccountNumber)
	}
}
