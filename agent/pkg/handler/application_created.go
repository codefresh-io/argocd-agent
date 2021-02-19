package handler

import (
	"github.com/codefresh-io/argocd-listener/agent/pkg/argo"
	"github.com/codefresh-io/argocd-listener/agent/pkg/codefresh"
	"github.com/codefresh-io/argocd-listener/agent/pkg/scheduler"
	"github.com/codefresh-io/argocd-listener/agent/pkg/store"
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

func (applicationCreatedHandler *ApplicationCreatedHandler) Handle(application argo.ArgoApplication) error {
	if store.GetStore().Codefresh.SyncMode != codefresh.ContinueSync {
		// ignore handling if autosync disabled
		return nil
	}

	api := codefresh.GetInstance()
	err := api.CreateEnvironment(application.Metadata.Name, application.Spec.Project, application.Metadata.Name)
	if err != nil {
		return err
	}

	scheduler.HandleNewApplications([]string{application.Metadata.Name})

	return nil
}
