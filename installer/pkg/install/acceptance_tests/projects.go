package acceptance_tests

import (
	"github.com/codefresh-io/argocd-listener/agent/pkg/argo"
	"github.com/codefresh-io/argocd-listener/installer/pkg/install"
)

type ProjectAcceptanceTest struct {
}

func (acceptanceTest *ProjectAcceptanceTest) Check(argoOptions *install.ArgoOptions) error {
	_, err := argo.GetProjectsWithCredentialsFromStorage()
	return err
}

func (acceptanceTest *ProjectAcceptanceTest) GetMessage() string {
	return "checking argocd projects accessibility..."
}
