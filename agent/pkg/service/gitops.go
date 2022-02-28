package service

import (
	"errors"
	"github.com/codefresh-io/argocd-listener/agent/pkg/api/argo"
	"github.com/codefresh-io/argocd-listener/agent/pkg/infra/git/provider"
	"github.com/codefresh-io/argocd-listener/agent/pkg/infra/logger"
	"github.com/codefresh-io/argocd-listener/agent/pkg/util"
	argoSdk "github.com/codefresh-io/argocd-sdk/pkg/api"
	codefreshSdk "github.com/codefresh-io/go-sdk/pkg/codefresh"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

type Gitops interface {
	MarkEnvAsRemoved(obj interface{}) (error, *codefreshSdk.Environment)
	HandleNewApplications(applications []string) []*EnvironmentWrapper
	ExtractNewApplication(application string) (*EnvironmentWrapper, error)
}

type EnvironmentWrapper struct {
	Environment codefreshSdk.Environment
	Commit      provider.ResourceCommit
}

type gitops struct {
	argoApi             argo.ArgoAPI
	argoResourceService ArgoResourceService
	envTransformer      ETransformer
	sharding            util.Sharding
}

func NewGitopsService() Gitops {
	return &gitops{
		argoApi:             argo.GetInstance(),
		envTransformer:      GetEnvTransformerInstance(argo.GetInstance()),
		argoResourceService: NewArgoResourceService(),
	}
}

func (gitops *gitops) MarkEnvAsRemoved(obj interface{}) (error, *codefreshSdk.Environment) {
	//disable for now
	//envTransformer := env.GetEnvTransformerInstance(argo.GetInstance())
	//err, env := envTransformer.PrepareEnvironment(obj.(*unstructured.Unstructured).Object)
	//if err != nil {
	//	return err, env
	//}
	//
	//env.HealthStatus = "Deleted"
	//
	//err = events.GetRolloutEventHandlerInstance().Handle(env)
	//
	//if err != nil {
	//	return err, nil
	//}

	return nil, nil
}

func (gitops *gitops) ExtractNewApplication(application string) (*EnvironmentWrapper, error) {
	applicationObj, err := gitops.argoApi.GetApplication(application)
	if err != nil {
		return nil, err
	}

	process := gitops.sharding.ShouldBeProcessed(&unstructured.Unstructured{Object: applicationObj})
	if !process {
		return nil, errors.New("should not be processed by this shard")
	}

	var app argoSdk.ArgoApplication

	util.Convert(applicationObj, &app)

	err, historyId := gitops.argoResourceService.ResolveHistoryId(app.Status.History, app.Status.OperationState.SyncResult.Revision, app.Metadata.Name)
	if err != nil {
		return nil, err
	}

	err, envWrapper := gitops.envTransformer.PrepareEnvironment(app, historyId)
	if err != nil {
		return nil, err
	}
	return envWrapper, nil
}

func (gitops *gitops) HandleNewApplications(applications []string) []*EnvironmentWrapper {
	var apps []*EnvironmentWrapper
	for _, application := range applications {
		newApp, err := gitops.ExtractNewApplication(application)
		if err != nil {
			logger.GetLogger().Errorf("Failed to handle new gitops application %v, reason: %v", application, err)
			continue
		}
		logger.GetLogger().Infof("Detect new gitops application %s, initiate initialization", application)
		apps = append(apps, newApp)
	}
	return apps
}
