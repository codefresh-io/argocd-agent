package acceptance

import (
	"github.com/codefresh-io/argocd-listener/agent/pkg/api/argo"
	"github.com/codefresh-io/argocd-listener/agent/pkg/infra/store"
	"github.com/codefresh-io/argocd-listener/installer/pkg/dictionary"
	"github.com/codefresh-io/argocd-listener/installer/pkg/install/entity"
)

type ArgoCredentialsAcceptanceTest struct {
	argoApi            argo.ArgoAPI
	unathorizedArgoApi argo.UnauthorizedApi
}

func (acceptanceTest *ArgoCredentialsAcceptanceTest) check(argoOptions *entity.ArgoOptions) error {
	if acceptanceTest.argoApi == nil {
		acceptanceTest.argoApi = argo.GetInstance()
	}
	if acceptanceTest.unathorizedArgoApi == nil {
		acceptanceTest.unathorizedArgoApi = argo.GetUnauthorizedApiInstance()
	}

	var err error
	token := argoOptions.Token
	if token == "" {
		token, err = acceptanceTest.unathorizedArgoApi.GetToken(argoOptions.Username, argoOptions.Password, argoOptions.Host)
		if err == nil {
			store.SetArgo(token, argoOptions.Host, "", "")
		}
	} else {
		store.SetArgo(token, argoOptions.Host, "", "")
		err = acceptanceTest.argoApi.CheckToken()
	}
	return err
}

func (acceptanceTest *ArgoCredentialsAcceptanceTest) getMessage() string {
	return dictionary.CheckArgoCredentials
}

func (acceptanceTest *ArgoCredentialsAcceptanceTest) failure() bool {
	return true
}
