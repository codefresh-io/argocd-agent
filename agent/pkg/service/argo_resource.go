package service

import (
	"github.com/codefresh-io/argocd-listener/agent/pkg/api/codefresh"
	argoSdk "github.com/codefresh-io/argocd-sdk/pkg/api"
	codefreshSdk "github.com/codefresh-io/go-sdk/pkg/codefresh"
	"github.com/thoas/go-funk"
)

type (
	argoResourceService struct {
	}

	Resource struct {
		Status string
		Name   string
		Commit codefreshSdk.Commit
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
	IdentifyChangedResources([]Resource, codefreshSdk.Commit) []Resource
	AdaptArgoProjects(projects []argoSdk.ProjectItem) []codefresh.AgentProject
	AdaptArgoApplications(applications []argoSdk.ApplicationItem) []codefresh.AgentApplication
}

// NewArgoResourceService new instance of service
func NewArgoResourceService() ArgoResourceService {
	return &argoResourceService{}
}

// IdentifyChangedResources understand which resources changed during current rollout
func (argoResourceService *argoResourceService) IdentifyChangedResources(resources []Resource, commit codefreshSdk.Commit) []Resource {
	result := funk.Filter(resources, func(resource Resource) bool {
		return resource.Status == OutOfSync
	})
	result = funk.Map(result.([]Resource), func(resource Resource) Resource {
		resource.Commit = commit
		return resource
	})
	return result.([]Resource)
}

func (argoResourceService *argoResourceService) AdaptArgoApplications(applications []argoSdk.ApplicationItem) []codefresh.AgentApplication {
	var result = make([]codefresh.AgentApplication, 0)

	for _, item := range applications {
		namespace := item.Spec.Destination.Namespace

		if namespace == "" {
			namespace = "-"
		}

		server := item.Spec.Destination.Server
		if server == "" {
			server = item.Spec.Destination.Name
		}

		newItem := codefresh.AgentApplication{
			Name:      item.Metadata.Name,
			UID:       item.Metadata.UID,
			Project:   item.Spec.Project,
			Server:    server,
			Namespace: namespace,
		}
		result = append(result, newItem)
	}

	return result
}

func (argoResourceService *argoResourceService) AdaptArgoProjects(projects []argoSdk.ProjectItem) []codefresh.AgentProject {
	var result = make([]codefresh.AgentProject, 0)

	for _, item := range projects {
		newItem := codefresh.AgentProject{
			Name: item.Metadata.Name,
			UID:  item.Metadata.UID,
		}
		result = append(result, newItem)
	}

	return result
}
