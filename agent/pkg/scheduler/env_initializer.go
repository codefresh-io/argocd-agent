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
	"github.com/jasonlvhit/gocron"
)

const envInitializationTime uint64 = 25

func isNewEnv(existingEnvs []store.Environment, newEnv codefreshSdk.CFEnvironment) bool {
	for _, env := range existingEnvs {
		if env.Name == newEnv.Metadata.Name {
			return false
		}
		continue
	}
	return true
}

func extractNewApplication(application string) (*service.EnvironmentWrapper, error) {
	applicationObj, err := argo.GetInstance().GetApplication(application)
	if err != nil {
		return nil, err
	}

	var app argoSdk.ArgoApplication

	util.Convert(applicationObj, &app)

	envTransformer := env2.GetEnvTransformerInstance(argo.GetInstance())

	err, historyId := service.NewArgoResourceService().ResolveHistoryId(app.Status.History, app.Status.OperationState.SyncResult.Revision, app.Metadata.Name)
	if err != nil {
		return nil, err
	}

	err, envWrapper := envTransformer.PrepareEnvironment(app, historyId)
	if err != nil {
		return nil, err
	}
	return envWrapper, nil
}

func handleNewApplications(applications []string) {
	for _, application := range applications {
		newApp, err := extractNewApplication(application)
		if err != nil {
			logger.GetLogger().Errorf("Failed to handle new gitops application %v, reason: %v", application, err)
			continue
		}
		logger.GetLogger().Infof("Detect new gitops application %s, initiate initialization", application)
		err = events.GetRolloutEventHandlerInstance().Handle(newApp)
		if err != nil {
			logger.GetLogger().Errorf("Failed to send environment, reason %v", err)
		}
	}
}

func handleEnvDifference() {
	storeInst := store.GetStore()
	var newEnvs []store.Environment
	var applications []string
	envs, _ := codefresh.GetInstance().GetEnvironments()
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

	handleNewApplications(applications)

}

// StartEnvInitializer start lister about new environments and update their statuses
func StartEnvInitializer() {
	job := gocron.Every(envInitializationTime).Seconds().Do(handleEnvDifference)

	if job != nil {
		err := job.Error()

		if err != "" {
			panic("Cant run env changes job because " + err)
		}
	}

	go gocron.Start()
}
