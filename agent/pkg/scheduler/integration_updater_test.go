package scheduler

import (
	codefreshSdk "github.com/codefresh-io/go-sdk/pkg/codefresh"
	"testing"
)

var UpdateIntegration func() error

type MockCodefreshApi struct {
}

func (api *MockCodefreshApi) GetDefaultGitContext() (error, *codefreshSdk.ContextPayload) {
	panic("implement me")
}

func (api *MockCodefreshApi) DeleteEnvironment(name string) error {
	panic("implement me")
}

func (api *MockCodefreshApi) SendResources(kind string, items interface{}, amount int) error {
	panic("implement me")
}

func (api *MockCodefreshApi) SendEvent(name string, props map[string]string) error {
	panic("implement me")
}

func (api *MockCodefreshApi) HeartBeat(error string) error {
	panic("implement me")
}

func (api *MockCodefreshApi) GetEnvironments() ([]codefreshSdk.CFEnvironment, error) {
	panic("implement me")
}

func (api *MockCodefreshApi) CreateIntegration(name string, host string, username string, password string, token string, serverVersion string, provider string, clusterName string) error {
	panic("implement me")
}

func (api *MockCodefreshApi) UpdateIntegration(name string, host string, username string, password string, token string, serverVersion string, provider string, clusterName string) error {
	return UpdateIntegration()
}

func (api *MockCodefreshApi) SendEnvironment(environment codefreshSdk.Environment) (map[string]interface{}, error) {
	panic("implement me")
}

func (api *MockCodefreshApi) CreateEnvironment(name string, project string, application string) error {
	panic("implement me")
}

func (api *MockCodefreshApi) SendApplicationResources(resources *codefreshSdk.ApplicationResources) error {
	panic("implement me")
}

func TestIntegrationUpdaterScheduler(t *testing.T) {

	integrationUpdaterScheduler := integrationUpdaterScheduler{codefreshApi: &MockCodefreshApi{}}
	time := integrationUpdaterScheduler.getTime()
	if time != "@every 100s" {
		t.Error("Wrong scheduling time")
	}

	ch := make(chan string)

	UpdateIntegration = func() error {
		ch <- "ok"
		return nil
	}

	go integrationUpdaterScheduler.updateIntegrationTask()

	result := <-ch
	if result != "ok" {
		t.Error("UpdateIntegration should be called")
	}
}
