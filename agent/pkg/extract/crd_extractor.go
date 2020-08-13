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
	env := transform.PrepareEnvironment(applicationObj)
	return &env, nil
}
