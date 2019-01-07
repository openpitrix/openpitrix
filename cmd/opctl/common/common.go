// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package common

import "fmt"

func Error(err error, msg string) {
	panic(fmt.Sprintf("%s error: %+v", msg, err))
}

type Param struct {
	Shorthand string      `yaml:"shorthand,omitempty"`
	Help      string      `yaml:"help,omitempty"`
	Type      string      `yaml:"type,omitempty"`
	Default   interface{} `yaml:"default,omitempty"`
}

type Cmd struct {
	Action      string           `yaml:"action,omitempty"`
	Request     string           `yaml:"request,omitempty"`
	Description string           `yaml:"description,omitempty"`
	Service     string           `yaml:"service,omitempty"`
	Path        map[string]Param `yaml:"path,omitempty"`
	Query       map[string]Param `yaml:"query,omitempty"`
	Body        map[string]Param `yaml:"body,omitempty"`
	Insecurity  bool             `yaml:"insecurity,omitempty"`
}

type Cmds []Cmd

func (cs Cmds) Len() int {
	return len(cs)
}
func (cs Cmds) Swap(i, j int) {
	cs[i], cs[j] = cs[j], cs[i]
}
func (cs Cmds) Less(i, j int) bool {
	if cs[i].Service == cs[j].Service {
		return cs[i].Action < cs[j].Action
	}
	return cs[i].Service < cs[j].Service
}
