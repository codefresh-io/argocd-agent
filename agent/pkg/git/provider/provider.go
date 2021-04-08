package provider

import (
	codefreshSdk "github.com/codefresh-io/go-sdk/pkg/codefresh"
)

// GitProvider interface for git providers
type GitProvider interface {
	// GetManifestRepoInfo information about manifest repo that include author, commits and other info
	GetManifestRepoInfo(repoUrl string, revision string) (error, *codefreshSdk.Gitops)
	// GetCommitByRevision retrieve git commit by sha
	GetCommitByRevision(repoUrl string, revision string) (error, *codefreshSdk.Commit)
}
