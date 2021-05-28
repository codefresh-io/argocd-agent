package acceptance

import (
	"github.com/codefresh-io/argocd-listener/agent/pkg/infra/store"
	"github.com/codefresh-io/argocd-listener/installer/pkg/dictionary"
	"github.com/codefresh-io/argocd-listener/installer/pkg/install/entity"
	argoSdk "github.com/codefresh-io/argocd-sdk/pkg/api"
	"testing"
)

var _ = func() bool {
	testing.Init()
	return true
}()

type ArgoCredsMockArgoApi struct {
}

func (api *ArgoCredsMockArgoApi) CreateDefaultApp() error {
	return nil
}

func (api *ArgoCredsMockArgoApi) CheckToken() error {
	panic("implement me")
}

func (api *ArgoCredsMockArgoApi) GetResourceTree(applicationName string) (*argoSdk.ResourceTree, error) {
	panic("implement me")
}

func (api *ArgoCredsMockArgoApi) GetResourceTreeAll(applicationName string) (interface{}, error) {
	panic("implement me")
}

func (api *ArgoCredsMockArgoApi) GetManagedResources(applicationName string) (*argoSdk.ManagedResource, error) {
	panic("implement me")
}

func (api *ArgoCredsMockArgoApi) GetVersion() (string, error) {
	panic("implement me")
}

func (api *ArgoCredsMockArgoApi) GetApplicationsWithCredentialsFromStorage() ([]argoSdk.ApplicationItem, error) {
	panic("implement me")
}

func (api *ArgoCredsMockArgoApi) GetProjectsWithCredentialsFromStorage() ([]argoSdk.ProjectItem, error) {
	return []argoSdk.ProjectItem{}, nil
}

func (api *ArgoCredsMockArgoApi) GetApplication(application string) (map[string]interface{}, error) {
	panic("implement me")
}

func (api *ArgoCredsMockArgoApi) GetClusters() ([]argoSdk.ClusterItem, error) {
	panic("implement me")
}

func (api *ArgoCredsMockArgoApi) GetApplications() ([]argoSdk.ApplicationItem, error) {
	panic("implement me")
}

func (api *ArgoCredsMockArgoApi) GetRepositories() ([]argoSdk.RepositoryItem, error) {
	panic("implement me")
}

type MockUnathourizedArgoApi struct {
}

func (api *MockUnathourizedArgoApi) GetApplications(token string, host string) ([]argoSdk.ApplicationItem, error) {
	return nil, nil
}

func (api *MockUnathourizedArgoApi) GetToken(username string, password string, host string) (string, error) {
	return "token", nil
}

func TestArgoCredsFailure(t *testing.T) {
	test := &ArgoCredentialsAcceptanceTest{}

	result := test.failure()
	if !result {
		t.Error("Should fail with error")
	}
}

func TestArgoCredsGetMessage(t *testing.T) {
	test := &ArgoCredentialsAcceptanceTest{}

	result := test.getMessage()

	if result != dictionary.CheckArgoCredentials {
		t.Error("Message is incorrect")
	}
}

func TestArgoCredsCheckWithoutToken(t *testing.T) {
	test := &ArgoCredentialsAcceptanceTest{
		argoApi:            &MockArgoApi{},
		unathorizedArgoApi: &MockUnathourizedArgoApi{},
	}

	err := test.check(&entity.ArgoOptions{})

	if err != nil {
		t.Error("Should be executed without an error")
	}

	tk := store.GetStore().Argo.Token
	if tk != "token" {
		t.Error("Wrong token")
	}

}

func TestArgoCredsCheckWithToken(t *testing.T) {
	test := &ArgoCredentialsAcceptanceTest{
		argoApi:            &MockArgoApi{},
		unathorizedArgoApi: &MockUnathourizedArgoApi{},
	}

	err := test.check(&entity.ArgoOptions{
		Token: "test",
	})

	if err != nil {
		t.Error("Should be executed without an error")
	}

}
