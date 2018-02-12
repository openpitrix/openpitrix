package logger

import (
	"fmt"
	"log"
	"os"
	"strings"
	"sync/atomic"
)

type Level uint32

const (
	PanicLevel Level = iota
	FatalLevel
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
	case FatalLevel:
		return "fatal"
	case PanicLevel:
		return "panic"
	}

	return "unknown"
}

func StringToLevel(level string) Level {
	switch level {
	case "panic":
		return PanicLevel
	case "fatal":
		return FatalLevel
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

func Info(v ...interface{}) {
	logger.Info(v...)
}

func Infof(format string, v ...interface{}) {
	logger.Infof(format, v...)
}

func Debug(v ...interface{}) {
	logger.Debug(v...)
}

func Debugf(format string, v ...interface{}) {
	logger.Debugf(format, v...)
}

func Warn(v ...interface{}) {
	logger.Warning(v...)
}

func Warnf(format string, v ...interface{}) {
	logger.Warningf(format, v...)
}

func Warning(v ...interface{}) {
	logger.Warning(v...)
}

func Warningf(format string, v ...interface{}) {
	logger.Warningf(format, v...)
}

func Error(v ...interface{}) {
	logger.Error(v...)
}

func Errorf(format string, v ...interface{}) {
	logger.Errorf(format, v...)
}

func Fatal(v ...interface{}) {
	logger.Fatal(v...)
}

func Fatalf(format string, v ...interface{}) {
	logger.Fatalf(format, v...)
}

func Panic(v ...interface{}) {
	logger.Panic(v...)
}

func Panicf(format string, v ...interface{}) {
	logger.Panicf(format, v...)
}

func SetLevelByString(level string) {
	logger.SetLevelByString(level)
}

func Disable() {
	SetLevelByString("fatal")
}

func NewLogger() *Logger {
	l := log.New(os.Stdout, "", log.LstdFlags|log.Lmicroseconds|log.Lshortfile)
	return &Logger{Level: InfoLevel, Logger: l}
}

type Logger struct {
	Level  Level
	Logger *log.Logger
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
	return "[" + strings.ToUpper(level.String()) + "] " + output
}

func (logger *Logger) log(level Level, args ...interface{}) {
	if logger.level() < level {
		return
	}
	logger.Logger.Output(4, logger.formatOutput(level, fmt.Sprint(args...)))
}

func (logger *Logger) logf(level Level, format string, args ...interface{}) {
	if logger.level() < level {
		return
	}
	logger.Logger.Output(4, logger.formatOutput(level, fmt.Sprintf(format, args...)))
}

func (logger *Logger) Debugf(format string, args ...interface{}) {
	logger.logf(DebugLevel, format, args...)
}

func (logger *Logger) Infof(format string, args ...interface{}) {
	logger.logf(InfoLevel, format, args...)
}

func (logger *Logger) Warnf(format string, args ...interface{}) {
	logger.logf(WarnLevel, format, args...)
}

func (logger *Logger) Warningf(format string, args ...interface{}) {
	logger.logf(WarnLevel, format, args...)
}

func (logger *Logger) Errorf(format string, args ...interface{}) {
	logger.logf(ErrorLevel, format, args...)
}

func (logger *Logger) Fatalf(format string, args ...interface{}) {
	logger.logf(FatalLevel, format, args...)
}

func (logger *Logger) Panicf(format string, args ...interface{}) {
	logger.logf(PanicLevel, format, args...)
}

func (logger *Logger) Debug(args ...interface{}) {
	logger.log(DebugLevel, args...)
}

func (logger *Logger) Info(args ...interface{}) {
	logger.log(InfoLevel, args...)
}

func (logger *Logger) Warn(args ...interface{}) {
	logger.log(WarnLevel, args...)
}

func (logger *Logger) Warning(args ...interface{}) {
	logger.log(WarnLevel, args...)
}

func (logger *Logger) Error(args ...interface{}) {
	logger.log(ErrorLevel, args...)
}

func (logger *Logger) Fatal(args ...interface{}) {
	logger.log(FatalLevel, args...)
}

func (logger *Logger) Panic(args ...interface{}) {
	logger.log(PanicLevel, args...)
}
