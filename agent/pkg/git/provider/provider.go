package provider

import (
	"github.com/codefresh-io/argocd-listener/agent/pkg/store"
	codefreshSdk "github.com/codefresh-io/go-sdk/pkg/codefresh"
)

type GitProvider interface {
	GetManifestRepoInfo(repoUrl string, revision string) (error, *codefreshSdk.Gitops)
	GetCommitByRevision(repoUrl string, revision string) (error, *codefreshSdk.Commit)
}

func Get() GitProvider {
	gitConfig := store.GetStore().Git
	if gitConfig.Type == "git.github" {
		return NewGithubProvider()
	}

	if gitConfig.Type == "git.bitbucket" {
		return NewBitbucketProvider()
	}

	return nil
}
