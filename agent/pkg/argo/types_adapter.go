package argo

import (
	codefresh2 "github.com/codefresh-io/argocd-listener/agent/pkg/codefresh"
)

func AdaptArgoApplications(applications []ApplicationItem) []codefresh2.AgentApplication {
	var result []codefresh2.AgentApplication

	for _, item := range applications {
		newItem := codefresh2.AgentApplication{
			Name:    item.Metadata.Name,
			UID:     item.Metadata.UID,
			Project: item.Spec.Project,
		}
		result = append(result, newItem)
	}

	return result
}

func AdaptArgoProjects(projects []ProjectItem) []codefresh2.AgentProject {
	var result []codefresh2.AgentProject

	for _, item := range projects {
		newItem := codefresh2.AgentProject{
			Name: item.Metadata.Name,
			UID:  item.Metadata.UID,
		}
		result = append(result, newItem)
	}

	return result
}
