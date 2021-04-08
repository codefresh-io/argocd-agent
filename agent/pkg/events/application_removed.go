package events

import (
	argoSdk "github.com/codefresh-io/argocd-sdk/pkg/api"
)

type ApplicationRemovedHandler struct {
}

var applicationRemovedHandler *ApplicationRemovedHandler

func GetApplicationRemovedHandlerInstance() *ApplicationRemovedHandler {
	if applicationRemovedHandler != nil {
		return applicationRemovedHandler
	}
	applicationRemovedHandler = &ApplicationRemovedHandler{}
	return applicationRemovedHandler
}

func (applicationRemovedHandler *ApplicationRemovedHandler) Handle(application argoSdk.ArgoApplication) error {
	return nil
}
