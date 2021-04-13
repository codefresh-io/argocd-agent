package service

import (
	"github.com/codefresh-io/argocd-listener/agent/pkg/api/argo"
	"github.com/codefresh-io/argocd-listener/agent/pkg/transform/env"
	codefreshSdk "github.com/codefresh-io/go-sdk/pkg/codefresh"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

type Gitops interface {
	MarkEnvAsRemoved(obj interface{}) (error, *codefreshSdk.Environment)
}

type gitops struct {
}

func NewGitopsService() Gitops {
	return &gitops{}
}

func (gitops *gitops) MarkEnvAsRemoved(obj interface{}) (error, *codefreshSdk.Environment) {
	envTransformer := env.GetEnvTransformerInstance(argo.GetInstance())
	err, env := envTransformer.PrepareEnvironment(obj.(*unstructured.Unstructured).Object)
	if err != nil {
		return err, env
	}

	env.HealthStatus = "Deleted"

	//err = events.GetRolloutEventHandlerInstance().Handle(env)

	//if err != nil {
	//	return err, nil
	//}

	return nil, nil
}
