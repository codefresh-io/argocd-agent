package events

import (
	"encoding/json"
	"errors"
	"github.com/codefresh-io/argocd-listener/agent/pkg/api/argo"
	"github.com/codefresh-io/argocd-listener/agent/pkg/api/codefresh"
	"github.com/codefresh-io/argocd-listener/agent/pkg/infra/logger"
	"github.com/codefresh-io/argocd-listener/agent/pkg/service"
	"github.com/codefresh-io/argocd-listener/agent/pkg/transform"
	"github.com/codefresh-io/argocd-listener/agent/pkg/util"
	argo2 "github.com/codefresh-io/argocd-sdk/pkg/api"
	codefreshSdk "github.com/codefresh-io/go-sdk/pkg/codefresh"
)

// RolloutHandler handle rollout event, rollout meat that some new release appear or state of current release is changed
type RolloutHandler struct {
	codefreshApi                   codefresh.CodefreshApi
	argoApi                        argo.ArgoAPI
	argoResourceService            service.ArgoResourceService
	applicationResourceTransformer transform.Transformer
}

var rolloutHandler *RolloutHandler

// GetRolloutEventHandlerInstance get singleton instance of rollout handler
func GetRolloutEventHandlerInstance() EventHandler {
	if rolloutHandler != nil {
		return rolloutHandler
	}
	rolloutHandler = &RolloutHandler{
		codefreshApi:                   codefresh.GetInstance(),
		argoApi:                        argo.GetInstance(),
		argoResourceService:            service.NewArgoResourceService(),
		applicationResourceTransformer: transform.GetApplicationResourcesTransformer(),
	}
	return rolloutHandler
}

func convert(resources interface{}) []service.Resource {
	manifestResourcesJson, _ := json.Marshal(resources)

	var manifestResourcesStruct []service.Resource

	_ = json.Unmarshal(manifestResourcesJson, &manifestResourcesStruct)
	return manifestResourcesStruct
}

// Handle handle rollout event , process and store info in codefresh
func (rolloutHandler *RolloutHandler) Handle(rollout interface{}) error {
	envWrapper := rollout.(*service.EnvironmentWrapper)
	_, err := rolloutHandler.codefreshApi.SendEnvironment(envWrapper.Environment)
	if err != nil {
		return err
	}
	env := envWrapper.Environment

	resources, err := rolloutHandler.argoApi.GetResourceTreeAll(env.Name)
	if err != nil {
		return err
	}

	app, err := rolloutHandler.argoApi.GetApplication(env.Name)
	if err != nil {
		return err
	}

	var newApp argo2.ArgoApplication

	util.Convert(app, &newApp)

	statuses, ok := app["status"].(map[string]interface{})
	if !ok {
		return errors.New("Failed to parse data from retrieved application, app : " + env.Name)
	}

	manifestResources, ok := statuses["resources"].([]interface{})
	if !ok {
		return errors.New("Failed to parse data from retrieved application, app : " + env.Name)
	}

	manifestResourcesStruct := convert(manifestResources)

	result := rolloutHandler.argoResourceService.IdentifyChangedResources(newApp, manifestResourcesStruct, envWrapper.Commit, env.HistoryId, env.Date)

	appResources := rolloutHandler.applicationResourceTransformer.Transform(service.ResourcesWrapper{
		ResourcesTree:     resources.([]interface{}),
		ManifestResources: result,
	})
	if appResources != nil {
		applicationResources := &codefreshSdk.ApplicationResources{
			Name:      env.Name,
			HistoryId: env.HistoryId,
			Revision:  env.SyncRevision,
			Resources: appResources,
		}
		logger.GetLogger().Infof("Send application resources for app %s with history %v", env.Name, env.HistoryId)
		err = rolloutHandler.codefreshApi.SendApplicationResources(applicationResources)
	} else {
		logger.GetLogger().Infof("Skip send application resources for app %s with history %v, because resources not exists", env.Name, env.HistoryId)
	}

	return err
}
