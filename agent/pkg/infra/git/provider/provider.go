package provider

import (
	"github.com/codefresh-io/argocd-listener/agent/pkg/infra/store"
	"github.com/codefresh-io/argocd-listener/agent/pkg/service"
	codefreshSdk "github.com/codefresh-io/go-sdk/pkg/codefresh"
)

const (
	GitlabContextType = "git.gitlab"
)

// GitProvider interface for git providers
type GitProvider interface {
	// GetManifestRepoInfo information about manifest repo that include author, commits and other info
	GetManifestRepoInfo(repoUrl string, revision string) (error, *codefreshSdk.Gitops)
	// GetCommitByRevision retrieve git commit by sha
	GetCommitByRevision(repoUrl string, revision string) (error, *service.ResourceCommit)
}

func GetGitProvider() GitProvider {
	context := store.GetStore().Git.Context
	if context.Spec.Type == GitlabContextType {
		return NewGitlabProvider()
	}
	return NewGithubProvider()
}
