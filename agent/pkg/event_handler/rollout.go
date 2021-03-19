package event_handler

import (
	"github.com/codefresh-io/argocd-listener/agent/pkg/codefresh"
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

	//resources, err := argo.GetInstance().GetResourceTreeAll(env.Name)
	//if err != nil {
	//	return err
	//}

	//appResources := transform.GetApplicationResourcesTransformer().Transform(resources)
	//if appResources != nil {
	//	applicationResources := &codefreshSdk.ApplicationResources{
	//		Name:      env.Name,
	//		HistoryId: env.HistoryId,
	//		Revision:  env.SyncRevision,
	//		Resources: appResources,
	//	}
	//	logger.GetLogger().Infof("Send application resources for app %s with history %v", env.Name, env.HistoryId)
	//	err = codefresh.GetInstance().SendApplicationResources(applicationResources)
	//} else {
	//	logger.GetLogger().Infof("Skip send application resources for app %s with history %v, because resources not exists", env.Name, env.HistoryId)
	//}

	return err
}
