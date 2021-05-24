package startup

import (
	"github.com/codefresh-io/argocd-listener/agent/pkg/api/argo"
	"github.com/codefresh-io/argocd-listener/agent/pkg/api/codefresh"
	"github.com/codefresh-io/argocd-listener/agent/pkg/infra/store"
)

type InputStore struct {
	input *Input
}

func NewInputStore(input *Input) *InputStore {
	return &InputStore{input}
}

func (inputStore *InputStore) Store() error {
	input := inputStore.input
	if input.argoToken == "" {
		token, err := argo.GetUnauthorizedApiInstance().GetToken(input.argoUsername, input.argoPassword, input.argoHost)
		if err != nil {
			return err
		}
		store.SetArgo(token, input.argoHost, input.argoUsername, input.argoPassword)
	} else {
		store.SetArgo(input.argoToken, input.argoHost, "", "")
	}

	store.SetSyncOptions(input.syncMode, input.applications)

	store.SetCodefresh(input.codefreshHost, input.codefreshToken, input.codefreshIntegrationName)

	store.SetEnvOptions(input.envName)

	if input.agentVersion != "" {
		store.SetAgent(input.agentVersion)
	}

	if input.password != "" {
		store.SetGit(input.password)
	}

	if input.gitIntegration != "" {
		err, gitContext := codefresh.GetInstance().GetGitContextByName(input.gitIntegration)
		if err == nil {
			store.SetGitContext(*gitContext)
		}
	}

	return nil

}
