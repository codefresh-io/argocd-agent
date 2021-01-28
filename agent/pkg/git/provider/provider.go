package provider

import (
	"github.com/codefresh-io/argocd-listener/agent/pkg/codefresh"
	"github.com/codefresh-io/argocd-listener/agent/pkg/git"
)

type GitProvider interface {
	GetManifestRepoInfo(repoUrl string, revision string) (error, *git.Gitops)
	GetCommitByRevision(repoUrl string, revision string) (error, *codefresh.Commit)
}
