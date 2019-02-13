// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package apigateway

import (
	"net/http"
	"strings"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"

	pbam "openpitrix.io/iam/pkg/pb/am"
	"openpitrix.io/openpitrix/pkg/client/am"
	"openpitrix.io/openpitrix/pkg/constants"
	"openpitrix.io/openpitrix/pkg/gerr"
	"openpitrix.io/openpitrix/pkg/logger"
	"openpitrix.io/openpitrix/pkg/sender"
	"openpitrix.io/openpitrix/pkg/util/ctxutil"
	"openpitrix.io/openpitrix/pkg/util/jwtutil"
	"openpitrix.io/openpitrix/pkg/util/stringutil"
)

const (
	Authorization = "Authorization"
)

var (
	amClient, _ = am.NewClient()
)

func httpAuth(mux *runtime.ServeMux, key string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {

		if req.URL.Path == "/v1/oauth2/token" {
			// skip auth sender
			mux.ServeHTTP(w, req)
			return
		}

		var err error
		ctx := req.Context()
		_, outboundMarshaler := runtime.MarshalerForRequest(mux, req)

		auth := strings.SplitN(req.Header.Get(Authorization), " ", 2)
		if auth[0] != "Bearer" {
			err = gerr.New(ctx, gerr.Unauthenticated, gerr.ErrorAuthFailure)
			runtime.HTTPError(ctx, mux, outboundMarshaler, w, req, err)
			return
		}
		s, err := jwtutil.Validate(key, auth[1])
		if err != nil {
			if err == jwtutil.ErrExpired {
				err = gerr.New(ctx, gerr.Unauthenticated, gerr.ErrorAccessTokenExpired)
			} else {
				err = gerr.New(ctx, gerr.Unauthenticated, gerr.ErrorAuthFailure)
			}
			runtime.HTTPError(ctx, mux, outboundMarshaler, w, req, err)
			return
		}

		if stringutil.StringIn(s.UserId, constants.InternalUsers) {
			// TODO: internal user should move into iam
			s.AccessPath = s.GetAccessPath()
			s.OwnerPath = s.GetOwnerPath()
		} else {
			v, err := amClient.CanDo(ctx, &pbam.CanDoRequest{
				UserId:    s.UserId,
				Url:       req.URL.Path,
				UrlMethod: req.Method,
			})
			if err != nil {
				logger.Error(ctx, "Sender [%+v] cannot [%s] [%s], err: %+v", s, req.Method, req.URL.Path, err)
				err = gerr.New(ctx, gerr.PermissionDenied, gerr.ErrorPermissionDenied)
				runtime.HTTPError(ctx, mux, outboundMarshaler, w, req, err)
				return
			}
			s.AccessPath = sender.OwnerPath(v.AccessPath)
			s.OwnerPath = sender.OwnerPath(v.OwnerPath)
			s.UserId = v.UserId
		}

		req.Header.Set(ctxutil.SenderKey, s.ToJson())
		req.Header.Del(Authorization)

		mux.ServeHTTP(w, req)
	})
}
