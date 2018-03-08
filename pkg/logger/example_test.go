// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package logger_test

import "openpitrix.io/openpitrix/pkg/logger"

func ExampleLogger() {
	logger.Debug("debug log, should ignore by default")
	logger.Info("info log, should visable")
	logger.Infof("format [%d]", 111)
	logger.SetLevelByString("debug")
	logger.Debug("debug log, now it becomes visible")
	// Output:
}
