package scheduler

import (
	"github.com/codefresh-io/argocd-listener/agent/pkg/heartbeat"
	"github.com/robfig/cron/v3"
)

func StartHeartBeat() {
	c := cron.New()
	_, _ = c.AddFunc("@every 8s", heartbeat.HeartBeatTask)
	c.Start()
}
