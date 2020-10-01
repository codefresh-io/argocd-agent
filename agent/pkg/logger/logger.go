package logger

import (
	"github.com/golang/glog"
)

type Logger struct {
}

func GetLogger() *Logger {
	return &Logger{}
}

func (log *Logger) Info(msg string) {
	glog.Info(msg)
}

func (log *Logger) Infof(format string, args ...interface{}) {
	glog.Infof(format, args)
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
