package questionnaire

import (
	"github.com/codefresh-io/argocd-listener/installer/pkg/install/entity"
	codefreshSdk "github.com/codefresh-io/go-sdk/pkg/codefresh"
	"testing"
)

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
	panic("implement me")
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

func TestGitContextQuestionnaire_AskAboutGitContext(t *testing.T) {
	questionnaire := GitContextQuestionnaire{codefreshApi: &MockCodefreshApi{}}
	options := &entity.InstallCmdOptions{}
	options.Git.Integration = "test"
	err := questionnaire.AskAboutGitContext(options)
	if err != nil {
		t.Error("Ask about context should fail without any error")
	}
}

func TestGitContextQuestionnaire_AskAboutGitContextWithoutIntegration(t *testing.T) {
	questionnaire := GitContextQuestionnaire{codefreshApi: &MockCodefreshApi{}}
	options := &entity.InstallCmdOptions{}
	err := questionnaire.AskAboutGitContext(options)
	if err != nil {
		t.Error("Ask about context should fail without any error")
	}
}
