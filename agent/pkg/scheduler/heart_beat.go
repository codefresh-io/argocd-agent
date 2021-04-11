package scheduler

import (
	"github.com/codefresh-io/argocd-listener/agent/pkg/service"
	"github.com/robfig/cron/v3"
)

func StartHeartBeat() {
	c := cron.New()
	_, _ = c.AddFunc("@every 8s", service.HeartBeatTask)
	c.Start()
}
