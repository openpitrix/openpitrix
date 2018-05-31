// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package gerr

import (
	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"openpitrix.io/openpitrix/pkg/logger"
	"openpitrix.io/openpitrix/pkg/pb"
)

const EN = "en"
const DefaultLocale = EN

func newStatus(code codes.Code, err error, errMsg ErrorMessage, a ...interface{}) *status.Status {
	locale := DefaultLocale

	s := status.New(code, errMsg.Message(locale, a...))

	errorDetail := &pb.ErrorDetail{ErrorName: errMsg.Name}
	if err != nil {
		errorDetail.Cause = err.Error()
	}

	sd, e := s.WithDetails(errorDetail)
	if e == nil {
		return sd
	} else {
		logger.Error("%+v", errors.WithStack(e))
	}
	return s
}

func New(code codes.Code, errMsg ErrorMessage, a ...interface{}) error {
	return newStatus(code, nil, errMsg, a...).Err()
}

func NewWithDetail(code codes.Code, err error, errMsg ErrorMessage, a ...interface{}) error {
	return newStatus(code, err, errMsg, a...).Err()
}
