package git

import (
	"github.com/google/go-github/github"
	"testing"
)

var _ = func() bool {
	testing.Init()
	return true
}()

func TestExtractRepoAndOwnerFromUrl(t *testing.T) {
	urls := []string{
		"ssh://user@host.xz/path/owner/repo.git/",
		"ssh://user@host.xz:8080/path/owner/repo.git/",
		"git://host.xz/path/owner/repo.git/",
		"git://host.xz:8080/path/owner/repo.git/",
		"http://host.xz/path/owner/repo.git/",
		"http://host.xz:8080/path/owner/repo.git/",
		"ftp://host.xz/path/owner/repo.git/",
		"ftp://host.xz:8080/path/owner/repo.git/",
	}

	for _, url := range urls {
		err, owner, repo := extractRepoAndOwnerFromUrl(url)
		if err != nil {
			t.Errorf("'ExtractRepoAndOwnerFromUrl' check error failed, error: %v", err.Error())
		}
		if owner != "owner" {
			t.Errorf("'ExtractRepoAndOwnerFromUrl check owner' failed, expected '%v', got '%v'", "owner", owner)
		}
		if repo != "repo" {
			t.Errorf("'ExtractRepoAndOwnerFromUrl check repo' failed, expected '%v', got '%v'", "repo", repo)
		}
	}
}

func TestGetInstance(t *testing.T) {
	err, api := GetInstance("https://github.com/owner/repo")
	_ = api
	if err != nil {
		t.Errorf("'GetInstance' check error failed, error: %v", err.Error())
	}
}

func TestGetComittersByCommits(t *testing.T) {
	api := api{
		Token:  "",
		Client: nil,
		Owner:  "",
		Repo:   "",
		Ctx:    nil,
	}

	login := "Ivan"
	avatar := "Link"

	commits := []*github.RepositoryCommit{
		{
			SHA:    nil,
			Commit: nil,
			Author: &github.User{
				Login:             &login,
				ID:                nil,
				NodeID:            nil,
				AvatarURL:         &avatar,
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
			},
			Committer:   nil,
			Parents:     nil,
			HTMLURL:     nil,
			URL:         nil,
			CommentsURL: nil,
			Stats:       nil,
			Files:       nil,
		},
	}

	err, result := api.GetComittersByCommits(commits)
	if err != nil {
		t.Error("Shouldnt fail during retrieve committers from commit")
	}

	if len(result) != 1 {
		t.Error("Should retrieve only one comitter from commits")
	}

	if result[0].Name != login {
		t.Error("Wrong login name")
	}

	if result[0].Avatar != avatar {
		t.Error("Wrong avatar")
	}

}
