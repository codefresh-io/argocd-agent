package scheduler

import (
	"github.com/codefresh-io/argocd-listener/agent/pkg/api/codefresh"
	"github.com/codefresh-io/argocd-listener/agent/pkg/infra/logger"
	"github.com/codefresh-io/argocd-listener/agent/pkg/infra/store"
)

type integrationUpdaterScheduler struct {
	codefreshApi codefresh.CodefreshApi
}

func GetIntegrationUpdatedScheduler() Scheduler {
	return &integrationUpdaterScheduler{codefreshApi: codefresh.GetInstance()}
}

func (integrationUpdaterScheduler *integrationUpdaterScheduler) updateIntegrationTask() {
	storeData := store.GetStore()

	err := integrationUpdaterScheduler.codefreshApi.UpdateIntegration(storeData.Codefresh.Integration, storeData.Argo.Host,
		"", "", storeData.Argo.Token, "", "", "")

	if err != nil {
		logger.GetLogger().Errorf("Failed to update integration, reason %v", err)
	}
}

func (integrationUpdaterScheduler *integrationUpdaterScheduler) getTime() string {
	return "@every 100s"
}

func (integrationUpdaterScheduler *integrationUpdaterScheduler) getFunc() func() {
	return integrationUpdaterScheduler.updateIntegrationTask
}

func (integrationUpdaterScheduler *integrationUpdaterScheduler) Run() {
	run(integrationUpdaterScheduler)
}
