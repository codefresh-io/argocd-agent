package acceptance

import (
	"github.com/codefresh-io/argocd-listener/installer/pkg/dictionary"
	"github.com/codefresh-io/argocd-listener/installer/pkg/install/entity"
	argoSdk "github.com/codefresh-io/argocd-sdk/pkg/api"
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

func (api *MockArgoApi) GetApplication(application string) (map[string]interface{}, error) {
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

func (api *MockArgoApi) GetApplicationsWithCredentialsFromStorage() ([]argoSdk.ApplicationItem, error) {
	return []argoSdk.ApplicationItem{}, nil
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

func (api *MockArgoApi) CreateDefaultApp() error {
	return nil
}

type MockPrompt struct {
}

func (p *MockPrompt) InputWithDefault(target *string, label string, defaultValue string) error {
	*target = defaultValue
	return nil
}

func (p *MockPrompt) InputPassword(target *string, label string) error {
	return nil
}

func (p *MockPrompt) Input(target *string, label string) error {
	return nil
}

func (p *MockPrompt) Confirm(label string) (error, bool) {
	return nil, false
}

func (p *MockPrompt) Multiselect(items []string, label string) (error, []string) {
	return nil, nil
}

func (p *MockPrompt) Select(items []string, label string) (error, string) {
	return nil, dictionary.StopInstallation
}

func TestEmptyResultOfApplications(t *testing.T) {
	test := &ApplicationAcceptanceTest{
		argoApi: &MockArgoApi{},
	}
	result := test.check(&entity.ArgoOptions{})

	if result == nil {
		t.Errorf("Acceptance test should be fail with error")
		return
	}

	if result.Error() != "could not access your application in argocd, check credentials and whether you have an application set-up" {
		t.Errorf("Acceptance test should be fail with error \"failed to retrieve applications, check token permissions or applications existence\", actual: %s", result.Error())
	}
}

func TestFailureCase(t *testing.T) {
	test := &ApplicationAcceptanceTest{
		argoApi: &MockArgoApi{},
		prompt:  &MockPrompt{},
	}
	result := test.failure()

	if result == false {
		t.Errorf("Failure should stop an installation")
		return
	}
}
