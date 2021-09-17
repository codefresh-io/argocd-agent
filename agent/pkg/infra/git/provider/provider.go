package provider

import (
	"github.com/codefresh-io/argocd-listener/agent/pkg/infra/store"
	codefreshSdk "github.com/codefresh-io/go-sdk/pkg/codefresh"
	"time"
)

const (
	GitlabContextType = "git.gitlab"
)

// GitProvider interface for git providers
type GitProvider interface {
	// GetManifestRepoInfo information about manifest repo that include author, commits and other info
	GetManifestRepoInfo(repoUrl string, revision string) (error, *codefreshSdk.Gitops)
	// GetCommitByRevision retrieve git commit by sha
	GetCommitByRevision(repoUrl string, revision string) (error, *ResourceCommit)
}

type (
	ResourceCommit struct {
		Time    *time.Time `json:"time,omitempty"`
		Message *string    `json:"message,omitempty"`
		Avatar  *string    `json:"avatar,omitempty"`
		Sha     *string    `json:"sha,omitempty"`
		Link    *string    `json:"link,omitempty"`
	}
)

func GetGitProvider() GitProvider {
	context := store.GetStore().Git.Context
	if context.Spec.Type == GitlabContextType {
		return NewGitlabProvider()
	}
	return NewGithubProvider()
}
