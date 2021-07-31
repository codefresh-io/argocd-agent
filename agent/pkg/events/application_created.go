package events

import (
	"github.com/codefresh-io/argocd-listener/agent/pkg/api/codefresh"
	"github.com/codefresh-io/argocd-listener/agent/pkg/infra/logger"
	"github.com/codefresh-io/argocd-listener/agent/pkg/infra/store"
	"github.com/codefresh-io/argocd-listener/agent/pkg/service"
	argoSdk "github.com/codefresh-io/argocd-sdk/pkg/api"
)

type ApplicationCreatedHandler struct {
	rolloutEventHandler EventHandler
}

var applicationCreatedHandler *ApplicationCreatedHandler

func GetApplicationCreatedHandlerInstance() *ApplicationCreatedHandler {
	if applicationCreatedHandler != nil {
		return applicationCreatedHandler
	}
	applicationCreatedHandler = &ApplicationCreatedHandler{
		rolloutEventHandler: GetRolloutEventHandlerInstance(),
	}
	return applicationCreatedHandler
}

func (applicationCreatedHandler *ApplicationCreatedHandler) Handle(application argoSdk.ArgoApplication) error {
	if store.GetStore().Codefresh.SyncMode != codefresh.ContinueSync {
		// ignore handling if autosync disabled
		return nil
	}

	err := codefresh.GetInstance().CreateEnvironment(application.Metadata.Name, application.Spec.Project, application.Metadata.Name)
	if err != nil {
		return err
	}

	apps := service.NewGitopsService().HandleNewApplications([]string{application.Metadata.Name})

	for _, application := range apps {
		err := applicationCreatedHandler.rolloutEventHandler.Handle(application)
		if err != nil {
			logger.GetLogger().Errorf("Failed to send environment, reason %v", err)
		}
	}

	return nil
}
