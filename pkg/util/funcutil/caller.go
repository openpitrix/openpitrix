// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package funcutil

import (
	"regexp"
	"runtime"
)

var (
	reInit    = regexp.MustCompile(`init·\d+$`) // main.init·1
	reClosure = regexp.MustCompile(`func·\d+$`) // main.func·001
)

// caller types:
// runtime.goexit
// runtime.main
// main.init
// main.main
// main.init·1 -> main.init
// main.func·001 -> main.func
// code.google.com/p/gettext-go/gettext.TestCallerName
// ...
func CallerName(skip int) string {
	pc, _, _, ok := runtime.Caller(skip)
	if !ok {
		return ""
	}
	name := runtime.FuncForPC(pc).Name()
	if reInit.MatchString(name) {
		return reInit.ReplaceAllString(name, "init")
	}
	if reClosure.MatchString(name) {
		return reClosure.ReplaceAllString(name, "func")
	}
	return name
}
