package scheduler

import (
	"github.com/codefresh-io/argocd-listener/agent/pkg/api/argo"
	"github.com/codefresh-io/argocd-listener/agent/pkg/api/codefresh"
	"github.com/codefresh-io/argocd-listener/agent/pkg/events"
	"github.com/codefresh-io/argocd-listener/agent/pkg/infra/logger"
	"github.com/codefresh-io/argocd-listener/agent/pkg/infra/store"
	"github.com/codefresh-io/argocd-listener/agent/pkg/service"
	env2 "github.com/codefresh-io/argocd-listener/agent/pkg/transform/env"
	"github.com/codefresh-io/argocd-listener/agent/pkg/util"
	argoSdk "github.com/codefresh-io/argocd-sdk/pkg/api"
	codefreshSdk "github.com/codefresh-io/go-sdk/pkg/codefresh"
)

type envInitializerScheduler struct {
	codefreshApi        codefresh.CodefreshApi
	rolloutEventHandler events.EventHandler
	argoApi             argo.ArgoAPI
	argoResourceService service.ArgoResourceService
	envTransformer      env2.ETransformer
}

func GetEnvInitializerScheduler() Scheduler {
	return &envInitializerScheduler{
		codefreshApi:        codefresh.GetInstance(),
		rolloutEventHandler: events.GetRolloutEventHandlerInstance(),
		argoApi:             argo.GetInstance(),
		envTransformer:      env2.GetEnvTransformerInstance(argo.GetInstance()),
	}
}

func (envInitializer *envInitializerScheduler) getTime() string {
	return "@every 100s"
}

func (envInitializer *envInitializerScheduler) getFunc() func() {
	return envInitializer.handleEnvDifference
}

//Run start lister about new environments and update their statuses
func (envInitializer *envInitializerScheduler) Run() {
	run(envInitializer)
}

func isNewEnv(existingEnvs []store.Environment, newEnv codefreshSdk.CFEnvironment) bool {
	for _, env := range existingEnvs {
		if env.Name == newEnv.Metadata.Name {
			return false
		}
		continue
	}
	return true
}

func (envInitializer *envInitializerScheduler) extractNewApplication(application string) (*service.EnvironmentWrapper, error) {
	applicationObj, err := envInitializer.argoApi.GetApplication(application)
	if err != nil {
		return nil, err
	}

	var app argoSdk.ArgoApplication

	util.Convert(applicationObj, &app)

	err, historyId := envInitializer.argoResourceService.ResolveHistoryId(app.Status.History, app.Status.OperationState.SyncResult.Revision, app.Metadata.Name)
	if err != nil {
		return nil, err
	}

	err, envWrapper := envInitializer.envTransformer.PrepareEnvironment(app, historyId)
	if err != nil {
		return nil, err
	}
	return envWrapper, nil
}

func (envInitializer *envInitializerScheduler) handleNewApplications(applications []string) {
	for _, application := range applications {
		newApp, err := envInitializer.extractNewApplication(application)
		if err != nil {
			logger.GetLogger().Errorf("Failed to handle new gitops application %v, reason: %v", application, err)
			continue
		}
		logger.GetLogger().Infof("Detect new gitops application %s, initiate initialization", application)
		err = envInitializer.rolloutEventHandler.Handle(newApp)
		if err != nil {
			logger.GetLogger().Errorf("Failed to send environment, reason %v", err)
		}
	}
}

func (envInitializer *envInitializerScheduler) handleEnvDifference() {
	storeInst := store.GetStore()
	var newEnvs []store.Environment
	var applications []string
	envs, _ := envInitializer.codefreshApi.GetEnvironments()
	for _, env := range envs {
		if env.Spec.Type != "argo" {
			continue
		}

		newEnv := store.Environment{
			Name: env.Metadata.Name,
		}
		newEnvs = append(newEnvs, newEnv)

		if isNewEnv(storeInst.Environments, env) {
			applications = append(applications, env.Spec.Application)
		}
	}

	store.SetEnvironments(newEnvs)

	envInitializer.handleNewApplications(applications)

}
