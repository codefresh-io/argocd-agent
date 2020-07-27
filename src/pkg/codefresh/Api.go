package codefresh

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"log"
	"net/http"
)

func buildHttpClient() *http.Client {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	return &http.Client{Transport: tr}
}

func SendEnvironment(environment Environment) map[string]interface{} {
	client := buildHttpClient()

	bytesRepresentation, err := json.Marshal(environment)
	if err != nil {
		log.Fatalln(err)
	}

	req, err := http.NewRequest("POST", "https://g.codefresh.io/api/environments-v2/argo/events", bytes.NewBuffer(bytesRepresentation))
	req.Header.Add("Authorization", "Bearer 5f1e81bbedd7b52b9a0fa94a.ef807a95c2a7032281b0f3c3970fcb6b")
	req.Header.Add("Content-Type", "application/json")

	resp, err := client.Do(req)

	var result map[string]interface{}

	json.NewDecoder(resp.Body).Decode(&result)

	defer resp.Body.Close()

	return result
}
