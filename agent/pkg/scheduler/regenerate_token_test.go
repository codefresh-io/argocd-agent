package scheduler

import (
	"github.com/codefresh-io/argocd-listener/agent/pkg/infra/store"
	argoSdk "github.com/codefresh-io/argocd-sdk/pkg/api"
	"testing"
)

var _ = func() bool {
	testing.Init()
	return true
}()

type MockUnathourizedArgoApi struct {
}

func (api *MockUnathourizedArgoApi) GetApplications(token string, host string) ([]argoSdk.ApplicationItem, error) {
	return nil, nil
}

func (api *MockUnathourizedArgoApi) GetToken(username string, password string, host string) (string, error) {
	return "token", nil
}

func TestRegenerateTokenScheduler(t *testing.T) {

	store.SetArgo("", "", "test", "test")

	regenerateTokenScheduler := regenerateTokenScheduler{argoApi: &MockUnathourizedArgoApi{}}
	regenerateTokenScheduler.regenerateToken()

	tk := store.GetStore().Argo.Token
	if tk == "" {
		t.Error("Token should be regenerated")
	}

}

func TestRegenerateTokenCreateInstance(t *testing.T) {
	regenerateTokenScheduler := GetRegenerateTokenScheduler()
	if regenerateTokenScheduler == nil {
		t.Error("Cant initialize regenerate token scheduler")
	}
}

func TestRegenerateTokenExecutionTime(t *testing.T) {
	regenerateTokenScheduler := GetRegenerateTokenScheduler()
	time := regenerateTokenScheduler.getTime()
	if time != "@every 30m" {
		t.Error("Wrong schedule time")
	}
}
