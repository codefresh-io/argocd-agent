package argo

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	store2 "github.com/codefresh-io/argocd-listener/agent/pkg/store"
	"net/http"
)

type ArgoApi interface {
	GetApplicationsWithCredentialsFromStorage() ([]ApplicationItem, error)
	GetResourceTree(applicationName string) (*ResourceTree, error)
	GetResourceTreeAll(applicationName string) (interface{}, error)
	GetManagedResources(applicationName string) (*ManagedResource, error)
}

type Api struct {
	Token string
	Host  string
}

var api *Api

func GetInstance() *Api {
	if api != nil {
		return api
	}

	argoConfig := store2.GetStore().Argo
	api = &Api{
		Token: argoConfig.Token,
		Host:  argoConfig.Host,
	}
	return api
}

func buildHttpClient() *http.Client {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	return &http.Client{Transport: tr}
}

func GetToken(username string, password string, host string) (string, error) {

	client := buildHttpClient()

	message := map[string]interface{}{
		"username": username,
		"password": password,
	}

	bytesRepresentation, err := json.Marshal(message)
	if err != nil {
		return "", errors.New("application error, cant retrieve argo token")
	}

	resp, err := client.Post(host+"/api/v1/session", "application/json", bytes.NewBuffer(bytesRepresentation))
	if err != nil {
		return "", err
	}

	if resp.StatusCode == 401 {
		return "", errors.New("cant retrieve argocd token, permission denied")
	}

	var result map[string]interface{}

	err = json.NewDecoder(resp.Body).Decode(&result)

	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	return result["token"].(string), nil
}

func (api *Api) CheckToken() error {
	client := buildHttpClient()
	req, err := http.NewRequest("GET", api.Host+"/api/v1/account", nil)

	if err != nil {
		return err
	}

	req.Header.Add("Authorization", "Bearer "+api.Token)
	resp, err := client.Do(req)

	if err != nil {
		return err
	}

	var result map[string]interface{}

	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&result)

	if err != nil {
		return err
	}

	return nil
}

func (api *Api) GetResourceTree(applicationName string) (*ResourceTree, error) {
	client := buildHttpClient()

	req, err := http.NewRequest("GET", api.Host+"/api/v1/applications/"+applicationName+"/resource-tree", nil)

	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", "Bearer "+api.Token)
	resp, err := client.Do(req)

	if err != nil {
		return nil, err
	}

	var result *ResourceTree

	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&result)

	if err != nil {
		return nil, err
	}

	return result, nil
}

//  TODO: refactor
func (api *Api) GetResourceTreeAll(applicationName string) (interface{}, error) {
	client := buildHttpClient()

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/v1/applications/%s/resource-tree", api.Host, applicationName), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", api.Token))
	resp, err := client.Do(req)

	if err != nil {
		return nil, err
	}

	var result interface{}

	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&result)

	if err != nil {
		return nil, err
	}

	return result.(map[string]interface{})["nodes"], nil
}

func (api *Api) GetManagedResources(applicationName string) (*ManagedResource, error) {
	token := store2.GetStore().Argo.Token
	host := store2.GetStore().Argo.Host

	client := buildHttpClient()

	req, err := http.NewRequest("GET", host+"/api/v1/applications/"+applicationName+"/managed-resources", nil)
	req.Header.Add("Authorization", "Bearer "+token)
	resp, err := client.Do(req)

	if err != nil {
		return nil, err
	}

	var result ManagedResource

	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&result)

	if err != nil {
		return nil, err
	}

	return &result, nil
}

func GetProjects(token string, host string) ([]ProjectItem, error) {
	client := buildHttpClient()

	req, err := http.NewRequest("GET", host+"/api/v1/projects", nil)
	req.Header.Add("Authorization", "Bearer "+token)
	resp, err := client.Do(req)

	if err != nil {
		return nil, err
	}

	var result Project

	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&result)

	if err != nil {
		return nil, err
	}

	return result.Items, nil
}

func GetProjectsWithCredentialsFromStorage() ([]ProjectItem, error) {
	token := store2.GetStore().Argo.Token
	host := store2.GetStore().Argo.Host

	return GetProjects(token, host)
}

func GetApplication(application string) (map[string]interface{}, error) {
	token := store2.GetStore().Argo.Token
	host := store2.GetStore().Argo.Host

	client := buildHttpClient()

	var result map[string]interface{}

	req, err := http.NewRequest("GET", host+"/api/v1/applications/"+application, nil)
	req.Header.Add("Authorization", "Bearer "+token)
	resp, err := client.Do(req)

	if resp.StatusCode != 200 {
		// TODO: add error handling and move it to common place
		return nil, errors.New(fmt.Sprintf("Failed to retrieve application, reason %v", resp.Status))
	}

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&result)

	if err != nil {
		return nil, err
	}

	return result, nil
}

func (api *Api) GetApplicationsWithCredentialsFromStorage() ([]ApplicationItem, error) {
	return GetApplications(api.Token, api.Host)
}

func GetApplications(token string, host string) ([]ApplicationItem, error) {

	client := buildHttpClient()

	req, err := http.NewRequest("GET", host+"/api/v1/applications", nil)
	req.Header.Add("Authorization", "Bearer "+token)
	resp, err := client.Do(req)

	if err != nil {
		return nil, err
	}

	var result Application

	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&result)

	if err != nil {
		return nil, err
	}

	return result.Items, nil
}
