package transform

import (
	"encoding/json"
	"github.com/codefresh-io/argocd-listener/agent/pkg/service"
	"testing"
)

var _ = func() bool {
	testing.Init()
	return true
}()

func TestApplicationResourcesTransformer(t *testing.T) {

	data := make([]interface{}, 1)

	item0 := make(map[string]interface{})
	item0["group"] = "group"
	item0["resourceVersion"] = "resourceVersion"
	item0["version"] = "version"
	item0["networkingInfo"] = "networkingInfo"
	item0["important"] = "important"

	data[0] = item0

	wrapper := service.ResourcesWrapper{
		ResourcesTree:     data,
		ManifestResources: nil,
	}

	result := GetApplicationResourcesTransformer().Transform(wrapper)

	transformationResult := result.([]interface{})

	if len(transformationResult) != 1 {
		t.Errorf("Not correct amount of transformed elements")
	}

	elemToMatch := transformationResult[0].(map[string]interface{})

	if len(elemToMatch) != 4 {
		t.Errorf("Garbage not removed during transformation")
	}

	if elemToMatch["important"] != "important" {
		t.Errorf("We lost important key")
	}

}

func TestApplicationResourcesTransformerInCaseNilInput(t *testing.T) {

	result := GetApplicationResourcesTransformer().Transform(nil)

	if result != nil {
		t.Errorf("Result should be nil")
	}
}

func convert(resources interface{}) []*service.Resource {
	manifestResourcesJson, _ := json.Marshal(resources)

	var manifestResourcesStruct []*service.Resource

	_ = json.Unmarshal(manifestResourcesJson, &manifestResourcesStruct)
	return manifestResourcesStruct
}

func TestApplicationResourcesTransformerInCaseManifestResourcesNotIncludeKind(t *testing.T) {

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

	wrapper := service.ResourcesWrapper{
		ResourcesTree:     data,
		ManifestResources: convert(manifestResources),
	}

	result := GetApplicationResourcesTransformer().Transform(wrapper)

	transformationResult := result.([]interface{})

	if len(transformationResult) != 1 {
		t.Errorf("Not correct amount of transformed elements")
	}

	elemToMatch := transformationResult[0].(map[string]interface{})

	if len(elemToMatch) != 5 {
		t.Errorf("Garbage not removed during transformation")
	}

	if elemToMatch["status"] == "" {
		t.Errorf("Status should not be found")
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
	item0["kind"] = "Service"

	data[0] = item0

	manifestResources := make([]interface{}, 1)

	mitem0 := make(map[string]interface{})
	mitem0["name"] = "test"
	mitem0["status"] = "OutOfSync"
	mitem0["kind"] = "Service"
	manifestResources[0] = mitem0

	wrapper := service.ResourcesWrapper{
		ResourcesTree:     data,
		ManifestResources: convert(manifestResources),
	}

	result := GetApplicationResourcesTransformer().Transform(wrapper)

	transformationResult := result.([]interface{})

	if len(transformationResult) != 1 {
		t.Errorf("Not correct amount of transformed elements")
	}

	elemToMatch := transformationResult[0].(map[string]interface{})

	if len(elemToMatch) != 10 {
		t.Errorf("Garbage not removed during transformation")
	}

	if elemToMatch["status"] != "OutOfSync" {
		t.Errorf("Status should be OutOfSync")
	}
}
