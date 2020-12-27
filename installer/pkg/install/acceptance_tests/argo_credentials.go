package acceptance_tests

import (
	"github.com/codefresh-io/argocd-listener/agent/pkg/argo"
	"github.com/codefresh-io/argocd-listener/agent/pkg/store"
	"github.com/codefresh-io/argocd-listener/installer/pkg/install"
)

type ArgoCredentialsAcceptanceTest struct {
}

func (acceptanceTest *ArgoCredentialsAcceptanceTest) Check(argoOptions *install.ArgoOptions) error {
	var err error
	token := argoOptions.Token
	if token == "" {
		token, err = argo.GetToken(argoOptions.Username, argoOptions.Password, argoOptions.Host)
		if err == nil {
			store.SetArgo(token, argoOptions.Host)
		}
	} else {
		store.SetArgo(token, argoOptions.Host)
		err = argo.GetInstance().CheckToken()
	}
	return err
}

func (acceptanceTest *ArgoCredentialsAcceptanceTest) GetMessage() string {
	return "checking argocd credentials..."
}
