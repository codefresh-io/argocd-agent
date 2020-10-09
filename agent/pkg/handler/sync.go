package handler

import (
	"github.com/codefresh-io/argocd-listener/agent/pkg/argo"
	"github.com/codefresh-io/argocd-listener/agent/pkg/codefresh"
	"github.com/codefresh-io/argocd-listener/agent/pkg/logger"
	"github.com/codefresh-io/argocd-listener/agent/pkg/store"
	"github.com/codefresh-io/argocd-listener/agent/pkg/util"
	"strings"
)

type SyncHandler struct {
}

var syncHandler *SyncHandler

func GetSyncHandlerInstance() *SyncHandler {
	if syncHandler != nil {
		return syncHandler
	}
	syncHandler = &SyncHandler{}
	return syncHandler
}

func (syncHandler *SyncHandler) Handle() error {
	syncMode := store.GetStore().Codefresh.SyncMode

	if syncMode == codefresh.None {
		logger.GetLogger().Info("Skip run sync handler because ")
		return nil
	}

	if syncMode == codefresh.OneTimeSync {
		applications, err := argo.GetApplicationsWithCredentialsFromStorage()
		if err != nil {
			return err
		}
		api := codefresh.GetInstance()
		for _, application := range applications {
			err = api.CreateEnvironment(application.Metadata.Name, application.Spec.Project, application.Metadata.Name)
			if err != nil {
				logger.GetLogger().Errorf("Failed to create environment, reason %v", err)
			}
		}
	}

	if syncMode == codefresh.SelectSync {
		selectedApps := store.GetStore().Codefresh.ApplicationsForSync
		logger.GetLogger().Infof("Start sync applications: %v", strings.Join(selectedApps, ","))

		applications, err := argo.GetApplicationsWithCredentialsFromStorage()
		if err != nil {
			return err
		}
		api := codefresh.GetInstance()
		for _, application := range applications {
			if util.Contains(selectedApps, application.Metadata.Name) {
				err = api.CreateEnvironment(application.Metadata.Name, application.Spec.Project, application.Metadata.Name)
				if err != nil {
					logger.GetLogger().Errorf("Failed to create environment, reason %v", err)
				}
			}
		}
	}

	return nil
}
