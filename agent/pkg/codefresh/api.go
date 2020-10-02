package codefresh

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/codefresh-io/argocd-listener/agent/pkg/logger"
	"github.com/codefresh-io/argocd-listener/agent/pkg/store"
	"github.com/guregu/null"
	"net/http"
	"strings"
)

type Api struct {
	Token       string
	Host        string
	Integration string
}

var api *Api

func GetInstance() *Api {
	if api != nil {
		return api
	}

	codefreshConfig := store.GetStore().Codefresh
	api = &Api{
		Token:       codefreshConfig.Token,
		Host:        codefreshConfig.Host,
		Integration: codefreshConfig.Integration,
	}
	return api
}

func (a *Api) GetDefaultGitContext() (error, *ContextPayload) {
	var result ContextPayload

	err := a.requestAPI(&requestOptions{
		method: "GET",
		path:   "/contexts/git/default",
	}, &result)
	if err != nil {
		return err, nil
	}
	// todo - add logs
	return nil, &result
}

func (a *Api) SendEnvironment(environment Environment) (map[string]interface{}, error) {

	logger.GetLogger().Infof("Successfully sent environment \"%v\" update to codefresh, services count %v", environment.Name, len(environment.Activities))

	var result map[string]interface{}
	err := a.requestAPI(&requestOptions{method: "POST", path: "/environments-v2/argo/events", body: environment}, &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (a *Api) SendResources(kind string, items interface{}) error {
	if items == nil {
		return nil
	}

	err := a.requestAPI(&requestOptions{
		method: "POST",
		path:   fmt.Sprintf("/argo-agent/%s", a.Integration),
		body:   &AgentState{Kind: kind, Items: items},
	}, nil)
	if err != nil {
		return err
	}

	logger.GetLogger().Infof("Successfully sent type: \"%s\" to codefresh", kind)

	return nil
}

func (a *Api) SendEvent(name string, props map[string]string) error {
	event := CodefreshEvent{Event: name, Props: props}

	err := a.requestAPI(&requestOptions{
		method: "POST",
		path:   "/gitops/system/events",
		body:   event,
	}, nil)
	if err != nil {
		return err
	}

	return nil
}

func (a *Api) HeartBeat(error string) error {
	var body interface{}

	if error != "" {
		body = Heartbeat{
			Error: error,
		}
	}

	err := a.requestAPI(&requestOptions{
		method: "POST",
		path:   fmt.Sprintf("/argo-agent/%s/heartbeat", a.Integration),
		body:   body,
	}, nil)
	if err != nil {
		return err
	}

	return nil
}

func (a *Api) GetEnvironments() ([]CFEnvironment, error) {
	var result MongoCFEnvWrapper
	err := a.requestAPI(&requestOptions{
		method: "GET",
		path:   "/environments-v2?plain=true&isEnvironment=false",
	}, &result)
	if err != nil {
		return nil, err
	}

	return result.Docs, nil
}

func prepareIntegration(name string, host string, username string, password string, token string) IntegrationPayloadData {
	payloadData := IntegrationPayloadData{
		Name: name,
		Url:  host,
	}

	if username != "" {
		payloadData.Username = null.NewString(username, true)
	}

	if password != "" {
		payloadData.Password = null.NewString(password, true)
	}

	if token != "" {
		payloadData.Token = null.NewString(token, true)
	}
	return payloadData
}

func (a *Api) CreateIntegration(name string, host string, username string, password string, token string) error {

	err := a.requestAPI(&requestOptions{
		method: "POST",
		path:   "/argo",
		body: &IntegrationPayload{
			Type: "argo-cd",
			Data: prepareIntegration(name, host, username, password, token),
		},
	}, nil)
	if err != nil {
		return err
	}

	return nil
}

func (a *Api) UpdateIntegration(name string, host string, username string, password string, token string) error {
	err := a.requestAPI(&requestOptions{
		method: "PUT",
		path:   fmt.Sprintf("/argo/%s", name),
		body: &IntegrationPayload{
			Type: "argo-cd",
			Data: prepareIntegration(name, host, username, password, token),
		},
	}, nil)
	if err != nil {
		return err
	}

	return nil
}

func (a *Api) GetIntegrations() ([]*IntegrationPayload, error) {
	var result []*IntegrationPayload

	err := a.requestAPI(&requestOptions{
		method: "GET",
		path:   "/argo",
	}, &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (a *Api) GetIntegrationByName(name string) (*IntegrationPayload, error) {
	var result IntegrationPayload

	err := a.requestAPI(&requestOptions{
		method: "GET",
		path:   fmt.Sprintf("/argo/%s", name),
	}, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (a *Api) DeleteIntegrationByName(name string) error {
	err := a.requestAPI(&requestOptions{
		method: "DELETE",
		path:   fmt.Sprintf("/argo/%s", name),
	}, nil)
	if err != nil {
		return err
	}

	return nil
}

func (a *Api) requestAPI(opt *requestOptions, target interface{}) error {

	var body []byte
	finalURL := fmt.Sprintf("%s%s", a.Host+"/api", opt.path)
	if opt.qs != nil {
		finalURL += a.getQs(opt.qs)
	}

	if opt.body != nil {
		body, _ = json.Marshal(opt.body)
	}

	request, err := http.NewRequest(opt.method, finalURL, bytes.NewBuffer(body))

	if err != nil {
		return err
	}

	request.Header.Set("Authorization", "Bearer "+a.Token)
	request.Header.Set("Content-Type", "application/json")

	response, err := a.buildHttpClient().Do(request)

	if err != nil {
		return err
	}

	defer response.Body.Close()

	if response.StatusCode < 200 || response.StatusCode > 299 {
		cfError := &CodefreshError{}
		err = json.NewDecoder(response.Body).Decode(cfError)

		if err != nil {
			return err
		}

		cfError.URL = finalURL

		return cfError
	}

	if target == nil {
		return nil
	}

	err = json.NewDecoder(response.Body).Decode(target)

	if err != nil {
		return err
	}

	return nil
}

func (a *Api) getQs(qs map[string]string) string {
	var arr []string
	for k, v := range qs {
		arr = append(arr, fmt.Sprintf("%s=%s", k, v))
	}
	return "?" + strings.Join(arr, "&")
}

func (a *Api) CreateEnvironment(name string, project string, application string) error {
	err := a.requestAPI(&requestOptions{
		method: "POST",
		path:   "/environments-v2",
		body: &EnvironmentPayload{
			Version: "1.0",
			Metadata: EnvironmentMetadata{
				Name: name,
			},
			Spec: EnvironmentSpec{
				Type:        "argo",
				Context:     a.Integration,
				Project:     project,
				Application: application,
			},
		},
	}, nil)
	if err != nil {
		return err
	}

	return nil
}

func (a *Api) DeleteEnvironment(name string) error {
	err := a.requestAPI(&requestOptions{
		method: "DELETE",
		path:   fmt.Sprintf("/environments-v2/%s", name),
	}, nil)
	if err != nil {
		return err
	}

	return nil
}

func (a *Api) buildHttpClient() *http.Client {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	return &http.Client{Transport: tr}
}
