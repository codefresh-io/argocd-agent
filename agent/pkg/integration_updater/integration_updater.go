package integration_updater

import (
	"github.com/codefresh-io/argocd-listener/agent/pkg/codefresh"
	"github.com/codefresh-io/argocd-listener/agent/pkg/logger"
	"github.com/codefresh-io/argocd-listener/agent/pkg/store"
	"github.com/codefresh-io/argocd-listener/installer/pkg/holder"
)

func UpdateIntegrationTask() {
	storeData := store.GetStore()
	holder.ApiHolder = *codefresh.GetInstance()

	err := holder.ApiHolder.UpdateIntegration(storeData.Codefresh.Integration, storeData.Argo.Host,
		"", "", storeData.Argo.Token, "", "", "")

	if err != nil {
		logger.GetLogger().Errorf("Failed to update integration, reason %v", err)
	}
}
