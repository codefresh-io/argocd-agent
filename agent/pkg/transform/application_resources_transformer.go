package transform

import (
	"github.com/codefresh-io/argocd-listener/agent/pkg/service"
)

// ApplicationResourcesTransformer handler for normalize application resources
type ApplicationResourcesTransformer struct {
}

var applicationResourcesTransformer *ApplicationResourcesTransformer

// GetApplicationResourcesTransformer singleton for get ApplicationResourcesTransformer instance
func GetApplicationResourcesTransformer() Transformer {
	if applicationResourcesTransformer == nil {
		applicationResourcesTransformer = new(ApplicationResourcesTransformer)
	}
	return applicationResourcesTransformer
}

func lookForRelatedManifestResource(appElem interface{}, resources []service.Resource) *service.Resource {
	for _, resource := range resources {
		appItem := appElem.(map[string]interface{})

		if resource.Kind != "" && resource.Name != "" && (resource.Name == appItem["name"]) && (resource.Kind == appItem["kind"]) {
			return &resource
		}
	}
	return nil
}

// Transform convert income data into argo resource
func (applicationResourcesTransformer *ApplicationResourcesTransformer) Transform(data interface{}) interface{} {

	if data == nil {
		return nil
	}

	dataObj, ok := data.(service.ResourcesWrapper)
	if !ok {
		return nil
	}
	resourcestree := dataObj.ResourcesTree

	for _, elem := range resourcestree {
		item := elem.(map[string]interface{})
		delete(item, "group")

		manifestResource := lookForRelatedManifestResource(item, dataObj.ManifestResources)
		if manifestResource != nil {
			item["status"] = manifestResource.Status
			item["commit"] = manifestResource.Commit
			item["updateAt"] = manifestResource.UpdatedAt
			item["historyId"] = manifestResource.HistoryId
		}
	}
	return resourcestree
}
