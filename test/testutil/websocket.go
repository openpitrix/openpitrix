// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package testutil

import (
	"context"
	"log"
	"net/url"

	"github.com/gorilla/websocket"

	"openpitrix.io/openpitrix/pkg/client/config"
	"openpitrix.io/openpitrix/pkg/topic"
	"openpitrix.io/openpitrix/pkg/util/jsonutil"
)

func GetIoClient(conf *ClientConfig) *IoClient {
	tokenSource, err := config.GetTokenSource(context.Background(), conf.ConfigPath)
	if err != nil {
		log.Fatal(err)
	}
	token, err := tokenSource.Token()
	if err != nil {
		log.Fatal(err)
	}

	endpoint := conf.GetEndpoint()
	var scheme = "ws"
	if endpoint.Scheme == "https" {
		scheme = "wss"
	}
	u := url.URL{
		Scheme:   scheme,
		Host:     endpoint.Host,
		Path:     "/v1/io",
		RawQuery: "sid=" + token.AccessToken,
	}
	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		panic(err)
	}
	return &IoClient{c}
}

type IoClient struct {
	conn *websocket.Conn
}

func (i *IoClient) ReadMessage() topic.Message {
	var msg topic.Message
	_, content, err := i.conn.ReadMessage()
	if err != nil {
		panic(err)
	}
	err = jsonutil.Decode(content, &msg)
	if err != nil {
		panic(err)
	}
	return msg
}
