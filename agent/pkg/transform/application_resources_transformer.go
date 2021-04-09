package transform

import "github.com/codefresh-io/argocd-listener/agent/pkg/argo"

type ApplicationResourcesTransformer struct {
}

var applicationResourcesTransformer *ApplicationResourcesTransformer

func GetApplicationResourcesTransformer() Transformer {
	if applicationResourcesTransformer == nil {
		applicationResourcesTransformer = new(ApplicationResourcesTransformer)
	}
	return applicationResourcesTransformer
}

func lookForRelatedManifestResource(name string, resources []interface{}) map[string]interface{} {
	for _, elem := range resources {
		item := elem.(map[string]interface{})
		if item["name"] == name {
			return item
		}
	}
	return nil
}

func (applicationResourcesTransformer *ApplicationResourcesTransformer) Transform(data interface{}) interface{} {

	if data == nil {
		return nil
	}

	dataObj, ok := data.(argo.ResourcesWrapper)
	if !ok {
		return nil
	}
	resourcestree := dataObj.ResourcesTree

	for _, elem := range resourcestree {
		item := elem.(map[string]interface{})
		delete(item, "group")
		delete(item, "resourceVersion")
		delete(item, "version")
		delete(item, "networkingInfo")

		manifestResource := lookForRelatedManifestResource(item["name"].(string), dataObj.ManifestResources)
		if manifestResource != nil {
			item["status"] = manifestResource["status"]
		}
	}
	return data
}
