package service

import (
	"github.com/codefresh-io/go-sdk/pkg/codefresh"
	"testing"
)

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
