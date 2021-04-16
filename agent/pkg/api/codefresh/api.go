package codefresh

import (
	"crypto/tls"
	"github.com/codefresh-io/argocd-listener/agent/pkg/api/argo"
	"github.com/codefresh-io/argocd-listener/agent/pkg/infra/logger"
	"github.com/codefresh-io/argocd-listener/agent/pkg/infra/store"
	codefreshSdk "github.com/codefresh-io/go-sdk/pkg/codefresh"
	"net/http"
)

type Api struct {
	codefreshApi codefreshSdk.Codefresh
	Integration  string
}

type CodefreshApi interface {
	CreateEnvironment(name string, project string, application string) error
	GetDefaultGitContext() (error, *codefreshSdk.ContextPayload)
	DeleteEnvironment(name string) error
	SendResources(kind string, items interface{}, amount int) error
	SendEvent(name string, props map[string]string) error
	HeartBeat(error string) error
	GetEnvironments() ([]codefreshSdk.CFEnvironment, error)
	CreateIntegration(name string, host string, username string, password string, token string, serverVersion string, provider string, clusterName string) error
	UpdateIntegration(name string, host string, username string, password string, token string, serverVersion string, provider string, clusterName string) error
	SendEnvironment(environment codefreshSdk.Environment) (map[string]interface{}, error)
	SendApplicationResources(resources *codefreshSdk.ApplicationResources) error
}

var api *Api

func GetInstance() *Api {
	if api != nil {
		return api
	}

	codefreshConfig := store.GetStore().Codefresh
	api = &Api{
		codefreshApi: BuildCodefreshSdk(codefreshConfig.Token, codefreshConfig.Host),
		Integration:  codefreshConfig.Integration,
	}

	return api
}

func BuildCodefreshSdk(token string, host string) codefreshSdk.Codefresh {
	return codefreshSdk.New(&codefreshSdk.ClientOptions{
		Auth: codefreshSdk.AuthOptions{
			Token: token,
		},
		Debug:  false,
		Host:   host,
		Client: buildHttpClient(),
	})
}

func buildHttpClient() *http.Client {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	return &http.Client{Transport: tr}
}

func (a *Api) GetDefaultGitContext() (error, *codefreshSdk.ContextPayload) {
	return a.codefreshApi.Contexts().GetDefaultGitContext()
}

func (a *Api) GetGitContexts() (error, *[]codefreshSdk.ContextPayload) {
	return a.codefreshApi.Contexts().GetGitContexts()
}

func (a *Api) GetGitContextByName(name string) (error, *codefreshSdk.ContextPayload) {
	return a.codefreshApi.Contexts().GetGitContextByName(name)
}

func (a *Api) SendResources(kind string, items interface{}, amount int) error {

	logger.GetLogger().Infof("Trying sent resources with type: \"%s\" to codefresh, amount: \"%v\"", kind, amount)

	err := a.codefreshApi.Argo().SendResources(kind, items, amount, a.Integration)

	if err != nil {
		return err
	}

	logger.GetLogger().Infof("Successfully sent type: \"%s\" to codefresh", kind)

	return nil
}

func (a *Api) SendEvent(name string, props map[string]string) error {
	return a.codefreshApi.Gitops().SendEvent(name, props)
}

func (a *Api) HeartBeat(error string) error {
	agentConfig := store.GetStore().Agent
	return a.codefreshApi.Argo().HeartBeat(error, agentConfig.Version, a.Integration)
}

func (a *Api) GetEnvironments() ([]codefreshSdk.CFEnvironment, error) {
	return a.codefreshApi.Gitops().GetEnvironments()
}

func (a *Api) CreateIntegration(name string, host string, username string, password string, token string, serverVersion string, provider string, clusterName string) error {
	payloadData := prepareIntegration(name, host, username, password, token, serverVersion, provider, clusterName)

	return a.codefreshApi.Argo().CreateIntegration(payloadData)
}

func (a *Api) UpdateIntegration(name string, host string, username string, password string, token string, serverVersion string, provider string, clusterName string) error {
	return a.codefreshApi.Argo().UpdateIntegration(name, prepareIntegration(name, host, username, password, token, serverVersion, provider, clusterName))
}

func (a *Api) SendEnvironment(environment codefreshSdk.Environment) (map[string]interface{}, error) {

	logger.GetLogger().Infof("Successfully sent environment \"%v\" update to codefresh, services count %v", environment.Name, len(environment.Activities))

	return a.codefreshApi.Gitops().SendEnvironment(environment)
}

func (a *Api) CreateEnvironment(name string, project string, application string) error {
	err := a.codefreshApi.Gitops().CreateEnvironment(name, project, application, a.Integration)
	if err != nil {
		return err
	}

	logger.GetLogger().Infof("Successfully create gitops application with name %s and application %s", name, application)

	return nil
}

func (a *Api) DeleteEnvironment(name string) error {
	return a.codefreshApi.Gitops().DeleteEnvironment(name)
}

func (a *Api) SendApplicationResources(resources *codefreshSdk.ApplicationResources) error {
	return a.codefreshApi.Gitops().SendApplicationResources(resources)
}

func prepareIntegration(name string, host string, username string, password string, token string, serverVersion string, provider string, clusterName string) codefreshSdk.IntegrationPayloadData {
	payloadData := codefreshSdk.IntegrationPayloadData{
		Name: name,
		Url:  host,
	}

	if token == "" {
		token = store.GetStore().Argo.Token
	}

	argoApi := argo.GetInstance()

	applications, _ := argoApi.GetApplications()
	clusters, _ := argoApi.GetClusters()
	repositories, _ := argoApi.GetRepositories()

	payloadData.Clusters = codefreshSdk.IntegrationItem{Amount: len(applications)}
	payloadData.Applications = codefreshSdk.IntegrationItem{Amount: len(clusters)}
	payloadData.Repositories = codefreshSdk.IntegrationItem{Amount: len(repositories)}

	if username != "" {
		payloadData.Username = &username
	}

	if password != "" {
		payloadData.Password = &password
	}

	if token != "" {
		payloadData.Token = &token
	}

	if serverVersion != "" {
		payloadData.ServerVersion = &serverVersion
	}

	if clusterName != "" {
		payloadData.ClusterName = &clusterName
	}

	if provider != "" {
		payloadData.Provider = &provider
	}

	return payloadData
}
