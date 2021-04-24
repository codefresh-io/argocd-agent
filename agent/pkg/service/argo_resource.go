package service

import (
	"fmt"
	"github.com/codefresh-io/argocd-listener/agent/pkg/api/codefresh"
	"github.com/codefresh-io/argocd-listener/agent/pkg/infra/logger"
	argoSdk "github.com/codefresh-io/argocd-sdk/pkg/api"
	codefreshSdk "github.com/codefresh-io/go-sdk/pkg/codefresh"
	"github.com/thoas/go-funk"
	"sort"
	"strings"
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

	ApplicationWrapper struct {
		Application argoSdk.ArgoApplication
		HistoryId   int64
	}
)

const (
	ChangedResourceKey = "configured"
)

// ArgoResourceService service for process argo resources
type ArgoResourceService interface {
	IdentifyChangedResources(argoSdk.ArgoApplication, []Resource, codefreshSdk.Commit) []Resource
	AdaptArgoProjects(projects []argoSdk.ProjectItem) []codefresh.AgentProject
	AdaptArgoApplications(applications []argoSdk.ApplicationItem) []codefresh.AgentApplication
	ResolveHistoryId(historyList []argoSdk.ApplicationHistoryItem, revision string, name string) (error, int64)
}

// NewArgoResourceService new instance of service
func NewArgoResourceService() ArgoResourceService {
	return &argoResourceService{}
}

// IdentifyChangedResources understand which resources changed during current rollout
func (argoResourceService *argoResourceService) IdentifyChangedResources(application argoSdk.ArgoApplication, serviceResources []Resource, commit codefreshSdk.Commit) []Resource {
	result := funk.Filter(application.Status.OperationState.SyncResult.Resources, func(resource argoSdk.SyncResultResource) bool {
		return strings.Contains(resource.Message, ChangedResourceKey)
	})
	syncResultResources := result.([]argoSdk.SyncResultResource)
	result = funk.Map(syncResultResources, func(syncResultResource argoSdk.SyncResultResource) Resource {
		resource := funk.Find(serviceResources, func(resource Resource) bool {
			return syncResultResource.Name == resource.Name && syncResultResource.Kind == resource.Kind
		}).(Resource)
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

func (argoResourceService *argoResourceService) ResolveHistoryId(historyList []argoSdk.ApplicationHistoryItem, revision string, name string) (error, int64) {
	if historyList == nil {
		logger.GetLogger().Errorf("can`t find history id for application %s, because history list is empty", name)
		return nil, -1
	}

	sort.Slice(historyList, func(i, j int) bool {
		return historyList[i].Id > historyList[j].Id
	})

	for _, item := range historyList {
		if item.Revision == revision {
			return nil, item.Id
		}
	}
	return fmt.Errorf("can`t find history id for application %s", name), 0
}
