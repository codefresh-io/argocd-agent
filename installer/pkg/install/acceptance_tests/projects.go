package acceptance_tests

import (
	"errors"
	"github.com/codefresh-io/argocd-listener/agent/pkg/argo"
	"github.com/codefresh-io/argocd-listener/installer/pkg/install"
)

type ProjectAcceptanceTest struct {
	argoApi argo.ArgoApi
}

func (acceptanceTest *ProjectAcceptanceTest) Check(argoOptions *install.ArgoOptions) error {
	projects, err := acceptanceTest.argoApi.GetProjectsWithCredentialsFromStorage()
	if err != nil {
		return err
	}
	if len(projects) == 0 {
		return errors.New("failed to retrieve projects, check token permissions or projects existence ")
	}
	return nil
}

func (acceptanceTest *ProjectAcceptanceTest) GetMessage() string {
	return "checking argocd projects accessibility..."
}
