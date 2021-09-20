package codefresh

import (
	codefreshSdk "github.com/codefresh-io/go-sdk/pkg/codefresh"
	"testing"
)

var _ = func() bool {
	testing.Init()
	return true
}()

type _gitops struct {
}

func (gitops *_gitops) CreateEnvironment(name string, project string, application string, integration string, namespace string) error {
	return nil
}
func (gitops *_gitops) SendEnvironment(environment codefreshSdk.Environment) (map[string]interface{}, error) {
	return nil, nil
}

func (gitops *_gitops) DeleteEnvironment(name string) error {
	return nil
}

func (gitops *_gitops) GetEnvironments() ([]codefreshSdk.CFEnvironment, error) {
	return nil, nil
}

func (gitops *_gitops) SendEvent(name string, props map[string]string) error {
	return nil
}

func (gitops *_gitops) SendApplicationResources(resources *codefreshSdk.ApplicationResources) error {
	return nil
}

type _argo struct {
}

func (a *_argo) CreateIntegration(integration codefreshSdk.IntegrationPayloadData) error {
	return nil
}

func (a *_argo) UpdateIntegration(name string, integration codefreshSdk.IntegrationPayloadData) error {
	return nil
}

func (a *_argo) GetIntegrations() ([]*codefreshSdk.IntegrationPayload, error) {
	return nil, nil
}

func (a *_argo) GetIntegrationByName(name string) (*codefreshSdk.IntegrationPayload, error) {
	return nil, nil
}

func (a *_argo) DeleteIntegrationByName(name string) error {
	return nil
}

func (a *_argo) HeartBeat(error string, version string, integration string) error {
	return nil
}

func (a *_argo) SendResources(kind string, items interface{}, amount int, integration string) error {
	return nil
}

type context struct {
}

func (ctx *context) GetGitContexts() (error, *[]codefreshSdk.ContextPayload) {
	return nil, nil
}

func (ctx *context) GetGitContextByName(name string) (error, *codefreshSdk.ContextPayload) {
	return nil, nil
}

func (ctx *context) GetDefaultGitContext() (error, *codefreshSdk.ContextPayload) {
	return nil, nil
}

type codefresh struct {
}

func (c *codefresh) Pipelines() codefreshSdk.IPipelineAPI {
	return nil
}

func (c *codefresh) Tokens() codefreshSdk.ITokenAPI {
	return nil
}

func (c *codefresh) RuntimeEnvironments() codefreshSdk.IRuntimeEnvironmentAPI {
	return nil
}

func (c *codefresh) Workflows() codefreshSdk.IWorkflowAPI {
	return nil
}

func (c *codefresh) Progresses() codefreshSdk.IProgressAPI {
	return nil
}

func (c *codefresh) Clusters() codefreshSdk.IClusterAPI {
	return nil
}

func (c *codefresh) Contexts() codefreshSdk.IContextAPI {
	return &context{}
}

func (c *codefresh) Argo() codefreshSdk.ArgoAPI {
	return &_argo{}
}

func (c *codefresh) Gitops() codefreshSdk.GitopsAPI {
	return &_gitops{}
}

func TestApi_GetDefaultGitContext(t *testing.T) {
	api := Api{
		codefreshApi: &codefresh{},
		Integration:  "Test",
	}

	err, _ := api.GetDefaultGitContext()
	if err != nil {
		t.Error("Should be able retrieve context")
	}
}

func TestApi_GetGitContexts(t *testing.T) {
	api := Api{
		codefreshApi: &codefresh{},
		Integration:  "Test",
	}

	err, _ := api.GetGitContexts()
	if err != nil {
		t.Error("Should be able retrieve contexts")
	}
}

func TestApi_SendResources(t *testing.T) {
	api := Api{
		codefreshApi: &codefresh{},
		Integration:  "Test",
	}
	payload := make(map[string]interface{}, 0)

	err := api.SendResources("test", payload, 0)
	if err != nil {
		t.Error("Should be able send resources without error")
	}
}

func TestApi_SendEvent(t *testing.T) {
	api := Api{
		codefreshApi: &codefresh{},
		Integration:  "Test",
	}
	payload := make(map[string]string, 0)

	err := api.SendEvent("test", payload)
	if err != nil {
		t.Error("Should be able send event without error")
	}
}

func TestApi_HeartBeat(t *testing.T) {
	api := Api{
		codefreshApi: &codefresh{},
		Integration:  "Test",
	}
	err := api.HeartBeat("test")
	if err != nil {
		t.Error("Should be able send heartbeat without error")
	}
}

func TestApi_GetEnvironments(t *testing.T) {
	api := Api{
		codefreshApi: &codefresh{},
		Integration:  "Test",
	}
	_, err := api.GetEnvironments()
	if err != nil {
		t.Error("Should be able get environments without error")
	}
}

func TestApi_CreateIntegration(t *testing.T) {
	api := Api{
		codefreshApi: &codefresh{},
		Integration:  "Test",
	}
	err := api.CreateIntegration("test", "host", "un", "pass", "tk", "v1", "gitops", "cluster")
	if err != nil {
		t.Error("Should be able create integration without error")
	}
}

func TestApi_UpdateIntegration(t *testing.T) {
	api := Api{
		codefreshApi: &codefresh{},
		Integration:  "Test",
	}
	err := api.UpdateIntegration("test", "host", "un", "pass", "tk", "v1", "gitops", "cluster")
	if err != nil {
		t.Error("Should be able update integration without error")
	}
}

func TestApi_SendEnvironment(t *testing.T) {
	api := Api{
		codefreshApi: &codefresh{},
		Integration:  "Test",
	}
	_, err := api.SendEnvironment(codefreshSdk.Environment{
		Gitops:       codefreshSdk.Gitops{},
		FinishedAt:   "",
		HealthStatus: "",
		SyncStatus:   "",
		HistoryId:    0,
		SyncRevision: "",
		Name:         "",
		Activities:   nil,
		Resources:    nil,
		RepoUrl:      "",
		Commit:       codefreshSdk.Commit{},
		SyncPolicy:   codefreshSdk.SyncPolicy{},
		Date:         "",
		ParentApp:    "",
	})
	if err != nil {
		t.Error("Should be able send environment without error")
	}
}

func TestApi_CreateEnvironment(t *testing.T) {
	api := Api{
		codefreshApi: &codefresh{},
		Integration:  "Test",
	}
	err := api.CreateEnvironment("test", "test", "test", "test")
	if err != nil {
		t.Error("Should be able create environment without error")
	}
}

func TestApi_DeleteEnvironment(t *testing.T) {
	api := Api{
		codefreshApi: &codefresh{},
		Integration:  "Test",
	}
	err := api.DeleteEnvironment("test")
	if err != nil {
		t.Error("Should be able delete environment without error")
	}
}

func TestApi_SendApplicationResources(t *testing.T) {
	api := Api{
		codefreshApi: &codefresh{},
		Integration:  "Test",
	}
	err := api.SendApplicationResources(nil)
	if err != nil {
		t.Error("Should be able send application resources without error")
	}
}
