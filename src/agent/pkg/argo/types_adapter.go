package argo

import (
	"github.com/codefresh-io/argocd-listener/src/agent/pkg/codefresh"
)

func AdaptArgoApplications(applications []ApplicationItem) []codefresh.AgentApplication {
	var result []codefresh.AgentApplication

	for _, item := range applications {
		newItem := codefresh.AgentApplication{
			Name:    item.Metadata.Name,
			UID:     item.Metadata.UID,
			Project: item.Spec.Project,
		}
		result = append(result, newItem)
	}

	return result
}

func AdaptArgoProjects(projects []ProjectItem) []codefresh.AgentProject {
	var result []codefresh.AgentProject

	for _, item := range projects {
		newItem := codefresh.AgentProject{
			Name: item.Metadata.Name,
			UID:  item.Metadata.UID,
		}
		result = append(result, newItem)
	}

	return result
}