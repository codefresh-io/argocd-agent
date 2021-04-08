package acceptance

import (
	"github.com/codefresh-io/argocd-listener/agent/pkg/argo"
	"github.com/codefresh-io/argocd-listener/agent/pkg/store"
	"github.com/codefresh-io/argocd-listener/installer/pkg/install/entity"
)

type ArgoCredentialsAcceptanceTest struct {
}

func (acceptanceTest *ArgoCredentialsAcceptanceTest) check(argoOptions *entity.ArgoOptions) error {
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

func (acceptanceTest *ArgoCredentialsAcceptanceTest) getMessage() string {
	return "checking argocd credentials..."
}
