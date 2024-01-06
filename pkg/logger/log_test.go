package logger

import (
	"testing"
)

func TestLogger(t *testing.T) {
	AccessLog := "/tmp/access.log"
	ErrorLog := "/tmp/error.log"

	Init(AccessLog, ErrorLog)

	Info("Hello World")
	Error("Bad World")
	Infof("%v, %v, %v", "First", "Second", "Third")
	Errorf("%v, %v, %v", "Third", "Second", "First")
}
