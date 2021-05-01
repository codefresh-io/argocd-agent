package git

import (
	"context"
	"encoding/base64"
	"errors"
	"github.com/bradleyfalzon/ghinstallation"
	"github.com/codefresh-io/argocd-listener/agent/pkg/infra/store"
	codefreshSdk "github.com/codefresh-io/go-sdk/pkg/codefresh"
	"github.com/google/go-github/github"
	"github.com/whilp/git-urls"
	"golang.org/x/oauth2"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

type api struct {
	Token  string
	Client *github.Client
	Owner  string
	Repo   string
	Ctx    context.Context
}

type Api interface {
	GetCommitBySha(sha string) (error, *github.RepositoryCommit)
	GetUserByUsername(username string) (error, *github.User)
	GetCommitsBySha(sha string) (error, []*github.RepositoryCommit)
	GetComittersByCommits(commits []*github.RepositoryCommit) (error, []codefreshSdk.User)
	GetIssuesAndPrsByCommits(commits []*github.RepositoryCommit) (error, []codefreshSdk.Annotation, []codefreshSdk.Annotation)
}

var githubApi *api

func GetInstance(repoUrl string) (error, Api) {
	err, owner, repo := extractRepoAndOwnerFromUrl(repoUrl)
	if err != nil {
		return err, nil
	}
	if githubApi != nil {
		githubApi.Owner = owner
		githubApi.Repo = repo
		return nil, githubApi
	}
	gitConfig := store.GetStore().Git
	err, ctx := getGitHttpClient(&store.GetStore().Git.Context)
	if err != nil {
		return err, nil
	}
	client := github.NewClient(ctx)

	githubApi = &api{
		Token:  gitConfig.Token,
		Ctx:    context.Background(),
		Client: client,
		Owner:  owner,
		Repo:   repo,
	}
	return nil, githubApi
}

func getGitHttpClient(gitContext *codefreshSdk.ContextPayload) (error, *http.Client) {
	auth := gitContext.Spec.Data.Auth
	if gitContext.Spec.Type == "git.github" {
		ctx := context.Background()
		ts := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: gitContext.Spec.Data.Auth.Password},
		)
		return nil, oauth2.NewClient(ctx, ts)
	}
	if gitContext.Spec.Type == "git.github-app" {
		appId, err := strconv.Atoi(auth.AppId)
		if err != nil {
			return err, nil
		}
		installationId, err := strconv.Atoi(auth.InstallationId)
		if err != nil {
			return err, nil
		}
		key, err := base64.URLEncoding.DecodeString(auth.PrivateKey)
		if err != nil {
			return err, nil
		}
		transport, err := ghinstallation.New(http.DefaultTransport, int64(appId), int64(installationId), key)
		if err != nil {
			return err, nil
		}
		return nil, &http.Client{Transport: transport}
	}
	return errors.New("Cant handle unknown git type"), nil
}

func extractRepoAndOwnerFromUrl(repoUrl string) (error, string, string) {
	u, err := giturls.Parse(repoUrl)
	if err != nil {
		return err, "", ""
	}

	// from url like this -> https://github.com/codefresh-io/argocd-agent.git/
	// to array like this -> string[]{"codefresh-io","argocd-agent.git",""}
	urlParts := strings.Split(u.Path, "/")
	filteredUrlParts := []string{}

	// removing all empty strings from array
	for _, part := range urlParts {
		if part != "" {
			filteredUrlParts = append(filteredUrlParts, part)
		}
	}
	if len(filteredUrlParts) > 1 {
		var re = regexp.MustCompile(`\.git$`)

		// getting the last element from array like -> string[]{"codefresh-io","argocd-agent.git"}
		// and removing unnecessary part of string
		// result will be "argocd-agent"
		repo := re.ReplaceAllString(filteredUrlParts[len(filteredUrlParts)-1], "")

		//getting the penultimate element from array like -> string[]{"codefresh-io","argocd-agent.git"}
		// result will be "codefresh-io"
		owner := filteredUrlParts[len(filteredUrlParts)-2]
		return nil, owner, repo
	}
	return nil, "", ""
}

func (a *api) GetCommitBySha(sha string) (error, *github.RepositoryCommit) {
	revisionCommit, _, err := a.Client.Repositories.GetCommit(a.Ctx, a.Owner, a.Repo, sha)
	if err != nil {
		return err, nil
	}
	return nil, revisionCommit
}

func (a *api) GetUserByUsername(username string) (error, *github.User) {
	user, _, err := a.Client.Users.Get(a.Ctx, username)
	if err != nil {
		return err, nil
	}
	return nil, user
}

func (a *api) GetCommitsBySha(sha string) (error, []*github.RepositoryCommit) {
	revisionCommit, _, err := a.Client.Repositories.GetCommit(a.Ctx, a.Owner, a.Repo, sha)
	if err != nil {
		return err, nil
	}
	return nil, []*github.RepositoryCommit{revisionCommit}
}

func (a *api) GetComittersByCommits(commits []*github.RepositoryCommit) (error, []codefreshSdk.User) {
	comitters := []codefreshSdk.User{}
	comittersSet := make(map[string]bool)
	for _, commit := range commits {
		author := commit.Author
		if author == nil {
			continue
		}
		_, exists := comittersSet[*author.Login]
		if exists != true {
			comittersSet[*author.Login] = true
			comitters = append(comitters, codefreshSdk.User{
				Name:   *author.Login,
				Avatar: *author.AvatarURL,
			})
		}
	}

	return nil, comitters
}

func (a *api) GetIssuesAndPrsByCommits(commits []*github.RepositoryCommit) (error, []codefreshSdk.Annotation, []codefreshSdk.Annotation) {
	allPullRequests, _, err := a.Client.PullRequests.List(a.Ctx, a.Owner, a.Repo, &github.PullRequestListOptions{State: "all"})
	if err != nil {
		return err, nil, nil
	}

	issues := []codefreshSdk.Annotation{}
	pullRequests := []codefreshSdk.Annotation{}

	for _, pr := range allPullRequests {
		mergeCommitSHA := pr.MergeCommitSHA
		if mergeCommitSHA == nil {
			continue
		}
		for _, commit := range commits {
			if commit.SHA == nil {
				continue
			}
			if *commit.SHA == *mergeCommitSHA {
				issue, _, err := a.Client.Issues.Get(a.Ctx, a.Owner, a.Repo, *pr.Number)
				if err != nil {
					return err, nil, nil
				}

				pullRequests = append(pullRequests, codefreshSdk.Annotation{
					Key:   *pr.Title,
					Value: *pr.HTMLURL,
				})
				issues = append(issues, codefreshSdk.Annotation{
					Key:   *issue.Title,
					Value: *issue.HTMLURL,
				})
			}
		}
	}
	return nil, issues, pullRequests
}
