// Copyright 2017 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package config

import (
	"flag"
	"fmt"
	"os"

	"github.com/golang/glog"
)

func init() {
	// avoid glog warning
	// and skip -h or -help flag when call flag.Parse()
	flag.CommandLine.Usage = func() {
		// fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		// flag.PrintDefaults()
	}
	flag.CommandLine.Init(os.Args[0], flag.ContinueOnError)
	flag.CommandLine.Parse([]string{})
}

func (p *Config) ActiveGlogFlags() {
	p.Glog.ActiveFlags()
}

func (p *Glog) ActiveFlags() {
	flag.CommandLine.Set("logtostderr", fmt.Sprintf("%v", p.LogToStderr))
	flag.CommandLine.Set("alsologtostderr", fmt.Sprintf("%v", p.AlsoLogTostderr))
	flag.CommandLine.Set("stderrthreshold", p.StderrThreshold)
	flag.CommandLine.Set("log_dir", p.LogDir)

	flag.CommandLine.Set("log_backtrace_at", p.LogBacktraceAt)
	flag.CommandLine.Set("v", fmt.Sprintf("%v", p.V))
	flag.CommandLine.Set("vmodule", p.VModule)

	if p.LogDir != "" {
		os.MkdirAll(p.LogDir, 0666)
	}
	if p.CopyStandardLogTo != "" {
		glog.CopyStandardLogTo(p.CopyStandardLogTo)
	}
}
