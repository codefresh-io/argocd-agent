package transform

import (
	"github.com/codefresh-io/argocd-listener/agent/pkg/argo"
	codefresh2 "github.com/codefresh-io/argocd-listener/agent/pkg/codefresh"
)

func AdaptArgoApplications(applications []argo.ApplicationItem) []codefresh2.AgentApplication {
	var result = make([]codefresh2.AgentApplication, 0)

	for _, item := range applications {
		namespace := item.Spec.Destination.Namespace

		if namespace == "" {
			namespace = "-"
		}

		server := item.Spec.Destination.Server
		if server == "" {
			server = item.Spec.Destination.Name
		}

		newItem := codefresh2.AgentApplication{
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

func AdaptArgoProjects(projects []argo.ProjectItem) []codefresh2.AgentProject {
	var result = make([]codefresh2.AgentProject, 0)

	for _, item := range projects {
		newItem := codefresh2.AgentProject{
			Name: item.Metadata.Name,
			UID:  item.Metadata.UID,
		}
		result = append(result, newItem)
	}

	return result
}
