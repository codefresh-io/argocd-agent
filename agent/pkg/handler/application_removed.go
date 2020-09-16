package handler

import (
	"github.com/codefresh-io/argocd-listener/agent/pkg/argo"
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

func (applicationRemovedHandler *ApplicationRemovedHandler) Handle(application argo.ArgoApplication) error {
	return nil
}
