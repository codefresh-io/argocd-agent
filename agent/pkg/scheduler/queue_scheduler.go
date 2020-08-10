package scheduler

import (
	"github.com/codefresh-io/argocd-listener/agent/pkg/util"
	"github.com/jasonlvhit/gocron"
)

func task() {
	queue := util.Get()
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
