package scheduler

import "github.com/robfig/cron/v3"

type Scheduler interface {
	getTime() string
	getFunc() func()
	Run()
}

func run(scheduler Scheduler) {
	c := cron.New()
	_, _ = c.AddFunc(scheduler.getTime(), scheduler.getFunc())
	c.Start()
}
