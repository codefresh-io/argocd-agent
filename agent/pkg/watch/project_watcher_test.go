package watch

import (
	argoSdk "github.com/codefresh-io/argocd-sdk/pkg/api"
	codefreshSdk "github.com/codefresh-io/go-sdk/pkg/codefresh"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"testing"
)

var _ = func() bool {
	testing.Init()
	return true
}()

var sendResourcesP func(len int) error

type PMockCodefreshApi struct {
}

func (api *PMockCodefreshApi) GetDefaultGitContext() (error, *codefreshSdk.ContextPayload) {
	panic("implement me")
}

func (api *PMockCodefreshApi) DeleteEnvironment(name string) error {
	panic("implement me")
}

func (api *PMockCodefreshApi) SendResources(kind string, items interface{}, amount int) error {
	return sendResourcesP(amount)
}

func (api *PMockCodefreshApi) SendEvent(name string, props map[string]string) error {
	panic("implement me")
}

func (api *PMockCodefreshApi) HeartBeat(error string) error {
	panic("implement me")
}

func (api *PMockCodefreshApi) GetEnvironments() ([]codefreshSdk.CFEnvironment, error) {
	panic("implement me")
}

func (api *PMockCodefreshApi) CreateIntegration(name string, host string, username string, password string, token string, serverVersion string, provider string, clusterName string) error {
	panic("implement me")
}

func (api *PMockCodefreshApi) UpdateIntegration(name string, host string, username string, password string, token string, serverVersion string, provider string, clusterName string) error {
	panic("implement me")
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
	metadata := struct {
		Name string `json:"name"`
	}{Name: "test"}
	return nil, &[]codefreshSdk.ContextPayload{
		{
			Metadata: metadata,
		},
	}
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
	panic("implement me")
}

func (api *PMockArgoApi) GetManagedResources(applicationName string) (*argoSdk.ManagedResource, error) {
	panic("implement me")
}

func (api *PMockArgoApi) GetVersion() (string, error) {
	panic("implement me")
}

func (api *PMockArgoApi) GetProjectsWithCredentialsFromStorage() ([]argoSdk.ProjectItem, error) {
	projects := make([]argoSdk.ProjectItem, 0)
	projects = append(projects, argoSdk.ProjectItem{
		Metadata: argoSdk.ProjectMetadata{
			Name: "Test",
		},
	})

	return projects, nil
}

func (api *PMockArgoApi) GetApplication(application string) (map[string]interface{}, error) {
	panic("implement me")
}

func (api *PMockArgoApi) GetDefaultGitContext() (error, *codefreshSdk.ContextPayload) {
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

func (api *PMockArgoApi) GetClusters() ([]argoSdk.ClusterItem, error) {
	panic("implement me")
}

func (api *PMockArgoApi) GetApplications() ([]argoSdk.ApplicationItem, error) {
	panic("implement me")
}

func (api *PMockArgoApi) GetRepositories() ([]argoSdk.RepositoryItem, error) {
	panic("implement me")
}

func TestProjectWatcherDeleteEvent(t *testing.T) {

	projectwatcher := projectWatcher{
		codefreshApi:    &PMockCodefreshApi{},
		informer:        nil,
		informerFactory: nil,
		argoApi:         &PMockArgoApi{},
	}

	obj := &unstructured.Unstructured{}

	sendResourcesP = func(len int) error {
		if len != 1 {
			t.Error("Unable watch delete events")
		}
		return nil
	}

	projectwatcher.delete(obj)

}

func TestProjectWatcherCreateEvent(t *testing.T) {
	projectwatcher := projectWatcher{
		codefreshApi:    &PMockCodefreshApi{},
		informer:        nil,
		informerFactory: nil,
		argoApi:         &PMockArgoApi{},
	}

	obj := &unstructured.Unstructured{}

	sendResourcesP = func(len int) error {
		if len != 1 {
			t.Error("Unable watch create events")
		}
		return nil
	}

	projectwatcher.add(obj)

}

func TestNewProjectWatcher(t *testing.T) {

	var watcher *projectWatcher
	watcher, err := NewProjectWatcher()

	if watcher.codefreshApi == nil || watcher.informer == nil || watcher.informerFactory == nil || watcher.argoApi == nil {
		t.Error("Missing watcher properties")
	}

	if err != nil {
		t.Error(err)
	}
}
