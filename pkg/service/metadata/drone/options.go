// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package drone

import (
	"openpitrix.io/openpitrix/pkg/constants"
)

type Options struct {
	DbPath string
	Id     string
	Port   int
}

func NewDefaultOptions() *Options {
	return &Options{
		Id:   MakeDroneId(""),
		Port: constants.DroneServicePort,
	}
}

func WithDbPath(path string) func(opt *Options) {
	return func(opt *Options) {
		opt.DbPath = path
	}
}
func WithDrondId(id string) func(opt *Options) {
	return func(opt *Options) {
		opt.Id = id
	}
}

func WithListenPort(port int) func(opt *Options) {
	return func(opt *Options) {
		opt.Port = port
	}
}

func (p *Options) Clone() *Options {
	var q = *p
	return &q
}
