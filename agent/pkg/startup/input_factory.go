package startup

import (
	"encoding/base64"
	"encoding/json"
	"github.com/codefresh-io/argocd-listener/agent/pkg/codefresh"
	"os"
)

type Input struct {
	argoHost                 string
	argoToken                string
	argoUsername             string
	argoPassword             string
	codefreshToken           string
	codefreshHost            string
	codefreshIntegrationName string
	applications             []string
	agentVersion             string
	gitIntegration           string
	password                 string
	syncMode                 string
}

type InputFactory struct {
}

func NewInputFactory() *InputFactory {
	return &InputFactory{}
}

func (inputFactory *InputFactory) Create() *Input {

	argoHost, _ := os.LookupEnv("ARGO_HOST")
	argoToken, _ := os.LookupEnv("ARGO_TOKEN")
	argoUsername, _ := os.LookupEnv("ARGO_USERNAME")
	argoPassword, _ := os.LookupEnv("ARGO_PASSWORD")
	codefreshToken, _ := os.LookupEnv("CODEFRESH_TOKEN")
	codefreshHost, codefreshHostExistance := os.LookupEnv("CODEFRESH_HOST")
	if !codefreshHostExistance {
		codefreshHost = "https://g.codefresh.io"
	}
	codefreshIntegrationName, _ := os.LookupEnv("CODEFRESH_INTEGRATION")

	var applications []string
	syncMode, _ := os.LookupEnv("SYNC_MODE")
	if syncMode == codefresh.SelectSync {
		applicationsToSyncEncodedJson, _ := os.LookupEnv("APPLICATIONS_FOR_SYNC")
		applicationsToSyncJson, _ := base64.StdEncoding.DecodeString(applicationsToSyncEncodedJson)
		_ = json.Unmarshal(applicationsToSyncJson, &applications)
	}

	agentVersion, _ := os.LookupEnv("AGENT_VERSION")

	gitIntegration, _ := os.LookupEnv("CODEFRESH_GIT_INTEGRATION")

	password, _ := os.LookupEnv("GIT_PASSWORD")

	return &Input{
		argoHost:                 argoHost,
		argoToken:                argoToken,
		argoUsername:             argoUsername,
		argoPassword:             argoPassword,
		codefreshToken:           codefreshToken,
		codefreshHost:            codefreshHost,
		codefreshIntegrationName: codefreshIntegrationName,
		applications:             applications,
		agentVersion:             agentVersion,
		gitIntegration:           gitIntegration,
		password:                 password,
		syncMode:                 syncMode,
	}
}