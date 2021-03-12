package git

import (
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
