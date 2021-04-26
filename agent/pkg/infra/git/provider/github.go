package provider

import (
	git2 "github.com/codefresh-io/argocd-listener/agent/pkg/infra/git"
	client2 "github.com/codefresh-io/argocd-listener/agent/pkg/infra/git/client"
	"github.com/codefresh-io/argocd-listener/agent/pkg/service"
	codefreshSdk "github.com/codefresh-io/go-sdk/pkg/codefresh"
)

type (
	Github struct {
	}
)

var github *Github

func NewGithubProvider() GitProvider {
	if github == nil {
		github = &Github{}
	}
	return github
}

func (github *Github) GetCommitByRevision(repoUrl string, revision string) (error, *service.ResourceCommit) {
	err, gitClient := git2.GetInstance(repoUrl)
	cachedGithub := client2.New(gitClient)
	if err != nil {
		return err, nil
	}
	err, commit := cachedGithub.GetCommitBySha(revision)
	if err != nil {
		return err, nil
	}

	result := &service.ResourceCommit{
		Message: commit.Commit.Message,
		Sha:     &revision,
	}

	if commit.Author != nil {
		result.Avatar = commit.Author.AvatarURL
	} else {
		err, usr := cachedGithub.GetUserByUsername(*commit.Commit.Author.Name)
		if err == nil && usr.AvatarURL != nil {
			result.Avatar = usr.AvatarURL
		}
	}

	result.Link = commit.HTMLURL

	return nil, result
}

func (github *Github) GetManifestRepoInfo(repoUrl string, revision string) (error, *codefreshSdk.Gitops) {
	defaultGitInfo := codefreshSdk.Gitops{
		Comitters: []codefreshSdk.User{},
		Prs:       []codefreshSdk.Annotation{},
		Issues:    []codefreshSdk.Annotation{},
	}
	err, gitClient := git2.GetInstance(repoUrl)
	cachedGithub := client2.New(gitClient)
	if err != nil {
		return err, &defaultGitInfo
	}

	err, commits := cachedGithub.GetCommitsBySha(revision)
	if err != nil {
		return err, &defaultGitInfo
	}

	err, committers := gitClient.GetComittersByCommits(commits)
	if err != nil {
		return err, &defaultGitInfo
	}

	err, _, prs := gitClient.GetIssuesAndPrsByCommits(commits)
	if err != nil {
		return err, &defaultGitInfo
	}

	gitInfo := codefreshSdk.Gitops{
		Comitters: committers,
		Prs:       prs,
		Issues:    []codefreshSdk.Annotation{},
	}

	return nil, &gitInfo
}
