// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package app

type Config struct {
	Raw         string
	Type        string    `json:"type,omitempty"`
	Properties  []*Config `json:"properties,omitempty"`
	Key         string    `json:"key,omitempty"`
	Description string    `json:"description,omitempty"`
	Required    bool      `json:"required,omitempty"`

	Default interface{} `json:"default,omitempty"`
	Pattern string      `json:"pattern,omitempty"`

	Limits            map[string][]string `json:"limits,omitempty"`
	AllowedOperations []string            `json:"allowed_operations,omitempty"`

	Port int `json:"port,omitempty"`

	Range []int `json:"range,omitempty"`
	Min   int   `json:"min,omitempty"`
	Max   int   `json:"max,omitempty"`
	Step  int   `json:"step,omitempty"`

	Changeable  bool   `json:"changeable,omitempty"`
	Separator   string `json:"separator,omitempty"`
	Multichoice bool   `json:"multichoice,omitempty"`
}
