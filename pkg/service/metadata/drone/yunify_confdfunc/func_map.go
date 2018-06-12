// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package yunify_confdfunc

import (
	"text/template"
)

func MakeCustomFuncMap() template.FuncMap {
	var m = template.FuncMap{}

	// yunify's confd add some new functions:
	// https://github.com/yunify/confd/commit/af06aaae59e96492d9bcc71342ca12f0af6f114d
	// https://github.com/yunify/confd/blob/master/resource/template/template_funcs.go

	m["filter"] = Filter
	m["toJson"] = ToJson
	m["toYaml"] = ToYaml

	m["min"] = min
	m["max"] = max

	// yunify's confd replaced some functions:
	// https://github.com/yunify/confd/commit/94860c6d62c4fc74d441b40fa8107742534cf6ee

	m["add"] = func(a, b interface{}) (interface{}, error) { return DoArithmetic(a, b, '+') }
	m["div"] = func(a, b interface{}) (interface{}, error) { return DoArithmetic(a, b, '/') }
	m["mul"] = func(a, b interface{}) (interface{}, error) { return DoArithmetic(a, b, '*') }
	m["sub"] = func(a, b interface{}) (interface{}, error) { return DoArithmetic(a, b, '-') }
	m["eq"] = eq
	m["ne"] = ne
	m["gt"] = gt
	m["ge"] = ge
	m["lt"] = lt
	m["le"] = le
	m["mod"] = mod

	return m
}
