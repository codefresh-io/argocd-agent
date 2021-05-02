package events

import (
	"github.com/codefresh-io/argocd-listener/agent/pkg/api/codefresh"
	"github.com/codefresh-io/argocd-listener/agent/pkg/service"
	argoSdk "github.com/codefresh-io/argocd-sdk/pkg/api"
	codefreshSdk "github.com/codefresh-io/go-sdk/pkg/codefresh"
	"testing"
)

var _ = func() bool {
	testing.Init()
	return true
}()

type PMockCodefreshApi struct {
}

func (api *PMockCodefreshApi) GetDefaultGitContext() (error, *codefreshSdk.ContextPayload) {
	panic("implement me")
}

func (api *PMockCodefreshApi) DeleteEnvironment(name string) error {
	panic("implement me")
}

func (api *PMockCodefreshApi) SendResources(kind string, items interface{}, amount int) error {
	panic("implement me")
}

func (api *PMockCodefreshApi) SendEvent(name string, props map[string]string) error {
	panic("implement me")
}

func (api *PMockCodefreshApi) HeartBeat(error string) error {
	panic("implement me")
}

func (api *PMockCodefreshApi) GetEnvironments() ([]codefreshSdk.CFEnvironment, error) {
	metadata := struct {
		Name string `json:"name"`
	}{Name: "test"}

	cfEnv := codefreshSdk.CFEnvironment{
		Metadata: metadata,
		Spec: struct {
			Type        string `json:"type"`
			Application string `json:"application"`
		}{Type: "argo", Application: "app"},
	}
	return []codefreshSdk.CFEnvironment{cfEnv}, nil
}

func (api *PMockCodefreshApi) CreateIntegration(name string, host string, username string, password string, token string, serverVersion string, provider string, clusterName string) error {
	panic("implement me")
}

func (api *PMockCodefreshApi) UpdateIntegration(name string, host string, username string, password string, token string, serverVersion string, provider string, clusterName string) error {
	return nil
}

func (api *PMockCodefreshApi) SendEnvironment(environment codefreshSdk.Environment) (map[string]interface{}, error) {
	return nil, nil
}

func (api *PMockCodefreshApi) CreateEnvironment(name string, project string, application string) error {
	panic("implement me")
}

func (api *PMockCodefreshApi) SendApplicationResources(resources *codefreshSdk.ApplicationResources) error {
	panic("implement me")
}

func (api *PMockCodefreshApi) GetGitContextByName(name string) (error, *codefreshSdk.ContextPayload) {
	return nil, nil
}

func (api *PMockCodefreshApi) GetGitContexts() (error, *[]codefreshSdk.ContextPayload) {
	panic("implement me")
}

type PMockArgoApi struct {
}

func (api *PMockArgoApi) CheckToken() error {
	panic("implement me")
}

func (api *PMockArgoApi) GetResourceTree(applicationName string) (*argoSdk.ResourceTree, error) {
	panic("implement me")
}

func (api *PMockArgoApi) GetResourceTreeAll(applicationName string) (interface{}, error) {
	return make([]interface{}, 0), nil
}

func (api *PMockArgoApi) GetManagedResources(applicationName string) (*argoSdk.ManagedResource, error) {
	panic("implement me")
}

func (api *PMockArgoApi) GetVersion() (string, error) {
	panic("implement me")
}

func (api *PMockArgoApi) GetProjectsWithCredentialsFromStorage() ([]argoSdk.ProjectItem, error) {
	panic("implement me")
}

func (api *PMockArgoApi) GetApplication(application string) (map[string]interface{}, error) {

	resources := make([]interface{}, 0)
	status := make(map[string]interface{})
	status["resources"] = resources

	result := make(map[string]interface{})
	result["status"] = status

	return result, nil
}

func (api *PMockArgoApi) GetDefaultGitContext() (error, *codefreshSdk.ContextPayload) {
	panic("implement me")
}

func (api *PMockArgoApi) GetClusters() ([]argoSdk.ClusterItem, error) {
	panic("implement me")
}

func (api *PMockArgoApi) GetApplications() ([]argoSdk.ApplicationItem, error) {
	panic("implement me")
}

func (api *PMockArgoApi) GetRepositories() ([]argoSdk.RepositoryItem, error) {
	panic("implement me")
}
func (api *PMockArgoApi) GetApplicationsWithCredentialsFromStorage() ([]argoSdk.ApplicationItem, error) {
	applications := make([]argoSdk.ApplicationItem, 0)
	applications = append(applications, argoSdk.ApplicationItem{
		Metadata: argoSdk.ApplicationMetadata{
			Name: "Test",
		},
		Spec: argoSdk.ApplicationSpec{
			Project: "Test-Project",
		},
	})

	return applications, nil
}

type MockApplicationResourceTransformer struct {
}

func (transformer *MockApplicationResourceTransformer) Transform(data interface{}) interface{} {
	return nil
}

type MockArgoResourceService struct {
}

func (argoResourceService *MockArgoResourceService) IdentifyChangedResources(app argoSdk.ArgoApplication, resources []service.Resource, commit service.ResourceCommit, historyId int64, updateAt string) []service.Resource {
	return []service.Resource{}
}

func (argoResourceService *MockArgoResourceService) AdaptArgoProjects(projects []argoSdk.ProjectItem) []codefresh.AgentProject {
	panic("implement me")
}

func (argoResourceService *MockArgoResourceService) AdaptArgoApplications(applications []argoSdk.ApplicationItem) []codefresh.AgentApplication {
	panic("implement me")
}

func (argoResourceService *MockArgoResourceService) ResolveHistoryId(historyList []argoSdk.ApplicationHistoryItem, revision string, name string) (error, int64) {
	return nil, 1
}

func TestRolloutHandler(t *testing.T) {

	rolloutHandler := RolloutHandler{
		codefreshApi:                   &PMockCodefreshApi{},
		argoApi:                        &PMockArgoApi{},
		argoResourceService:            &MockArgoResourceService{},
		applicationResourceTransformer: &MockApplicationResourceTransformer{},
	}

	wrapper := &service.EnvironmentWrapper{
		Environment: codefreshSdk.Environment{},
		Commit:      service.ResourceCommit{},
	}

	err := rolloutHandler.Handle(wrapper)
	if err != nil {
		t.Error("Rollout should be handler without error")
	}
}
