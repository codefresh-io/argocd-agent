package newrelic

import (
	"fmt"
	"github.com/codefresh-io/argocd-listener/agent/pkg/infra/logger"
	"github.com/codefresh-io/argocd-listener/agent/pkg/infra/store"
	newrelic "github.com/newrelic/go-agent"
)

type NewrelicApp interface {
	RecordCustomEvent(eventType string, params map[string]interface{}) error
}

type newrelicApi struct {
	api newrelic.Application
}

var api *newrelicApi

func GetInstance() NewrelicApp {
	if api != nil {
		return api
	}

	storeState := store.GetStore()
	newRelicLicense := storeState.NewRelic.Key
	envName := storeState.Env.Name

	if newRelicLicense != "" {
		logger.GetLogger().Infof("Initialize newrelic for env %s", envName)
		cfg := newrelic.NewConfig(fmt.Sprintf("argo-agent[%s]", envName), newRelicLicense)
		newrelicApp, err := newrelic.NewApplication(cfg)
		api = &newrelicApi{
			api: newrelicApp,
		}

		logger.GetLogger().Errorf("Initialize newrelic error %s", err.Error())
	}

	return api
}

func (a *newrelicApi) RecordCustomEvent(eventType string, params map[string]interface{}) error {
	if api == nil {
		err := a.api.RecordCustomEvent(eventType, params)
		if err != nil {
			logger.GetLogger().Errorf("Newrelic RecordCustomEvent \"%s\" error %s", eventType, err.Error())
			return err
		}
	}
	return nil
}
