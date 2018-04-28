// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package libconfd

import (
	"text/template"
)

type Options func(*Config)

func (p *Config) applyOptions(opts ...Options) *Config {
	for _, fn := range opts {
		fn(p)
	}
	return p
}

func WithOnetimeMode() Options {
	return func(opt *Config) {
		opt.Onetime = true
	}
}

func WithIntervalMode() Options {
	return func(opt *Config) {
		opt.Onetime = false
		opt.Watch = false
	}
}

func WithInterval(interval int) Options {
	return func(opt *Config) {
		opt.Interval = interval
	}
}

func WithWatchMode() Options {
	return func(opt *Config) {
		opt.Onetime = false
		opt.Watch = true
	}
}

func WithFuncMap(maps ...template.FuncMap) Options {
	return func(opt *Config) {
		if opt.FuncMap == nil {
			opt.FuncMap = make(template.FuncMap)
		}
		for _, m := range maps {
			for k, fn := range m {
				opt.FuncMap[k] = fn
			}
		}
	}
}

func WithAbsKeyAdjuster(fn func(absKey string) (realKey string)) Options {
	return func(opt *Config) {
		opt.HookAbsKeyAdjuster = fn
	}
}

func WithFuncMapUpdater(fn func(m template.FuncMap, basefn *TemplateFunc)) Options {
	return func(opt *Config) {
		opt.FuncMapUpdater = fn
	}
}

func WithHookOnCheckCmdDone(fn func(trName, cmd string, err error)) Options {
	return func(opt *Config) {
		opt.HookOnCheckCmdDone = fn
	}
}

func WithHookOnReloadCmdDone(fn func(trName, cmd string, err error)) Options {
	return func(opt *Config) {
		opt.HookOnReloadCmdDone = fn
	}
}

func WithHookOnUpdateDone(fn func(trName string, err error)) Options {
	return func(opt *Config) {
		opt.HookOnUpdateDone = fn
	}
}
