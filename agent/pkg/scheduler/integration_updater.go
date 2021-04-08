package scheduler

import (
	"github.com/codefresh-io/argocd-listener/agent/pkg/codefresh"
	"github.com/codefresh-io/argocd-listener/agent/pkg/logger"
	"github.com/codefresh-io/argocd-listener/agent/pkg/store"
	"github.com/robfig/cron/v3"
)

func updateIntegrationTask() {
	storeData := store.GetStore()

	err := codefresh.GetInstance().UpdateIntegration(storeData.Codefresh.Integration, storeData.Argo.Host,
		"", "", storeData.Argo.Token, "", "", "")

	if err != nil {
		logger.GetLogger().Errorf("Failed to update integration, reason %v", err)
	}
}

func StartUpdateIntegration() {
	c := cron.New()
	_, _ = c.AddFunc("@every 10s", updateIntegrationTask) // time???
	c.Start()
}
