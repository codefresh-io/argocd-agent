package client

import (
	"fmt"
	git2 "github.com/codefresh-io/argocd-listener/agent/pkg/infra/git"
	"github.com/google/go-github/github"
)

type CachedGithub struct {
	GitClient      git2.Api
	commitBySha    map[string]github.RepositoryCommit
	commitsBySha   map[string][]github.RepositoryCommit
	userByUsername map[string]github.User
}

var cachedGithub *CachedGithub

var cacheSize = 1000

func New(gitClient git2.Api) *CachedGithub {
	if cachedGithub == nil {
		commitBySha := make(map[string]github.RepositoryCommit)
		commitsBySha := make(map[string][]github.RepositoryCommit)
		userByUsername := make(map[string]github.User)

		cachedGithub = &CachedGithub{GitClient: gitClient, commitBySha: commitBySha, commitsBySha: commitsBySha, userByUsername: userByUsername}
	}
	return cachedGithub
}

func (cachedGithub *CachedGithub) GetCommitBySha(revision string) (error, *github.RepositoryCommit) {
	key := fmt.Sprintf("revision-%s", revision)

	commit, exist := cachedGithub.commitBySha[key]
	if exist {
		return nil, &commit
	}

	err, commitBySha := cachedGithub.GitClient.GetCommitBySha(revision)
	if err != nil {
		return err, nil
	}

	if len(cachedGithub.commitBySha) > cacheSize {
		cachedGithub.commitBySha = make(map[string]github.RepositoryCommit)
	}

	cachedGithub.commitBySha[key] = *commitBySha

	return nil, commitBySha
}

func (cachedGithub *CachedGithub) GetCommitsBySha(revision string) (error, []*github.RepositoryCommit) {
	key := fmt.Sprintf("revision-commits-%s", revision)

	cachedCommits, exist := cachedGithub.commitsBySha[key]
	if exist {
		result := make([]*github.RepositoryCommit, 0)

		for _, commit := range cachedCommits {
			result = append(result, &commit)
		}

		return nil, result
	}

	err, commits := cachedGithub.GitClient.GetCommitsBySha(revision)
	if err != nil {
		return err, nil
	}

	if len(cachedGithub.commitsBySha) > cacheSize {
		cachedGithub.commitsBySha = make(map[string][]github.RepositoryCommit)
	}

	result := make([]github.RepositoryCommit, 0)

	for _, commit := range commits {
		result = append(result, *commit)
	}

	cachedGithub.commitsBySha[key] = result

	return nil, commits
}

func (cachedGithub *CachedGithub) GetUserByUsername(username string) (error, *github.User) {
	key := fmt.Sprintf("user-%s", username)

	cachedUser, exist := cachedGithub.userByUsername[key]
	if exist {
		return nil, &cachedUser
	}

	err, user := cachedGithub.GitClient.GetUserByUsername(username)
	if err != nil {
		return err, nil
	}

	if len(cachedGithub.userByUsername) > cacheSize {
		cachedGithub.userByUsername = make(map[string]github.User)
	}

	cachedGithub.userByUsername[key] = *user

	return nil, user
}
