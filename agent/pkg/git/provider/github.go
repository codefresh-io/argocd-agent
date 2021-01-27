package provider

import (
	"github.com/codefresh-io/argocd-listener/agent/pkg/codefresh"
	"github.com/codefresh-io/argocd-listener/agent/pkg/git"
	"github.com/codefresh-io/argocd-listener/agent/pkg/git/client"
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

func (github *Github) GetCommitByRevision(repoUrl string, revision string) (error, *codefresh.Commit) {
	err, gitClient := git.GetInstance(repoUrl)
	cachedGithub := client.New(gitClient)
	if err != nil {
		return err, nil
	}
	err, commit := cachedGithub.GetCommitBySha(revision)
	if err != nil {
		return err, nil
	}

	result := &codefresh.Commit{
		Message: commit.Commit.Message,
	}

	if commit.Author != nil {
		result.Avatar = commit.Author.AvatarURL
	} else {
		err, usr := cachedGithub.GetUserByUsername(*commit.Commit.Author.Name)
		if err == nil && usr.AvatarURL != nil {
			result.Avatar = usr.AvatarURL
		}
	}

	return nil, result
}

func (github *Github) GetManifestRepoInfo(repoUrl string, revision string) (error, *git.Gitops) {
	defaultGitInfo := git.Gitops{
		Comitters: []git.User{},
		Prs:       []git.Annotation{},
		Issues:    []git.Annotation{},
	}
	err, gitClient := git.GetInstance(repoUrl)
	cachedGithub := client.New(gitClient)
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

	gitInfo := git.Gitops{
		Comitters: committers,
		Prs:       prs,
		Issues:    []git.Annotation{},
	}

	return nil, &gitInfo
}
