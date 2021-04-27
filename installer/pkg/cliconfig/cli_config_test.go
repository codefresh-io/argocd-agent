package cliconfig

import "testing"

var _ = func() bool {
	testing.Init()
	return true
}()

func TestGetCurrentConfig(t *testing.T) {
	config := `contexts:
  local:
    type: APIKey
    name: local
    url: 'http://local.codefresh.io'
    token: token
    beta: false
    onPrem: true
current-context: local`
	cliConfig := cliConfig{}
	configItem, err := cliConfig.getCurrentConfig([]byte(config))
	if err != nil {
		t.Error("Get current config should be executed without error")
	}

	if configItem == nil || configItem.Name != "local" {
		t.Error("Got issue during retrieve config")
	}
}
