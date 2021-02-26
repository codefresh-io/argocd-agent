package transform

import (
	argoSdk "github.com/codefresh-io/argocd-sdk/pkg/api"
	"testing"
)

type MockArgoApi struct {
}

func (m MockArgoApi) GetApplicationsWithCredentialsFromStorage() ([]argoSdk.ApplicationItem, error) {
	panic("implement me")
}

func (m MockArgoApi) GetVersion() (string, error) {
	panic("implement me")
}

func (m MockArgoApi) GetProjectsWithCredentialsFromStorage() ([]argoSdk.ProjectItem, error) {
	panic("implement me")
}

func (m MockArgoApi) GetApplication(application string) (map[string]interface{}, error) {
	panic("implement me")
}

func (m MockArgoApi) GetResourceTree(applicationName string) (*argoSdk.ResourceTree, error) {
	var nodes = make([]argoSdk.Node, 0)
	nodes = append(nodes, argoSdk.Node{
		Kind: "Deploy",
		Uid:  "Uid",
		Health: argoSdk.Health{
			Status: "Health",
		},
	})

	nodes = append(nodes, argoSdk.Node{
		Kind: "Deploy",
		Uid:  "Uid2",
		Health: argoSdk.Health{
			Status: "Unhealth",
		},
	})

	return &argoSdk.ResourceTree{
		Nodes: nodes,
	}, nil
}

func (m MockArgoApi) GetResourceTreeAll(applicationName string) (interface{}, error) {
	panic("implement me")
}

func (m MockArgoApi) GetManagedResources(applicationName string) (*argoSdk.ManagedResource, error) {
	liveState := "{\"kind\":\"Service\",\"metadata\":{ \"name\":\"test-api\",\"namespace\":\"andrii\",\"uid\":\"46263671-f290-11ea-8d49-42010a8001b0\"},\"spec\":{ \"template\": { \"spec\": { \"containers\":[{\"image\":\"andriicodefresh/test:v7\",\"name\":\"test-api\"}] } }, \"clusterIP\":\"10.27.251.224\",\"ports\":[{\"port\":80,\"protocol\":\"TCP\",\"targetPort\":1700}]}}"

	var resourceItems = make([]argoSdk.ManagedResourceItem, 0)
	resourceItems = append(resourceItems, argoSdk.ManagedResourceItem{
		Kind:        "Deployment",
		TargetState: "",
		LiveState:   liveState,
		Name:        "Test",
	})

	return &argoSdk.ManagedResource{
		Items: resourceItems,
	}, nil
}

func TestPrepareEnvironment(t *testing.T) {

	envTransformer := GetEnvTransformerInstance(MockArgoApi{})

	services, err := envTransformer.prepareEnvironmentActivity("test")
	if err != nil {
		t.Error(err)
	}

	if len(services) != 1 {
		t.Errorf("We should prepare 1 services for send to codefresh")
	}

}
