package provider

import (
	"fmt"
	"github.com/codefresh-io/argocd-listener/agent/pkg/store"
	codefreshSdk "github.com/codefresh-io/go-sdk/pkg/codefresh"
	bitbucketSdk "github.com/ktrysmt/go-bitbucket"
	"strings"
)

type (
	Bitbucket struct {
		client *bitbucketSdk.Client
	}
)

var bitbucket *Bitbucket

func NewBitbucketProvider() GitProvider {
	if bitbucket == nil {
		gitConfig := store.GetStore().Git
		bitbucket = &Bitbucket{client: bitbucketSdk.NewOAuthbearerToken(gitConfig.Token)}

	}
	return bitbucket
}

func (bitbucket *Bitbucket) GetCommitByRevision(repoUrl string, revision string) (error, *codefreshSdk.Commit) {
	repoUrl = "https://bitbucket.org/vadpasseka/iremember"
	// repo url example = https://github.com/andrii-codefresh/argo
	parts := strings.Split(repoUrl, "/")

	commit, err := bitbucket.client.Repositories.Commits.GetCommit(&bitbucketSdk.CommitsOptions{
		Owner:    parts[3],
		RepoSlug: parts[4],
		Revision: revision,
	})

	fmt.Println(commit)

	return err, nil
}

func (bitbucket *Bitbucket) GetManifestRepoInfo(repoUrl string, revision string) (error, *codefreshSdk.Gitops) {
	defaultGitInfo := codefreshSdk.Gitops{
		Comitters: []codefreshSdk.User{},
		Prs:       []codefreshSdk.Annotation{},
		Issues:    []codefreshSdk.Annotation{},
	}
	//err, gitClient := git.GetInstance(repoUrl)
	//cachedGithub := client.New(gitClient)
	//if err != nil {
	//	return err, &defaultGitInfo
	//}
	//
	//err, commits := cachedGithub.GetCommitsBySha(revision)
	//if err != nil {
	//	return err, &defaultGitInfo
	//}
	//
	//err, committers := gitClient.GetComittersByCommits(commits)
	//if err != nil {
	//	return err, &defaultGitInfo
	//}
	//
	//err, _, prs := gitClient.GetIssuesAndPrsByCommits(commits)
	//if err != nil {
	//	return err, &defaultGitInfo
	//}
	//
	//gitInfo := codefreshSdk.Gitops{
	//	Comitters: committers,
	//	Prs:       prs,
	//	Issues:    []codefreshSdk.Annotation{},
	//}

	return nil, &defaultGitInfo
}
