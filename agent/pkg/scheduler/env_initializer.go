package scheduler

import (
	"github.com/codefresh-io/argocd-listener/agent/pkg/api/argo"
	"github.com/codefresh-io/argocd-listener/agent/pkg/api/codefresh"
	"github.com/codefresh-io/argocd-listener/agent/pkg/events"
	"github.com/codefresh-io/argocd-listener/agent/pkg/infra/logger"
	"github.com/codefresh-io/argocd-listener/agent/pkg/infra/store"
	"github.com/codefresh-io/argocd-listener/agent/pkg/service"
	codefreshSdk "github.com/codefresh-io/go-sdk/pkg/codefresh"
)

type envInitializerScheduler struct {
	codefreshApi        codefresh.CodefreshApi
	rolloutEventHandler events.EventHandler
	argoApi             argo.ArgoAPI
	argoResourceService service.ArgoResourceService
	gitopsService       service.Gitops
	envTransformer      service.ETransformer
}

func GetEnvInitializerScheduler() Scheduler {
	return &envInitializerScheduler{
		codefreshApi:        codefresh.GetInstance(),
		rolloutEventHandler: events.GetRolloutEventHandlerInstance(),
		argoApi:             argo.GetInstance(),
		envTransformer:      service.GetEnvTransformerInstance(argo.GetInstance()),
		argoResourceService: service.NewArgoResourceService(),
		gitopsService:       service.NewGitopsService(),
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
		if env.Name == newEnv.Metadata.Name && env.Context == newEnv.Spec.Context {
			return false
		}
		continue
	}
	return true
}

func (envInitializer *envInitializerScheduler) extractNewApplication(application string) (*service.EnvironmentWrapper, error) {
	return envInitializer.gitopsService.ExtractNewApplication(application)
}

func (envInitializer *envInitializerScheduler) handleNewApplications(applications []string) {
	apps := envInitializer.gitopsService.HandleNewApplications(applications)

	for _, application := range apps {
		err := envInitializer.rolloutEventHandler.Handle(application)
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
			Name:    env.Metadata.Name,
			Context: env.Spec.Context,
		}
		newEnvs = append(newEnvs, newEnv)

		if isNewEnv(storeInst.Environments, env) {
			applications = append(applications, env.Spec.Application)
		}
	}

	store.SetEnvironments(newEnvs)

	envInitializer.handleNewApplications(applications)

}
