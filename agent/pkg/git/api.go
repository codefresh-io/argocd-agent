package git

import (
	"context"
	"github.com/codefresh-io/argocd-listener/agent/pkg/store"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
	"net/url"
	"strings"
)

type Api struct {
	Token  string
	Client *github.Client //github.Client
	Owner  string
	Repo   string
	Ctx    context.Context
}

var api *Api

func GetInstance(repoUrl string) (error, *Api) {
	if api != nil {
		return nil, api
	}
	gitConfig := store.GetStore().Git
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: gitConfig.Token},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)
	err, owner, repo  := _extractRepoAndOwnerFromUrl(repoUrl)
	if err != nil {
		return err, nil
	}

	api = &Api{
		Token:  gitConfig.Token,
		Ctx:    ctx,
		Client: client,
		Owner:  owner,
		Repo:   repo,
	}
	return nil, api
}

func _extractRepoAndOwnerFromUrl(repoUrl string) (error, string, string) {
	u, err  := url.Parse(repoUrl)
	if err != nil {
		return err, "", ""
	}
	urlParts := strings.Split(u.Path, "/")
	return nil, urlParts[len(urlParts)-2], urlParts[len(urlParts)-1]
}


func (a *Api) GetCommitsBySha(sha string) (error, []*github.RepositoryCommit) {
	revisionCommit, _, err := api.Client.Repositories.GetCommit(api.Ctx, api.Owner, api.Repo, sha)
	if err != nil {
		return err, nil
	}
	commits := []*github.RepositoryCommit{revisionCommit}
	if len(revisionCommit.Parents) > 0 {
		for i := 0; i < len(revisionCommit.Parents); i++ {
			commitInfo, _, err := api.Client.Repositories.GetCommit(api.Ctx, api.Owner, api.Repo, *revisionCommit.Parents[i].SHA)
			commits = append(commits, commitInfo)
			if err != nil {
				return err, nil
			}
		}
	}

	return nil, commits
}

func (a *Api) GetComittersByCommits(commits []*github.RepositoryCommit) (error, []User) {
	comitters := []User{}
	comittersSet := make(map[string]bool)
	for _, commit := range commits {
		author := commit.Author
		if author == nil {
			continue
		}
		_, exists := comittersSet[*author.Login]
		if exists != true {
			comittersSet[*author.Login] = true
			comitters = append(comitters, User{
				Name:   *author.Login,
				Avatar: *author.AvatarURL,
			})
		}
	}

	return nil, comitters
}

func (a *Api) GetIssuesAndPrsByCommits(commits []*github.RepositoryCommit) (error, []Annotation, []Annotation) {
	allPullRequests, _, err := api.Client.PullRequests.List(api.Ctx, api.Owner, api.Repo, &github.PullRequestListOptions{State: "all"})
	if err != nil {
		return err, nil, nil
	}

	issues := []Annotation{}
	pullRequests := []Annotation{}

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

				pullRequests = append(pullRequests, Annotation{
					Key:   *pr.Title,
					Value: *pr.URL,
				})
				issues = append(issues, Annotation{
					Key:   *issue.Title,
					Value: *issue.URL,
				})
			}
		}
	}
	return nil, issues, pullRequests
}

