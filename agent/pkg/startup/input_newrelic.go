package startup

import (
	"github.com/codefresh-io/argocd-listener/agent/pkg/infra/newrelic"
	"github.com/codefresh-io/argocd-listener/agent/pkg/infra/store"
)

type InputNewrelic struct {
	input *Input
}

func NewInputNewrelic(input *Input) *InputNewrelic {
	return &InputNewrelic{input}
}

func (InputNewrelic *InputNewrelic) Init() error {
	storeState := store.GetStore()
	newRelicLicense := storeState.NewRelic.Key
	envName := storeState.Env.Name

	if newRelicLicense != "" {
		return newrelic.Init(newRelicLicense, envName)
	}
	return nil
}
