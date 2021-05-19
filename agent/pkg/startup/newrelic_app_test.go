package startup

import (
	"testing"
)

var _ = func() bool {
	testing.Init()
	return true
}()

func TestInitNewRelic(t *testing.T) {
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
	err := NewNewrelicApp(input).Init()
	if err == nil {
		t.Error("TestInitNewRelic should be failed")
	}
}
