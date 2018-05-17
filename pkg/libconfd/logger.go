// Copyright 2018 <chaishushan{AT}gmail.com>. All rights reserved.
// Use of this source code is governed by a Apache-style
// license that can be found in the LICENSE file.

// copy from https://github.com/chai2010/logger

package libconfd

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"sync/atomic"
)

var pkgLogger atomic.Value

func init() {
	SetLogger(NewStdLogger(os.Stderr, "", "", 0))
}

// NewStdLogger create new logger based on std log.
// If level is empty string, use WARN as the default level.
// If flag is zore, use 'log.LstdFlags|log.Lshortfile' as the default flag.
// Level: DEBUG < INFO < WARN < ERROR < PANIC < FATAL
func NewStdLogger(out io.Writer, prefix, level string, flag int) Logger {
	return newStdLogger(out, prefix, level, flag)
}

func GetLogger() Logger {
	return pkgLogger.Load().(Logger)
}

func SetLogger(new Logger) (old Logger) {
	if x := pkgLogger.Load(); x != nil {
		old = x.(Logger)
	}
	pkgLogger.Store(new)
	return
}

// Logger interface
//
// See https://github.com/chai2010/logger
type Logger interface {
	Assert(condition bool, v ...interface{})
	Assertln(condition bool, v ...interface{})
	Assertf(condition bool, format string, v ...interface{})
	Debug(v ...interface{})
	Debugln(v ...interface{})
	Debugf(format string, v ...interface{})
	Info(v ...interface{})
	Infoln(v ...interface{})
	Infof(format string, v ...interface{})
	Warning(v ...interface{})
	Warningln(v ...interface{})
	Warningf(format string, v ...interface{})
	Error(v ...interface{})
	Errorln(v ...interface{})
	Errorf(format string, v ...interface{})
	Panic(v ...interface{})
	Panicln(v ...interface{})
	Panicf(format string, v ...interface{})
	Fatal(v ...interface{})
	Fatalln(v ...interface{})
	Fatalf(format string, v ...interface{})

	// Level: DEBUG < INFO < WARN < ERROR < PANIC < FATAL
	GetLevel() string
	SetLevel(new string) (old string)
}

type logLevelType uint32

const (
	logUnknownLevel logLevelType = iota // invalid
	logDebugLevel
	logInfoLevel
	logWarnLevel
	logErrorLevel
	logPanicLevel
	logFatalLevel
)

func (level logLevelType) Valid() bool {
	return level >= logDebugLevel && level <= logFatalLevel
}

func newLogLevel(name string) logLevelType {
	switch strings.ToUpper(name) {
	case "DEBUG", "":
		return logDebugLevel
	case "INFO":
		return logInfoLevel
	case "WARN":
		return logWarnLevel
	case "ERROR":
		return logErrorLevel
	case "PANIC":
		return logPanicLevel
	case "FATAL":
		return logFatalLevel
	}
	return logUnknownLevel
}

func (level logLevelType) String() string {
	switch level {
	case logDebugLevel:
		return "DEBUG"
	case logInfoLevel:
		return "INFO"
	case logWarnLevel:
		return "WARN"
	case logErrorLevel:
		return "ERROR"
	case logPanicLevel:
		return "PANIC"
	case logFatalLevel:
		return "FATAL"
	}
	return "INVALID"
}

type stdLogger struct {
	level logLevelType
	*log.Logger
}

func newStdLogger(out io.Writer, prefix, level string, flag int) *stdLogger {
	if flag == 0 {
		flag = log.LstdFlags | log.Lshortfile
	}
	if level == "" {
		level = "INFO"
	}

	p := &stdLogger{Logger: log.New(out, prefix, flag)}
	p.SetLevel(level)
	return p
}

func (p *stdLogger) getLevel() logLevelType {
	return logLevelType(atomic.LoadUint32((*uint32)(&p.level)))
}
func (p *stdLogger) setLevel(level logLevelType) logLevelType {
	return logLevelType(atomic.SwapUint32((*uint32)(&p.level), uint32(level)))
}

func (p *stdLogger) getLevelName() string {
	return p.getLevel().String()
}
func (p *stdLogger) setLevelByName(levelName string) string {
	level := newLogLevel(levelName)
	if !level.Valid() {
		panic("invalid level: " + levelName)
	}
	return p.setLevel(level).String()
}

func (p *stdLogger) GetLevel() string {
	return p.getLevel().String()
}
func (p *stdLogger) SetLevel(new string) (old string) {
	return p.setLevelByName(new)
}

