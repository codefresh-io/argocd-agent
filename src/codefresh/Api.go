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

	req, err := http.NewRequest("POST", "http://local.codefresh.io/api/environments-v2/argo", bytes.NewBuffer(bytesRepresentation))
	req.Header.Add("Authorization", "Bearer 5f1a9baf6adbd676fa3bd6df.6556ecb5b8b4f8fd0ca51155c02af06f")
	req.Header.Add("Content-Type", "application/json")

	resp, err := client.Do(req)

	var result map[string]interface{}

	json.NewDecoder(resp.Body).Decode(&result)

	defer resp.Body.Close()

	return result
}
