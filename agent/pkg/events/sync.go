package events

import (
	"github.com/codefresh-io/argocd-listener/agent/pkg/api/argo"
	"github.com/codefresh-io/argocd-listener/agent/pkg/api/codefresh"
	"github.com/codefresh-io/argocd-listener/agent/pkg/infra/logger"
	"github.com/codefresh-io/argocd-listener/agent/pkg/infra/store"
	"github.com/codefresh-io/argocd-listener/agent/pkg/util"
	"strings"
)

type SyncHandler struct {
	codefreshApi codefresh.CodefreshApi

	argoApi argo.ArgoAPI
}

var syncHandler *SyncHandler

func GetSyncHandlerInstance(codefreshApi codefresh.CodefreshApi, argoApi argo.ArgoAPI) *SyncHandler {
	if syncHandler != nil {
		return syncHandler
	}
	syncHandler = &SyncHandler{
		codefreshApi,
		argoApi,
	}
	return syncHandler
}

func (syncHandler *SyncHandler) Handle() error {
	syncMode := store.GetStore().Codefresh.SyncMode

	if syncMode == codefresh.None {
		logger.GetLogger().Info("Skip run sync handler because ")
		return nil
	}

	if syncMode == codefresh.OneTimeSync {
		applications, err := syncHandler.argoApi.GetApplicationsWithCredentialsFromStorage()
		if err != nil {
			return err
		}
		for _, application := range applications {
			err = syncHandler.codefreshApi.CreateEnvironment(application.Metadata.Name, application.Spec.Project, application.Metadata.Name, application.Spec.Destination.Namespace)
			if err != nil {
				logger.GetLogger().Errorf("Failed to create environment, reason %v", err)
			}
		}
	}

	if syncMode == codefresh.SelectSync {
		selectedApps := store.GetStore().Codefresh.ApplicationsForSync
		logger.GetLogger().Infof("Start sync applications: %v", strings.Join(selectedApps, ","))

		applications, err := syncHandler.argoApi.GetApplicationsWithCredentialsFromStorage()
		if err != nil {
			return err
		}
		for _, application := range applications {
			if util.Contains(selectedApps, application.Metadata.Name) {
				err = syncHandler.codefreshApi.CreateEnvironment(application.Metadata.Name, application.Spec.Project, application.Metadata.Name, application.Spec.Destination.Namespace)
				if err != nil {
					logger.GetLogger().Errorf("Failed to create environment, reason %v", err)
				}
			}
		}
	}

	return nil
}
