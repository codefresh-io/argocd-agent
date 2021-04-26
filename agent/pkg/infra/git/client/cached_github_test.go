package client

import (
	codefreshSdk "github.com/codefresh-io/go-sdk/pkg/codefresh"
	"github.com/google/go-github/github"
	"testing"
)

var _ = func() bool {
	testing.Init()
	return true
}()

type githubApi struct {
}

func (api *githubApi) GetCommitBySha(sha string) (error, *github.RepositoryCommit) {

	return nil, &github.RepositoryCommit{
		SHA:         &sha,
		Commit:      nil,
		Author:      nil,
		Committer:   nil,
		Parents:     nil,
		HTMLURL:     nil,
		URL:         nil,
		CommentsURL: nil,
		Stats:       nil,
		Files:       nil,
	}
}

func (api *githubApi) GetUserByUsername(username string) (error, *github.User) {
	panic("Not implemented yet")
}

func (api *githubApi) GetCommitsBySha(sha string) (error, []*github.RepositoryCommit) {
	panic("Not implemented yet")
}

func (api *githubApi) GetComittersByCommits(commits []*github.RepositoryCommit) (error, []codefreshSdk.User) {
	panic("Not implemented yet")
}

func (api *githubApi) GetIssuesAndPrsByCommits(commits []*github.RepositoryCommit) (error, []codefreshSdk.Annotation, []codefreshSdk.Annotation) {
	panic("Not implemented yet")
}

func TestGetCommitBySha(t *testing.T) {
	cachedGithub := New(&githubApi{})
	sha := "revision"
	err, commit := cachedGithub.GetCommitBySha(sha)
	if err != nil {
		t.Error("Should retrieve commit without error")
	}

	if *commit.SHA != sha {
		t.Error("Sha is wrong")
	}
}
