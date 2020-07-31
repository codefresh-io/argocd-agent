package codefresh

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	store2 "github.com/codefresh-io/argocd-listener/agent/pkg/store"
	"net/http"
)

func buildHttpClient() *http.Client {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	return &http.Client{Transport: tr}
}

func requestAPI(opt *requestOptions, target interface{}) error {

	codefresh := store2.GetStore().Codefresh

	var body []byte
	finalURL := fmt.Sprintf("%s%s", codefresh.Host+"/api", opt.path)

	if opt.body != nil {
		body, _ = json.Marshal(opt.body)
	}

	request, err := http.NewRequest(opt.method, finalURL, bytes.NewBuffer(body))

	if err != nil {
		return err
	}

	request.Header.Set("Authorization", "Bearer "+codefresh.Token)
	request.Header.Set("Content-Type", "application/json")

	response, err := buildHttpClient().Do(request)

	if err != nil {
		return err
	}

	if response.StatusCode < 200 || response.StatusCode > 299 {
		cfError := &CodefreshError{}
		err = json.NewDecoder(response.Body).Decode(cfError)

		if err != nil {
			return err
		}

		return fmt.Errorf("%d: %s", response.StatusCode, cfError.Message)
	}

	err = json.NewDecoder(response.Body).Decode(target)

	if err != nil {
		return err
	}

	return nil
}

func SendEnvironment(environment Environment) (map[string]interface{}, error) {
	var result map[string]interface{}
	err := requestAPI(&requestOptions{method: "POST", path: "/environments-v2/argo/events", body: environment}, result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func SendResources(kind string, items interface{}) error {
	err := requestAPI(&requestOptions{
		method: "POST",
		path:   "/argo-agent/argo-demo", //TODO: change to integration name from store
		body:   &AgentState{Kind: kind, Items: items},
	}, nil)
	if err != nil {
		return err
	}

	return nil
}
