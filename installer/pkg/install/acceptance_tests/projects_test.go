package acceptance_tests

import (
	"github.com/codefresh-io/argocd-listener/agent/pkg/argo"
	"github.com/codefresh-io/argocd-listener/installer/pkg/install"
	"testing"
)

type PrjMockArgoApi struct {
}

func (api *PrjMockArgoApi) GetResourceTree(applicationName string) (*argo.ResourceTree, error) {
	panic("implement me")
}

func (api *PrjMockArgoApi) GetResourceTreeAll(applicationName string) (interface{}, error) {
	panic("implement me")
}

func (api *PrjMockArgoApi) GetManagedResources(applicationName string) (*argo.ManagedResource, error) {
	panic("implement me")
}

func (api *PrjMockArgoApi) GetVersion() (string, error) {
	panic("implement me")
}

func (api *PrjMockArgoApi) GetApplicationsWithCredentialsFromStorage() ([]argo.ApplicationItem, error) {
	panic("implement me")
}

func (api *PrjMockArgoApi) GetProjectsWithCredentialsFromStorage() ([]argo.ProjectItem, error) {
	return []argo.ProjectItem{}, nil
}

func TestEmptyResultOfProjects(t *testing.T) {
	test := &ProjectAcceptanceTest{
		argoApi: &PrjMockArgoApi{},
	}
	result := test.Check(&install.ArgoOptions{})

	if result == nil {
		t.Errorf("Project test should be fail with error")
		return
	}

	if result.Error() != "could not access your project in argocd, check credentials and whether you have an project set-up" {
		t.Errorf("Acceptance test should be fail with error \"failed to retrieve projects, check token permissions or projects existence\", actual: %s", result.Error())
	}
}
