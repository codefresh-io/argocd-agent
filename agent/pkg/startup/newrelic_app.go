package startup

import (
	"github.com/codefresh-io/argocd-listener/agent/pkg/infra/store"
	"github.com/codefresh-io/go-infra/pkg/logger"
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
		cfg := newrelic.NewConfig("nomios[kubernetes]", newRelicLicense)
		var err error
		nrApp, err = newrelic.NewApplication(cfg)
		if nil != err {
			log.WithError(err).Error("failed to setup New Relic agent with provided license")
			return err
		}
		log.Debug("setting New Relic agent hook for Logrus logging")
		nrHook := logger.NewNewRelicLogrusHook(nrApp, []log.Level{log.ErrorLevel, log.FatalLevel, log.PanicLevel})
		log.AddHook(nrHook)
	}
	return nil
}
