package event_handler

import (
	"github.com/codefresh-io/argocd-listener/agent/pkg/argo"
	"github.com/codefresh-io/argocd-listener/agent/pkg/codefresh"
	"github.com/codefresh-io/argocd-listener/agent/pkg/transform"
	codefreshSdk "github.com/codefresh-io/go-sdk/pkg/codefresh"
)

type RolloutHandler struct {
}

var rolloutHandler *RolloutHandler

func GetRolloutEventHandlerInstance() EventHandler {
	if rolloutHandler != nil {
		return rolloutHandler
	}
	rolloutHandler = &RolloutHandler{}
	return rolloutHandler
}

func (rolloutHandler *RolloutHandler) Handle(rollout interface{}) error {
	env := rollout.(*codefreshSdk.Environment)
	_, err := codefresh.GetInstance().SendEnvironment(*env)
	if err != nil {
		return err
	}

	resources, err := argo.GetInstance().GetResourceTreeAll(env.Name)
	if err != nil {
		return err
	}

	appResources := transform.GetApplicationResourcesTransformer().Transform(resources)
	if appResources != nil {
		applicationResources := &codefreshSdk.ApplicationResources{
			Name:      env.Name,
			HistoryId: env.HistoryId,
			Revision:  env.SyncRevision,
			Resources: appResources,
		}

		err = codefresh.GetInstance().SendApplicationResources(applicationResources)
	}

	return err
}
