package scheduler

import (
	"fmt"
	"github.com/codefresh-io/argocd-listener/agent/pkg/api/codefresh"
	"github.com/codefresh-io/argocd-listener/agent/pkg/infra/git/provider"
	"github.com/codefresh-io/argocd-listener/agent/pkg/infra/store"
	"github.com/codefresh-io/argocd-listener/agent/pkg/service"
	argoSdk "github.com/codefresh-io/argocd-sdk/pkg/api"
	codefreshSdk "github.com/codefresh-io/go-sdk/pkg/codefresh"
	"testing"
	"time"
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
			Context     string `json:"context"`
		}{Type: "argo", Application: "app", Context: "test"},
	}
	return []codefreshSdk.CFEnvironment{cfEnv}, nil
}

func (api *PMockCodefreshApi) CreateIntegration(name string, host string, username string, password string, token string, serverVersion string, provider string, clusterName string) error {
	panic("implement me")
}

func (api *PMockCodefreshApi) UpdateIntegration(name string, host string, username string, password string, token string, serverVersion string, provider string, clusterName string) error {
	return UpdateIntegration()
}

func (api *PMockCodefreshApi) SendEnvironment(environment codefreshSdk.Environment) (map[string]interface{}, error) {
	panic("implement me")
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

type MockArgoApi struct {
}

func (api *MockArgoApi) CreateDefaultApp() error {
	return nil
}

func (api *MockArgoApi) CheckToken() error {
	panic("implement me")
}

func (api *MockArgoApi) GetResourceTree(applicationName string) (*argoSdk.ResourceTree, error) {
	panic("implement me")
}

func (api *MockArgoApi) GetResourceTreeAll(applicationName string) (interface{}, error) {
	panic("implement me")
}

func (api *MockArgoApi) GetManagedResources(applicationName string) (*argoSdk.ManagedResource, error) {
	panic("implement me")
}

func (api *MockArgoApi) GetVersion() (string, error) {
	panic("implement me")
}

func (api *MockArgoApi) GetProjectsWithCredentialsFromStorage() ([]argoSdk.ProjectItem, error) {
	panic("implement me")
}

func (api *MockArgoApi) GetApplication(application string) (map[string]interface{}, error) {
	return make(map[string]interface{}), nil
}

func (api *MockArgoApi) GetDefaultGitContext() (error, *codefreshSdk.ContextPayload) {
	panic("implement me")
}

func (api *MockArgoApi) GetClusters() ([]argoSdk.ClusterItem, error) {
	panic("implement me")
}

func (api *MockArgoApi) GetApplications() ([]argoSdk.ApplicationItem, error) {
	panic("implement me")
}

func (api *MockArgoApi) GetRepositories() ([]argoSdk.RepositoryItem, error) {
	panic("implement me")
}
func (api *MockArgoApi) GetApplicationsWithCredentialsFromStorage() ([]argoSdk.ApplicationItem, error) {
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

var rolloutHandlerFunc func(payload interface{}) error

type MockRolloutEventHandler struct {
}

func (rollout *MockRolloutEventHandler) Handle(payload interface{}) error {
	return rolloutHandlerFunc(payload)
}

type MockEnvTransformer struct {
}

func (transformer *MockEnvTransformer) PrepareEnvironment(app argoSdk.ArgoApplication, historyId int64) (error, *service.EnvironmentWrapper) {
	return nil, &service.EnvironmentWrapper{
		Environment: codefreshSdk.Environment{},
		Commit:      provider.ResourceCommit{},
	}
}

type MockArgoResourceService struct {
}

func (argoResourceService *MockArgoResourceService) IdentifyChangedResources(app argoSdk.ArgoApplication, resources []service.Resource, commit provider.ResourceCommit, historyId int64, updateAt string) []*service.Resource {
	panic("implement me")
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

type GitopsService struct {
}

func (gitopsService *GitopsService) MarkEnvAsRemoved(obj interface{}) (error, *codefreshSdk.Environment) {
	return nil, nil
}

func (gitopsService *GitopsService) HandleNewApplications(applications []string) []*service.EnvironmentWrapper {
	return []*service.EnvironmentWrapper{
		&service.EnvironmentWrapper{
			Environment: codefreshSdk.Environment{},
			Commit:      provider.ResourceCommit{},
		},
	}
}

func (gitopsService *GitopsService) ExtractNewApplication(application string) (*service.EnvironmentWrapper, error) {
	return nil, nil
}

func TestNewEnv(t *testing.T) {

	envs := []store.Environment{
		{
			Name: "test",
		},
	}

	metadata := struct {
		Name string `json:"name"`
	}{Name: "test"}

	cfEnv := codefreshSdk.CFEnvironment{
		Metadata: metadata,
	}

	newEnv := isNewEnv(envs, cfEnv)
	if newEnv {
		t.Error("Should return that env already exist")
	}
}

func TestHandleEnvDifference(t *testing.T) {

	envInitializerScheduler := envInitializerScheduler{
		codefreshApi:        &PMockCodefreshApi{},
		rolloutEventHandler: &MockRolloutEventHandler{},
		argoApi:             &MockArgoApi{},
		argoResourceService: &MockArgoResourceService{},
		envTransformer:      &MockEnvTransformer{},
		gitopsService:       &GitopsService{},
	}

	result := make(chan string)

	rolloutHandlerFunc = func(payload interface{}) error {
		result <- "ok"
		return nil
	}

	go envInitializerScheduler.handleEnvDifference()

	select {
	case ret := <-result:
		fmt.Println(ret)
	case <-time.After(10 * time.Second):
		t.Error("Rollout func should be called")
	}

}

func TestCreateInstance(t *testing.T) {
	envInitializerScheduler := GetEnvInitializerScheduler()
	if envInitializerScheduler == nil {
		t.Error("Cant initialize env initializer")
	}
}

func TestExecutionTime(t *testing.T) {
	envInitializerScheduler := GetEnvInitializerScheduler()
	time := envInitializerScheduler.getTime()
	if time != "@every 100s" {
		t.Error("Wrong schedule time")
	}
}
