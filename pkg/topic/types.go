// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package topic

import "github.com/gorilla/websocket"

type UserMessage struct {
	UserId  string
	Message Message
}

type Message struct {
	// Type: optional create/delete/update
	Type        MessageType `json:"type,omitempty"`
	ResourceSet []Resource  `json:"resource_set,omitempty"`
}
type Resource struct {
	ResourceType     string  `json:"rtype,omitempty"`
	ResourceId       string  `json:"rid,omitempty"`
	Status           *string `json:"status,omitempty"`
	TransitionStatus *string `json:"tstatus,omitempty"`
}

type Receiver struct {
	UserId string
	Conn   *websocket.Conn
}

type MessageType string

const (
	Create MessageType = "create"
	Update MessageType = "update"
	Delete MessageType = "delete"
)

const topicPrefix = "topic"
