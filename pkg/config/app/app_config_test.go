// Copyright 2017 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package app_config

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
	cfg, err := LoadConfig()
	tAssert(t, err == nil, err)

	var pb_const pb.Const
	tAssert(t, cfg.App.Host == pb_const.GetAppHost())
	tAssert(t, cfg.App.Port == int(pb_const.GetAppPort()))

	tCheckDefaultGlog(t)
	tCheckDefultDB(t)
}
func TestConfig_envChanged(t *testing.T) {
	tChangeGlogConfig()
	tCheckGlogChanged(t)

	tChangeServiceHostAndPort()
	tCheckServiceHostAndPort(t)

	tChangeDBConfig()
	tCheckDBChanged(t)
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

func tCheckDefultDB(t *testing.T) {
	cfg, err := LoadConfig()
	tAssert(t, err == nil, err)

	var pb_const pb.Const
	tAssert(t, cfg.DB.Type == pb_const.GetDbType())
	tAssert(t, cfg.DB.Host == pb_const.GetDbHost())
	tAssert(t, cfg.DB.Port == int(pb_const.GetDbPort()))
	tAssert(t, cfg.DB.Encoding == pb_const.GetDbEncoding())
	tAssert(t, cfg.DB.Engine == pb_const.GetDbEngine())
	tAssert(t, cfg.DB.DbName == pb_const.GetDbDbname())

	var pb_default pb.Default
	tAssert(t, cfg.DB.RootName == pb_default.GetDbAdminName())
	tAssert(t, cfg.DB.RootPassword == pb_default.GetDbAdminPassword())
	tAssert(t, cfg.DB.UserName == pb_default.GetDbAppUserName())
	tAssert(t, cfg.DB.UserPassword == pb_default.GetDbAppUserPassword())
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
	os.Setenv("OPENPITRIX_CONFIG_APP_HOST", "host")
	os.Setenv("OPENPITRIX_CONFIG_APP_PORT", "1234")
}
func tCheckServiceHostAndPort(t *testing.T) {
	cfg, err := LoadConfig()
	tAssert(t, err == nil, err)

	tAssert(t, cfg.App.Host == "host")
	tAssert(t, cfg.App.Port == 1234)
}

func tChangeDBConfig() {
	os.Setenv("OPENPITRIX_CONFIG_DB_TYPE", "mysql-new")
	os.Setenv("OPENPITRIX_CONFIG_DB_HOST", "mysql-host")
	os.Setenv("OPENPITRIX_CONFIG_DB_PORT", "3307")
	os.Setenv("OPENPITRIX_CONFIG_DB_ENCODING", "ascii")
	os.Setenv("OPENPITRIX_CONFIG_DB_ENGINE", "MyISAM")
	os.Setenv("OPENPITRIX_CONFIG_DB_DBNAME", "dbname")

	os.Setenv("OPENPITRIX_CONFIG_DB_ROOTNAME", "admin")
	os.Setenv("OPENPITRIX_CONFIG_DB_ROOTPASSWORD", "admin-123")
	os.Setenv("OPENPITRIX_CONFIG_DB_USERNAME", "user")
	os.Setenv("OPENPITRIX_CONFIG_DB_USERPASSWORD", "user-123")
}
func tCheckDBChanged(t *testing.T) {
	cfg, err := LoadConfig()
	tAssert(t, err == nil, err)

	tAssert(t, cfg.DB.Type == "mysql-new")
	tAssert(t, cfg.DB.Host == "mysql-host")
	tAssert(t, cfg.DB.Port == 3307)
	tAssert(t, cfg.DB.Encoding == "ascii")
	tAssert(t, cfg.DB.Engine == "MyISAM")
	tAssert(t, cfg.DB.DbName == "dbname")

	tAssert(t, cfg.DB.RootName == "admin")
	tAssert(t, cfg.DB.RootPassword == "admin-123")
	tAssert(t, cfg.DB.UserName == "user")
	tAssert(t, cfg.DB.UserPassword == "user-123")
}
