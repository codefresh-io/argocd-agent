package acceptance_tests

import (
	"github.com/codefresh-io/argocd-listener/installer/pkg/install"
	argoSdk "github.com/codefresh-io/argocd-sdk/pkg/api"
	"testing"
)

type MockArgoApi struct {
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

func TestEmptyResultOfApplications(t *testing.T) {
	test := &ApplicationAcceptanceTest{
		argoApi: &MockArgoApi{},
	}
	result := test.Check(&install.ArgoOptions{})

	if result == nil {
		t.Errorf("Acceptance test should be fail with error")
		return
	}

	if result.Error() != "could not access your application in argocd, check credentials and whether you have an application set-up" {
		t.Errorf("Acceptance test should be fail with error \"failed to retrieve applications, check token permissions or applications existence\", actual: %s", result.Error())
	}
}
