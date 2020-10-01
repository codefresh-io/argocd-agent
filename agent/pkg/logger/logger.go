package logger

import "github.com/golang/glog"

type Logger struct {
}

func GetLogger() *Logger {
	return &Logger{}
}

func (log *Logger) Info(format string, args ...interface{}) {
	glog.Infof(format, args)
}

func (log *Logger) Error(format string, args ...interface{}) {
	glog.Errorf(format, args)
}

func (log *Logger) ErrorE(err error) {
	glog.Error(err.Error())
}
