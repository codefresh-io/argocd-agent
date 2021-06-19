package acceptance

import (
	"errors"
	"fmt"
	"github.com/codefresh-io/argocd-listener/agent/pkg/api/argo"
	"github.com/codefresh-io/argocd-listener/agent/pkg/infra/store"
	"github.com/codefresh-io/argocd-listener/installer/pkg/dictionary"
	"github.com/codefresh-io/argocd-listener/installer/pkg/install/entity"
	"github.com/codefresh-io/argocd-listener/installer/pkg/prompt"
)

type ArgoAccessibilityAcceptanceTest struct {
	unathorizedArgoApi argo.UnauthorizedApi
	prompt             prompt.Prompt
}

func (acceptanceTest *ArgoAccessibilityAcceptanceTest) check(argoOptions *entity.ArgoOptions) (error, bool) {
	store.SetArgoHost(argoOptions.Host)
	if acceptanceTest.unathorizedArgoApi == nil {
		acceptanceTest.unathorizedArgoApi = argo.GetUnauthorizedApiInstance()
	}

	_, err := acceptanceTest.unathorizedArgoApi.GetVersion(argoOptions.Host)
	if err != nil {
		return errors.New(fmt.Sprintf(dictionary.CouldNotAccessToArgocdServer, argoOptions.Host)), false
	}
	return nil, false
}

func (acceptanceTest *ArgoAccessibilityAcceptanceTest) getMessage() string {
	return dictionary.CheckArgoServerAccessability
}

func (acceptanceTest *ArgoAccessibilityAcceptanceTest) failure(argoOptions *entity.ArgoOptions) bool {
	options := []string{
		dictionary.StopInstallation,
		dictionary.ContinueInstallationBehindFirewall,
	}
	err, result := acceptanceTest.prompt.Select(options, fmt.Sprintf(dictionary.CouldNotAccessToArgocdServer, argoOptions.Host))
	if err != nil {
		return true
	}

	if result == dictionary.StopInstallation {
		return true
	}

	if result == dictionary.ContinueInstallationBehindFirewall {
		argoOptions.FailFast = true
		return false
	}

	return true
}
