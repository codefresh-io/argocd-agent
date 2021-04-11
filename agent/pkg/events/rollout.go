package events

import (
	"errors"
	"github.com/codefresh-io/argocd-listener/agent/pkg/api/argo"
	"github.com/codefresh-io/argocd-listener/agent/pkg/api/codefresh"
	"github.com/codefresh-io/argocd-listener/agent/pkg/infra/logger"
	"github.com/codefresh-io/argocd-listener/agent/pkg/transform"
	codefreshSdk "github.com/codefresh-io/go-sdk/pkg/codefresh"
)

// RolloutHandler handle rollout event, rollout meat that some new release appear or state of current release is changed
type RolloutHandler struct {
}

var rolloutHandler *RolloutHandler

// GetRolloutEventHandlerInstance get singleton instance of rollout handler
func GetRolloutEventHandlerInstance() EventHandler {
	if rolloutHandler != nil {
		return rolloutHandler
	}
	rolloutHandler = &RolloutHandler{}
	return rolloutHandler
}

// Handle handle rollout event , process and store info in codefresh
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

	app, err := argo.GetInstance().GetApplication(env.Name)
	if err != nil {
		return err
	}

	statuses, ok := app["status"].(map[string]interface{})
	if !ok {
		return errors.New("Failed to parse data from retrieved application, app : " + env.Name)
	}

	manifestResources, ok := statuses["resources"].([]interface{})
	if !ok {
		return errors.New("Failed to parse data from retrieved application, app : " + env.Name)
	}

	appResources := transform.GetApplicationResourcesTransformer().Transform(argo.ResourcesWrapper{
		ResourcesTree:     resources.([]interface{}),
		ManifestResources: manifestResources,
	})
	if appResources != nil {
		applicationResources := &codefreshSdk.ApplicationResources{
			Name:      env.Name,
			HistoryId: env.HistoryId,
			Revision:  env.SyncRevision,
			Resources: appResources,
		}
		logger.GetLogger().Infof("Send application resources for app %s with history %v", env.Name, env.HistoryId)
		err = codefresh.GetInstance().SendApplicationResources(applicationResources)
	} else {
		logger.GetLogger().Infof("Skip send application resources for app %s with history %v, because resources not exists", env.Name, env.HistoryId)
	}

	return err
}
