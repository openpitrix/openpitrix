// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package logger

import (
	"fmt"
	"runtime"
	"strings"
	"sync/atomic"
	"time"
)

type Level uint32

const (
	CriticalLevel Level = iota
	ErrorLevel
	WarnLevel
	InfoLevel
	DebugLevel
)

func (level Level) String() string {
	switch level {
	case DebugLevel:
		return "debug"
	case InfoLevel:
		return "info"
	case WarnLevel:
		return "warning"
	case ErrorLevel:
		return "error"
	case CriticalLevel:
		return "critical"
	}

	return "unknown"
}

func StringToLevel(level string) Level {
	switch level {
	case "critical":
		return CriticalLevel
	case "error":
		return ErrorLevel
	case "warn", "warning":
		return WarnLevel
	case "debug":
		return DebugLevel
	case "info":
		return InfoLevel
	}
	return InfoLevel
}

var logger = NewLogger()

func Info(format string, v ...interface{}) {
	logger.Info(format, v...)
}

func Debug(format string, v ...interface{}) {
	logger.Debug(format, v...)
}

func Warn(format string, v ...interface{}) {
	logger.Warn(format, v...)
}

func Error(format string, v ...interface{}) {
	logger.Error(format, v...)
}

func Critical(format string, v ...interface{}) {
	logger.Critical(format, v...)
}

func Fatalf(format string, v ...interface{}) {
	logger.Critical(format, v...)
}

func SetLevelByString(level string) {
	logger.SetLevelByString(level)
}

func NewLogger() *Logger {
	return &Logger{Level: InfoLevel}
}

type Logger struct {
	Level Level
}

func (logger *Logger) level() Level {
	return Level(atomic.LoadUint32((*uint32)(&logger.Level)))
}

func (logger *Logger) SetLevel(level Level) {
	atomic.StoreUint32((*uint32)(&logger.Level), uint32(level))
}

func (logger *Logger) SetLevelByString(level string) {
	logger.SetLevel(StringToLevel(level))
}

func (logger *Logger) formatOutput(level Level, output string) string {
	now := time.Now().Format("2006-01-02 15:04:05.99999")
	_, file, line, ok := runtime.Caller(4)
	if !ok {
		file = "???"
		line = 0
	}
	// short file name
	for i := len(file) - 1; i > 0; i-- {
		if file[i] == '/' {
			file = file[i+1:]
			break
		}
	}
	// 2018-03-27 02:08:44.93894 -INFO- Api service start http://openpitrix-api-gateway:9100 (main.go:44)
	return fmt.Sprintf("%s -%s- %s (%s:%d)", now, strings.ToUpper(level.String()), output, file, line)
}

func (logger *Logger) logf(level Level, format string, args ...interface{}) {
	if logger.level() < level {
		return
	}
	fmt.Println(logger.formatOutput(level, fmt.Sprintf(format, args...)))
}

func (logger *Logger) Debug(format string, args ...interface{}) {
	logger.logf(DebugLevel, format, args...)
}

func (logger *Logger) Info(format string, args ...interface{}) {
	logger.logf(InfoLevel, format, args...)
}

func (logger *Logger) Warn(format string, args ...interface{}) {
	logger.logf(WarnLevel, format, args...)
}

func (logger *Logger) Error(format string, args ...interface{}) {
	logger.logf(ErrorLevel, format, args...)
}

func (logger *Logger) Critical(format string, args ...interface{}) {
	logger.logf(CriticalLevel, format, args...)
}
