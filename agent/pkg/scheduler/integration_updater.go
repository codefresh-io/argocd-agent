package scheduler

import (
	"github.com/codefresh-io/argocd-listener/agent/pkg/service"
	"github.com/robfig/cron/v3"
)

func StartUpdateIntegration() {
	gitopsService := service.NewGitopsService()
	c := cron.New()
	_, _ = c.AddFunc("@every 10s", gitopsService.UpdateIntegration) // time???
	c.Start()
}
