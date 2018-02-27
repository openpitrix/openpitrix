package logger

import "testing"

func TestPrint(t *testing.T) {
	Debug(222)
	Info(111)
	Infof("format [%d]", 111)
	SetLevelByString("debug")
	Info(111)
	Debug(222)
	Debug(222)
}
