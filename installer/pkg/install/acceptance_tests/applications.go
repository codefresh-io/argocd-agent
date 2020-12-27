package acceptance_tests

import (
	"github.com/codefresh-io/argocd-listener/agent/pkg/argo"
	"github.com/codefresh-io/argocd-listener/installer/pkg/install"
)

type ApplicationAcceptanceTest struct {
}

func (acceptanceTest *ApplicationAcceptanceTest) Check(argoOptions *install.ArgoOptions) error {
	_, err := argo.GetInstance().GetApplicationsWithCredentialsFromStorage()
	return err
}

func (acceptanceTest *ApplicationAcceptanceTest) GetMessage() string {
	return "checking argocd applications accessibility..."
}
