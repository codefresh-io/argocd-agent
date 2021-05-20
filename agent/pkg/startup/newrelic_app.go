package startup

import (
	"fmt"
	"github.com/codefresh-io/argocd-listener/agent/pkg/infra/store"
	newrelic "github.com/newrelic/go-agent"
)

var nrApp newrelic.Application

type NewRelicApp struct {
	input *Input
}

func NewNewrelicApp(input *Input) *NewRelicApp {
	return &NewRelicApp{input}
}

func (NewRelicApp *NewRelicApp) Init() error {
	storeState := store.GetStore()
	newRelicLicense := storeState.NewRelic.Key
	envName := storeState.Env.Name

	if newRelicLicense != "" {
		cfg := newrelic.NewConfig(fmt.Sprintf("argo-agent[%s]", envName), newRelicLicense)
		_, err := newrelic.NewApplication(cfg)
		return err
	}
	return nil
}
