package acceptance_tests

import (
	"errors"
	"github.com/codefresh-io/argocd-listener/agent/pkg/argo"
	"github.com/codefresh-io/argocd-listener/installer/pkg/install"
)

type ProjectAcceptanceTest struct {
	argoApi argo.ArgoAPI
}

func (acceptanceTest *ProjectAcceptanceTest) Check(argoOptions *install.ArgoOptions) error {
	if acceptanceTest.argoApi == nil {
		acceptanceTest.argoApi = argo.GetInstance()
	}

	projects, err := acceptanceTest.argoApi.GetProjectsWithCredentialsFromStorage()
	if err != nil {
		return err
	}
	if len(projects) == 0 {
		return errors.New("could not access your project in argocd, check credentials and whether you have an project set-up")
	}
	return nil
}

func (acceptanceTest *ProjectAcceptanceTest) GetMessage() string {
	return "checking argocd projects accessibility..."
}
