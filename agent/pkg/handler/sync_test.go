package handler

import (
	"github.com/codefresh-io/argocd-listener/agent/pkg/codefresh"
	"github.com/codefresh-io/argocd-listener/agent/pkg/store"
	argoSdk "github.com/codefresh-io/argocd-sdk/pkg/api"
	"testing"
)

var createdEnv []string

type MockArgoApi struct {
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
	panic("implement me")
}

type MockCodefreshApi struct {
}

func (api *MockCodefreshApi) CreateEnvironment(name string, project string, application string) error {
	createdEnv = append(createdEnv, name)
	return nil
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
