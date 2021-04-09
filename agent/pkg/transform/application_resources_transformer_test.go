package transform

import (
	"github.com/codefresh-io/argocd-listener/agent/pkg/argo"
	"testing"
)

func TestApplicationResourcesTransformer(t *testing.T) {

	data := make([]interface{}, 1)

	item0 := make(map[string]interface{})
	item0["group"] = "group"
	item0["resourceVersion"] = "resourceVersion"
	item0["version"] = "version"
	item0["networkingInfo"] = "networkingInfo"
	item0["important"] = "important"

	data[0] = item0

	wrapper := argo.ResourcesWrapper{
		ResourcesTree:     data,
		ManifestResources: nil,
	}

	result := GetApplicationResourcesTransformer().Transform(wrapper)

	transformationResult := result.([]interface{})

	if len(transformationResult) != 1 {
		t.Errorf("Not correct amount of transformed elements")
	}

	elemToMatch := transformationResult[0].(map[string]interface{})

	if len(elemToMatch) != 1 {
		t.Errorf("Garbage not removed during transformation")
	}

	if elemToMatch["important"] != "important" {
		t.Errorf("We lost important key")
	}

	envTransformer := GetEnvTransformerInstance(&MockArgoApi{})
	if envTransformer.argoApi == nil {
		t.Errorf("Should export argoApi in struct")
	}
}

func TestApplicationResourcesTransformerInCaseNilInput(t *testing.T) {

	result := GetApplicationResourcesTransformer().Transform(nil)

	if result != nil {
		t.Errorf("Result should be nil")
	}
}

func TestApplicationResourcesTransformerInCaseManifestResourcesIncludeSyncStatus(t *testing.T) {

	data := make([]interface{}, 1)

	item0 := make(map[string]interface{})
	item0["group"] = "group"
	item0["resourceVersion"] = "resourceVersion"
	item0["version"] = "version"
	item0["networkingInfo"] = "networkingInfo"
	item0["important"] = "important"
	item0["name"] = "test"

	data[0] = item0

	manifestResources := make([]interface{}, 1)

	mitem0 := make(map[string]interface{})
	mitem0["name"] = "test"
	mitem0["status"] = "OutOfSync"
	manifestResources[0] = mitem0

	wrapper := argo.ResourcesWrapper{
		ResourcesTree:     data,
		ManifestResources: manifestResources,
	}

	result := GetApplicationResourcesTransformer().Transform(wrapper)

	transformationResult := result.([]interface{})

	if len(transformationResult) != 1 {
		t.Errorf("Not correct amount of transformed elements")
	}

	elemToMatch := transformationResult[0].(map[string]interface{})

	if len(elemToMatch) != 3 {
		t.Errorf("Garbage not removed during transformation")
	}

	if elemToMatch["status"] != "OutOfSync" {
		t.Errorf("We lost status key")
	}
}
