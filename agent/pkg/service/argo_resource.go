package service

import (
	"github.com/codefresh-io/go-sdk/pkg/codefresh"
	"github.com/thoas/go-funk"
)

type (
	argoResourceService struct {
	}

	Resource struct {
		Status string
		Name   string
		Commit codefresh.Commit
		Kind   string
	}

	ResourcesWrapper struct {
		ResourcesTree     []interface{}
		ManifestResources []Resource
	}
)

const (
	OutOfSync = "OutOfSync"
)

// ArgoResourceService service for process argo resources
type ArgoResourceService interface {
	IdentifyChangedResources([]Resource, codefresh.Commit) []Resource
}

// NewArgoResourceService new instance of service
func NewArgoResourceService() ArgoResourceService {
	return &argoResourceService{}
}

// IdentifyChangedResources understand which resources changed during current rollout
func (argoResourceService *argoResourceService) IdentifyChangedResources(resources []Resource, commit codefresh.Commit) []Resource {
	result := funk.Filter(resources, func(resource Resource) bool {
		return resource.Status == OutOfSync
	})
	result = funk.Map(result.([]Resource), func(resource Resource) Resource {
		resource.Commit = commit
		return resource
	})
	return result.([]Resource)
}
