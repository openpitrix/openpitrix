// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package gerr

import (
	"fmt"

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
		errorDetail.Cause = fmt.Sprintf("%+v", err)
		logger.Error("%+v", err)
	}

	sd, e := s.WithDetails(errorDetail)
	if e == nil {
		return sd
	} else {
		logger.Error("%+v", errors.WithStack(e))
	}
	return s
}

func ClearErrorCause(err error) error {
	if e, ok := status.FromError(err); ok {
		details := e.Details()
		if len(details) > 0 {
			detail := details[0]
			if d, ok := detail.(*pb.ErrorDetail); ok {
				d.Cause = ""
				// clear detail
				proto := e.Proto()
				proto.Details = proto.Details[:0]
				e = status.FromProto(proto)
				e, _ := e.WithDetails(d)
				return e.Err()
			}
		}
	}
	return err
}

type GRPCError interface {
	error
	GRPCStatus() *status.Status
}

func New(code codes.Code, errMsg ErrorMessage, a ...interface{}) GRPCError {
	return newStatus(code, nil, errMsg, a...).Err().(GRPCError)
}

func NewWithDetail(code codes.Code, err error, errMsg ErrorMessage, a ...interface{}) GRPCError {
	return newStatus(code, err, errMsg, a...).Err().(GRPCError)
}

func IsGRPCError(err error) bool {
	if e, ok := err.(GRPCError); ok && e != nil {
		return true
	}
	return false
}
