package acceptance

import (
	"errors"
	"github.com/codefresh-io/argocd-listener/agent/pkg/argo"
	"github.com/codefresh-io/argocd-listener/installer/pkg/install/entity"
)

type ApplicationAcceptanceTest struct {
	argoApi argo.ArgoAPI
}

func (acceptanceTest *ApplicationAcceptanceTest) check(argoOptions *entity.ArgoOptions) error {
	if acceptanceTest.argoApi == nil {
		acceptanceTest.argoApi = argo.GetInstance()
	}

	applications, err := acceptanceTest.argoApi.GetApplicationsWithCredentialsFromStorage()
	if err != nil {
		return err
	}
	if len(applications) == 0 {
		return errors.New("could not access your application in argocd, check credentials and whether you have an application set-up")
	}
	return err
}

func (acceptanceTest *ApplicationAcceptanceTest) getMessage() string {
	return "checking argocd applications accessibility..."
}
