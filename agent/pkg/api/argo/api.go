package argo

import (
	"crypto/tls"
	"github.com/codefresh-io/argocd-listener/agent/pkg/infra/store"
	argoSdk "github.com/codefresh-io/argocd-sdk/pkg/api"
	"net/http"
)

// ArgoAPI responsible for proxy calls for argosdk that implement argo api
type ArgoAPI interface {
	GetApplicationsWithCredentialsFromStorage() ([]argoSdk.ApplicationItem, error)
	GetResourceTree(applicationName string) (*argoSdk.ResourceTree, error)
	GetResourceTreeAll(applicationName string) (interface{}, error)
	GetManagedResources(applicationName string) (*argoSdk.ManagedResource, error)
	GetProjectsWithCredentialsFromStorage() ([]argoSdk.ProjectItem, error)
	GetApplication(application string) (map[string]interface{}, error)
	CheckToken() error
	GetClusters() ([]argoSdk.ClusterItem, error)
	GetApplications() ([]argoSdk.ApplicationItem, error)
	GetRepositories() ([]argoSdk.RepositoryItem, error)
	CreateDefaultApp() error
}

type argoAPI struct {
	sdk argoSdk.Argo
}

type unauthorizedApi struct {
}

type UnauthorizedApi interface {
	GetApplications(token string, host string) ([]argoSdk.ApplicationItem, error)
	GetToken(username string, password string, host string) (string, error)
	GetVersion(host string) (string, error)
}

var api *argoAPI
var unauthorizedArgoApi *unauthorizedApi

// GetInstance build and provide as singleton new instance of ArgoAPI interface
func GetInstance() ArgoAPI {
	if api != nil {
		return api
	}

	argoConfig := store.GetStore().Argo
	api = &argoAPI{
		sdk: buildArgoSdk(argoConfig.Token, argoConfig.Host),
	}
	return api
}

// We need to change api instance after after regeneration argocd token
func ResetInstance() {
	argoConfig := store.GetStore().Argo
	*api = argoAPI{
		sdk: buildArgoSdk(argoConfig.Token, argoConfig.Host),
	}
}

// GetUnauthorizedApiInstance build and provide singleton for unathorized argo api
func GetUnauthorizedApiInstance() UnauthorizedApi {
	if unauthorizedArgoApi != nil {
		return unauthorizedArgoApi
	}

	unauthorizedArgoApi = &unauthorizedApi{}
	return unauthorizedArgoApi
}

func buildArgoSdk(token string, host string) argoSdk.Argo {
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

// GetToken retrieve argocd token use basic auth in before use ArgoAPI interface
func (api *unauthorizedApi) GetToken(username string, password string, host string) (string, error) {
	return argoSdk.GetToken(username, password, host)
}

// CheckToken validate argocd token
func (api *argoAPI) CheckToken() error {
	return api.sdk.Auth().CheckToken()
}

// GetResourceTree retrieve argo application resources tree, include deployment, services and so on
func (api *argoAPI) GetResourceTree(applicationName string) (*argoSdk.ResourceTree, error) {
	return api.sdk.Application().GetResourceTree(applicationName)
}

// GetResourceTreeAll deprecated , should be user GetResourceTree instead with retrieval only nodes field
func (api *argoAPI) GetResourceTreeAll(applicationName string) (interface{}, error) {
	return api.sdk.Application().GetResourceTreeAll(applicationName)
}

// GetVersion get argocd server version
func (api *unauthorizedApi) GetVersion(host string) (string, error) {
	sdk := buildArgoSdk("", host)
	return sdk.Version().GetVersion()
}

// GetManagedResources
func (api *argoAPI) GetManagedResources(applicationName string) (*argoSdk.ManagedResource, error) {
	return api.sdk.Application().GetManagedResources(applicationName)
}

// GetClusters get argocd connected clusters
func (api *argoAPI) GetClusters() ([]argoSdk.ClusterItem, error) {
	return api.sdk.Clusters().GetClusters()
}

// GetApplications get argocd applications
func (api *argoAPI) GetApplications() ([]argoSdk.ApplicationItem, error) {
	return api.sdk.Application().GetApplications()
}

// GetRepositories get argocd connected repositories
func (api *argoAPI) GetRepositories() ([]argoSdk.RepositoryItem, error) {
	return api.sdk.Repository().GetRepositories()
}

// GetProjectsWithCredentialsFromStorage retrieve projects use credentials from storage that we init during startup
func (api *argoAPI) GetProjectsWithCredentialsFromStorage() ([]argoSdk.ProjectItem, error) {
	token := store.GetStore().Argo.Token
	host := store.GetStore().Argo.Host
	sdk := buildArgoSdk(token, host)
	return sdk.Project().GetProjects()
}

// GetApplication get detailed application information
func (api *argoAPI) GetApplication(application string) (map[string]interface{}, error) {
	return api.sdk.Application().GetApplication(application)
}

func (api *argoAPI) CreateDefaultApp() error {
	var requestOptions argoSdk.CreateApplicationOpt
	requestOptions.Metadata.Name = "cf-guestbook"
	requestOptions.Spec.Project = "default"
	requestOptions.Spec.Destination.Name = ""
	requestOptions.Spec.Destination.Namespace = ""
	requestOptions.Spec.Destination.Server = "https://kubernetes.default.svc"
	requestOptions.Spec.Source.RepoURL = "https://github.com/argoproj/argocd-example-apps.git"
	requestOptions.Spec.Source.Path = "guestbook"
	requestOptions.Spec.Source.TargetRevision = "HEAD"
	return api.sdk.Application().CreateApplication(requestOptions)
}

// GetApplicationsWithCredentialsFromStorage get detailed application information use credentials from storage that we init during startup
func (api *argoAPI) GetApplicationsWithCredentialsFromStorage() ([]argoSdk.ApplicationItem, error) {
	token := store.GetStore().Argo.Token
	host := store.GetStore().Argo.Host
	sdk := buildArgoSdk(token, host)
	return sdk.Application().GetApplications()
}

// GetApplications get applications with token as param, without init API interface
func (api *unauthorizedApi) GetApplications(token string, host string) ([]argoSdk.ApplicationItem, error) {
	sdk := buildArgoSdk(token, host)
	return sdk.Application().GetApplications()
}
