// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package test

import (
	"net/url"

	"github.com/gorilla/websocket"

	"openpitrix.io/openpitrix/pkg/topic"
	"openpitrix.io/openpitrix/pkg/util/jsonutil"
)

func GetIoClient(conf *ClientConfig, userId string) *IoClient {
	u := url.URL{
		Scheme:   "ws",
		Host:     conf.Host,
		Path:     conf.BasePath + "v1/io",
		RawQuery: "uid=" + userId,
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
