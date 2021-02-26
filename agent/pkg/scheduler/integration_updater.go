package scheduler

import (
	"github.com/codefresh-io/argocd-listener/agent/pkg/integration_updater"
	"github.com/robfig/cron/v3"
)

func StartUpdateIntegration() {
	c := cron.New()
	_, _ = c.AddFunc("@every 10s", integration_updater.UpdateIntegrationTask) // time???
	c.Start()
}
