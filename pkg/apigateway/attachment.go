// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package apigateway

import (
	"context"
	"net/http"
	"strings"

	attachmentclient "openpitrix.io/openpitrix/pkg/client/attachment"
	"openpitrix.io/openpitrix/pkg/logger"
	"openpitrix.io/openpitrix/pkg/pb"
)

func ServeAttachments(prefix string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		path := req.URL.Path
		path = strings.TrimPrefix(path, prefix)
		params := strings.SplitN(path, "/", 2)
		if len(params) != 2 {
			logger.Error(nil, "%+v", params)
			w.WriteHeader(http.StatusForbidden)
			return
		}
		attachmentId := params[0]
		filename := params[1]

		attachmentServiceClient, err := attachmentclient.NewAttachmentServiceClient()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			logger.Error(nil, "Cannot connect attachment manager: %+v", err)
			return
		}
		r := pb.GetAttachmentRequest{
			AttachmentId: attachmentId,
			Filename:     filename,
		}
		res, err := attachmentServiceClient.GetAttachment(context.Background(), &r)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			logger.Error(nil, "Cannot get attachment: %+v", err)
			return
		}
		if len(res.Content) == 0 {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		w.Header().Set("Etag", res.Etag)

		if match := req.Header.Get("If-None-Match"); match != "" {
			if strings.Contains(match, res.Etag) {
				w.WriteHeader(http.StatusNotModified)
				return
			}
		}

		w.Write(res.Content)

		return
	})
}
