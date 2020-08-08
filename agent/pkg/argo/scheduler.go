package argo

import (
	"github.com/jasonlvhit/gocron"
)

func task() {
	queue := Get()
	queue.Notify()
}

func Schedule() {

	job := gocron.Every(1).Minute().Do(task)

	if job != nil {
		err := job.Error()

		if err != "" {
			panic("Cant start job because " + err)
		}
	}

	<-gocron.Start()

}
