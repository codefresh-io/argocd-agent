package newrelic

import (
	"fmt"
	"github.com/codefresh-io/argocd-listener/agent/pkg/infra/logger"
	"github.com/codefresh-io/argocd-listener/agent/pkg/infra/store"
	"github.com/codefresh-io/argocd-listener/agent/pkg/util"
	nr "github.com/newrelic/go-agent"
)

const (
	EventGetUsers         = "GetUsers"
	EventGetIssues        = "GetIssues"
	EventGetCommitBySha   = "GetCommitBySha"
	EventGetCommitsBySha  = "GetCommitsBySha"
	EventListPullRequests = "ListPullRequests"
)

type EventParams struct {
	Owner string `json:"owner"`
	Repo  string `json:"repo"`
}

type NewrelicApp interface {
	RecordCustomEvent(eventType string, params EventParams) error
}

type newrelicApp struct {
	api nr.Application
}

var app *newrelicApp

func GetInstance() NewrelicApp {
	if app != nil {
		return app
	}

	storeState := store.GetStore()
	newRelicLicense := storeState.NewRelic.Key
	envName := storeState.Env.Name

	if newRelicLicense != "" {
		logger.GetLogger().Infof("Initialize newrelic for env %s", envName)
		cfg := nr.NewConfig(fmt.Sprintf("argo-agent[%s]", envName), newRelicLicense)
		nrApp, err := nr.NewApplication(cfg)
		app = &newrelicApp{
			api: nrApp,
		}

		logger.GetLogger().Errorf("Initialize newrelic error %s", err.Error())
	}

	return app
}

func (a *newrelicApp) RecordCustomEvent(eventType string, params EventParams) error {
	if app == nil {
		var nrParams map[string]interface{}
		util.Convert(params, &nrParams)
		err := a.api.RecordCustomEvent(fmt.Sprintf("GitopsDashboardGit::%s", eventType), nrParams)
		if err != nil {
			logger.GetLogger().Errorf("Newrelic RecordCustomEvent \"%s\" error %s", eventType, err.Error())
			return err
		}
	}
	return nil
}
