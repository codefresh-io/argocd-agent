package newrelic

import (
	"fmt"
	"github.com/codefresh-io/argocd-listener/agent/pkg/infra/logger"
	nr "github.com/newrelic/go-agent"
)

var newrelicApp nr.Application

func Init(newRelicLicense string, envName string) error {
	var err error
	if newRelicLicense != "" {
		logger.GetLogger().Infof("Initialize newrelic for env %s", envName)
		cfg := nr.NewConfig(fmt.Sprintf("argo-agent[%s]", envName), newRelicLicense)
		newrelicApp, err = nr.NewApplication(cfg)
		return err
	}
	return nil
}

func RecordCustomEvent(eventType string, params map[string]interface{}) error {
	if newrelicApp == nil {
		err := newrelicApp.RecordCustomEvent(fmt.Sprintf("GitopsDashboardGit::%s", eventType), params)
		if err != nil {
			logger.GetLogger().Errorf("Newrelic RecordCustomEvent \"%s\" error %s", eventType, err.Error())
			return err
		}
	}
	return nil
}
