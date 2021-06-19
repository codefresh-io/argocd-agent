package acceptance

import (
	"errors"
	"github.com/codefresh-io/argocd-listener/agent/pkg/api/argo"
	"github.com/codefresh-io/argocd-listener/installer/pkg/dictionary"
	"github.com/codefresh-io/argocd-listener/installer/pkg/install/entity"
	"github.com/codefresh-io/argocd-listener/installer/pkg/prompt"
)

type ApplicationAcceptanceTest struct {
	argoApi argo.ArgoAPI
	prompt  prompt.Prompt
}

func (acceptanceTest *ApplicationAcceptanceTest) check(argoOptions *entity.ArgoOptions) (error, bool) {
	if argoOptions.FailFast {
		return nil, true
	}

	if acceptanceTest.argoApi == nil {
		acceptanceTest.argoApi = argo.GetInstance()
	}

	applications, err := acceptanceTest.argoApi.GetApplicationsWithCredentialsFromStorage()
	if err != nil {
		return err, false
	}
	if len(applications) == 0 {
		return errors.New("could not access your application in argocd, check credentials and whether you have an application set-up"), false
	}
	return err, false
}

func (acceptanceTest *ApplicationAcceptanceTest) getMessage() string {
	return dictionary.CheckArgoApplicationsAccessability
}

func (acceptanceTest *ApplicationAcceptanceTest) failure(argoOptions *entity.ArgoOptions) bool {
	options := []string{
		dictionary.StopInstallation,
		dictionary.ContinueInstallation,
		dictionary.SetupDemoApplication,
	}
	err, result := acceptanceTest.prompt.Select(options, dictionary.NoApplication)
	if err != nil {
		return true
	}

	if result == dictionary.StopInstallation {
		return true
	}

	if result == dictionary.ContinueInstallation {
		return false
	}

	if result == dictionary.SetupDemoApplication {
		_ = acceptanceTest.argoApi.CreateDefaultApp()
		return false
	}

	return true
}
