package startup

import (
	"testing"
)

var _ = func() bool {
	testing.Init()
	return true
}()

func TestInitNewRelicWithKey(t *testing.T) {
	input := &Input{
		argoHost:                 "http://argo-host",
		argoToken:                "1234",
		argoUsername:             "",
		argoPassword:             "",
		codefreshToken:           "token",
		codefreshHost:            "http://cf-host",
		newRelicLicense:          "key",
		codefreshIntegrationName: "test",
		applications:             nil,
		agentVersion:             "124",
		gitIntegration:           "integration",
		password:                 "",
		syncMode:                 "AUTOSYNC",
	}
	NewInputStore(input).Store()
	err := NewInputNewrelic(input).Init()
	if err == nil {
		t.Error("TestInitNewRelic should be failed")
	}
}

func TestInitNewRelicWithoutKey(t *testing.T) {
	input := &Input{
		argoHost:                 "http://argo-host",
		argoToken:                "1234",
		argoUsername:             "",
		argoPassword:             "",
		codefreshToken:           "token",
		codefreshHost:            "http://cf-host",
		newRelicLicense:          "",
		codefreshIntegrationName: "test",
		applications:             nil,
		agentVersion:             "124",
		gitIntegration:           "integration",
		password:                 "",
		syncMode:                 "AUTOSYNC",
	}
	NewInputStore(input).Store()
	err := NewInputNewrelic(input).Init()
	if err != nil {
		t.Error("TestInitNewRelic should not be failed")
	}
}
