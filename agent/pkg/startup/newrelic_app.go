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

func (NewRelicApp *NewRelicApp) Init() {
	newRelicLicense := store.GetStore().NewRelic.Key

	if newRelicLicense != "" {
		log.Debug("setting New Relic agent")
		cfg := newrelic.NewConfig("argo-agent[kubernetes]", newRelicLicense)
		var err error
		nrApp, err = newrelic.NewApplication(cfg)
		if nil != err {
			log.WithError(err).Error("failed to setup New Relic agent with provided license")
		}
	}
}
