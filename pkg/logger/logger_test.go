// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package logger

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func readBuf(buf *bytes.Buffer) string {
	str := buf.String()
	buf.Reset()
	return str
}

func TestLogger(t *testing.T) {
	buf := new(bytes.Buffer)
	SetOutput(buf)

	Debug("debug log, should ignore by default")
	assert.Empty(t, readBuf(buf))

	Info("info log, should visable")
	assert.Contains(t, readBuf(buf), "info log, should visable")

	Info("format [%d]", 111)
	assert.Contains(t, readBuf(buf), "format [111]")

	SetLevelByString("debug")
	Debug("debug log, now it becomes visible")
	assert.Contains(t, readBuf(buf), "debug log, now it becomes visible")

	logger = NewLogger()
	logger.SetPrefix("(prefix)").SetSuffix("(suffix)").SetOutput(buf)

	logger.Warn("log_content")
	log := readBuf(buf)
	assert.Regexp(t, " -WARNING- \\(prefix\\)log_content \\(testing.go:\\d+\\)\\(suffix\\)", log)
	t.Log(log)

	logger.HideCallstack()
	logger.Warn("log_content")
	log = readBuf(buf)
	assert.Regexp(t, " -WARNING- \\(prefix\\)log_content\\(suffix\\)", log)
	t.Log(log)
}
