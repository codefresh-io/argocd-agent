package argo

import (
	"crypto/tls"
	"github.com/codefresh-io/argocd-listener/agent/pkg/store"
	argoSdk "github.com/codefresh-io/argocd-sdk/pkg/api"
	"net/http"
)

type ArgoApi interface {
	GetApplicationsWithCredentialsFromStorage() ([]argoSdk.ApplicationItem, error)
	GetResourceTree(applicationName string) (*argoSdk.ResourceTree, error)
	GetResourceTreeAll(applicationName string) (interface{}, error)
	GetManagedResources(applicationName string) (*argoSdk.ManagedResource, error)
	GetVersion() (string, error)
	GetProjectsWithCredentialsFromStorage() ([]argoSdk.ProjectItem, error)
	GetApplication(application string) (map[string]interface{}, error)
}

type Api struct {
	sdk argoSdk.Argo
}

type UnauthorizedApi struct {
}

var api *Api
var unauthorizedApi *UnauthorizedApi

func GetInstance() *Api {
	if api != nil {
		return api
	}

	argoConfig := store.GetStore().Argo
	api = &Api{
		sdk: BuildArgoSdk(argoConfig.Token, argoConfig.Host),
	}
	return api
}

func GetUnauthorizedApiInstance() *UnauthorizedApi {
	if unauthorizedApi != nil {
		return unauthorizedApi
	}

	unauthorizedApi = &UnauthorizedApi{}
	return unauthorizedApi
}

func BuildArgoSdk(token string, host string) argoSdk.Argo {
	return argoSdk.New(&argoSdk.ClientOptions{
		Auth: argoSdk.AuthOptions{
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

func GetToken(username string, password string, host string) (string, error) {
	return argoSdk.GetToken(username, password, host)
}

func (api *Api) CheckToken() error {
	return api.sdk.Auth().CheckToken()
}

func (api *Api) GetResourceTree(applicationName string) (*argoSdk.ResourceTree, error) {
	return api.sdk.Application().GetResourceTree(applicationName)
}

func (api *Api) GetResourceTreeAll(applicationName string) (interface{}, error) {
	return api.sdk.Application().GetResourceTreeAll(applicationName)
}

func (api *Api) GetVersion() (string, error) {
	return api.sdk.Version().GetVersion()
}

func (api *Api) GetManagedResources(applicationName string) (*argoSdk.ManagedResource, error) {
	return api.sdk.Application().GetManagedResources(applicationName)
}

func (api *Api) GetProjectsWithCredentialsFromStorage() ([]argoSdk.ProjectItem, error) {
	token := store.GetStore().Argo.Token
	host := store.GetStore().Argo.Host
	sdk := BuildArgoSdk(token, host)
	return sdk.Project().GetProjects()
}

func (api *Api) GetApplication(application string) (map[string]interface{}, error) {
	return api.sdk.Application().GetApplication(application)
}

func (api *Api) GetApplicationsWithCredentialsFromStorage() ([]argoSdk.ApplicationItem, error) {
	token := store.GetStore().Argo.Token
	host := store.GetStore().Argo.Host
	sdk := BuildArgoSdk(token, host)
	return sdk.Application().GetApplications()
}

func (api *UnauthorizedApi) GetApplications(token string, host string) ([]argoSdk.ApplicationItem, error) {
	sdk := BuildArgoSdk(token, host)
	return sdk.Application().GetApplications()
}
