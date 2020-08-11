package scheduler

import (
	"github.com/codefresh-io/argocd-listener/agent/pkg/heartbeat"
	"github.com/jasonlvhit/gocron"
)

var HeartBeatInterval uint64 = 5

func StartHeartBeat() {
	job := gocron.Every(HeartBeatInterval).Second().Do(heartbeat.HeartBeatTask)

	if job != nil {
		err := job.Error()

		if err != "" {
			panic("Cant heartbeat job because " + err)
		}
	}

	go gocron.Start()
}
