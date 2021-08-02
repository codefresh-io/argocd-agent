package logger

import (
	"flag"
	"github.com/golang/glog"
	"os"
)

const (
	DebugLevel = "debug"
)

type Logger struct {
	loglevel string
}

var logger *Logger

func GetLogger() *Logger {
	if logger == nil {
		_ = flag.Set("logtostderr", "true")
		_ = flag.Set("stderrthreshold", "WARNING")
		_ = flag.Set("v", "2")
		// NOTE: This next line is key you have to call flag.Parse() for the command line
		// options or "flags" that are defined in the glog module to be picked up.
		flag.Parse()
		logger = &Logger{loglevel: os.Getenv("LOG_LEVEL")}
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

func (log *Logger) Debug(msg string) {
	if log.loglevel == DebugLevel {
		glog.Info(msg)
	}
}

func (log *Logger) ErrorE(err error) {
	glog.Error(err.Error())
}
