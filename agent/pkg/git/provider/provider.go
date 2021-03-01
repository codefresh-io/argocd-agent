package provider

import (
	codefreshSdk "github.com/codefresh-io/go-sdk/pkg/codefresh"
)

type GitProvider interface {
	GetManifestRepoInfo(repoUrl string, revision string) (error, *codefreshSdk.Gitops)
	GetCommitByRevision(repoUrl string, revision string) (error, *codefreshSdk.Commit)
}
