package acceptance_tests

import (
	"github.com/codefresh-io/argocd-listener/installer/pkg/install"
	argoSdk "github.com/codefresh-io/argocd-sdk/pkg/api"
	"testing"
)

var _ = func() bool {
	testing.Init()
	return true
}()

type PrjMockArgoApi struct {
}

func (api *PrjMockArgoApi) CheckToken() error {
	panic("implement me")
}

func (api *PrjMockArgoApi) GetResourceTree(applicationName string) (*argoSdk.ResourceTree, error) {
	panic("implement me")
}

func (api *PrjMockArgoApi) GetResourceTreeAll(applicationName string) (interface{}, error) {
	panic("implement me")
}

func (api *PrjMockArgoApi) GetManagedResources(applicationName string) (*argoSdk.ManagedResource, error) {
	panic("implement me")
}

func (api *PrjMockArgoApi) GetVersion() (string, error) {
	panic("implement me")
}

func (api *PrjMockArgoApi) GetApplicationsWithCredentialsFromStorage() ([]argoSdk.ApplicationItem, error) {
	panic("implement me")
}

func (api *PrjMockArgoApi) GetProjectsWithCredentialsFromStorage() ([]argoSdk.ProjectItem, error) {
	return []argoSdk.ProjectItem{}, nil
}

func (api *PrjMockArgoApi) GetApplication(application string) (map[string]interface{}, error) {
	panic("implement me")
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