func (p *stdLogger) Assert(condition bool, v ...interface{}) {
	if l := logDebugLevel; p.getLevel() <= l && !condition {
		p.Output(2, "[ASSERT] "+fmt.Sprint(v...))
		os.Exit(1)
	}
}
func (p *stdLogger) Assertln(condition bool, v ...interface{}) {
	if l := logDebugLevel; p.getLevel() <= l && !condition {
		p.Output(2, "[ASSERT] "+fmt.Sprintln(v...))
		os.Exit(1)
	}
}
func (p *stdLogger) Assertf(condition bool, format string, v ...interface{}) {
	if l := logDebugLevel; p.getLevel() <= l && !condition {
		p.Output(2, "[ASSERT] "+fmt.Sprintf(format, v...))
		os.Exit(1)
	}
}

func (p *stdLogger) Debug(v ...interface{}) {
	if l := logDebugLevel; p.getLevel() <= l {
		p.Output(2, "["+l.String()+"] "+fmt.Sprint(v...))
	}
}
func (p *stdLogger) Debugln(v ...interface{}) {
	if l := logDebugLevel; p.getLevel() <= l {
		p.Output(2, "["+l.String()+"] "+fmt.Sprintln(v...))
	}
}
func (p *stdLogger) Debugf(format string, v ...interface{}) {
	if l := logDebugLevel; p.getLevel() <= l {
		p.Output(2, "["+l.String()+"] "+fmt.Sprintf(format, v...))
	}
}

func (p *stdLogger) Info(v ...interface{}) {
	if l := logInfoLevel; p.getLevel() <= l {
		p.Output(2, "["+l.String()+"] "+fmt.Sprint(v...))
	}
}
func (p *stdLogger) Infoln(v ...interface{}) {
	if l := logInfoLevel; p.getLevel() <= l {
		p.Output(2, "["+l.String()+"] "+fmt.Sprintln(v...))
	}
}
func (p *stdLogger) Infof(format string, v ...interface{}) {
	if l := logInfoLevel; p.getLevel() <= l {
		p.Output(2, "["+l.String()+"] "+fmt.Sprintf(format, v...))
	}
}

func (p *stdLogger) Warning(v ...interface{}) {
	if l := logWarnLevel; p.getLevel() <= l {
		p.Output(2, "["+l.String()+"] "+fmt.Sprint(v...))
	}
}
func (p *stdLogger) Warningln(v ...interface{}) {
	if l := logWarnLevel; p.getLevel() <= l {
		p.Output(2, "["+l.String()+"] "+fmt.Sprintln(v...))
	}
}
func (p *stdLogger) Warningf(format string, v ...interface{}) {
	if l := logWarnLevel; p.getLevel() <= l {
		p.Output(2, "["+l.String()+"] "+fmt.Sprintf(format, v...))
	}
}

func (p *stdLogger) Error(v ...interface{}) {
	if l := logErrorLevel; p.getLevel() <= l {
		p.Output(2, "["+l.String()+"] "+fmt.Sprint(v...))
	}
}
func (p *stdLogger) Errorln(v ...interface{}) {
	if l := logErrorLevel; p.getLevel() <= l {
		p.Output(2, "["+l.String()+"] "+fmt.Sprintln(v...))
	}
}
func (p *stdLogger) Errorf(format string, v ...interface{}) {
	if l := logErrorLevel; p.getLevel() <= l {
		p.Output(2, "["+l.String()+"] "+fmt.Sprintf(format, v...))
	}
}

func (p *stdLogger) Panic(v ...interface{}) {
	s := fmt.Sprint(v...)
	if l := logPanicLevel; p.getLevel() <= l {
		p.Output(2, "["+l.String()+"] "+s)
	}
	panic(s)
}
func (p *stdLogger) Panicln(v ...interface{}) {
	s := fmt.Sprintln(v...)
	if l := logPanicLevel; p.getLevel() <= l {
		p.Output(2, "["+l.String()+"] "+s)
	}
	panic(s)
}
func (p *stdLogger) Panicf(format string, v ...interface{}) {
	s := fmt.Sprintf(format, v...)
	if l := logPanicLevel; p.getLevel() <= l {
		p.Output(2, "["+l.String()+"] "+s)
	}
	panic(s)
}

func (p *stdLogger) Fatal(v ...interface{}) {
	const l = logFatalLevel
	p.Output(2, "["+l.String()+"] "+fmt.Sprint(v...))
	os.Exit(1)
}
func (p *stdLogger) Fatalln(v ...interface{}) {
	const l = logFatalLevel
	p.Output(2, "["+l.String()+"] "+fmt.Sprintln(v...))
	os.Exit(1)
}
func (p *stdLogger) Fatalf(format string, v ...interface{}) {
	const l = logFatalLevel
	p.Output(2, "["+l.String()+"] "+fmt.Sprintf(format, v...))
	os.Exit(1)
}
