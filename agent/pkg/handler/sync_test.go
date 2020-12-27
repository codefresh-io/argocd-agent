package handler

import (
	"github.com/codefresh-io/argocd-listener/agent/pkg/argo"
	"github.com/codefresh-io/argocd-listener/agent/pkg/codefresh"
	"github.com/codefresh-io/argocd-listener/agent/pkg/store"
	"testing"
)

var createdEnv []string

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

type MockCodefreshApi struct {
}

func (api *MockCodefreshApi) CreateEnvironment(name string, project string, application string) error {
	createdEnv = append(createdEnv, name)
	return nil
}

func (api *MockArgoApi) GetApplicationsWithCredentialsFromStorage() ([]argo.ApplicationItem, error) {
	applications := make([]argo.ApplicationItem, 0)
	applications = append(applications, argo.ApplicationItem{
		Metadata: argo.ApplicationMetadata{
			Name: "Test",
		},
		Spec: argo.ApplicationSpec{
			Project: "Test-Project",
		},
	})

	return applications, nil
}

func TestSyncWithNoneMode(t *testing.T) {

	createdEnv = make([]string, 0)

	store.SetSyncOptions(codefresh.None, []string{})

	syncHandler := GetSyncHandlerInstance(&MockCodefreshApi{}, &MockArgoApi{})
	err := syncHandler.Handle()

	if err != nil {
		t.Error(err)
	}

	if len(createdEnv) > 0 {
		t.Errorf("Envs should not be created during NONE mode")
	}

}

func TestSyncWithOneTimeSync(t *testing.T) {

	createdEnv = make([]string, 0)

	store.SetSyncOptions(codefresh.OneTimeSync, []string{})

	syncHandler := GetSyncHandlerInstance(&MockCodefreshApi{}, &MockArgoApi{})

	err := syncHandler.Handle()
	if err != nil {
		t.Error(err)
	}

	if len(createdEnv) != 1 {
		t.Errorf("Single env should be created during OneTimeSync mode")
	}

}

func TestSyncWithSelectSync(t *testing.T) {

	createdEnv = make([]string, 0)

	applications := []string{
		"Test",
	}

	store.SetSyncOptions(codefresh.SelectSync, applications)

	syncHandler := GetSyncHandlerInstance(&MockCodefreshApi{}, &MockArgoApi{})

	err := syncHandler.Handle()
	if err != nil {
		t.Error(err)
	}

	if len(createdEnv) != 1 {
		t.Errorf("Single env should be created during Select sync mode")
	}

}

func TestSyncWithSelectSyncWithNonExistApplication(t *testing.T) {

	createdEnv = make([]string, 0)

	applications := []string{
		"Test2",
	}

	store.SetSyncOptions(codefresh.SelectSync, applications)

	syncHandler := GetSyncHandlerInstance(&MockCodefreshApi{}, &MockArgoApi{})

	err := syncHandler.Handle()
	if err != nil {
		t.Error(err)
	}

	if len(createdEnv) != 0 {
		t.Errorf("Zero envs should be created during Select sync mode")
	}

}
