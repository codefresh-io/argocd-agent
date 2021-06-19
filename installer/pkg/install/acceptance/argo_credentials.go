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

func (acceptanceTest *ArgoCredentialsAcceptanceTest) check(argoOptions *entity.ArgoOptions) (error, bool) {
	if argoOptions.FailFast {
		return nil, true
	}

	var err error
	token := argoOptions.Token
	if token == "" {
		if acceptanceTest.unathorizedArgoApi == nil {
			acceptanceTest.unathorizedArgoApi = argo.GetUnauthorizedApiInstance()
		}
		token, err = acceptanceTest.unathorizedArgoApi.GetToken(argoOptions.Username, argoOptions.Password, argoOptions.Host)
		if err == nil {
			store.SetArgo(token, argoOptions.Host, "", "")
		}
	} else {
		store.SetArgo(token, argoOptions.Host, "", "")
		if acceptanceTest.argoApi == nil {
			acceptanceTest.argoApi = argo.GetInstance()
		}
		err = acceptanceTest.argoApi.CheckToken()
	}
	return err, false
}

func (acceptanceTest *ArgoCredentialsAcceptanceTest) getMessage() string {
	return dictionary.CheckArgoCredentials
}

func (acceptanceTest *ArgoCredentialsAcceptanceTest) failure(argoOptions *entity.ArgoOptions) bool {
	return true
}
