package kube

import (
	"testing"
)

func TestNew(t *testing.T) {
	options := Options{
		ContextName:      "Context",
		Namespace:        "Namespace",
		PathToKubeConfig: "kubeConfigPath",
		InCluster:        false,
	}
	_, err := New(&options)

	if err == nil {
		t.Error("Should throw error for kubeConfigPath")
	}
}
