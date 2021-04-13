package service

import (
	codefreshSdk "github.com/codefresh-io/go-sdk/pkg/codefresh"
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
