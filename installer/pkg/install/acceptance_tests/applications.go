package acceptance_tests

import (
	"errors"
	"github.com/codefresh-io/argocd-listener/agent/pkg/argo"
	"github.com/codefresh-io/argocd-listener/installer/pkg/install"
)

type ApplicationAcceptanceTest struct {
	argoApi argo.ArgoApi
}

func (acceptanceTest *ApplicationAcceptanceTest) Check(argoOptions *install.ArgoOptions) error {
	if acceptanceTest.argoApi == nil {
		acceptanceTest.argoApi = argo.GetInstance()
	}

	applications, err := acceptanceTest.argoApi.GetApplicationsWithCredentialsFromStorage()
	if err != nil {
		return err
	}
	if len(applications) == 0 {
		return errors.New("failed to retrieve applications, check token permissions or applications existence ")
	}
	return err
}

func (acceptanceTest *ApplicationAcceptanceTest) GetMessage() string {
	return "checking argocd applications accessibility..."
}
