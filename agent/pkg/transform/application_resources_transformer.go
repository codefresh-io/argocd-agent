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
	if data == nil {
		return nil
	}
	resources, ok := data.([]interface{})
	if !ok {
		return nil
	}
	for _, elem := range resources {
		item := elem.(map[string]interface{})
		delete(item, "group")
		delete(item, "resourceVersion")
		delete(item, "version")
		delete(item, "networkingInfo")
	}
	return data
}
