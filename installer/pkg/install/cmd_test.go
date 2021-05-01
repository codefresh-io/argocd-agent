package install

import (
	"github.com/codefresh-io/argocd-listener/installer/pkg/install/entity"
	argoSdk "github.com/codefresh-io/argocd-sdk/pkg/api"
	codefreshSdk "github.com/codefresh-io/go-sdk/pkg/codefresh"
	"testing"
)

var _ = func() bool {
	testing.Init()
	return true
}()

type MockArgoApi struct {
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
	return "v1", nil
}

func (api *MockArgoApi) GetProjectsWithCredentialsFromStorage() ([]argoSdk.ProjectItem, error) {
	panic("implement me")
}

func (api *MockArgoApi) GetApplication(application string) (map[string]interface{}, error) {
	panic("implement me")
}

func (api *MockArgoApi) GetDefaultGitContext() (error, *codefreshSdk.ContextPayload) {
	panic("implement me")
}

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
	return nil
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

func (api *MockArgoApi) GetClusters() ([]argoSdk.ClusterItem, error) {
	panic("implement me")
}

func (api *MockArgoApi) GetApplications() ([]argoSdk.ApplicationItem, error) {
	panic("implement me")
}

func (api *MockArgoApi) GetRepositories() ([]argoSdk.RepositoryItem, error) {
	panic("implement me")
}

func TestEnsureIntegration(t *testing.T) {
	cmdOptions := entity.InstallCmdOptions{
		Kube: entity.Kube{},
		Argo: entity.ArgoOptions{},
		Codefresh: struct {
			Host                   string
			Token                  string
			Integration            string
			Suffix                 string
			SyncMode               string
			ApplicationsForSync    string
			ApplicationsForSyncArr []string
			Provider               string
		}{},
		Git: struct {
			Integration string
			Password    string
		}{},
		Host: struct {
			HttpProxy  string
			HttpsProxy string
		}{},
		Agent: struct {
			Version string
		}{},
	}

	installIntegration := InstallIntegration{
		argoApi:      &MockArgoApi{},
		codefreshApi: &MockCodefreshApi{},
	}
	err := installIntegration.ensureIntegration(&cmdOptions, "tet")
	if err != nil {
		t.Error("Should be executed without error")
	}
}
