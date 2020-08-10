package codefresh

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type Api struct {
	Token       string
	Host        string
	Integration string
}

func (a *Api) SendEnvironment(environment Environment) (map[string]interface{}, error) {
	var result map[string]interface{}
	err := a.requestAPI(&requestOptions{method: "POST", path: "/environments-v2/argo/events", body: environment}, result)
	if err != nil {
		return nil, err
	}
	fmt.Println(fmt.Sprintf("Send environment to codefresh %v", environment))
	return result, nil
}

func (a *Api) SendResources(kind string, items interface{}) error {
	err := a.requestAPI(&requestOptions{
		method: "POST",
		path:   fmt.Sprintf("/argo-agent/%s", a.Integration),
		body:   &AgentState{Kind: kind, Items: items},
	}, nil)
	if err != nil {
		return err
	}

	return nil
}

func (a *Api) CreateIntegration(name string, host string, username string, password string, ensure bool) error {
	qs := make(map[string]string)

	if ensure {
		qs["ensure"] = "true"
	}

	err := a.requestAPI(&requestOptions{
		method: "POST",
		path:   "/argo",
		qs:     qs,
		body: &IntegrationPayload{
			Type: "argo-cd",
			Data: IntegrationPayloadData{
				Name:     name,
				Url:      host,
				Username: username,
				Password: password,
			},
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

func (a *Api) buildHttpClient() *http.Client {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	return &http.Client{Transport: tr}
}
