package newrelic

import (
	"errors"
	"fmt"
	"github.com/codefresh-io/argocd-listener/agent/pkg/infra/logger"
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
	Init(newRelicLicense, envName string) error
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
	return &newrelicApp{}
}

func (a *newrelicApp) Init(newRelicLicense, envName string) error {
	if newRelicLicense != "" && envName != "" {
		logger.GetLogger().Infof("Initialize newrelic for env %s", envName)
		cfg := nr.NewConfig(fmt.Sprintf("argo-agent[%s]", envName), newRelicLicense)
		nrApp, err := nr.NewApplication(cfg)

		if err != nil {
			return err
		}
		a.api = nrApp
		return nil
	}

	return errors.New("failed to initiate new relic")
}

func (a *newrelicApp) RecordCustomEvent(eventType string, params EventParams) error {
	if a.api == nil {
		logger.GetLogger().Infof("failed to record new relic event because new relic app is not initialized")
		return errors.New("failed to record event because new relic app is not initialized")
	}
	var nrParams map[string]interface{}
	util.Convert(params, &nrParams)
	err := a.api.RecordCustomEvent(fmt.Sprintf("GitopsDashboardGit%s", eventType), nrParams)
	if err != nil {
		logger.GetLogger().Errorf("Newrelic RecordCustomEvent \"%s\" error %s", eventType, err.Error())
		return err
	}
	logger.GetLogger().Infof("Successfully record new relic event %s and params %v", eventType, params)
	return nil
}
