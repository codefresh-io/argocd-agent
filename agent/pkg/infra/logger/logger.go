package logger

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strconv"

	"github.com/golang/glog"
	"github.com/sergi/go-diff/diffmatchpatch"
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
		glog.Info(fmt.Sprintf("[debug] %s", msg))
	}
}

func (log *Logger) Debugf(format string, args ...interface{}) {
	if log.loglevel == DebugLevel {
		glog.Infof(fmt.Sprintf("[debug] %s", format), args...)
	}
}

func (log *Logger) Diff(obj1 interface{}, obj2 interface{}) {
	if log.loglevel == DebugLevel {
		printResourceDiff, _ := os.LookupEnv("PRINT_RESOURCE_DIFF")
		if truthy, err := strconv.ParseBool(printResourceDiff); err == nil && truthy {
			log.printDiff(obj1, obj2)
		}
	}
}

func (log *Logger) ErrorE(err error) {
	glog.Error(err.Error())
}

func (log *Logger) printDiff(oldItem interface{}, newItem interface{}) error {
	if oldItem == nil || newItem == nil {
		log.Debug("[DIFF] ignore diff view because one of the entities is nil")
		return nil
	}

	prevState, err := json.Marshal(oldItem)
	if err != nil {
		return err
	}
	newState, err := json.Marshal(newItem)
	if err != nil {
		return err
	}

	dmp := diffmatchpatch.New()
	diffs := dmp.DiffMain(string(prevState), string(newState), false)
	log.Debugf("[DIFF] %s", dmp.DiffPrettyText(diffs))

	return nil
}
