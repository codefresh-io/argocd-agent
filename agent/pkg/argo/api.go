package argo

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	store2 "github.com/codefresh-io/argocd-listener/agent/pkg/store"
	"log"
	"net/http"
)

func buildHttpClient() *http.Client {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	return &http.Client{Transport: tr}
}

func GetToken(username string, password string, host string) string {

	client := buildHttpClient()

	message := map[string]interface{}{
		"username": username,
		"password": password,
	}

	bytesRepresentation, err := json.Marshal(message)
	if err != nil {
		log.Fatalln(err)
	}

	resp, err := client.Post(host+"/api/v1/session", "application/json", bytes.NewBuffer(bytesRepresentation))
	if err != nil {
		log.Fatalln(err)
	}
	var result map[string]interface{}

	json.NewDecoder(resp.Body).Decode(&result)

	defer resp.Body.Close()

	return result["token"].(string)
}

func GetResourceTree(applicationName string) (*ResourceTree, error) {
	token := store2.GetStore().Argo.Token
	host := store2.GetStore().Argo.Host

	client := buildHttpClient()

	req, err := http.NewRequest("GET", host+"/api/v1/applications/"+applicationName+"/resource-tree", nil)
	req.Header.Add("Authorization", "Bearer "+token)
	resp, err := client.Do(req)

	if err != nil {
		log.Fatalln(err)
	}

	var result *ResourceTree

	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&result)

	if err != nil {
		return nil, err
	}

	return result, nil
}

func GetManagedResources(applicationName string) ManagedResource {
	token := store2.GetStore().Argo.Token
	host := store2.GetStore().Argo.Host

	client := buildHttpClient()

	req, err := http.NewRequest("GET", host+"/api/v1/applications/"+applicationName+"/managed-resources", nil)
	req.Header.Add("Authorization", "Bearer "+token)
	resp, err := client.Do(req)

	if err != nil {
		log.Fatalln(err)
	}

	var result ManagedResource

	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&result)

	if err != nil {
		//return nil, err
	}

	return result
}

func GetProjects() []ProjectItem {
	token := store2.GetStore().Argo.Token
	host := store2.GetStore().Argo.Host

	client := buildHttpClient()

	req, err := http.NewRequest("GET", host+"/api/v1/projects", nil)
	req.Header.Add("Authorization", "Bearer "+token)
	resp, err := client.Do(req)

	if err != nil {
		log.Fatalln(err)
	}

	var result Project

	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&result)

	if err != nil {
		//return nil, err
	}

	return result.Items
}

func GetApplications() []ApplicationItem {
	token := store2.GetStore().Argo.Token
	host := store2.GetStore().Argo.Host

	client := buildHttpClient()

	req, err := http.NewRequest("GET", host+"/api/v1/applications", nil)
	req.Header.Add("Authorization", "Bearer "+token)
	resp, err := client.Do(req)

	if err != nil {
		log.Fatalln(err)
	}

	var result Application

	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&result)

	if err != nil {
		//return nil, err
	}

	return result.Items
}
