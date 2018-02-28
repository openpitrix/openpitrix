// Copyright 2017 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package sender

import (
	"context"
	"encoding/json"
	"net/http"

	context2 "golang.org/x/net/context"
	"google.golang.org/grpc/metadata"

	"openpitrix.io/openpitrix/pkg/logger"
)

const senderKey = "sender"

type Info struct {
	UserId string `json="user_id"`
}

func GetSenderFromContext(ctx context.Context) *Info {
	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		logger.Debugf("%+v", md[senderKey])
		sender := Info{}
		err := json.Unmarshal([]byte(md[senderKey][0]), &sender)
		if err != nil {
			panic(err)
		}
		return &sender
	}
	return nil
}

func AuthUserInfo(authKey string) *Info {
	logger.Debugf("got auth key: %+v", authKey)
	// TODO: validate auth key && get user info from db
	return &Info{UserId: "system"}
}

func ServeMuxSetSender(_ context2.Context, request *http.Request) metadata.MD {
	md := metadata.MD{}
	authKey := request.Header.Get("X-Auth-Key")
	user := AuthUserInfo(authKey)
	userJson, _ := json.Marshal(&user)
	md["sender"] = []string{string(userJson)}
	return md
}
