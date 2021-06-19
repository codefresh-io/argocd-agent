package acceptance

import (
	"errors"
	"github.com/codefresh-io/argocd-listener/agent/pkg/api/argo"
	"github.com/codefresh-io/argocd-listener/installer/pkg/install/entity"
)

type ProjectAcceptanceTest struct {
	argoApi argo.ArgoAPI
}

func (acceptanceTest *ProjectAcceptanceTest) check(argoOptions *entity.ArgoOptions) (error, bool) {
	if argoOptions.FailFast {
		return nil, true
	}
	if acceptanceTest.argoApi == nil {
		acceptanceTest.argoApi = argo.GetInstance()
	}

	projects, err := acceptanceTest.argoApi.GetProjectsWithCredentialsFromStorage()
	if err != nil {
		return err, false
	}
	if len(projects) == 0 {
		return errors.New("could not access your project in argocd, check credentials and whether you have an project set-up"), false
	}
	return nil, false
}

func (acceptanceTest *ProjectAcceptanceTest) getMessage() string {
	return "checking argocd projects accessibility..."
}

func (acceptanceTest *ProjectAcceptanceTest) failure(argoOptions *entity.ArgoOptions) bool {
	return true
}
