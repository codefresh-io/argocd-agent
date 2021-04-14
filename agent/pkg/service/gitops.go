package service

import (
	"github.com/codefresh-io/argocd-listener/agent/pkg/api/codefresh"
	"github.com/codefresh-io/argocd-listener/agent/pkg/infra/logger"
	"github.com/codefresh-io/argocd-listener/agent/pkg/infra/store"
	codefreshSdk "github.com/codefresh-io/go-sdk/pkg/codefresh"
)

type Gitops interface {
	MarkEnvAsRemoved(obj interface{}) (error, *codefreshSdk.Environment)
	UpdateIntegration()
}

type gitops struct {
}

func NewGitopsService() Gitops {
	return &gitops{}
}

func (gitops *gitops) UpdateIntegration() {
	storeData := store.GetStore()

	err := codefresh.GetInstance().UpdateIntegration(storeData.Codefresh.Integration, storeData.Argo.Host,
		"", "", storeData.Argo.Token, "", "", "")

	if err != nil {
		logger.GetLogger().Errorf("Failed to update integration, reason %v", err)
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
