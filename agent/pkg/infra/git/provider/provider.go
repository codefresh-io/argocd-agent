package provider

import (
	"github.com/codefresh-io/argocd-listener/agent/pkg/service"
	codefreshSdk "github.com/codefresh-io/go-sdk/pkg/codefresh"
)

// GitProvider interface for git providers
type GitProvider interface {
	// GetManifestRepoInfo information about manifest repo that include author, commits and other info
	GetManifestRepoInfo(repoUrl string, revision string) (error, *codefreshSdk.Gitops)
	// GetCommitByRevision retrieve git commit by sha
	GetCommitByRevision(repoUrl string, revision string) (error, *service.ResourceCommit)
}
