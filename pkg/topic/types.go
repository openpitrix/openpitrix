// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package topic

import "github.com/gorilla/websocket"

type model interface {
	GetTopicResource() Resource
}

type userMessage struct {
	UserId  string
	Message Message
}

type Message struct {
	// Type: optional create/delete/update
	Type     messageType `json:"type,omitempty"`
	Resource Resource    `json:"resource,omitempty"`
}

type Resource struct {
	ResourceType     string  `json:"rtype,omitempty"`
	ResourceId       string  `json:"rid,omitempty"`
	Status           *string `json:"status,omitempty"`
	TransitionStatus *string `json:"tstatus,omitempty"`
}

func NewResource(rtype, rid string) Resource {
	return Resource{
		ResourceType:     rtype,
		ResourceId:       rid,
		Status:           nil,
		TransitionStatus: nil,
	}
}

func (r Resource) SetStatus(status string) Resource {
	r.Status = &status
	return r
}

func (r Resource) SetTransitionStatus(transitionStatus string) Resource {
	r.TransitionStatus = &transitionStatus
	return r
}

func (r Resource) GetTopicResource() Resource { return r }

type receiver struct {
	UserId string
	Conn   *websocket.Conn
}

type messageType string

const (
	Create messageType = "create"
	Update messageType = "update"
	Delete messageType = "delete"
)

const topicPrefix = "topic"
