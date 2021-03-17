package transform

import "testing"

func TestApplicationResourcesTransformer(t *testing.T) {
	data := make([]interface{}, 1)

	item0 := make(map[string]interface{})
	item0["group"] = "group"
	item0["resourceVersion"] = "resourceVersion"
	item0["version"] = "version"
	item0["namespace"] = "namespace"
	item0["networkingInfo"] = "networkingInfo"
	item0["important"] = "important"

	data[0] = item0

	result := GetApplicationResourcesTransformer().Transform(data)

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

	envTransformer := GetEnvTransformerInstance(MockArgoApi{})
	if envTransformer.argoApi == nil {
		t.Errorf("Should export argoApi in struct")
	}
}
