package transform

import (
	argoSdk "github.com/codefresh-io/argocd-sdk/pkg/api"
	"testing"
)

var _ = func() bool {
	testing.Init()
	return true
}()

func TestAdaptArgoApplications(t *testing.T) {
	var apps []argoSdk.ApplicationItem
	agentApps := AdaptArgoApplications(apps)
	if len(agentApps) != 0 {
		t.Error("Wrong result")
	}
}
