package transform

type ApplicationResourcesTransformer struct {
}

var applicationResourcesTransformer *ApplicationResourcesTransformer

func GetApplicationResourcesTransformer() Transformer {
	if applicationResourcesTransformer == nil {
		applicationResourcesTransformer = new(ApplicationResourcesTransformer)
	}
	return applicationResourcesTransformer
}

func (applicationResourcesTransformer *ApplicationResourcesTransformer) Transform(data interface{}) interface{} {
	for _, elem := range data.([]interface{}) {
		item := elem.(map[string]interface{})
		delete(item, "group")
		delete(item, "resourceVersion")
		delete(item, "version")
		delete(item, "namespace")
		delete(item, "networkingInfo")
	}
	return data
}
