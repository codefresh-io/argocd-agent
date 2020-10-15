package extract

import (
	"github.com/codefresh-io/argocd-listener/agent/pkg/argo"
	"github.com/codefresh-io/argocd-listener/agent/pkg/codefresh"
	"github.com/codefresh-io/argocd-listener/agent/pkg/transform"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
)

func ExtractNewApplication(application string) (*codefresh.Environment, error) {
	applicationObj, err := argo.GetApplication(application)
	if err != nil {
		return nil, err
	}

	envTransformer := transform.GetEnvTransformerInstance(argo.GetInstance())

	err, env := envTransformer.PrepareEnvironment(applicationObj)
	if err != nil {
		return nil, err
	}
	return env, nil
}
