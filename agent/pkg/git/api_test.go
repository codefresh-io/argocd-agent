package git

import (
	"testing"
)

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
			t.Errorf("'ExtractRepoAndOwnerFromUrl' failed, error: %v", err.Error())
		}
		if owner != "owner" {
			t.Errorf("'ExtractRepoAndOwnerFromUrl' failed, expected '%v', got '%v'", "owner", owner)
		}
		if repo != "repo.git" {
			t.Errorf("'ExtractRepoAndOwnerFromUrl' failed, expected '%v', got '%v'", "repo.git", repo)
		}
	}
}
