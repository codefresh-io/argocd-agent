package startup

import (
	"github.com/codefresh-io/argocd-listener/agent/pkg/infra/store"
	newrelic "github.com/newrelic/go-agent"
	log "github.com/sirupsen/logrus"
)

var nrApp newrelic.Application

type NewRelicApp struct {
	input *Input
}

func NewNewrelicApp(input *Input) *NewRelicApp {
	return &NewRelicApp{input}
}

func (NewRelicApp *NewRelicApp) Init() error {
	newRelicLicense := store.GetStore().NewRelic.Key

	if newRelicLicense != "" {
		log.Debug("setting New Relic agent")
		cfg := newrelic.NewConfig("argo-agent[kubernetes]", newRelicLicense)
		_, err := newrelic.NewApplication(cfg)
		return err
	}
	return nil
}
