package scheduler

import (
	"github.com/codefresh-io/argocd-listener/agent/pkg/service"
	"github.com/robfig/cron/v3"
)

func StartHeartBeat() {
	hb := service.New()

	c := cron.New()
	_, _ = c.AddFunc("@every 8s", hb.HeartBeatTask)
	c.Start()
}
