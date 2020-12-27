package acceptance_tests

import (
	"github.com/codefresh-io/argocd-listener/agent/pkg/argo"
	"github.com/codefresh-io/argocd-listener/installer/pkg/install"
	"testing"
)

type MockArgoApi struct {
}

func (api *MockArgoApi) GetResourceTree(applicationName string) (*argo.ResourceTree, error) {
	panic("implement me")
}

func (api *MockArgoApi) GetResourceTreeAll(applicationName string) (interface{}, error) {
	panic("implement me")
}

func (api *MockArgoApi) GetManagedResources(applicationName string) (*argo.ManagedResource, error) {
	panic("implement me")
}

func (api *MockArgoApi) GetVersion() (string, error) {
	panic("implement me")
}

func (api *MockArgoApi) GetProjectsWithCredentialsFromStorage() ([]argo.ProjectItem, error) {
	panic("implement me")
}

func (api *MockArgoApi) GetApplicationsWithCredentialsFromStorage() ([]argo.ApplicationItem, error) {
	return []argo.ApplicationItem{}, nil
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

	if result.Error() != "failed to retrieve applications, check token permissions or applications existence " {
		t.Errorf("Acceptance test should be fail with error \"failed to retrieve applications, check token permissions or applications existence\", actual: %s", result.Error())
	}
}
