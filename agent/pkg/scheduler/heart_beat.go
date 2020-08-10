package scheduler

import (
	"fmt"
	"github.com/codefresh-io/argocd-listener/agent/pkg/codefresh"
	"github.com/jasonlvhit/gocron"
)

var HeartBeatInterval uint64 = 5

func heartBeatTask() {
	err := codefresh.GetInstance().HeartBeat()
	if err != nil {
		fmt.Println(err)
	}
}

func StartHeartBeat() {
	job := gocron.Every(HeartBeatInterval).Second().Do(heartBeatTask)

	if job != nil {
		err := job.Error()

		if err != "" {
			panic("Cant heartbeat job because " + err)
		}
	}

	<-gocron.Start()
}
