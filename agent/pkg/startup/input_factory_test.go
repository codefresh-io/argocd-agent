package startup

import (
	"os"
	"testing"
)

func TestInputFactory(t *testing.T) {
	argoHost := "http://host"
	argoToken := "token"
	argoUsername := "username"
	argoPassword := "password"
	codefreshToken := "token"
	envName := "name"
	newRelicLicense := "key"
	codefreshHost := "http://cf-host"
	codefreshIntegration := "integration"
	syncMode := "sync-mode"
	agentVersion := "v1"
	gitIntegration := "git-integration"
	gitPassword := "1234"

	os.Setenv("ARGO_HOST", argoHost)
	os.Setenv("ARGO_TOKEN", argoToken)
	os.Setenv("ARGO_USERNAME", argoUsername)
	os.Setenv("ARGO_PASSWORD", argoPassword)
	os.Setenv("CODEFRESH_TOKEN", codefreshToken)
	os.Setenv("NEWRELIC_LICENSE_KEY", newRelicLicense)
	os.Setenv("CODEFRESH_HOST", codefreshHost)
	os.Setenv("CODEFRESH_INTEGRATION", codefreshIntegration)
	os.Setenv("SYNC_MODE", syncMode)
	os.Setenv("AGENT_VERSION", agentVersion)
	os.Setenv("CODEFRESH_GIT_INTEGRATION", gitIntegration)
	os.Setenv("GIT_PASSWORD", gitPassword)
	os.Setenv("ENV_NAME", envName)
	os.Setenv("CREATE_INTEGRATION_IF_NOT_EXIST", "something-wrong")

	input := NewInputFactory().Create()

	if input.argoHost != argoHost {
		t.Error("ARGO_HOST env variables not parsed to structure")
	}

	if input.argoToken != argoToken {
		t.Error("ARGO_TOKEN env variables not parsed to structure")
	}

	if input.argoUsername != argoUsername {
		t.Error("ARGO_USERNAME env variables not parsed to structure")
	}

	if input.argoPassword != argoPassword {
		t.Error("ARGO_PASSWORD env variables not parsed to structure")
	}

	if input.codefreshToken != codefreshToken {
		t.Error("CODEFRESH_TOKEN env variables not parsed to structure")
	}

	if input.newRelicLicense != newRelicLicense {
		t.Error("NEWRELIC_LICENSE_KEY env variables not parsed to structure")
	}

	if input.envName != envName {
		t.Error("ENV_NAME env variables not parsed to structure")
	}

	if input.codefreshIntegrationName != codefreshIntegration {
		t.Error("CODEFRESH_INTEGRATION env variables not parsed to structure")
	}

	if input.syncMode != syncMode {
		t.Error("SYNC_MODE env variables not parsed to structure")
	}

	if input.agentVersion != agentVersion {
		t.Error("AGENT_VERSION env variables not parsed to structure")
	}

	if input.gitIntegration != gitIntegration {
		t.Error("CODEFRESH_GIT_INTEGRATION env variables not parsed to structure")
	}

	if input.password != gitPassword {
		t.Error("GIT_PASSWORD env variables not parsed to structure")
	}

	if input.createIntegrationIfNotExist != false {
		t.Error("CREATE_INTEGRATION_IF_NOT_EXIST env variables not parsed to structure")
	}

}
