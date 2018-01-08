// Copyright 2017 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package api_config

import (
	"fmt"
	"os"
	"testing"

	pb "openpitrix.io/openpitrix/pkg/service.pb"
)

func tAssert(tb testing.TB, condition bool, a ...interface{}) {
	tb.Helper()
	if !condition {
		if msg := fmt.Sprint(a...); msg != "" {
			tb.Fatal("Assert failed: " + msg)
		} else {
			tb.Fatal("Assert failed")
		}
	}
}

func tAssertf(tb testing.TB, condition bool, format string, a ...interface{}) {
	tb.Helper()
	if !condition {
		if msg := fmt.Sprintf(format, a...); msg != "" {
			tb.Fatal("Assert failed: " + msg)
		} else {
			tb.Fatal("Assert failed")
		}
	}
}

func TestConfig_default(t *testing.T) {
	tCheckDefaultServiceHostAndPort(t)
	tCheckDefaultGlog(t)
}
func TestConfig_envChanged(t *testing.T) {
	tChangeGlogConfig()
	tCheckGlogChanged(t)

	tChangeServiceHostAndPort()
	tCheckServiceHostAndPort(t)
}

func tCheckDefaultServiceHostAndPort(t *testing.T) {
	cfg, err := LoadConfig()
	tAssert(t, err == nil, err)

	var pb_const pb.Const
	tAssert(t, cfg.Api.Host == pb_const.GetApiHost())
	tAssert(t, cfg.Api.Port == int(pb_const.GetApiPort()))
	tAssert(t, cfg.App.Host == pb_const.GetAppHost())
	tAssert(t, cfg.App.Port == int(pb_const.GetAppPort()))
	tAssert(t, cfg.Runtime.Host == pb_const.GetRuntimeHost())
	tAssert(t, cfg.Runtime.Port == int(pb_const.GetRuntimePort()))
	tAssert(t, cfg.Cluster.Host == pb_const.GetClusterHost())
	tAssert(t, cfg.Cluster.Port == int(pb_const.GetClusterPort()))
	tAssert(t, cfg.Repo.Host == pb_const.GetRepoHost())
	tAssert(t, cfg.Repo.Port == int(pb_const.GetRepoPort()))
}

func tCheckDefaultGlog(t *testing.T) {
	cfg, err := LoadConfig()
	tAssert(t, err == nil, err)

	var pb_default pb.Default
	tAssert(t, cfg.Glog.LogToStderr == pb_default.GetGlogLogToStderr())
	tAssert(t, cfg.Glog.AlsoLogToStderr == pb_default.GetGlogAlsoLogToStderr())
	tAssert(t, cfg.Glog.StderrThreshold == pb_default.GetGlogStderrThreshold())
	tAssert(t, cfg.Glog.LogDir == pb_default.GetGlogLogDir())
	tAssert(t, cfg.Glog.LogBacktraceAt == pb_default.GetGlogLogBacktraceAt())
	tAssert(t, cfg.Glog.V == int(pb_default.GetGlogV()))
	tAssert(t, cfg.Glog.VModule == pb_default.GetGlogVModule())
	tAssert(t, cfg.Glog.CopyStandardLogTo == pb_default.GetGlogCopyStandardLogTo())
}

func tChangeGlogConfig() {
	os.Setenv("OPENPITRIX_CONFIG_GLOG_LOGTOSTDERR", "true")
	os.Setenv("OPENPITRIX_CONFIG_GLOG_ALSOLOGTOSTDERR", "false")
	os.Setenv("OPENPITRIX_CONFIG_GLOG_STDERRTHRESHOLD", "ERROR-new")
	os.Setenv("OPENPITRIX_CONFIG_GLOG_LOGDIR", "dir-new")

	os.Setenv("OPENPITRIX_CONFIG_GLOG_LOGBACKTRACEAT", "trace-new")
	os.Setenv("OPENPITRIX_CONFIG_GLOG_V", "123")
	os.Setenv("OPENPITRIX_CONFIG_GLOG_VMODULE", "dir-new")

	os.Setenv("OPENPITRIX_CONFIG_GLOG_COPYSTANDARDLOGTO", "INFO-new")
}
func tCheckGlogChanged(t *testing.T) {
	cfg, err := LoadConfig()
	tAssert(t, err == nil, err)

	tAssert(t, cfg.Glog.LogToStderr == true)
	tAssert(t, cfg.Glog.AlsoLogToStderr == false)
	tAssert(t, cfg.Glog.StderrThreshold == "ERROR-new")
	tAssert(t, cfg.Glog.LogDir == "dir-new")
	tAssert(t, cfg.Glog.LogBacktraceAt == "trace-new")
	tAssert(t, cfg.Glog.V == 123)
	tAssert(t, cfg.Glog.VModule == "dir-new")
	tAssert(t, cfg.Glog.CopyStandardLogTo == "INFO-new")
}

func tChangeServiceHostAndPort() {
	os.Setenv("OPENPITRIX_CONFIG_API_HOST", "openpitrix-api-new")         // -new
	os.Setenv("OPENPITRIX_CONFIG_API_PORT", "19100")                      // +10000
	os.Setenv("OPENPITRIX_CONFIG_APP_HOST", "openpitrix-app-new")         // -new
	os.Setenv("OPENPITRIX_CONFIG_APP_PORT", "19101")                      // +10000
	os.Setenv("OPENPITRIX_CONFIG_RUNTIME_HOST", "openpitrix-runtime-new") // -new
	os.Setenv("OPENPITRIX_CONFIG_RUNTIME_PORT", "19102")                  // +10000
	os.Setenv("OPENPITRIX_CONFIG_CLUSTER_HOST", "openpitrix-cluster-new") // -new
	os.Setenv("OPENPITRIX_CONFIG_CLUSTER_PORT", "19103")                  // +10000
	os.Setenv("OPENPITRIX_CONFIG_REPO_HOST", "openpitrix-repo-new")       // -new
	os.Setenv("OPENPITRIX_CONFIG_REPO_PORT", "19104")                     // +10000
}
func tCheckServiceHostAndPort(t *testing.T) {
	cfg, err := LoadConfig()
	tAssert(t, err == nil, err)

	var pb_const pb.Const
	tAssert(t, cfg.Api.Host == pb_const.GetApiHost()+"-new")
	tAssert(t, cfg.Api.Port == int(pb_const.GetApiPort()+10000))
	tAssert(t, cfg.App.Host == pb_const.GetAppHost()+"-new")
	tAssert(t, cfg.App.Port == int(pb_const.GetAppPort())+10000)
	tAssert(t, cfg.Runtime.Host == pb_const.GetRuntimeHost()+"-new")
	tAssert(t, cfg.Runtime.Port == int(pb_const.GetRuntimePort())+10000)
	tAssert(t, cfg.Cluster.Host == pb_const.GetClusterHost()+"-new")
	tAssert(t, cfg.Cluster.Port == int(pb_const.GetClusterPort())+10000)
	tAssert(t, cfg.Repo.Host == pb_const.GetRepoHost()+"-new")
	tAssert(t, cfg.Repo.Port == int(pb_const.GetRepoPort())+10000)
}
