package integration_updater

import (
	"github.com/codefresh-io/argocd-listener/agent/pkg/argo"
	"github.com/codefresh-io/argocd-listener/agent/pkg/codefresh"
	"github.com/codefresh-io/argocd-listener/agent/pkg/logger"
	"github.com/codefresh-io/argocd-listener/agent/pkg/store"
	"github.com/codefresh-io/argocd-listener/installer/pkg/holder"
)

var heartbeatAmount = 0

func UpdateIntegrationTask() {
	storeData := store.GetStore()
	holder.ApiHolder = codefresh.Api{
		Token:       storeData.Codefresh.Token,       // +
		Host:        storeData.Codefresh.Host,        // +
		Integration: storeData.Codefresh.Integration, // +
	}
	argoToken := storeData.Argo.Token
	clusters, _ := argo.GetClusters(argoToken, storeData.Argo.Host)
	applications, _ := argo.GetApplications(argoToken, storeData.Argo.Host)
	repositories, _ := argo.GetRepositories(argoToken, storeData.Argo.Host)

	err := holder.ApiHolder.UpdateIntegration(storeData.Codefresh.Integration, storeData.Argo.Host,
		"", "", storeData.Argo.Token, "", "", "",
		len(clusters), len(applications), len(repositories))

	if err != nil {
		logger.GetLogger().Errorf("Failed to update integration, reason %v", err)
	}
}
