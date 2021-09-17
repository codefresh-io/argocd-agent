package provider

import (
	"github.com/codefresh-io/argocd-listener/agent/pkg/infra/git/provider/api"
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

func (github *Github) GetCommitByRevision(repoUrl string, revision string) (error, *ResourceCommit) {
	err, gitClient := api.GetInstance(repoUrl)
	cachedGithub := api.New(gitClient)
	if err != nil {
		return err, nil
	}
	err, commit := cachedGithub.GetCommitBySha(revision)
	if err != nil {
		return err, nil
	}

	result := &ResourceCommit{
		Message: commit.Commit.Message,
		Sha:     &revision,
	}

	if commit.Commit.Committer.Date != nil {
		result.Time = commit.Commit.Committer.Date
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
		Comitters: []codefreshSdk.GitopsUser{},
		Prs:       []codefreshSdk.Annotation{},
		Issues:    []codefreshSdk.Annotation{},
	}
	err, gitClient := api.GetInstance(repoUrl)
	cachedGithub := api.New(gitClient)
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
