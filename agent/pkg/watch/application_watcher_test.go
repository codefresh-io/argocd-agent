package watch

import (
	"github.com/codefresh-io/argocd-listener/agent/pkg/infra/queue"
	argoSdk "github.com/codefresh-io/argocd-sdk/pkg/api"
	codefreshSdk "github.com/codefresh-io/go-sdk/pkg/codefresh"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"testing"
)

var _ = func() bool {
	testing.Init()
	return true
}()

var sendResources func(len int) error

type MockCodefreshApi struct {
}

func (api *MockCodefreshApi) GetDefaultGitContext() (error, *codefreshSdk.ContextPayload) {
	panic("implement me")
}

func (api *MockCodefreshApi) DeleteEnvironment(name string) error {
	panic("implement me")
}

func (api *MockCodefreshApi) SendResources(kind string, items interface{}, amount int) error {
	return sendResources(amount)
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
	panic("implement me")
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
	panic("implement me")
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

func TestApplicationWatcherUpdateEvent(t *testing.T) {

	appwatcher := applicationWatcher{
		codefreshApi:    nil,
		itemQueue:       queue.GetInstance(),
		informer:        nil,
		informerFactory: nil,
	}

	syncResult := make(map[string]interface{})
	syncResult["Revision"] = "123"

	operationState := make(map[string]interface{})
	operationState["SyncResult"] = syncResult

	status := make(map[string]interface{})
	status["OperationState"] = operationState
	historyItem := make(map[string]interface{})
	historyItem["Revision"] = "123"
	historyItem["Id"] = 1

	status["History"] = []map[string]interface{}{
		historyItem,
	}

	payload := make(map[string]interface{})
	payload["Status"] = status

	obj := &unstructured.Unstructured{Object: payload}

	appwatcher.update(obj)

	size := queue.GetInstance().Size()

	if size != 1 {
		t.Error("Unable watch update events")
	}

}

func TestApplicationWatcherDeleteEvent(t *testing.T) {

	appwatcher := applicationWatcher{
		codefreshApi:    &MockCodefreshApi{},
		itemQueue:       queue.GetInstance(),
		informer:        nil,
		informerFactory: nil,
		argoApi:         &MockArgoApi{},
	}

	obj := &unstructured.Unstructured{}

	sendResources = func(len int) error {
		if len != 1 {
			t.Error("Unable watch delete events")
		}
		return nil
	}

	appwatcher.delete(obj)

}

func TestApplicationWatcherCreateEvent(t *testing.T) {
	q := queue.GetInstance().New()

	appwatcher := applicationWatcher{
		codefreshApi:    &MockCodefreshApi{},
		itemQueue:       q,
		informer:        nil,
		informerFactory: nil,
		argoApi:         &MockArgoApi{},
	}

	//Status.OperationState.SyncResult.Revision

	syncResult := make(map[string]interface{})
	syncResult["Revision"] = "123"

	operationState := make(map[string]interface{})
	operationState["SyncResult"] = syncResult

	status := make(map[string]interface{})
	status["OperationState"] = operationState
	historyItem := make(map[string]interface{})
	historyItem["Revision"] = "123"
	historyItem["Id"] = 1

	status["History"] = []map[string]interface{}{
		historyItem,
	}

	payload := make(map[string]interface{})
	payload["Status"] = status

	obj := &unstructured.Unstructured{Object: payload}

	sendResources = func(len int) error {
		if len != 1 {
			t.Error("Unable watch create events")
		}
		return nil
	}

	appwatcher.add(obj)

	size := q.Size()

	if size != 1 {
		t.Error("Unable watch update events")
	}

}
