package cfeventsender

import (
	"fmt"
	codefreshSdk "github.com/codefresh-io/go-sdk/pkg/codefresh"
	"testing"
	"time"
)

var _ = func() bool {
	testing.Init()
	return true
}()

var sendEventMock func(eventName string) error

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
	return sendEventMock(name)
}

func (api *MockCodefreshApi) HeartBeat(error string) error {
	panic("implement me")
}

func (api *MockCodefreshApi) GetEnvironments() ([]codefreshSdk.CFEnvironment, error) {
	panic("implement me")
}

func (api *MockCodefreshApi) UpdateIntegration(name string, host string, username string, password string, token string, serverVersion string, provider string, clusterName string) error {
	panic("implement me")
}

func (api *MockCodefreshApi) SendEnvironment(environment codefreshSdk.Environment) (map[string]interface{}, error) {
	panic("implement me")
}

func (api *MockCodefreshApi) CreateEnvironment(name string, project string, application string) error {
	return nil
}

func (api *MockCodefreshApi) SendApplicationResources(resources *codefreshSdk.ApplicationResources) error {
	panic("implement me")
}

func (api *MockCodefreshApi) CreateIntegration(name string, host string, username string, password string, token string, serverVersion string, provider string, clusterName string) error {
	panic("implement me")
}

func (api *MockCodefreshApi) GetGitContextByName(name string) (error, *codefreshSdk.ContextPayload) {
	return nil, nil
}

func (api *MockCodefreshApi) GetGitContexts() (error, *[]codefreshSdk.ContextPayload) {
	metadata := struct {
		Name string `json:"name"`
	}{Name: "test"}
	return nil, &[]codefreshSdk.ContextPayload{
		{
			Metadata: metadata,
		},
	}
}

func TestNew(t *testing.T) {
	client := New(EVENT_AGENT_UNINSTALL)
	if client.eventName != EVENT_AGENT_UNINSTALL {
		t.Errorf("'TestNew' failed, unexpected event name after init, expected '%v', got '%v'", EVENT_AGENT_UNINSTALL, client.eventName)
	}

	client = New(EVENT_AGENT_INSTALL)
	if client.eventName != EVENT_AGENT_INSTALL {
		t.Errorf("'TestNew' failed, must return existing state, expected '%v', got '%v'", EVENT_AGENT_INSTALL, client.eventName)
	}
}

func TestSuccess(t *testing.T) {
	result := make(chan string)

	sendEventMock = func(name string) error {
		if name != EVENT_AGENT_INSTALL {
			t.Error("Event is incorrect")
		}
		result <- "ok"
		return nil
	}

	cfApi := &MockCodefreshApi{}

	client := New(EVENT_AGENT_INSTALL).WithCustomCodefreshApi(cfApi)

	go client.Success("Something went wrong")

	select {
	case ret := <-result:
		fmt.Println(ret)
	case <-time.After(10 * time.Second):
		t.Error("Create integration func should be called")
	}

}

func TestFailed(t *testing.T) {
	result := make(chan string)

	sendEventMock = func(name string) error {
		if name != EVENT_AGENT_UNINSTALL {
			t.Error("Event is incorrect")
		}
		result <- "ok"
		return nil
	}

	cfApi := &MockCodefreshApi{}

	client := New(EVENT_AGENT_UNINSTALL).WithCustomCodefreshApi(cfApi)

	go client.Fail("Something went wrong")

	select {
	case ret := <-result:
		fmt.Println(ret)
	case <-time.After(10 * time.Second):
		t.Error("Create integration func should be called")
	}

}
