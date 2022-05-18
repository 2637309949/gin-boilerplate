package swagger

import (
	"context"
	"gin-boilerplate/comm/logger"
)

const (
	test = iota
	release
)

var swagMode = release

func isRelease() bool {
	return swagMode == release
}

// Println calls Output to print to the standard logger when release mode.
func Println(v ...interface{}) {
	if isRelease() {
		logger.Info(context.TODO(), v...)
	}
}

// Printf calls Output to print to the standard logger when release mode.
func Printf(format string, v ...interface{}) {
	if isRelease() {
		logger.Infof(context.TODO(), format, v...)
	}
}
