package startup

import (
	"github.com/codefresh-io/argocd-listener/agent/pkg/infra/store"
	"testing"
)

var _ = func() bool {
	testing.Init()
	return true
}()

func TestValidStore(t *testing.T) {
	input := &Input{
		argoHost:                 "http://argo-host",
		argoToken:                "1234",
		argoUsername:             "",
		argoPassword:             "",
		codefreshToken:           "token",
		newRelicLicense:          "key",
		envName:                  "name",
		codefreshHost:            "http://cf-host",
		codefreshIntegrationName: "",
		applications:             nil,
		agentVersion:             "124",
		gitIntegration:           "integration",
		password:                 "",
		syncMode:                 "AUTOSYNC",
	}

	err := NewInputStore(input).Store()

	argo := store.GetStore().Argo
	if argo.Token != input.argoToken || argo.Host != input.argoHost {
		t.Error("Failed to retrieve correct information from store about argo")
	}

	codefresh := store.GetStore().Codefresh
	if codefresh.SyncMode != input.syncMode {
		t.Error("Failed to retrieve correct information from store about sync mode")
	}

	if codefresh.Host != input.codefreshHost ||
		codefresh.Token != input.codefreshToken ||
		codefresh.Integration != input.codefreshIntegrationName {
		t.Error("Failed to retrieve correct information from store about codefresh creds")
	}

	newRelic := store.GetStore().NewRelic

	if newRelic.Key != input.newRelicLicense {
		t.Error("Failed to retrieve correct information from store about newrelic creds")
	}

	env := store.GetStore().Env
	if env.Name != input.envName {
		t.Error("Failed to retrieve correct information from store about env name")
	}

	agent := store.GetStore().Agent
	if agent.Version != input.agentVersion {
		t.Error("Failed to retrieve correct information from store about agent version")
	}

	if err != nil {
		t.Error("Validation should be passed, input is correct")
	}
}
