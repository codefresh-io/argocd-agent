package event_handler

import (
	"github.com/codefresh-io/argocd-listener/agent/pkg/codefresh"
	//"github.com/codefresh-io/argocd-listener/agent/pkg/scheduler"
	"github.com/codefresh-io/argocd-listener/agent/pkg/store"
	argoSdk "github.com/codefresh-io/argocd-sdk/pkg/api"
)

type ApplicationCreatedHandler struct {
}

var applicationCreatedHandler *ApplicationCreatedHandler

func GetApplicationCreatedHandlerInstance() *ApplicationCreatedHandler {
	if applicationCreatedHandler != nil {
		return applicationCreatedHandler
	}
	applicationCreatedHandler = &ApplicationCreatedHandler{}
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

	//scheduler.HandleNewApplications([]string{application.Metadata.Name})

	return nil
}