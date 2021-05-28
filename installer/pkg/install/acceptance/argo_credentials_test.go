package acceptance

import (
	"github.com/codefresh-io/argocd-listener/installer/pkg/dictionary"
	"testing"
)

func TestArgoCredsFailure(t *testing.T) {
	test := &ArgoCredentialsAcceptanceTest{}

	result := test.failure()
	if !result {
		t.Error("Should fail with error")
	}
}

func TestArgoCredsGetMessage(t *testing.T) {
	test := &ArgoCredentialsAcceptanceTest{}

	result := test.getMessage()

	if result != dictionary.CheckArgoCredentials {
		t.Error("Message is incorrect")
	}
}
