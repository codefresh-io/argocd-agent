package util

import (
	argoSdk "github.com/codefresh-io/argocd-sdk/pkg/api"
	"testing"
)

var _ = func() bool {
	testing.Init()
	return true
}()

func TestContains(t *testing.T) {
	arr := []string{"1", "2"}
	result := Contains(arr, "2")
	if !result {
		t.Error("Element should be found")
	}
}

func TestContainsFalse(t *testing.T) {
	arr := []string{"1", "2"}
	result := Contains(arr, "3")
	if result {
		t.Error("Element should be not found")
	}
}

func TestConvert(t *testing.T) {
	labels := map[string]interface{}{"app.kubernetes.io/instance": "apps-root"}
	envItem := map[string]interface{}{
		"metadata": struct {
			name   string
			labels map[string]interface{}
		}{
			labels: labels,
			name:   "test",
		},
	}

	var env argoSdk.ArgoApplication

	Convert(envItem, &env)

	if env.Metadata.Name != "test" {
		t.Error("Wrong environment name")
	}
}
