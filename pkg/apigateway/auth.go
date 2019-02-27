// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package apigateway

import (
	"net/http"
	"strings"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"

	"openpitrix.io/openpitrix/pkg/client/access"
	"openpitrix.io/openpitrix/pkg/gerr"
	"openpitrix.io/openpitrix/pkg/logger"
	"openpitrix.io/openpitrix/pkg/pb"
	"openpitrix.io/openpitrix/pkg/sender"
	"openpitrix.io/openpitrix/pkg/util/ctxutil"
	"openpitrix.io/openpitrix/pkg/util/jwtutil"
)

const (
	Authorization = "Authorization"
)

var (
	accessClient, _ = access.NewClient()
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

		v, err := accessClient.CanDo(ctx, &pb.CanDoRequest{
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

		logger.Debug(ctx, "CanDo: %+v", v)
		s.AccessPath = sender.OwnerPath(v.GetAccessPath())
		s.OwnerPath = sender.OwnerPath(v.GetOwnerPath())
		s.UserId = v.UserId

		req.Header.Set(ctxutil.SenderKey, s.ToJson())
		req.Header.Del(Authorization)

		mux.ServeHTTP(w, req)
	})
}
