package startup

import (
	"encoding/base64"
	"encoding/json"
	"github.com/codefresh-io/argocd-listener/agent/pkg/api/codefresh"
	"os"
	"strconv"
)

type Input struct {
	namespace                   string
	argoHost                    string
	argoToken                   string
	argoUsername                string
	argoPassword                string
	codefreshToken              string
	newRelicLicense             string
	envName                     string
	codefreshHost               string
	codefreshIntegrationName    string
	applications                []string
	agentVersion                string
	gitIntegration              string
	password                    string
	syncMode                    string
	createIntegrationIfNotExist bool
	numberOfShard               int
	replicas                    int
}

type InputFactory struct {
}

func NewInputFactory() *InputFactory {
	return &InputFactory{}
}

func (inputFactory *InputFactory) Create() *Input {

	createIntegrationIfNotExist, createIntegrationIfNotExistExistance := os.LookupEnv("CREATE_INTEGRATION_IF_NOT_EXIST")

	if !createIntegrationIfNotExistExistance {
		createIntegrationIfNotExist = "false"
	}

	argoHost, _ := os.LookupEnv("ARGO_HOST")
	argoToken, _ := os.LookupEnv("ARGO_TOKEN")
	argoUsername, _ := os.LookupEnv("ARGO_USERNAME")
	argoPassword, _ := os.LookupEnv("ARGO_PASSWORD")
	codefreshToken, _ := os.LookupEnv("CODEFRESH_TOKEN")
	newRelicLicense, _ := os.LookupEnv("NEWRELIC_LICENSE_KEY")
	envName, _ := os.LookupEnv("ENV_NAME")
	namespace, _ := os.LookupEnv("NAMESPACE")
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

	input := &Input{
		namespace:                namespace,
		argoHost:                 argoHost,
		argoToken:                argoToken,
		argoUsername:             argoUsername,
		argoPassword:             argoPassword,
		codefreshToken:           codefreshToken,
		newRelicLicense:          newRelicLicense,
		envName:                  envName,
		codefreshHost:            codefreshHost,
		codefreshIntegrationName: codefreshIntegrationName,
		applications:             applications,
		agentVersion:             agentVersion,
		gitIntegration:           gitIntegration,
		password:                 password,
		syncMode:                 syncMode,
	}

	createIntegrationIfNotExistBool, err := strconv.ParseBool(createIntegrationIfNotExist)
	if err != nil {
		input.createIntegrationIfNotExist = false
	} else {
		input.createIntegrationIfNotExist = createIntegrationIfNotExistBool
	}

	return input
}
