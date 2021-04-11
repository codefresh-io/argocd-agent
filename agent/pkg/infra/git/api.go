package git

import (
	"context"
	"github.com/codefresh-io/argocd-listener/agent/pkg/infra/store"
	codefreshSdk "github.com/codefresh-io/go-sdk/pkg/codefresh"
	"github.com/google/go-github/github"
	"github.com/whilp/git-urls"
	"golang.org/x/oauth2"
	"regexp"
	"strings"
)

type Api struct {
	Token  string
	Client *github.Client
	Owner  string
	Repo   string
	Ctx    context.Context
}

var api *Api

func GetInstance(repoUrl string) (error, *Api) {
	err, owner, repo := extractRepoAndOwnerFromUrl(repoUrl)
	if err != nil {
		return err, nil
	}
	if api != nil {
		api.Owner = owner
		api.Repo = repo
		return nil, api
	}
	gitConfig := store.GetStore().Git
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: gitConfig.Token},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	api = &Api{
		Token:  gitConfig.Token,
		Ctx:    ctx,
		Client: client,
		Owner:  owner,
		Repo:   repo,
	}
	return nil, api
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

func (a *Api) GetCommitBySha(sha string) (error, *github.RepositoryCommit) {
	revisionCommit, _, err := api.Client.Repositories.GetCommit(api.Ctx, api.Owner, api.Repo, sha)
	if err != nil {
		return err, nil
	}
	return nil, revisionCommit
}

func (a *Api) GetUserByUsername(username string) (error, *github.User) {
	user, _, err := api.Client.Users.Get(api.Ctx, username)
	if err != nil {
		return err, nil
	}
	return nil, user
}

func (a *Api) GetCommitsBySha(sha string) (error, []*github.RepositoryCommit) {
	revisionCommit, _, err := api.Client.Repositories.GetCommit(api.Ctx, api.Owner, api.Repo, sha)
	if err != nil {
		return err, nil
	}
	return nil, []*github.RepositoryCommit{revisionCommit}
}

func (a *Api) GetComittersByCommits(commits []*github.RepositoryCommit) (error, []codefreshSdk.User) {
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

func (a *Api) GetIssuesAndPrsByCommits(commits []*github.RepositoryCommit) (error, []codefreshSdk.Annotation, []codefreshSdk.Annotation) {
	allPullRequests, _, err := api.Client.PullRequests.List(api.Ctx, api.Owner, api.Repo, &github.PullRequestListOptions{State: "all"})
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
				issue, _, err := api.Client.Issues.Get(api.Ctx, api.Owner, api.Repo, *pr.Number)
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
