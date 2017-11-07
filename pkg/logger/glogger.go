// Copyright 2017 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package logger

import (
	"fmt"

	"github.com/golang/glog"
	_ "google.golang.org/grpc/grpclog/glogger"
)

type glogger struct{}

func (g *glogger) Info(args ...interface{}) {
	glog.InfoDepth(2, args...)
}

func (g *glogger) Infoln(args ...interface{}) {
	glog.InfoDepth(2, fmt.Sprintln(args...))
}

func (g *glogger) Infof(format string, args ...interface{}) {
	glog.InfoDepth(2, fmt.Sprintf(format, args...))
}

func (g *glogger) Warning(args ...interface{}) {
	glog.WarningDepth(2, args...)
}

func (g *glogger) Warningln(args ...interface{}) {
	glog.WarningDepth(2, fmt.Sprintln(args...))
}

func (g *glogger) Warningf(format string, args ...interface{}) {
	glog.WarningDepth(2, fmt.Sprintf(format, args...))
}

func (g *glogger) Error(args ...interface{}) {
	glog.ErrorDepth(2, args...)
}

func (g *glogger) Errorln(args ...interface{}) {
	glog.ErrorDepth(2, fmt.Sprintln(args...))
}

func (g *glogger) Errorf(format string, args ...interface{}) {
	glog.ErrorDepth(2, fmt.Sprintf(format, args...))
}

func (g *glogger) Fatal(args ...interface{}) {
	glog.FatalDepth(2, args...)
}

func (g *glogger) Fatalln(args ...interface{}) {
	glog.FatalDepth(2, fmt.Sprintln(args...))
}

func (g *glogger) Fatalf(format string, args ...interface{}) {
	glog.FatalDepth(2, fmt.Sprintf(format, args...))
}

func (g *glogger) V(l int) bool {
	return bool(glog.V(glog.Level(l)))
}
