package acceptance

import (
	"github.com/codefresh-io/argocd-listener/installer/pkg/dictionary"
	"github.com/codefresh-io/argocd-listener/installer/pkg/install/entity"
	argoSdk "github.com/codefresh-io/argocd-sdk/pkg/api"
	"testing"
)

type MMockPrompt struct {
}

func (p *MMockPrompt) InputWithDefault(target *string, label string, defaultValue string) error {
	*target = defaultValue
	return nil
}

func (p *MMockPrompt) InputPassword(target *string, label string) error {
	return nil
}

func (p *MMockPrompt) Input(target *string, label string) error {
	return nil
}

func (p *MMockPrompt) Confirm(label string) (error, bool) {
	return nil, false
}

func (p *MMockPrompt) Multiselect(items []string, label string) (error, []string) {
	return nil, nil
}

func (p *MMockPrompt) Select(items []string, label string) (error, string) {
	return nil, dictionary.StopInstallation
}

type MMockUnathourizedArgoApi struct {
}

func (api *MMockUnathourizedArgoApi) GetApplications(token string, host string) ([]argoSdk.ApplicationItem, error) {
	return nil, nil
}

func (api *MMockUnathourizedArgoApi) GetToken(username string, password string, host string) (string, error) {
	return "token", nil
}

func (api *MMockUnathourizedArgoApi) GetVersion(host string) (string, error) {
	return "v1", nil
}

func TestGetMessage(t *testing.T) {
	test := &ArgoAccessibilityAcceptanceTest{
		unathorizedArgoApi: nil,
		prompt:             nil,
	}

	msg := test.getMessage()

	if msg != dictionary.CheckArgoServerAccessability {
		t.Error("Wrong message")
	}
}

func TestArgoAccessibilityFailureWithStopInstallation(t *testing.T) {
	test := &ArgoAccessibilityAcceptanceTest{
		unathorizedArgoApi: nil,
		prompt:             &MMockPrompt{},
	}

	result := test.failure(&entity.ArgoOptions{})
	if !result {
		t.Error("Wrong failure result")
	}
}

func TestArgoAccessibilityFailureWithSpecificPrefix(t *testing.T) {
	test := &ArgoAccessibilityAcceptanceTest{
		unathorizedArgoApi: nil,
		prompt:             &MMockPrompt{},
	}

	result := test.failure(&entity.ArgoOptions{
		Host: "https://argocd-server",
	})

	if result {
		t.Error("Wrong failure result")
	}
}

func TestArgoAccessibilityCheck(t *testing.T) {
	test := &ArgoAccessibilityAcceptanceTest{
		unathorizedArgoApi: &MMockUnathourizedArgoApi{},
		prompt:             &MMockPrompt{},
	}

	result, _ := test.check(&entity.ArgoOptions{})
	if result != nil {
		t.Error("Wrong check result")
	}
}
