package logger

import (
	"flag"
	"github.com/golang/glog"
)

type Logger struct {
}

var logger *Logger

func GetLogger() *Logger {
	if logger == nil {
		// NOTE: This next line is key you have to call flag.Parse() for the command line
		// options or "flags" that are defined in the glog module to be picked up.
		flag.Parse()
		logger = &Logger{}
	}
	return logger
}

func (log *Logger) Info(msg string) {
	glog.Info(msg)
}

func (log *Logger) Infof(format string, args ...interface{}) {
	glog.Infof(format, args...)
}

func (log *Logger) Errorf(format string, args ...interface{}) {
	glog.Errorf(format, args...)
}

func (log *Logger) Error(msg string) {
	glog.Error(msg)
}

func (log *Logger) ErrorE(err error) {
	glog.Error(err.Error())
}
