package api

import (
	codefreshSdk "github.com/codefresh-io/go-sdk/pkg/codefresh"
	"github.com/google/go-github/github"
	"testing"
	"time"
)

var _ = func() bool {
	testing.Init()
	return true
}()

type mgithubApi struct {
}

func _getCommitTime() *time.Time {
	commitTimeLayout := "2006-01-02T15:04:05.000Z"
	commitTimeStr := "2014-11-12T11:45:26.371Z"
	commitTime, _ := time.Parse(commitTimeLayout, commitTimeStr)
	return &commitTime
}

func (api *mgithubApi) GetCommitBySha(sha string) (error, *github.RepositoryCommit) {

	return nil, &github.RepositoryCommit{
		SHA: &sha,
		Commit: &github.Commit{
			Committer: &github.CommitAuthor{
				Date: _getCommitTime(),
			},
		},
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

func (api *mgithubApi) GetUserByUsername(username string) (error, *github.User) {
	return nil, &github.User{
		Login:             &username,
		ID:                nil,
		NodeID:            nil,
		AvatarURL:         nil,
		HTMLURL:           nil,
		GravatarID:        nil,
		Name:              nil,
		Company:           nil,
		Blog:              nil,
		Location:          nil,
		Email:             nil,
		Hireable:          nil,
		Bio:               nil,
		PublicRepos:       nil,
		PublicGists:       nil,
		Followers:         nil,
		Following:         nil,
		CreatedAt:         nil,
		UpdatedAt:         nil,
		SuspendedAt:       nil,
		Type:              nil,
		SiteAdmin:         nil,
		TotalPrivateRepos: nil,
		OwnedPrivateRepos: nil,
		PrivateGists:      nil,
		DiskUsage:         nil,
		Collaborators:     nil,
		Plan:              nil,
		URL:               nil,
		EventsURL:         nil,
		FollowingURL:      nil,
		FollowersURL:      nil,
		GistsURL:          nil,
		OrganizationsURL:  nil,
		ReceivedEventsURL: nil,
		ReposURL:          nil,
		StarredURL:        nil,
		SubscriptionsURL:  nil,
		TextMatches:       nil,
		Permissions:       nil,
	}
}

func (api *mgithubApi) GetCommitsBySha(sha string) (error, []*github.RepositoryCommit) {
	return nil, []*github.RepositoryCommit{
		{
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
		},
		{
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
		},
	}
}

func (api *mgithubApi) GetComittersByCommits(commits []*github.RepositoryCommit) (error, []codefreshSdk.GitopsUser) {
	panic("Not implemented yet")
}

func (api *mgithubApi) GetIssuesAndPrsByCommits(commits []*github.RepositoryCommit) (error, []codefreshSdk.Annotation, []codefreshSdk.Annotation) {
	panic("Not implemented yet")
}

func TestGetCommitBySha(t *testing.T) {
	cachedGithub := New(&mgithubApi{})
	sha := "revision"
	err, commit := cachedGithub.GetCommitBySha(sha)
	if err != nil {
		t.Error("Should retrieve commit without error")
	}
	if commit.Commit.Committer.Date.String() != _getCommitTime().String() {
		t.Error("Commit date is wrong")
	}

	if *commit.SHA != sha {
		t.Error("Sha is wrong")
	}
}

func TestGetCommitsBySha(t *testing.T) {
	cachedGithub := New(&mgithubApi{})
	sha := "revision"
	err, commits := cachedGithub.GetCommitsBySha(sha)
	if err != nil {
		t.Error("Should retrieve commit without error")
	}

	if len(commits) != 2 {
		t.Error("Retrieve wrong number of commits")
	}

	if *commits[0].SHA != sha {
		t.Error("Sha is wrong")
	}
}

func TestGetUserByUsername(t *testing.T) {
	cachedGithub := New(&mgithubApi{})
	username := "test"
	err, user := cachedGithub.GetUserByUsername(username)
	if err != nil {
		t.Error("Should retrieve user without error")
	}

	if *user.Login != username {
		t.Error("Username is wrong")
	}
}
