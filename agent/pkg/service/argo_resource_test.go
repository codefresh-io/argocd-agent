package service

import (
	argoSdk "github.com/codefresh-io/argocd-sdk/pkg/api"
	"github.com/codefresh-io/go-sdk/pkg/codefresh"
	"testing"
)

var _ = func() bool {
	testing.Init()
	return true
}()

func TestArgoResourceIdentifyChangedResources(t *testing.T) {
	service := NewArgoResourceService()

	resources := make([]Resource, 0)
	resources = append(resources, Resource{
		Status: "OutOfSync",
		Name:   "test",
		Kind:   "Service",
	})
	resources = append(resources, Resource{
		Status: "Success",
		Name:   "test2",
		Kind:   "Deployment",
	})

	commitMessage := "Commit message"
	avatar := "avatar"

	changedResources := service.IdentifyChangedResources(resources, codefresh.Commit{
		Message: &commitMessage,
		Avatar:  &avatar,
	})

	if len(changedResources) != 1 {
		t.Error("We should identify only 1 changed resource")
	}

	if *changedResources[0].Commit.Message != commitMessage {
		t.Error("Commit message is incorrect")
	}

	if *changedResources[0].Commit.Avatar != avatar {
		t.Error("Avatar is incorrect")
	}
}

func TestAdaptArgoApplicationsEmptyState(t *testing.T) {
	svc := NewArgoResourceService()

	var apps []argoSdk.ApplicationItem
	agentApps := svc.AdaptArgoApplications(apps)
	if len(agentApps) != 0 {
		t.Error("Wrong result")
	}
}

func TestAdaptArgoApplicationsNonEmpty(t *testing.T) {
	svc := NewArgoResourceService()

	apps := make([]argoSdk.ApplicationItem, 0)
	apps = append(apps, argoSdk.ApplicationItem{
		Metadata: argoSdk.ApplicationMetadata{},
		Spec:     argoSdk.ApplicationSpec{},
	})

	agentApps := svc.AdaptArgoApplications(apps)
	if len(agentApps) != 1 {
		t.Error("Wrong result")
	}
}

func TestAdaptArgoProjectsNonEmpty(t *testing.T) {
	svc := NewArgoResourceService()

	projects := make([]argoSdk.ProjectItem, 0)
	projects = append(projects, argoSdk.ProjectItem{
		Metadata: argoSdk.ProjectMetadata{
			Name: "Test",
			UID:  "UUID",
		},
	})

	agentApps := svc.AdaptArgoProjects(projects)
	if len(agentApps) != 1 {
		t.Error("Wrong result")
	}
}

func TestAdaptArgoProjectsEmptyState(t *testing.T) {
	svc := NewArgoResourceService()

	var projects []argoSdk.ProjectItem
	agentProjects := svc.AdaptArgoProjects(projects)
	if len(agentProjects) != 0 {
		t.Error("Wrong result")
	}
}
